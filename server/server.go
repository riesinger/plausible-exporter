package server

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	s http.Server
}

func New() *Server {
	return &Server{}
}

func (srv *Server) ListenAndServe(listenAddress string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv.s = http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}
	log.Printf("Starting server, listening on %s", listenAddress)
	return srv.s.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server")
	if err := srv.s.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server stopped")
	return nil
}
