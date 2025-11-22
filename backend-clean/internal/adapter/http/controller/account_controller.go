package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/adapter/http/presenter"
	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

type AccountController struct {
	inputFactory  func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort
	outputFactory func() *presenter.AccountPresenter
	repoFactory   func() port.AccountRepository
}

func NewAccountController(
	inputFactory func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort,
	outputFactory func() *presenter.AccountPresenter,
	repoFactory func() port.AccountRepository,
) *AccountController {
	return &AccountController{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
		repoFactory:   repoFactory,
	}
}

func (c *AccountController) CreateOrGet(ctx echo.Context) error {
	var body openapi.ModelsCreateOrGetAccountRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	input, p := c.newIO()
	err := input.CreateOrGet(ctx.Request().Context(), account.OAuthAccountInput{
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
	return ctx.JSON(http.StatusOK, p.Response())
}

func (c *AccountController) GetByID(ctx echo.Context, accountID string) error {
	input, p := c.newIO()
	err := input.GetByID(ctx.Request().Context(), accountID)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Response())
}

func (c *AccountController) GetCurrent(ctx echo.Context) error {
	accountID, err := currentAccountID(ctx)
	if err != nil {
		return handleError(ctx, err)
	}
	input, p := c.newIO()
	err = input.GetByID(ctx.Request().Context(), accountID)
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Response())
}

func (c *AccountController) newIO() (port.AccountInputPort, *presenter.AccountPresenter) {
	output := c.outputFactory()
	input := c.inputFactory(c.repoFactory(), output)
	return input, output
}
