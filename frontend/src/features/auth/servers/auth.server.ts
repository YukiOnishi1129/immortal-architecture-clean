import "server-only";
import { getServerSession } from "next-auth/next";
import { authOptions } from "../lib/options";

export async function getSessionServer() {
  return getServerSession(authOptions);
}
