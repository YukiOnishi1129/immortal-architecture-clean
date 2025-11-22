// Package db implements gateway repositories.
package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	sqldb "immortal-architecture-clean/backend/internal/adapter/gateway/db/sqlc"
	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/template"
	"immortal-architecture-clean/backend/internal/port"
)

// TemplateRepository implements template persistence.
type TemplateRepository struct {
	pool    *pgxpool.Pool
	queries *sqldb.Queries
}

var _ port.TemplateRepository = (*TemplateRepository)(nil)

// NewTemplateRepository creates TemplateRepository.
func NewTemplateRepository(pool *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{
		pool:    pool,
		queries: sqldb.New(pool),
	}
}

// List returns templates by filters.
func (r *TemplateRepository) List(ctx context.Context, filters template.Filters) ([]template.WithUsage, error) {
	params := &sqldb.ListTemplatesParams{}
	if filters.OwnerID != nil && *filters.OwnerID != "" {
		if id, err := toUUID(*filters.OwnerID); err == nil {
			params.Column1 = id
		}
	}
	if filters.Query != nil && *filters.Query != "" {
		params.Column2 = *filters.Query
	}

	rows, err := queriesForContext(ctx, r.queries).ListTemplates(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]template.WithUsage, 0, len(rows))
	for _, row := range rows {
		fields, err := r.listFields(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, template.WithUsage{
			Template: template.Template{
				ID:        uuidToString(row.ID),
				Name:      row.Name,
				OwnerID:   uuidToString(row.OwnerID),
				UpdatedAt: timestamptzToTime(row.UpdatedAt),
				Fields:    fields,
			},
			Fields: fields,
			IsUsed: row.IsUsed,
		})
	}
	return result, nil
}

// Get returns a template with usage and fields.
func (r *TemplateRepository) Get(ctx context.Context, id string) (*template.WithUsage, error) {
	pgID, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).GetTemplateByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	fields, err := r.listFields(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	return &template.WithUsage{
		Template: template.Template{
			ID:        uuidToString(row.ID),
			Name:      row.Name,
			OwnerID:   uuidToString(row.OwnerID),
			UpdatedAt: timestamptzToTime(row.UpdatedAt),
			Fields:    fields,
		},
		Fields: fields,
		IsUsed: row.IsUsed,
	}, nil
}

// Create inserts a template.
func (r *TemplateRepository) Create(ctx context.Context, tpl template.Template) (*template.Template, error) {
	owner, err := toUUID(tpl.OwnerID)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).CreateTemplate(ctx, &sqldb.CreateTemplateParams{
		Name:    tpl.Name,
		OwnerID: owner,
	})
	if err != nil {
		return nil, err
	}
	return &template.Template{
		ID:        uuidToString(row.ID),
		Name:      row.Name,
		OwnerID:   uuidToString(row.OwnerID),
		UpdatedAt: timestamptzToTime(row.UpdatedAt),
	}, nil
}

// Update updates template name.
func (r *TemplateRepository) Update(ctx context.Context, tpl template.Template) (*template.Template, error) {
	pgID, err := toUUID(tpl.ID)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).UpdateTemplate(ctx, &sqldb.UpdateTemplateParams{
		ID:   pgID,
		Name: tpl.Name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	return &template.Template{
		ID:        uuidToString(row.ID),
		Name:      row.Name,
		OwnerID:   uuidToString(row.OwnerID),
		UpdatedAt: timestamptzToTime(row.UpdatedAt),
	}, nil
}

// Delete deletes a template.
func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	pgID, err := toUUID(id)
	if err != nil {
		return err
	}
	return queriesForContext(ctx, r.queries).DeleteTemplate(ctx, pgID)
}

// ReplaceFields replaces template fields.
func (r *TemplateRepository) ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error {
	pgID, err := toUUID(templateID)
	if err != nil {
		return err
	}
	q := queriesForContext(ctx, r.queries)
	if err := q.DeleteFieldsByTemplate(ctx, pgID); err != nil {
		return err
	}
	for idx, f := range fields {
		order := f.Order
		if order == 0 {
			order = idx + 1
		}
		if _, err := q.CreateField(ctx, &sqldb.CreateFieldParams{
			TemplateID: pgID,
			Label:      f.Label,
			Order:      int32(order), //nolint:gosec
			IsRequired: f.IsRequired,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *TemplateRepository) listFields(ctx context.Context, templateID pgtype.UUID) ([]template.Field, error) {
	rows, err := queriesForContext(ctx, r.queries).ListFieldsByTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}
	fields := make([]template.Field, 0, len(rows))
	for _, f := range rows {
		fields = append(fields, template.Field{
			ID:         uuidToString(f.ID),
			Label:      f.Label,
			Order:      int(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	return fields, nil
}
