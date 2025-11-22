import "server-only";
import {
  getAuthenticatedSessionServer,
  requireAuthServer,
} from "@/features/auth/servers/redirect.server";
import type { NoteFilters } from "@/features/note/types";
import { NoteResponseSchema } from "../../dto/note.dto";
import { noteService } from "../../service/note/note.service";

export async function getNoteByIdQuery(id: string) {
  await requireAuthServer();

  const note = await noteService.getNoteById(id);

  if (!note) {
    return null;
  }

  return NoteResponseSchema.parse(note);
}

export async function listNoteQuery(filters?: NoteFilters) {
  await requireAuthServer();

  const notes = await noteService.getNotes(filters);
  return notes.map((note) => NoteResponseSchema.parse(note));
}

export async function listMyNoteQuery(filters?: NoteFilters) {
  const session = await getAuthenticatedSessionServer();

  const notes = await noteService.getNotes({
    ...filters,
    ownerId: session.account.id,
  });

  return notes.map((note) => NoteResponseSchema.parse(note));
}
