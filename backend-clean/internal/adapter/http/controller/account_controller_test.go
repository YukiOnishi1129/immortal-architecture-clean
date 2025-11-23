package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	ctrlmock "immortal-architecture-clean/backend/internal/adapter/http/controller/mock"
	"immortal-architecture-clean/backend/internal/adapter/http/presenter"
	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/port"
)

func TestAccountController_CreateOrGet(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "[Success] create or get",
			body:       `{"email":"user@example.com","name":"Taro","provider":"google","providerAccountId":"pid"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "[Fail] bind error",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/accounts/auth", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.CreateOrGet(c)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK && (p.Response() == nil || p.Response().Email != "user@example.com") {
				t.Fatalf("presenter response not set: %+v", p.Response())
			}
		})
	}
}

func TestAccountController_GetCurrent(t *testing.T) {
	tests := []struct {
		name       string
		headerID   string
		wantStatus int
	}{
		{name: "[Success] with header", headerID: "acc-1", wantStatus: http.StatusOK},
		{name: "[Fail] missing header", headerID: "", wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/accounts/me", nil)
			if tt.headerID != "" {
				req.Header.Set("X-Account-ID", tt.headerID)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.GetCurrent(c)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAccountController_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		getErr     error
		wantStatus int
	}{
		{name: "[Success] get by id", wantStatus: http.StatusOK},
		{name: "[Fail] not found", getErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{GetErr: tt.getErr}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/accounts/acc-1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.GetByID(c, "acc-1")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
