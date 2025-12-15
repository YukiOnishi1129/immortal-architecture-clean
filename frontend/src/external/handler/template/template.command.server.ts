import "server-only";

import { withAuth } from "@/features/auth/servers/auth.guard";
import {
  type CreateTemplateRequest,
  CreateTemplateRequestSchema,
  type DeleteTemplateRequest,
  DeleteTemplateRequestSchema,
  TemplateResponseSchema,
  type UpdateTemplateByIdRequest,
  UpdateTemplateByIdRequestSchema,
} from "../../dto/template.dto";
import { templateService } from "../../service/template/template.service";

export async function createTemplateCommand(request: CreateTemplateRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = CreateTemplateRequestSchema.parse(request);
    const template = await templateService.createTemplate(accountId, validated);
    return TemplateResponseSchema.parse(template);
  });
}

export async function updateTemplateCommand(
  request: UpdateTemplateByIdRequest,
) {
  return withAuth(async ({ accountId }) => {
    try {
      const validated = UpdateTemplateByIdRequestSchema.parse(request);
      const { id, ...updateData } = validated;
      const template = await templateService.updateTemplate(
        id,
        accountId,
        updateData,
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
  });
}

export async function deleteTemplateCommand(request: DeleteTemplateRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = DeleteTemplateRequestSchema.parse(request);
    await templateService.deleteTemplate(validated.id, accountId);
    return { success: true };
  });
}
