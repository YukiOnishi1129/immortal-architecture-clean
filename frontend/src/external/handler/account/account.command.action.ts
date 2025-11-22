"use server";

import type {
  CreateOrGetAccountRequest,
  CreateOrGetAccountResponse,
  UpdateAccountRequest,
  UpdateAccountResponse,
} from "../../dto/account.dto";
import {
  createOrGetAccountCommand,
  updateAccountCommand,
} from "./account.command.server";

export async function createOrGetAccountCommandAction(
  request: CreateOrGetAccountRequest,
): Promise<CreateOrGetAccountResponse> {
  return createOrGetAccountCommand(request);
}

export async function updateAccountCommandAction(
  id: string,
  request: UpdateAccountRequest,
): Promise<UpdateAccountResponse> {
  return updateAccountCommand(id, request);
}
