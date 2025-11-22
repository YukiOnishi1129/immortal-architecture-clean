package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/port"
)

type NoteController struct {
	input port.NoteInputPort
}

func NewNoteController(input port.NoteInputPort) *NoteController {
	return &NoteController{input: input}
}

func (c *NoteController) List(ctx echo.Context, params openapi.NotesListNotesParams) error {
	var status *note.NoteStatus
	if params.Status != nil {
		s := note.NoteStatus(*params.Status)
		status = &s
	}
	filters := note.Filters{
		Status:     status,
		TemplateID: params.TemplateId,
		OwnerID:    params.OwnerId,
		Query:      params.Q,
	}
	notes, err := c.input.List(ctx.Request().Context(), filters)
	if err != nil {
		return handleError(ctx, err)
	}
	res := make([]openapi.ModelsNoteResponse, 0, len(notes))
	for _, n := range notes {
		res = append(res, toNoteResponse(n))
	}
	return ctx.JSON(http.StatusOK, res)
}

func (c *NoteController) Create(ctx echo.Context) error {
	var body openapi.ModelsCreateNoteRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	sections := []port.SectionInput{}
	if body.Sections != nil {
		for _, s := range *body.Sections {
			sections = append(sections, port.SectionInput{
				FieldID: s.FieldId,
				Content: s.Content,
			})
		}
	}
	n, err := c.input.Create(ctx.Request().Context(), port.NoteCreateInput{
		Title:      body.Title,
		TemplateID: body.TemplateId.String(),
		OwnerID:    body.OwnerId.String(),
		Sections:   sections,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toNoteResponse(*n))
}

func (c *NoteController) Delete(ctx echo.Context, noteID string) error {
	if err := c.input.Delete(ctx.Request().Context(), noteID, ""); err != nil { // TODO: ownerID from auth
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, openapi.ModelsSuccessResponse{Success: true})
}

func (c *NoteController) GetByID(ctx echo.Context, noteID string) error {
	n, err := c.input.Get(ctx.Request().Context(), noteID)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toNoteResponse(*n))
}

func (c *NoteController) Update(ctx echo.Context, noteID string) error {
	var body openapi.ModelsUpdateNoteRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	sections := make([]port.SectionUpdateInput, 0, len(body.Sections))
	for _, s := range body.Sections {
		sections = append(sections, port.SectionUpdateInput{
			SectionID: s.Id,
			Content:   s.Content,
		})
	}
	n, err := c.input.Update(ctx.Request().Context(), port.NoteUpdateInput{
		ID:       noteID,
		Title:    body.Title,
		OwnerID:  body.OwnerId,
		Sections: sections,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toNoteResponse(*n))
}

func (c *NoteController) Publish(ctx echo.Context, noteID string) error {
	n, err := c.input.ChangeStatus(ctx.Request().Context(), port.NoteStatusChangeInput{
		ID:      noteID,
		Status:  note.StatusPublish,
		OwnerID: "",
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toNoteResponse(*n))
}

func (c *NoteController) Unpublish(ctx echo.Context, noteID string) error {
	n, err := c.input.ChangeStatus(ctx.Request().Context(), port.NoteStatusChangeInput{
		ID:      noteID,
		Status:  note.StatusDraft,
		OwnerID: "",
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toNoteResponse(*n))
}
