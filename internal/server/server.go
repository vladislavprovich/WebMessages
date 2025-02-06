package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

const timeSeconds = 50

type Server struct {
	logger    *zap.Logger
	port      string
	staticDir string
	httpSrv   *http.Server
}

func NewServer(port, staticDir string, logger *zap.Logger) *Server {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	ws := NewWebSocket()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/ws", ws.WebSocketHandler)

	return &Server{
		logger:    logger,
		port:      port,
		staticDir: staticDir,
		httpSrv: &http.Server{
			Addr:              fmt.Sprintf(":%s", port),
			Handler:           r,
			ReadTimeout:       timeSeconds * time.Second,
			WriteTimeout:      timeSeconds * time.Second,
			IdleTimeout:       timeSeconds * time.Second,
			ReadHeaderTimeout: timeSeconds * time.Second,
		},
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	s.logger.Info("Starting server...", zap.String("address", "http://localhost"+addr))

	err := s.httpSrv.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal("Server failed to start", zap.Error(err))
		return err
	}

	s.logger.Info("Server started successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Stopping server...")
	return s.httpSrv.Shutdown(ctx)
}
