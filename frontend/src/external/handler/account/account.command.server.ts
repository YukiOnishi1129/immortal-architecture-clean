import "server-only";

import { withAuth } from "@/features/auth/servers/auth.guard";
import {
  type CreateOrGetAccountRequest,
  CreateOrGetAccountRequestSchema,
  type CreateOrGetAccountResponse,
  type UpdateAccountByIdRequest,
  UpdateAccountByIdRequestSchema,
  type UpdateAccountResponse,
} from "../../dto/account.dto";
import { accountService } from "../../service/account/account.service";

// NOTE: createOrGetAccountCommandはOAuth認証時に呼ばれるため、withAuthを適用しない
export async function createOrGetAccountCommand(
  request: CreateOrGetAccountRequest,
): Promise<CreateOrGetAccountResponse> {
  const validated = CreateOrGetAccountRequestSchema.parse(request);
  return accountService.createOrGet(validated);
}

export async function updateAccountCommand(
  request: UpdateAccountByIdRequest,
): Promise<UpdateAccountResponse> {
  return withAuth(async ({ accountId }) => {
    const validated = UpdateAccountByIdRequestSchema.parse(request);

    // Check if the user is updating their own account
    if (accountId !== validated.id) {
      throw new Error("Forbidden: Can only update your own account");
    }

    const { id, ...updateData } = validated;
    return accountService.update(id, updateData);
  });
}
