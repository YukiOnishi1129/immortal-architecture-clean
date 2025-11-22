package presenter

import (
	"context"
	"strings"
	"time"

	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

type AccountPresenter struct {
	account *openapi.ModelsAccountResponse
}

var _ port.AccountOutputPort = (*AccountPresenter)(nil)

func NewAccountPresenter() *AccountPresenter {
	return &AccountPresenter{}
}

func (p *AccountPresenter) PresentAccount(ctx context.Context, a *account.Account) error {
	var lastLogin time.Time
	if a.LastLoginAt != nil {
		lastLogin = *a.LastLoginAt
	}
	p.account = &openapi.ModelsAccountResponse{
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
	return nil
}

func (p *AccountPresenter) Response() *openapi.ModelsAccountResponse {
	return p.account
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
