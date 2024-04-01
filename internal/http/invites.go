package http

import (
	"encoding/json"
	"github.com/arkreddy21/eligos"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) inviteRoutes(r chi.Router) {
	r.Post("/create", s.handleInviteCreate)
	r.Post("/accept", s.handleInviteAccept)
	r.Post("/reject", s.handleInviteReject)
	r.Get("/get", s.handleGetInvites)
}

func (s *Server) handleInviteCreate(w http.ResponseWriter, r *http.Request) {
	var invite eligos.Invite
	err := json.NewDecoder(r.Body).Decode(&invite)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := s.UserService.GetUser(invite.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.InviteService.CreateInvite(&invite)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	wsPayload, err := json.Marshal(invite)
	if err != nil {
		return
	}
	s.hub.SendMessageToUser(user.Id, "invite", wsPayload)
}

func (s *Server) handleInviteAccept(w http.ResponseWriter, r *http.Request) {
	var invite eligos.Invite
	err := json.NewDecoder(r.Body).Decode(&invite)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := s.UserService.GetUser(invite.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.SpaceService.AddUserById(user.Id, invite.SpaceId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.InviteService.DeleteInviteById(invite.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) handleInviteReject(w http.ResponseWriter, r *http.Request) {
	var invite eligos.Invite
	err := json.NewDecoder(r.Body).Decode(&invite)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.InviteService.DeleteInviteById(invite.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) handleGetInvites(w http.ResponseWriter, r *http.Request) {
	email, ok := r.URL.Query()["email"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("email query param not found"))
		return
	}
	user, err := s.UserService.GetUser(email[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	invites, err := s.InviteService.GetInvitesByUser(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := json.Marshal(invites)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
