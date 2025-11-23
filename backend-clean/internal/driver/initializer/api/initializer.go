// Package initializer wires dependencies for the API server.
package initializer

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	// Allow frontend (localhost:3000) to call the API during development.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins(),
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	ac := httpcontroller.NewAccountController(accountInputFactory, accountOutputFactory, accountRepoFactory)
	nc := httpcontroller.NewNoteController(noteInputFactory, noteOutputFactory, noteRepoFactory, templateRepoFactory, txFactory)
	tc := httpcontroller.NewTemplateController(templateInputFactory, templateOutputFactory, templateRepoFactory, txFactory)
	server := httpcontroller.NewServer(ac, nc, tc)
	openapi.RegisterHandlers(e, server)

	return e, cleanup, nil
}

func allowedOrigins() []string {
	fromEnv := os.Getenv("CLIENT_ORIGIN")
	if strings.TrimSpace(fromEnv) == "" {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
	parts := strings.Split(fromEnv, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
	return origins
}
