package server

import (
	"context"
	"github.com/bwjson/StudyBuddy/configs"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *configs.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Server.Port,
			Handler:        handler,
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			MaxHeaderBytes: cfg.Server.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
