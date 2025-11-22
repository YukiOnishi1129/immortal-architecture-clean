package port

import (
	"context"

	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/domain/template"
)

type NoteInputPort interface {
	List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error)
	Get(ctx context.Context, id string) (*note.WithMeta, error)
	Create(ctx context.Context, input NoteCreateInput) (*note.WithMeta, error)
	Update(ctx context.Context, input NoteUpdateInput) (*note.WithMeta, error)
	ChangeStatus(ctx context.Context, input NoteStatusChangeInput) (*note.WithMeta, error)
	Delete(ctx context.Context, id, ownerID string) error
}

type NoteOutputPort interface {
	PresentNoteList(ctx context.Context, notes []note.WithMeta) error
	PresentNote(ctx context.Context, note *note.WithMeta) error
	PresentNoteDeleted(ctx context.Context) error
}

type NoteRepository interface {
	List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error)
	Get(ctx context.Context, id string) (*note.WithMeta, error)
	Create(ctx context.Context, n note.Note) (*note.Note, error)
	Update(ctx context.Context, n note.Note) (*note.Note, error)
	UpdateStatus(ctx context.Context, id string, status note.NoteStatus) (*note.Note, error)
	Delete(ctx context.Context, id string) error
	ReplaceSections(ctx context.Context, noteID string, sections []note.Section) error
}

type NoteCreateInput struct {
	Title      string
	TemplateID string
	OwnerID    string
	Sections   []SectionInput
}

type SectionInput struct {
	FieldID string
	Content string
}

type NoteUpdateInput struct {
	ID       string
	Title    string
	OwnerID  string
	Sections []SectionUpdateInput
}

type SectionUpdateInput struct {
	SectionID string
	Content   string
}

type NoteStatusChangeInput struct {
	ID      string
	OwnerID string
	Status  note.NoteStatus
}

// Helpers to alias domain types when needed
type NoteFilters = note.Filters
type NoteWithMeta = note.WithMeta
type TemplateFields = []template.Field
