package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arkreddy21/eligos"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// Hub keeps track of all active clients.
// It receives messages from every client
// and sends them only to the clients who need it
type Hub struct {
	// Registered clients.
	clients map[uuid.UUID]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan uuid.UUID
}

// WsMessage is to send/receive messages in a space
type WsMessage struct {
	Proto   string          `json:"proto"`
	Spaceid uuid.UUID       `json:"spaceid"`
	Payload json.RawMessage `json:"payload"`
}

// WsNotification is to send notifications to a user
type WsNotification struct {
	Proto   string          `json:"proto"`
	Userid  uuid.UUID       `json:"userid"`
	Payload json.RawMessage `json:"payload"`
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan uuid.UUID),
	}
}

func (h *Hub) run(s *Server) {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case clientid := <-h.unregister:
			if val, ok := h.clients[clientid]; ok {
				close(val.send)
				delete(h.clients, clientid)
			}
		case message := <-h.broadcast:
			var data WsMessage
			err := json.Unmarshal(message, &data)
			if err != nil {
				continue
			}
			res, err := handleRequest(data.Payload, data.Proto, s)
			if err != nil {
				continue
			}
			response, err := json.Marshal(WsMessage{
				Proto:   data.Proto,
				Spaceid: data.Spaceid,
				Payload: res,
			})
			users, err := s.SpaceService.GetUsersInSpace(data.Spaceid)
			if err != nil {
				continue
			}
			//TODO O(users x clients) not efficient
			for _, user := range *users {
				client, ok := h.clients[user.Id]
				if !ok {
					// user is not connected
					continue
				}
				//client.send <- response
				select {
				case client.send <- response:
				default:
					close(client.send)
					delete(h.clients, client.id)
				}
			}
		}
	}
}

// takes payload and proto from the websocket message and returns the appropriate payload to send back
func handleRequest(payload json.RawMessage, proto string, s *Server) (json.RawMessage, error) {
	switch proto {
	case "message":
		var m eligos.MessageWUser
		err := json.Unmarshal(payload, &m)
		if err != nil {
			return nil, err
		}
		message, err := s.MessageService.CreateMessage(m)
		if err != nil {
			return nil, err
		}
		response, err := json.Marshal(message)
		if err != nil {
			return nil, err
		}
		return response, nil
	default:
		return nil, fmt.Errorf("unknown proto")
	}
}

// SendMessageToUser sends message to a particular user outside of spaces.
// useful for sending notifications, invites etc
func (h *Hub) SendMessageToUser(userId uuid.UUID, proto string, payload []byte) {
	client, ok := h.clients[userId]
	if !ok {
		return
	}

	res, err := json.Marshal(WsNotification{
		Proto:   proto,
		Userid:  userId,
		Payload: payload,
	})
	if err != nil {
		return
	}

	select {
	case client.send <- res:
	default:
		close(client.send)
		delete(h.clients, client.id)
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	//user id
	id uuid.UUID

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the Hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c.id
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the Hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
