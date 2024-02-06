package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/arkreddy21/eligos"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	server *http.Server
	router *chi.Mux

	//websocket hub
	hub *Hub

	jwtKey []byte

	//database services
	UserService    eligos.UserServiceI
	SpaceService   eligos.SpaceServiceI
	MessageService eligos.MessageServiceI
}

func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	key, ok := os.LookupEnv("ELIGOSJWTKEY")
	if !ok {
		log.Fatal("ELIGOSJWTKEY env variable not set")
	}
	s.jwtKey = []byte(key)

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.router.Route("/api/auth", s.authRoutes)
	s.router.Get("/api/ws", s.handleWs)

	s.router.Group(func(r chi.Router) {
		r.Use(s.validateJwt)
		r.Get("/api/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
		r.Get("/api/user", s.handleUser)
		r.Route("/api/space", s.spaceRoutes)
	})

	//create a websocket hub
	s.hub = newHub()

	return s
}

func (s *Server) Open() {
	go s.hub.run(s)
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
