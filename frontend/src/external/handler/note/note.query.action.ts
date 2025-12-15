"use server";

import type {
  GetNoteByIdRequest,
  ListMyNoteRequest,
  ListNoteRequest,
} from "../../dto/note.dto";
import {
  getNoteByIdQuery,
  listMyNoteQuery,
  listNoteQuery,
} from "./note.query.server";

export async function getNoteByIdQueryAction(request: GetNoteByIdRequest) {
  return getNoteByIdQuery(request);
}

export async function listNoteQueryAction(request?: ListNoteRequest) {
  return listNoteQuery(request);
}

export async function listMyNoteQueryAction(request?: ListMyNoteRequest) {
  return listMyNoteQuery(request);
}
