import "server-only";
import { getAuthenticatedSessionServer } from "@/features/auth/servers/redirect.server";
import {
  CreateTemplateRequestSchema,
  TemplateResponseSchema,
  UpdateTemplateRequestSchema,
} from "../../dto/template.dto";
import { templateService } from "../../service/template/template.service";

export async function createTemplateCommand(request: unknown) {
  const session = await getAuthenticatedSessionServer();

  // Validate request
  const validated = CreateTemplateRequestSchema.parse(request);

  const template = await templateService.createTemplate(
    session.account.id,
    validated,
  );

  return TemplateResponseSchema.parse(template);
}

export async function updateTemplateCommand(id: string, request: unknown) {
  const session = await getAuthenticatedSessionServer();

  try {
    // Validate request
    const validated = UpdateTemplateRequestSchema.parse(request);

    // Update template
    const template = await templateService.updateTemplate(
      id,
      session.account.id,
      validated,
    );
    return TemplateResponseSchema.parse(template);
  } catch (error) {
    if (error instanceof Error) {
      if (error.message === "TEMPLATE_FIELD_IN_USE") {
        throw new Error(
          "テンプレートの項目は変更・削除できません。ノートで使用されています。",
        );
      }
      if (error.message === "TEMPLATE_STRUCTURE_LOCKED") {
        throw new Error(
          "テンプレートの項目は変更・削除できません。ノートで使用されています。",
        );
      }
    }
    throw error;
  }
}

export async function deleteTemplateCommand(id: string) {
  const session = await getAuthenticatedSessionServer();

  // Delete template
  await templateService.deleteTemplate(id, session.account.id);
  return { success: true };
}
