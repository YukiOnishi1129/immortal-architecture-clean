package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/template"
	"immortal-architecture-clean/backend/internal/port"
)

type TemplateController struct {
	input port.TemplateInputPort
}

func NewTemplateController(input port.TemplateInputPort) *TemplateController {
	return &TemplateController{input: input}
}

func (c *TemplateController) List(ctx echo.Context, params openapi.TemplatesListTemplatesParams) error {
	filters := template.Filters{
		Query:   params.Q,
		OwnerID: params.OwnerId,
	}
	templates, err := c.input.List(ctx.Request().Context(), filters)
	if err != nil {
		return handleError(ctx, err)
	}
	resp := make([]openapi.ModelsTemplateResponse, 0, len(templates))
	for _, t := range templates {
		resp = append(resp, toTemplateResponse(t))
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *TemplateController) Create(ctx echo.Context) error {
	var body openapi.ModelsCreateTemplateRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	fields := make([]template.Field, 0, len(body.Fields))
	for _, f := range body.Fields {
		fields = append(fields, template.Field{
			Label:      f.Label,
			Order:      int(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	tpl, err := c.input.Create(ctx.Request().Context(), port.TemplateCreateInput{
		Name:    body.Name,
		OwnerID: body.OwnerId.String(),
		Fields:  fields,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toTemplateResponse(*tpl))
}

func (c *TemplateController) Delete(ctx echo.Context, templateID string) error {
	if err := c.input.Delete(ctx.Request().Context(), templateID); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, openapi.ModelsSuccessResponse{Success: true})
}

func (c *TemplateController) GetByID(ctx echo.Context, templateID string) error {
	tpl, err := c.input.Get(ctx.Request().Context(), templateID)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toTemplateResponse(*tpl))
}

func (c *TemplateController) Update(ctx echo.Context, templateID string) error {
	var body openapi.ModelsUpdateTemplateRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	fields := make([]template.Field, 0, len(body.Fields))
	for _, f := range body.Fields {
		fields = append(fields, template.Field{
			ID:         valueOrEmpty(f.Id),
			Label:      f.Label,
			Order:      int(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	tpl, err := c.input.Update(ctx.Request().Context(), port.TemplateUpdateInput{
		ID:     templateID,
		Name:   body.Name,
		Fields: fields,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toTemplateResponse(*tpl))
}
