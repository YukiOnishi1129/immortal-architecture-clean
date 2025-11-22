// Package main starts the API server.
package main

import (
	"context"
	"log"

	"immortal-architecture-clean/backend/internal/driver/initializer"
)

func main() {
	ctx := context.Background()
	e, cleanup, err := initializer.BuildServer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}
	defer cleanup()

	addr := ":8080"
	log.Printf("starting HTTP server at %s\n", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
