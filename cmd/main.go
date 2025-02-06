package main

import (
	"context"
	"errors"
	"log"
	"messenger/cmd/config"
	"messenger/internal/server"
	"messenger/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

const timeOutSeconds = 10

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log := logger.NewLogger(cfg.Log.Level)
	log.Info("Starting server...")

	ws := server.NewWebSocket()
	log.Info("WebSocket created...", zap.Any("websocket", ws))

	srv := server.NewServer(cfg.Server.Port, cfg.Server.StaticDir, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err = srv.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error("Error starting server")
		}
	}()

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeOutSeconds*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error")
	}

	log.Info("Server gracefully stopped")
}
