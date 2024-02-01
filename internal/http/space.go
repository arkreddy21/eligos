package http

import (
	"encoding/json"
	"fmt"
	"github.com/arkreddy21/eligos"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (s *Server) spaceRoutes(r chi.Router) {
	r.Post("/create", s.handleCreateSpace)
	r.Post("/adduser", s.handleAddUserToSpace)
	// returns all spaces that a user belongs to
	r.Get("/spaces", s.handleGetSpaces)
}

func (s *Server) handleCreateSpace(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name   string
		Userid uuid.UUID
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if body.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("please provide a name"))
		return
	}
	err = s.SpaceService.CreateSpace(&eligos.Space{Name: body.Name}, body.Userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create space"))
		fmt.Println(err)
		return
	}
	response, err := json.Marshal(map[string]string{
		"status": "ok",
	})
	w.Write(response)
}

func (s *Server) handleAddUserToSpace(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email   string
		SpaceId uuid.UUID
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := s.UserService.GetUser(body.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user not found"))
		return
	}
	err = s.SpaceService.AddUserById(user.Id, body.SpaceId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to add user to space"))
		return
	}
	response, err := json.Marshal(map[string]string{
		"status": "ok",
	})
	w.Write(response)
}

// returns all spaces that a user belongs to
func (s *Server) handleGetSpaces(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)
	uid, _ := uuid.Parse(userId)
	spaces, err := s.UserService.GetSpaces(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse body"))
		return
	}
	response, _ := json.Marshal(*spaces)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
