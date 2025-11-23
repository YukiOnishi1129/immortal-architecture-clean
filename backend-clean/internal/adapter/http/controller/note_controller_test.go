package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	ctrlmock "immortal-architecture-clean/backend/internal/adapter/http/controller/mock"
	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/adapter/http/presenter"
	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/port"
)

func TestNoteController_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "[Success] create note",
			body:       `{"title":"Hello","templateId":"00000000-0000-0000-0000-000000000001","ownerId":"00000000-0000-0000-0000-000000000002","sections":[{"fieldId":"f1","content":"c1"}]}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "[Fail] bind error",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewNotePresenter()
			input := &ctrlmock.NoteInputStub{}
			ctrl := NewNoteController(
				func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
					input.Output = output
					return input
				},
				func() *presenter.NotePresenter { return p },
				func() port.NoteRepository { return nil },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.Create(c)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestNoteController_List(t *testing.T) {
	tests := []struct {
		name       string
		filters    openapi.NotesListNotesParams
		inErr      error
		wantStatus int
	}{
		{name: "[Success] list notes", filters: openapi.NotesListNotesParams{}, wantStatus: http.StatusOK},
		{name: "[Fail] repo error", filters: openapi.NotesListNotesParams{}, inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewNotePresenter()
			input := &ctrlmock.NoteInputStub{Notes: []note.WithMeta{{Note: note.Note{ID: "n1"}}}, Err: tt.inErr}
			ctrl := NewNoteController(
				func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
					input.Output = output
					return input
				},
				func() *presenter.NotePresenter { return p },
				func() port.NoteRepository { return nil },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodGet, "/api/notes", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.List(c, tt.filters)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestNoteController_Get(t *testing.T) {
	tests := []struct {
		name       string
		inErr      error
		wantStatus int
	}{
		{name: "[Success] get note", wantStatus: http.StatusOK},
		{name: "[Fail] not found", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewNotePresenter()
			input := &ctrlmock.NoteInputStub{Err: tt.inErr}
			ctrl := NewNoteController(
				func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
					input.Output = output
					return input
				},
				func() *presenter.NotePresenter { return p },
				func() port.NoteRepository { return nil },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodGet, "/api/notes/n1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.GetByID(c, "n1")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestNoteController_Update(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		params     openapi.NotesUpdateNoteParams
		inErr      error
		wantStatus int
	}{
		{name: "[Success] update note", body: `{"title":"New","sections":[{"id":"sec1","content":"c"}]}`, params: openapi.NotesUpdateNoteParams{OwnerId: "owner"}, wantStatus: http.StatusOK},
		{name: "[Fail] missing owner", body: `{"title":"New","sections":[{"id":"sec1","content":"c"}]}`, params: openapi.NotesUpdateNoteParams{OwnerId: ""}, wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewNotePresenter()
			input := &ctrlmock.NoteInputStub{Err: tt.inErr}
			ctrl := NewNoteController(
				func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
					input.Output = output
					return input
				},
				func() *presenter.NotePresenter { return p },
				func() port.NoteRepository { return nil },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodPut, "/api/notes/n1", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Update(c, "n1", tt.params)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestNoteController_Publish(t *testing.T) {
	tests := []struct {
		name       string
		ownerID    string
		wantStatus int
	}{
		{name: "[Success] publish note", ownerID: "owner", wantStatus: http.StatusOK},
		{name: "[Fail] publish missing owner", ownerID: "", wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewNotePresenter()
			input := &ctrlmock.NoteInputStub{}
			ctrl := NewNoteController(
				func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
					input.Output = output
					return input
				},
				func() *presenter.NotePresenter { return p },
				func() port.NoteRepository { return nil },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodPost, "/api/notes/n1/publish", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Publish(c, "n1", openapi.NotesPublishNoteParams{OwnerId: tt.ownerID})
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestNoteController_Unpublish(t *testing.T) {
	e := echo.New()
	p := presenter.NewNotePresenter()
	input := &ctrlmock.NoteInputStub{}
	ctrl := NewNoteController(
		func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
			input.Output = output
			return input
		},
		func() *presenter.NotePresenter { return p },
		func() port.NoteRepository { return nil },
		func() port.TemplateRepository { return nil },
		func() port.TxManager { return nil },
	)
	req := httptest.NewRequest(http.MethodPost, "/api/notes/n1/unpublish", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = ctrl.Unpublish(c, "n1", openapi.NotesUnpublishNoteParams{OwnerId: "owner"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestNoteController_Delete(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name       string
		ownerID    string
		wantStatus int
	}{
		{name: "[Success] delete", ownerID: "owner", wantStatus: http.StatusOK},
		{name: "[Fail] owner missing", ownerID: "", wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewNotePresenter()
			input := &ctrlmock.NoteInputStub{}
			ctrl := NewNoteController(
				func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
					input.Output = output
					return input
				},
				func() *presenter.NotePresenter { return p },
				func() port.NoteRepository { return nil },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodDelete, "/api/notes/n1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Delete(c, "n1", openapi.NotesDeleteNoteParams{OwnerId: tt.ownerID})
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
