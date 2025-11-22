import "server-only";

import { getAuthenticatedSessionServer } from "@/features/auth/servers/redirect.server";
import {
  CreateNoteRequestSchema,
  NoteResponseSchema,
  PublishNoteRequestSchema,
  UnpublishNoteRequestSchema,
  UpdateNoteRequestSchema,
} from "../../dto/note.dto";
import { noteService } from "../../service/note/note.service";

export async function createNoteCommand(request: unknown) {
  const session = await getAuthenticatedSessionServer();

  // リクエストのバリデーション
  const validated = CreateNoteRequestSchema.parse(request);

  const note = await noteService.createNote(session.account.id, validated);
  return NoteResponseSchema.parse(note);
}

export async function updateNoteCommand(id: string, request: unknown) {
  const _session = await getAuthenticatedSessionServer();

  // リクエストのバリデーション
  const validated = UpdateNoteRequestSchema.parse(request);

  const note = await noteService.updateNote(id, validated);
  return NoteResponseSchema.parse(note);
}

export async function publishNoteCommand(request: unknown) {
  const session = await getAuthenticatedSessionServer();

  // Validate request
  const validated = PublishNoteRequestSchema.parse(request);

  if (!session?.account?.id) {
    throw new Error("Account not found");
  }

  const note = await noteService.publishNote(validated.noteId);
  return NoteResponseSchema.parse(note);
}

export async function unpublishNoteCommand(request: unknown) {
  const session = await getAuthenticatedSessionServer();

  // Validate request
  const validated = UnpublishNoteRequestSchema.parse(request);

  if (!session?.account?.id) {
    throw new Error("Account not found");
  }

  const note = await noteService.unpublishNote(validated.noteId);
  return NoteResponseSchema.parse(note);
}

export async function deleteNoteCommand(id: string) {
  const session = await getAuthenticatedSessionServer();

  if (!session?.account?.id) {
    throw new Error("Account not found");
  }

  await noteService.deleteNote(id);

  return { success: true };
}
