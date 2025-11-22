// Package factory provides constructors for driver-level wiring.
package factory

import (
	"github.com/jackc/pgx/v5/pgxpool"

	gatewaydb "immortal-architecture-clean/backend/internal/adapter/gateway/db"
	"immortal-architecture-clean/backend/internal/port"
)

// NewAccountRepoFactory returns a factory that creates AccountRepository.
func NewAccountRepoFactory(pool *pgxpool.Pool) func() port.AccountRepository {
	return func() port.AccountRepository {
		return gatewaydb.NewAccountRepository(pool)
	}
}

// NewTemplateRepoFactory returns a factory that creates TemplateRepository.
func NewTemplateRepoFactory(pool *pgxpool.Pool) func() port.TemplateRepository {
	return func() port.TemplateRepository {
		return gatewaydb.NewTemplateRepository(pool)
	}
}

// NewNoteRepoFactory returns a factory that creates NoteRepository.
func NewNoteRepoFactory(pool *pgxpool.Pool) func() port.NoteRepository {
	return func() port.NoteRepository {
		return gatewaydb.NewNoteRepository(pool)
	}
}
