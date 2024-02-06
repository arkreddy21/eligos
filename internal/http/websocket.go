package http

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173"
	},
}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("token not provided"))
		return
	}
	token, err := jwt.ParseWithClaims(keys[0], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("user unauthorized"))
		return
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	userid, err := uuid.Parse(claims.Issuer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := Client{hub: s.hub, conn: conn, id: userid, send: make(chan []byte, 256)}
	s.hub.register <- &client
	go client.writePump()
	go client.readPump()
}
