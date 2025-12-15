import "server-only";

import { getAuthenticatedSessionServer } from "@/features/auth/servers/redirect.server";

export async function withAuth<T>(
  handler: (ctx: { accountId: string }) => Promise<T>,
): Promise<T> {
  const session = await getAuthenticatedSessionServer();
  if (!session.account) {
    throw new Error("Unexpected: account should exist after authentication");
  }
  return handler({ accountId: session.account.id });
}
