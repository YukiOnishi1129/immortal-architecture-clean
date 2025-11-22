"use server";

import type { TemplateFilters } from "@/features/template/types";
import {
  getTemplateByIdQuery,
  listMyTemplatesQuery,
  listTemplatesQuery,
} from "./template.query.server";

export async function getTemplateByIdQueryAction(id: string) {
  return getTemplateByIdQuery(id);
}

export async function listTemplatesQueryAction(filters?: TemplateFilters) {
  return listTemplatesQuery(filters);
}

export async function listMyTemplatesQueryAction() {
  return listMyTemplatesQuery();
}
