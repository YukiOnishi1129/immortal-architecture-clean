import "server-only";

import { withAuth } from "@/features/auth/servers/auth.guard";
import {
  type CreateNoteRequest,
  CreateNoteRequestSchema,
  type DeleteNoteRequest,
  DeleteNoteRequestSchema,
  NoteResponseSchema,
  type PublishNoteRequest,
  PublishNoteRequestSchema,
  type UnpublishNoteRequest,
  UnpublishNoteRequestSchema,
  type UpdateNoteByIdRequest,
  UpdateNoteByIdRequestSchema,
} from "../../dto/note.dto";
import { noteService } from "../../service/note/note.service";

export async function createNoteCommand(request: CreateNoteRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = CreateNoteRequestSchema.parse(request);
    const note = await noteService.createNote(accountId, validated);
    return NoteResponseSchema.parse(note);
  });
}

export async function updateNoteCommand(request: UpdateNoteByIdRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = UpdateNoteByIdRequestSchema.parse(request);
    const { id, ...updateData } = validated;
    const note = await noteService.updateNote(id, accountId, updateData);
    return NoteResponseSchema.parse(note);
  });
}

export async function publishNoteCommand(request: PublishNoteRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = PublishNoteRequestSchema.parse(request);
    const note = await noteService.publishNote(validated.noteId, accountId);
    return NoteResponseSchema.parse(note);
  });
}

export async function unpublishNoteCommand(request: UnpublishNoteRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = UnpublishNoteRequestSchema.parse(request);
    const note = await noteService.unpublishNote(validated.noteId, accountId);
    return NoteResponseSchema.parse(note);
  });
}

export async function deleteNoteCommand(request: DeleteNoteRequest) {
  return withAuth(async ({ accountId }) => {
    const validated = DeleteNoteRequestSchema.parse(request);
    await noteService.deleteNote(validated.id, accountId);
    return { success: true };
  });
}
