package presenter

import (
	"context"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/port"
)

type NotePresenter struct {
	note      *openapi.ModelsNoteResponse
	notes     []openapi.ModelsNoteResponse
	deletedOK bool
}

var _ port.NoteOutputPort = (*NotePresenter)(nil)

func NewNotePresenter() *NotePresenter {
	return &NotePresenter{}
}

func (p *NotePresenter) PresentNoteList(ctx context.Context, notes []note.WithMeta) error {
	res := make([]openapi.ModelsNoteResponse, 0, len(notes))
	for _, n := range notes {
		res = append(res, toNoteResponse(n))
	}
	p.notes = res
	return nil
}

func (p *NotePresenter) PresentNote(ctx context.Context, n *note.WithMeta) error {
	resp := toNoteResponse(*n)
	p.note = &resp
	return nil
}

func (p *NotePresenter) PresentNoteDeleted(ctx context.Context) error {
	p.deletedOK = true
	return nil
}

func (p *NotePresenter) Note() *openapi.ModelsNoteResponse {
	return p.note
}

func (p *NotePresenter) Notes() []openapi.ModelsNoteResponse {
	return p.notes
}

func (p *NotePresenter) DeleteResponse() openapi.ModelsSuccessResponse {
	return openapi.ModelsSuccessResponse{Success: p.deletedOK}
}

func toNoteResponse(n note.WithMeta) openapi.ModelsNoteResponse {
	sections := make([]openapi.ModelsSection, 0, len(n.Sections))
	for _, s := range n.Sections {
		sections = append(sections, openapi.ModelsSection{
			Id:         s.Section.ID,
			FieldId:    s.Section.FieldID,
			FieldLabel: s.FieldLabel,
			Content:    s.Section.Content,
			IsRequired: s.IsRequired,
		})
	}
	return openapi.ModelsNoteResponse{
		Id:           n.Note.ID,
		Title:        n.Note.Title,
		TemplateId:   n.Note.TemplateID,
		TemplateName: n.TemplateName,
		OwnerId:      n.Note.OwnerID,
		Owner: openapi.ModelsAccountSummary{
			Id:        n.Note.OwnerID,
			FirstName: "",
			LastName:  "",
		},
		Status:    openapi.ModelsNoteStatus(n.Note.Status),
		Sections:  sections,
		CreatedAt: n.Note.CreatedAt,
		UpdatedAt: n.Note.UpdatedAt,
	}
}
