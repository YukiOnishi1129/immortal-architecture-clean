package controller

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/account"
	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/domain/template"
)

func handleError(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, domainerr.ErrNotFound):
		return ctx.JSON(http.StatusNotFound, openapi.ModelsNotFoundError{Code: openapi.ModelsNotFoundErrorCodeNOTFOUND, Message: err.Error()})
	case errors.Is(err, domainerr.ErrUnauthorized):
		return ctx.JSON(http.StatusForbidden, openapi.ModelsForbiddenError{Code: openapi.ModelsForbiddenErrorCodeFORBIDDEN, Message: err.Error()})
	case errors.Is(err, account.ErrInvalidEmail), errors.Is(err, account.ErrInvalidName):
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: err.Error()})
	case errors.Is(err, domainerr.ErrInvalidStatus) || errors.Is(err, domainerr.ErrInvalidStatusChange) || errors.Is(err, domainerr.ErrInvalidTemplateField):
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: err.Error()})
	default:
		return ctx.JSON(http.StatusInternalServerError, openapi.ModelsErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
}

func toAccountResponse(a *account.Account) openapi.ModelsAccountResponse {
	var lastLogin time.Time
	if a.LastLoginAt != nil {
		lastLogin = *a.LastLoginAt
	}
	return openapi.ModelsAccountResponse{
		Id:          a.ID,
		Email:       a.Email.String(),
		FirstName:   a.FirstName,
		LastName:    a.LastName,
		FullName:    strings.TrimSpace(a.FirstName + " " + a.LastName),
		Thumbnail:   strPtrOrNil(a.Thumbnail),
		LastLoginAt: lastLogin,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
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

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
