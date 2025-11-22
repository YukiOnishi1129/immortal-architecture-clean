import NextAuth from "next-auth";

import { authOptions } from "@/features/auth/lib/options";

export const handler = NextAuth(authOptions);
