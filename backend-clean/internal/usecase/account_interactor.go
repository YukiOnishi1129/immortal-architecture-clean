package usecase

import (
	"context"

	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

type AccountInteractor struct {
	repo port.AccountRepository
}

var _ port.AccountInputPort = (*AccountInteractor)(nil)

func NewAccountInteractor(repo port.AccountRepository) *AccountInteractor {
	return &AccountInteractor{repo: repo}
}

func (u *AccountInteractor) CreateOrGet(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	email, err := account.ParseEmail(input.Email)
	if err != nil {
		return nil, err
	}
	acc := account.Account{
		Email:             email,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Provider:          input.Provider,
		ProviderAccountID: input.ProviderAccountID,
		Thumbnail:         valueOrEmpty(input.Thumbnail),
	}
	if err := account.Validate(acc); err != nil {
		return nil, err
	}
	return u.repo.UpsertOAuthAccount(ctx, input)
}

func (u *AccountInteractor) GetByID(ctx context.Context, id string) (*account.Account, error) {
	return u.repo.GetByID(ctx, id)
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
