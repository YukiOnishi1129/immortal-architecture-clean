package usecase

import (
	"context"

	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/template"
	"immortal-architecture-clean/backend/internal/port"
)

type TemplateInteractor struct {
	repo port.TemplateRepository
	tx   port.TxManager
}

var _ port.TemplateInputPort = (*TemplateInteractor)(nil)

func NewTemplateInteractor(repo port.TemplateRepository, tx port.TxManager) *TemplateInteractor {
	return &TemplateInteractor{repo: repo, tx: tx}
}

func (u *TemplateInteractor) List(ctx context.Context, filters template.Filters) ([]template.WithUsage, error) {
	return u.repo.List(ctx, filters)
}

func (u *TemplateInteractor) Get(ctx context.Context, id string) (*template.WithUsage, error) {
	return u.repo.Get(ctx, id)
}

func (u *TemplateInteractor) Create(ctx context.Context, input port.TemplateCreateInput) (*template.WithUsage, error) {
	var createdID string
	if len(input.Fields) > 0 {
		fields, err := template.NormalizeAndValidate(input.Fields)
		if err != nil {
			return nil, err
		}
		input.Fields = fields
	}
	err := u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		tpl, err := u.repo.Create(txCtx, template.Template{
			Name:    input.Name,
			OwnerID: input.OwnerID,
		})
		if err != nil {
			return err
		}
		createdID = tpl.ID
		if len(input.Fields) > 0 {
			if err := u.repo.ReplaceFields(txCtx, tpl.ID, input.Fields); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.repo.Get(ctx, createdID)
}

func (u *TemplateInteractor) Update(ctx context.Context, input port.TemplateUpdateInput) (*template.WithUsage, error) {
	if input.Fields != nil {
		fields, err := template.NormalizeAndValidate(input.Fields)
		if err != nil {
			return nil, err
		}
		input.Fields = fields
	}
	err := u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, err := u.repo.Update(txCtx, template.Template{
			ID:   input.ID,
			Name: input.Name,
		})
		if err != nil {
			return err
		}
		if input.Fields != nil {
			if len(input.Fields) == 0 {
				return domainerr.ErrInvalidTemplateField
			}
			if err := u.repo.ReplaceFields(txCtx, input.ID, input.Fields); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.repo.Get(ctx, input.ID)
}

func (u *TemplateInteractor) Delete(ctx context.Context, id string) error {
	tpl, err := u.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if tpl.IsUsed {
		return errors.ErrTemplateInUse
	}
	return u.repo.Delete(ctx, id)
}
