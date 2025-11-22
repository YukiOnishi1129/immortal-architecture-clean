package controller

import (
	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
)

// Server implements the OpenAPI ServerInterface by delegating to domain-specific controllers.
type Server struct {
	account  *AccountController
	note     *NoteController
	template *TemplateController
}

func NewServer(ac *AccountController, nc *NoteController, tc *TemplateController) *Server {
	return &Server{account: ac, note: nc, template: tc}
}

// OpenAPI generated interface methods
func (s *Server) AccountsCreateOrGetAccount(ctx echo.Context) error {
	return s.account.CreateOrGet(ctx)
}

func (s *Server) AccountsGetCurrentAccount(ctx echo.Context) error {
	return s.account.GetCurrent(ctx)
}

func (s *Server) AccountsGetAccountById(ctx echo.Context, accountId string) error {
	return s.account.GetByID(ctx, accountId)
}

func (s *Server) NotesListNotes(ctx echo.Context, params openapi.NotesListNotesParams) error {
	return s.note.List(ctx, params)
}

func (s *Server) NotesCreateNote(ctx echo.Context) error {
	return s.note.Create(ctx)
}

func (s *Server) NotesDeleteNote(ctx echo.Context, noteId string) error {
	return s.note.Delete(ctx, noteId)
}

func (s *Server) NotesGetNoteById(ctx echo.Context, noteId string) error {
	return s.note.GetByID(ctx, noteId)
}

func (s *Server) NotesUpdateNote(ctx echo.Context, noteId string) error {
	return s.note.Update(ctx, noteId)
}

func (s *Server) NotesPublishNote(ctx echo.Context, noteId string) error {
	return s.note.Publish(ctx, noteId)
}

func (s *Server) NotesUnpublishNote(ctx echo.Context, noteId string) error {
	return s.note.Unpublish(ctx, noteId)
}

func (s *Server) TemplatesListTemplates(ctx echo.Context, params openapi.TemplatesListTemplatesParams) error {
	return s.template.List(ctx, params)
}

func (s *Server) TemplatesCreateTemplate(ctx echo.Context) error {
	return s.template.Create(ctx)
}

func (s *Server) TemplatesDeleteTemplate(ctx echo.Context, templateId string) error {
	return s.template.Delete(ctx, templateId)
}

func (s *Server) TemplatesGetTemplateById(ctx echo.Context, templateId string) error {
	return s.template.GetByID(ctx, templateId)
}

func (s *Server) TemplatesUpdateTemplate(ctx echo.Context, templateId string) error {
	return s.template.Update(ctx, templateId)
}
