// Package db implements gateway repositories.
package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	sqldb "immortal-architecture-clean/backend/internal/adapter/gateway/db/sqlc"
	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/port"
)

// AccountRepository implements account persistence.
type AccountRepository struct {
	pool    *pgxpool.Pool
	queries *sqldb.Queries
}

var _ port.AccountRepository = (*AccountRepository)(nil)

// NewAccountRepository creates AccountRepository.
func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		pool:    pool,
		queries: sqldb.New(pool),
	}
}

// UpsertOAuthAccount inserts or updates an OAuth account.
func (r *AccountRepository) UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	q := queriesForContext(ctx, r.queries)

	row, err := q.UpsertAccount(ctx, &sqldb.UpsertAccountParams{
		Email:             input.Email,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Provider:          input.Provider,
		ProviderAccountID: input.ProviderAccountID,
		Thumbnail:         pgNullableText(input.Thumbnail),
		LastLoginAt:       pgNullableTime(nil),
	})
	if err != nil {
		return nil, err
	}
	return toDomainAccount(row)
}

// GetByID fetches account by ID.
func (r *AccountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	q := queriesForContext(ctx, r.queries)
	uuid, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := q.GetAccountByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return toDomainAccount(row)
}

func toDomainAccount(a *sqldb.Account) (*account.Account, error) {
	var lastLogin *time.Time
	if a.LastLoginAt.Valid {
		t := timestamptzToTime(a.LastLoginAt)
		lastLogin = &t
	}
	email, err := account.ParseEmail(a.Email)
	if err != nil {
		return nil, err
	}
	return &account.Account{
		ID:                uuidToString(a.ID),
		Email:             email,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		IsActive:          a.IsActive,
		Provider:          a.Provider,
		ProviderAccountID: a.ProviderAccountID,
		Thumbnail:         nullableTextToString(a.Thumbnail),
		LastLoginAt:       lastLogin,
		CreatedAt:         timestamptzToTime(a.CreatedAt),
		UpdatedAt:         timestamptzToTime(a.UpdatedAt),
	}, nil
}
