"use server";

import type { NoteFilters } from "@/features/note/types";
import {
  getNoteByIdQuery,
  listMyNoteQuery,
  listNoteQuery,
} from "./note.query.server";

export async function getNoteByIdQueryAction(id: string) {
  return getNoteByIdQuery(id);
}

export async function listNoteQueryAction(filters?: NoteFilters) {
  return listNoteQuery(filters);
}

export async function listMyNoteQueryAction(filters?: NoteFilters) {
  return listMyNoteQuery(filters);
}
