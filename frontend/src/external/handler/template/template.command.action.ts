"use server";

import type {
  CreateTemplateRequest,
  UpdateTemplateRequest,
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
  id: string,
  request: UpdateTemplateRequest,
) {
  return updateTemplateCommand(id, request);
}

export async function deleteTemplateCommandAction(id: string) {
  return deleteTemplateCommand(id);
}
