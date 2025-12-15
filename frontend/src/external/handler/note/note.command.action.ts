"use server";

import type {
  CreateNoteRequest,
  DeleteNoteRequest,
  PublishNoteRequest,
  UnpublishNoteRequest,
  UpdateNoteByIdRequest,
} from "../../dto/note.dto";
import {
  createNoteCommand,
  deleteNoteCommand,
  publishNoteCommand,
  unpublishNoteCommand,
  updateNoteCommand,
} from "./note.command.server";

export async function createNoteCommandAction(request: CreateNoteRequest) {
  return createNoteCommand(request);
}

export async function updateNoteCommandAction(request: UpdateNoteByIdRequest) {
  return updateNoteCommand(request);
}

export async function publishNoteCommandAction(request: PublishNoteRequest) {
  return publishNoteCommand(request);
}

export async function unpublishNoteCommandAction(
  request: UnpublishNoteRequest,
) {
  return unpublishNoteCommand(request);
}

export async function deleteNoteCommandAction(request: DeleteNoteRequest) {
  return deleteNoteCommand(request);
}
