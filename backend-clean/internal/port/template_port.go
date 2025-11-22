package port

import (
	"context"

	"immortal-architecture-clean/backend/internal/domain/template"
)

type TemplateInputPort interface {
	List(ctx context.Context, filters template.Filters) error
	Get(ctx context.Context, id string) error
	Create(ctx context.Context, input TemplateCreateInput) error
	Update(ctx context.Context, input TemplateUpdateInput) error
	Delete(ctx context.Context, id, ownerID string) error
}

type TemplateOutputPort interface {
	PresentTemplateList(ctx context.Context, templates []template.WithUsage) error
	PresentTemplate(ctx context.Context, template *template.WithUsage) error
	PresentTemplateDeleted(ctx context.Context) error
}

type TemplateRepository interface {
	List(ctx context.Context, filters template.Filters) ([]template.WithUsage, error)
	Get(ctx context.Context, id string) (*template.WithUsage, error)
	Create(ctx context.Context, tpl template.Template) (*template.Template, error)
	Update(ctx context.Context, tpl template.Template) (*template.Template, error)
	Delete(ctx context.Context, id string) error
	ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error
}

type TemplateCreateInput struct {
	Name    string
	OwnerID string
	Fields  []template.Field
}

type TemplateUpdateInput struct {
	ID      string
	Name    string
	Fields  []template.Field
	OwnerID string
}
