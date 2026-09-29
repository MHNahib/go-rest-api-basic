package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MHNahib/rest-api/internal/config"
	"github.com/MHNahib/rest-api/internal/http/handlers/todos"
	"github.com/MHNahib/rest-api/internal/storage"
	"github.com/MHNahib/rest-api/internal/storage/sqlite"
	"github.com/MHNahib/rest-api/internal/utils/response"
)

func main() {
	fmt.Println("Bismillah")

	appConfig := config.MountConfig()

	fmt.Printf("app is on env: %s mode\n", appConfig.Env)

	database, err := sqlite.New(appConfig)

	if err != nil {
		slog.Error("cannot create database", "err", err.Error())
		os.Exit(1)
	}

	defer database.Db.Close()

	slog.Info("database is ready")

	server := createSever(appConfig, database)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("cannot start server", "err", err.Error())
		}
		slog.Info("server is ready on", "address", appConfig.Server.Address)
	}()
	slog.Info("server is ready on", "address", appConfig.Server.Address)

	<-done

	slog.Info("Shutting down the server!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", "error", err)
	}

	slog.Info("Server stopped")
}

func createSever(config *config.Config, database storage.Storage) *http.Server {
	server := http.Server{
		Addr:    config.Server.Address,
		Handler: appRouter(database),
	}
	return &server
}

func appRouter(database storage.Storage) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJson(w, http.StatusOK, nil, "success")
	})

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJson(w, http.StatusOK, nil, "success")
	})

	router.HandleFunc("GET /todo", todos.GetTodosHandler(database))
	router.HandleFunc("POST /todo", todos.CreateTodoHandler(database))
	router.HandleFunc("PUT /todo", todos.UpdateTodoHandler(database))
	router.HandleFunc("POST /todo/delete", todos.DeleteTodoHandler(database))

	return router
}
