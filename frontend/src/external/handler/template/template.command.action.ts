"use server";

import type {
  CreateTemplateRequest,
  DeleteTemplateRequest,
  UpdateTemplateByIdRequest,
} from "@/external/dto/template.dto";
import {
  createTemplateCommand,
  deleteTemplateCommand,
  updateTemplateCommand,
} from "./template.command.server";

export async function createTemplateCommandAction(
  request: CreateTemplateRequest,
) {
  return createTemplateCommand(request);
}

export async function updateTemplateCommandAction(
  request: UpdateTemplateByIdRequest,
) {
  return updateTemplateCommand(request);
}

export async function deleteTemplateCommandAction(
  request: DeleteTemplateRequest,
) {
  return deleteTemplateCommand(request);
}
