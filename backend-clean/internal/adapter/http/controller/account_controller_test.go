package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"immortal-architecture-clean/backend/internal/adapter/http/presenter"
	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

// stubAccountInput is a lightweight input port stub for controller tests.
type stubAccountInput struct {
	createErr error
	getErr    error
	output    port.AccountOutputPort
}

func (s *stubAccountInput) CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error {
	// simulate usecase setting presenter response
	if s.output != nil && s.createErr == nil {
		_ = s.output.PresentAccount(ctx, &account.Account{
			ID:        "acc-1",
			Email:     account.Email(input.Email),
			FirstName: input.FirstName,
			Provider:  input.Provider,
		})
	}
	return s.createErr
}

func (s *stubAccountInput) GetByID(ctx context.Context, id string) error {
	if s.output != nil && s.getErr == nil {
		_ = s.output.PresentAccount(ctx, &account.Account{
			ID:        id,
			Email:     account.Email("user@example.com"),
			FirstName: "Taro",
			Provider:  "google",
		})
	}
	return s.getErr
}

// Tests

func TestAccountController_CreateOrGet_Success(t *testing.T) {
	p := presenter.NewAccountPresenter()
	input := &stubAccountInput{}
	ctrl := NewAccountController(
		func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
			input.output = output
			return input
		},
		func() *presenter.AccountPresenter { return p },
		func() port.AccountRepository { return nil },
	)

	e := echo.New()
	body := `{"email":"user@example.com","name":"Taro","provider":"google","providerAccountId":"pid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/accounts/auth", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := ctrl.CreateOrGet(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if p.Response() == nil || p.Response().Email != "user@example.com" {
		t.Fatalf("presenter response not set: %+v", p.Response())
	}
}

func TestAccountController_CreateOrGet_BadRequest(t *testing.T) {
	ctrl := NewAccountController(
		func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
			return &stubAccountInput{output: output}
		},
		func() *presenter.AccountPresenter { return presenter.NewAccountPresenter() },
		func() port.AccountRepository { return nil },
	)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/accounts/auth", bytes.NewBufferString("not-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = ctrl.CreateOrGet(c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestAccountController_GetCurrent_NoHeader(t *testing.T) {
	ctrl := NewAccountController(
		func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
			return &stubAccountInput{output: output}
		},
		func() *presenter.AccountPresenter { return presenter.NewAccountPresenter() },
		func() port.AccountRepository { return nil },
	)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/accounts/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = ctrl.GetCurrent(c)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}
