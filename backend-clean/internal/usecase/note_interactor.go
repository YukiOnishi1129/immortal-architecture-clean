package usecase

import (
	"context"
	"strings"

	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/domain/service"
	"immortal-architecture-clean/backend/internal/domain/template"
	"immortal-architecture-clean/backend/internal/port"
)

type NoteInteractor struct {
	notes     port.NoteRepository
	templates port.TemplateRepository
	tx        port.TxManager
}

var _ port.NoteInputPort = (*NoteInteractor)(nil)

func NewNoteInteractor(notes port.NoteRepository, templates port.TemplateRepository, tx port.TxManager) *NoteInteractor {
	return &NoteInteractor{
		notes:     notes,
		templates: templates,
		tx:        tx,
	}
}

func (u *NoteInteractor) List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error) {
	return u.notes.List(ctx, filters)
}

func (u *NoteInteractor) Get(ctx context.Context, id string) (*note.WithMeta, error) {
	return u.notes.Get(ctx, id)
}

func (u *NoteInteractor) Create(ctx context.Context, input port.NoteCreateInput) (*note.WithMeta, error) {
	tpl, err := u.templates.Get(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	sections := buildSections("", tpl.Template.Fields, input.Sections)
	if err := note.ValidateNoteForCreate(input.Title, tpl.Template, sections); err != nil {
		return nil, err
	}

	var noteID string
	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		built, err := service.BuildNote(input.Title, tpl.Template, input.OwnerID, sections)
		if err != nil {
			return err
		}
		nn, err := u.notes.Create(txCtx, built)
		if err != nil {
			return err
		}
		noteID = nn.ID
		sectionsWithID := buildSections(noteID, tpl.Template.Fields, input.Sections)
		if err := note.ValidateSections(tpl.Template.Fields, sectionsWithID); err != nil {
			return err
		}
		return u.notes.ReplaceSections(txCtx, noteID, sectionsWithID)
	})
	if err != nil {
		return nil, err
	}
	return u.notes.Get(ctx, noteID)
}

func (u *NoteInteractor) Update(ctx context.Context, input port.NoteUpdateInput) (*note.WithMeta, error) {
	current, err := u.notes.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := note.ValidateNoteOwnership(current.Note.OwnerID, input.OwnerID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, note.ErrTitleRequired
	}

	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, err := u.notes.Update(txCtx, note.Note{
			ID:    input.ID,
			Title: input.Title,
		})
		if err != nil {
			return err
		}
		if input.Sections != nil {
			tpl, err := u.templates.Get(ctx, current.Note.TemplateID)
			if err != nil {
				return err
			}
			sections := buildSections(current.Note.ID, tpl.Template.Fields, input.Sections)
			if err := note.ValidateSections(tpl.Template.Fields, sections); err != nil {
				return err
			}
			if err := u.notes.ReplaceSections(txCtx, input.ID, sections); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.notes.Get(ctx, input.ID)
}

func (u *NoteInteractor) ChangeStatus(ctx context.Context, input port.NoteStatusChangeInput) (*note.WithMeta, error) {
	current, err := u.notes.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := note.ValidateNoteOwnership(current.Note.OwnerID, input.OwnerID); err != nil {
		return nil, err
	}
	if err := input.Status.Validate(); err != nil {
		return nil, err
	}
	// domain service handles owner check + transition rule
	if input.Status == note.StatusPublish {
		if err := service.CanPublish(current.Note, input.OwnerID); err != nil {
			return nil, err
		}
	} else {
		if err := service.CanUnpublish(current.Note, input.OwnerID); err != nil {
			return nil, err
		}
	}
	if err := note.CanChangeStatus(current.Note.Status, input.Status); err != nil {
		return nil, err
	}

	if _, err := u.notes.UpdateStatus(ctx, input.ID, input.Status); err != nil {
		return nil, err
	}
	return u.notes.Get(ctx, input.ID)
}

func (u *NoteInteractor) Delete(ctx context.Context, id, ownerID string) error {
	current, err := u.notes.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := note.ValidateNoteOwnership(current.Note.OwnerID, ownerID); err != nil {
		return err
	}
	return u.notes.Delete(ctx, id)
}

func buildSections(noteID string, templateFields []template.Field, inputs []port.SectionInput) []note.Section {
	if len(inputs) == 0 {
		sections := service.BuildSectionsFromTemplate(templateFields)
		for i := range sections {
			sections[i].NoteID = noteID
		}
		return sections
	}

	sections := make([]note.Section, 0, len(inputs))
	for _, s := range inputs {
		sections = append(sections, note.Section{
			FieldID: s.FieldID,
			NoteID:  noteID,
			Content: s.Content,
		})
	}
	return sections
}
