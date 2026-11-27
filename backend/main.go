package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/tekmen0/simple-calculator/backend/internal/httpapi"
)

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpapi.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}

	slog.Info("starting calculator API", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("calculator API stopped", "err", err)
		os.Exit(1)
	}
}
