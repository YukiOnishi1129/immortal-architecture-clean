package service

import (
	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/domain/template"
)

// BuildNote validates template and sections, then returns a Note ready for persistence.
func BuildNote(title string, tpl template.Template, ownerID string, sections []note.Section) (note.Note, error) {
	if err := note.ValidateNoteForCreate(title, tpl, sections); err != nil {
		return note.Note{}, err
	}
	if ownerID == "" {
		return note.Note{}, domainerr.ErrOwnerRequired
	}
	return note.Note{
		Title:      title,
		TemplateID: tpl.ID,
		OwnerID:    ownerID,
		Status:     note.StatusDraft,
		Sections:   sections,
	}, nil
}
