package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"motor-de-rol/backend/db"
	"motor-de-rol/backend/repository"
	"motor-de-rol/backend/transport"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	host = "localhost"
	port = 8080
)

func main() {
	database, err := db.Connect("motor-de-rol.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	characterRepo := repository.NewSQLiteCharacterRepo(database)
	charIdentityRepo := repository.NewSQLiteCharacterIdentityRepo(database)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	api := humachi.New(router, huma.DefaultConfig("Motor de Rol API", "1.0.0"))

	characterController := transport.NewCharacterController(characterRepo)
	characterController.Register(api)

	charIdentityController := transport.NewCharacterIdentityController(charIdentityRepo)
	charIdentityController.Register(api)

	printBanner()

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Server starting on http://%s", addr)
	log.Printf("API Documentation: http://%s/docs", addr)
	log.Printf("OpenAPI JSON: http://%s/openapi.json", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	sqlDB, dbErr := database.DB()
	if dbErr == nil {
		sqlDB.Close()
	}

	log.Println("Server stopped gracefully")
}

func printBanner() {
	fmt.Print(`
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║           MOTOR DE ROL - API SERVER                      ║
║                                                          ║
║   REST API:    http://localhost:8080                     ║
║   Docs:        http://localhost:8080/docs                ║
║   OpenAPI:     http://localhost:8080/openapi.json        ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
`)
}
