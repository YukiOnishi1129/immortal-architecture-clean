package usecase

import (
	"context"

	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

type AccountInteractor struct {
	repo   port.AccountRepository
	output port.AccountOutputPort
}

var _ port.AccountInputPort = (*AccountInteractor)(nil)

func NewAccountInteractor(repo port.AccountRepository, output port.AccountOutputPort) *AccountInteractor {
	return &AccountInteractor{repo: repo, output: output}
}

func (u *AccountInteractor) CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error {
	email, err := account.ParseEmail(input.Email)
	if err != nil {
		return err
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
		return err
	}
	a, err := u.repo.UpsertOAuthAccount(ctx, input)
	if err != nil {
		return err
	}
	return u.output.PresentAccount(ctx, a)
}

func (u *AccountInteractor) GetByID(ctx context.Context, id string) error {
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.output.PresentAccount(ctx, a)
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
