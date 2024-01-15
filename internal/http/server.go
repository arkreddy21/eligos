package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	router *chi.Mux
}

func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	return s
}

func (s *Server) Open() {
	fmt.Println("listening on port 4000")
	s.server = &http.Server{Addr: "0.0.0.0:4000", Handler: s.router}
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("unable to run http server: ", err)
		}
	}()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	return err
}
