package main

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v4"

	gatewaydb "immortal-architecture-clean/backend/internal/adapter/gateway/db"
	httpcontroller "immortal-architecture-clean/backend/internal/adapter/http/controller"
	openapi "immortal-architecture-clean/backend/internal/adapter/http/generated/openapi"
	driverdb "immortal-architecture-clean/backend/internal/driver/db"
	"immortal-architecture-clean/backend/internal/usecase"
)

func main() {
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := driverdb.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	txMgr := driverdb.NewTxManager(pool)

	// Repositories
	accountRepo := gatewaydb.NewAccountRepository(pool)
	templateRepo := gatewaydb.NewTemplateRepository(pool)
	noteRepo := gatewaydb.NewNoteRepository(pool)

	// UseCases
	accountUC := usecase.NewAccountInteractor(accountRepo)
	templateUC := usecase.NewTemplateInteractor(templateRepo, txMgr)
	noteUC := usecase.NewNoteInteractor(noteRepo, templateRepo, txMgr)

	// HTTP server wiring
	e := echo.New()
	ac := httpcontroller.NewAccountController(accountUC)
	nc := httpcontroller.NewNoteController(noteUC)
	tc := httpcontroller.NewTemplateController(templateUC)
	server := httpcontroller.NewServer(ac, nc, tc)
	openapi.RegisterHandlers(e, server)

	addr := ":8080"
	log.Printf("starting HTTP server at %s\n", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
