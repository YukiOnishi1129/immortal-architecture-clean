package presenter

import (
	"context"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/template"
	"immortal-architecture-clean/backend/internal/port"
)

type TemplatePresenter struct {
	template *openapi.ModelsTemplateResponse
	list     []openapi.ModelsTemplateResponse
	deleted  bool
}

var _ port.TemplateOutputPort = (*TemplatePresenter)(nil)

func NewTemplatePresenter() *TemplatePresenter {
	return &TemplatePresenter{}
}

func (p *TemplatePresenter) PresentTemplateList(ctx context.Context, templates []template.WithUsage) error {
	res := make([]openapi.ModelsTemplateResponse, 0, len(templates))
	for _, t := range templates {
		res = append(res, toTemplateResponse(t))
	}
	p.list = res
	return nil
}

func (p *TemplatePresenter) PresentTemplate(ctx context.Context, tpl *template.WithUsage) error {
	resp := toTemplateResponse(*tpl)
	p.template = &resp
	return nil
}

func (p *TemplatePresenter) PresentTemplateDeleted(ctx context.Context) error {
	p.deleted = true
	return nil
}

func (p *TemplatePresenter) Template() *openapi.ModelsTemplateResponse {
	return p.template
}

func (p *TemplatePresenter) Templates() []openapi.ModelsTemplateResponse {
	return p.list
}

func (p *TemplatePresenter) DeleteResponse() openapi.ModelsSuccessResponse {
	return openapi.ModelsSuccessResponse{Success: p.deleted}
}

func toTemplateResponse(t template.WithUsage) openapi.ModelsTemplateResponse {
	fields := make([]openapi.ModelsField, 0, len(t.Fields))
	for _, f := range t.Fields {
		fields = append(fields, openapi.ModelsField{
			Id:         f.ID,
			Label:      f.Label,
			Order:      int32(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	return openapi.ModelsTemplateResponse{
		Id:        t.Template.ID,
		Name:      t.Template.Name,
		OwnerId:   t.Template.OwnerID,
		Owner:     openapi.ModelsAccountSummary{Id: t.Template.OwnerID},
		Fields:    fields,
		IsUsed:    t.IsUsed,
		UpdatedAt: t.Template.UpdatedAt,
	}
}
