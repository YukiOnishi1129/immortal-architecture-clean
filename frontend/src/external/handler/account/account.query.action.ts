"use server";

import {
  getAccountByIdQuery,
  getCurrentAccountQuery,
} from "./account.query.server";

export async function getCurrentAccountQueryAction() {
  return getCurrentAccountQuery();
}

export async function getAccountByIdQueryAction(id: string) {
  return getAccountByIdQuery(id);
}
