import type { Account } from "@/features/accounts/types";

declare module "next-auth" {
  interface Session {
    account?: Account;
    error?: "RefreshTokenMissing" | "RefreshAccessTokenError";
  }

  interface User {
    id: string;
    account?: Account;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    account?: Account;
    accessToken?: string;
    refreshToken?: string;
    accessTokenExpires?: number;
    idToken?: string;
    error?: "RefreshTokenMissing" | "RefreshAccessTokenError";
  }
}

// Google OAuth profile type
export interface GoogleProfile {
  sub: string;
  name: string;
  given_name: string;
  family_name: string;
  picture?: string;
  email: string;
  email_verified: boolean;
  locale?: string;
}
