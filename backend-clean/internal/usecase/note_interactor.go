package usecase

import (
	"context"

	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
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
	template, err := u.templates.Get(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	sections := buildSections(template.Fields, input.Sections)

	var noteID string
	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		nn, err := u.notes.Create(txCtx, note.Note{
			Title:      input.Title,
			TemplateID: input.TemplateID,
			OwnerID:    input.OwnerID,
			Status:     note.StatusDraft,
		})
		if err != nil {
			return err
		}
		noteID = nn.ID
		return u.notes.ReplaceSections(txCtx, noteID, sections)
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
	if current.Note.OwnerID != input.OwnerID {
		return nil, domainerr.ErrUnauthorized
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
			sections := make([]note.Section, 0, len(input.Sections))
			for _, s := range input.Sections {
				sections = append(sections, note.Section{
					ID:      s.SectionID,
					Content: s.Content,
				})
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
	if current.Note.OwnerID != input.OwnerID {
		return nil, domainerr.ErrUnauthorized
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
	if current.Note.OwnerID != ownerID {
		return domainerr.ErrUnauthorized
	}
	return u.notes.Delete(ctx, id)
}

func buildSections(templateFields []template.Field, inputs []port.SectionInput) []note.Section {
	if len(inputs) == 0 {
		return service.BuildSectionsFromTemplate(templateFields)
	}

	sections := make([]note.Section, 0, len(inputs))
	for _, s := range inputs {
		sections = append(sections, note.Section{
			FieldID: s.FieldID,
			Content: s.Content,
		})
	}
	return sections
}
