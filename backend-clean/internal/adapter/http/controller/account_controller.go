package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

type AccountController struct {
	input port.AccountInputPort
}

func NewAccountController(input port.AccountInputPort) *AccountController {
	return &AccountController{input: input}
}

func (c *AccountController) CreateOrGet(ctx echo.Context) error {
	var body openapi.ModelsCreateOrGetAccountRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	acc, err := c.input.CreateOrGet(ctx.Request().Context(), account.OAuthAccountInput{
		Email:             body.Email,
		FirstName:         body.Name,
		LastName:          "",
		Provider:          body.Provider,
		ProviderAccountID: body.ProviderAccountId,
		Thumbnail:         body.Thumbnail,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toAccountResponse(acc))
}

func (c *AccountController) GetByID(ctx echo.Context, accountID string) error {
	acc, err := c.input.GetByID(ctx.Request().Context(), accountID)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, toAccountResponse(acc))
}
