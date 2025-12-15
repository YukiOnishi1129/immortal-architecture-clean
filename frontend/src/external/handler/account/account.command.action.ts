"use server";

import type {
  CreateOrGetAccountRequest,
  CreateOrGetAccountResponse,
  UpdateAccountByIdRequest,
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
  request: UpdateAccountByIdRequest,
): Promise<UpdateAccountResponse> {
  return updateAccountCommand(request);
}
