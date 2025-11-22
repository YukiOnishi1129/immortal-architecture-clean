"use server";

import {
  createNoteCommand,
  deleteNoteCommand,
  publishNoteCommand,
  unpublishNoteCommand,
  updateNoteCommand,
} from "./note.command.server";

export async function createNoteCommandAction(request: unknown) {
  return createNoteCommand(request);
}

export async function updateNoteCommandAction(id: string, request: unknown) {
  return updateNoteCommand(id, request);
}

export async function publishNoteCommandAction(request: unknown) {
  return publishNoteCommand(request);
}

export async function unpublishNoteCommandAction(request: unknown) {
  return unpublishNoteCommand(request);
}

export async function deleteNoteCommandAction(id: string) {
  return deleteNoteCommand(id);
}
