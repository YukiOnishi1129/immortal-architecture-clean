package port

import (
	"context"

	"immortal-architecture-clean/backend/internal/domain/account"
)

// AccountInputPort defines entrypoints for account use cases.
type AccountInputPort interface {
	CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error
	GetByID(ctx context.Context, id string) error
}

// AccountOutputPort converts account結果を外部向けに整形するための契約。
type AccountOutputPort interface {
	PresentAccount(ctx context.Context, account *account.Account) error
}

// AccountRepository abstracts persistence for accounts.
type AccountRepository interface {
	UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error)
	GetByID(ctx context.Context, id string) (*account.Account, error)
}
