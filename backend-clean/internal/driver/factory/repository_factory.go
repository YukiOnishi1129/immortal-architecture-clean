package factory

import (
	"github.com/jackc/pgx/v5/pgxpool"

	gatewaydb "immortal-architecture-clean/backend/internal/adapter/gateway/db"
	"immortal-architecture-clean/backend/internal/port"
)

func NewAccountRepoFactory(pool *pgxpool.Pool) func() port.AccountRepository {
	return func() port.AccountRepository {
		return gatewaydb.NewAccountRepository(pool)
	}
}

func NewTemplateRepoFactory(pool *pgxpool.Pool) func() port.TemplateRepository {
	return func() port.TemplateRepository {
		return gatewaydb.NewTemplateRepository(pool)
	}
}

func NewNoteRepoFactory(pool *pgxpool.Pool) func() port.NoteRepository {
	return func() port.NoteRepository {
		return gatewaydb.NewNoteRepository(pool)
	}
}
