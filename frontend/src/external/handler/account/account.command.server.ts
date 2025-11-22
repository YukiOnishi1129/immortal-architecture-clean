import "server-only";

import { getAuthenticatedSessionServer } from "@/features/auth/servers/redirect.server";
import {
  type CreateOrGetAccountRequest,
  CreateOrGetAccountRequestSchema,
  type CreateOrGetAccountResponse,
  type UpdateAccountRequest,
  UpdateAccountRequestSchema,
  type UpdateAccountResponse,
} from "../../dto/account.dto";
import { accountService } from "../../service/account/account.service";

export async function createOrGetAccountCommand(
  request: CreateOrGetAccountRequest,
): Promise<CreateOrGetAccountResponse> {
  const validated = CreateOrGetAccountRequestSchema.parse(request);

  return accountService.createOrGet(validated);
}

export async function updateAccountCommand(
  id: string,
  request: UpdateAccountRequest,
): Promise<UpdateAccountResponse> {
  const session = await getAuthenticatedSessionServer();

  // Check if the user is updating their own account
  if (session.account.id !== id) {
    throw new Error("Forbidden: Can only update your own account");
  }

  const validated = UpdateAccountRequestSchema.parse(request);

  return accountService.update(id, validated);
}
