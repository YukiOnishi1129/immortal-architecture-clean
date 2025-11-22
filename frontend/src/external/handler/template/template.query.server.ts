import "server-only";

import { getSessionServer } from "@/features/auth/servers/auth.server";
import {
  getAuthenticatedSessionServer,
  requireAuthServer,
} from "@/features/auth/servers/redirect.server";
import type { TemplateFilters } from "@/features/template/types";
import {
  TemplateDetailResponseSchema,
  TemplateResponseSchema,
} from "../../dto/template.dto";
import { templateService } from "../../service/template/template.service";

export async function getTemplateByIdQuery(id: string) {
  const template = await templateService.getTemplateById(id);

  if (!template) {
    return null;
  }

  return TemplateDetailResponseSchema.parse(template);
}

export async function listTemplatesQuery(filters?: TemplateFilters) {
  await requireAuthServer();

  // Get current user for onlyMyTemplates filter
  const session = await getSessionServer();

  const ownerFilter =
    filters?.onlyMyTemplates && session?.account.id
      ? session.account.id
      : filters?.ownerId;

  const templates = await templateService.getTemplates({
    ownerId: ownerFilter,
    q: filters?.q,
  });

  return templates.map((template) => TemplateResponseSchema.parse(template));
}

export async function listMyTemplatesQuery() {
  const session = await getAuthenticatedSessionServer();

  const templates = await templateService.getTemplates({
    ownerId: session.account.id,
  });

  return templates.map((template) => TemplateResponseSchema.parse(template));
}
