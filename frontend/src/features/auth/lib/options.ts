import type { NextAuthOptions } from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import { createOrGetAccountCommand } from "@/external/handler/account/account.command.server";
import { refreshGoogleTokenCommand } from "@/external/handler/auth/token.command.server";
import type { Account } from "@/features/account/types";
import type { GoogleProfile } from "@/features/auth/types/next-auth";

function buildProfileName(profile: GoogleProfile): string {
  if (profile.name) return profile.name;
  return `${profile.given_name || ""} ${profile.family_name || ""}`.trim();
}

export const authOptions: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID || "",
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || "",
      authorization: {
        params: {
          access_type: "offline",
          prompt: "consent",
          scope: "openid email profile",
        },
      },
      profile(profile) {
        return {
          id: profile.sub,
        };
      },
    }),
  ],
  callbacks: {
    async signIn({ user, account, profile }): Promise<boolean> {
      if (!account || account.provider !== "google") return false;

      try {
        const googleProfile = profile as GoogleProfile;
        const fullName = buildProfileName(googleProfile);

        const accountData = await createOrGetAccountCommand({
          email: googleProfile.email,
          name: fullName,
          provider: account.provider,
          providerAccountId: account.providerAccountId,
          thumbnail: googleProfile.picture,
        });

        user.account = accountData;
        return true;
      } catch (error) {
        console.error("Error during sign in:", error);
        return false;
      }
    },
    async jwt({ token, user, account }) {
      if (user?.account) {
        token.account = user.account as Account;
      }

      if (account) {
        token.accessToken = account.access_token ?? token.accessToken;
        token.refreshToken = account.refresh_token ?? token.refreshToken;
        token.idToken = account.id_token ?? token.idToken;
        token.accessTokenExpires = account.expires_at
          ? account.expires_at * 1000
          : Date.now() + 1000 * 60 * 60;
        token.error = undefined;
        return token;
      }

      if (
        token.accessToken &&
        typeof token.accessTokenExpires === "number" &&
        token.accessTokenExpires > Date.now() + 60_000
      ) {
        return token;
      }

      if (!token.refreshToken) {
        return {
          ...token,
          accessToken: undefined,
          idToken: undefined,
          error: "RefreshTokenMissing",
        };
      }

      try {
        const refreshed = await refreshGoogleTokenCommand({
          refreshToken: token.refreshToken,
        });

        return {
          ...token,
          accessToken: refreshed.accessToken ?? token.accessToken,
          idToken: refreshed.idToken ?? token.idToken,
          accessTokenExpires:
            refreshed.accessTokenExpires ?? Date.now() + 1000 * 60 * 60,
          error: undefined,
        };
      } catch (error) {
        console.error("Error refreshing access token:", error);
        return {
          ...token,
          accessToken: undefined,
          idToken: undefined,
          error: "RefreshAccessTokenError",
        };
      }
    },
    async session({ session, token }) {
      if (token.account) {
        session.account = token.account as Account;

        // アカウント情報から標準的なsession.userフィールドを設定
        session.user = {
          name: `${token.account.firstName} ${token.account.lastName}`,
          email: token.account.email,
          image: token.account.thumbnail ?? null,
        };
      }

      if (token.error) {
        session.error = token.error;
      }

      return session;
    },
  },
  pages: {
    signIn: "/login",
  },
  session: {
    strategy: "jwt",
  },
};
