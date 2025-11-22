// Package initializer wires dependencies for the API server.
package initializer

import (
	"context"
	"errors"
	"os"

	"github.com/labstack/echo/v4"

	httpcontroller "immortal-architecture-clean/backend/internal/adapter/http/controller"
	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	driverdb "immortal-architecture-clean/backend/internal/driver/db"
	"immortal-architecture-clean/backend/internal/driver/factory"
)

// BuildServer composes all dependencies and returns an Echo server and cleanup function.
func BuildServer(ctx context.Context) (*echo.Echo, func(), error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, func() {}, errors.New("DATABASE_URL is not set")
	}

	pool, err := driverdb.NewPool(ctx, dbURL)
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() {
		pool.Close()
	}

	txMgr := driverdb.NewTxManager(pool)

	accountRepoFactory := factory.NewAccountRepoFactory(pool)
	templateRepoFactory := factory.NewTemplateRepoFactory(pool)
	noteRepoFactory := factory.NewNoteRepoFactory(pool)
	txFactory := factory.NewTxFactory(txMgr)

	accountOutputFactory := factory.NewAccountOutputFactory()
	templateOutputFactory := factory.NewTemplateOutputFactory()
	noteOutputFactory := factory.NewNoteOutputFactory()

	accountInputFactory := factory.NewAccountInputFactory()
	templateInputFactory := factory.NewTemplateInputFactory()
	noteInputFactory := factory.NewNoteInputFactory()

	e := echo.New()
	ac := httpcontroller.NewAccountController(accountInputFactory, accountOutputFactory, accountRepoFactory)
	nc := httpcontroller.NewNoteController(noteInputFactory, noteOutputFactory, noteRepoFactory, templateRepoFactory, txFactory)
	tc := httpcontroller.NewTemplateController(templateInputFactory, templateOutputFactory, templateRepoFactory, txFactory)
	server := httpcontroller.NewServer(ac, nc, tc)
	openapi.RegisterHandlers(e, server)

	return e, cleanup, nil
}
