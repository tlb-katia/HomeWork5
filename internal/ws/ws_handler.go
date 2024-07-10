package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *Hub
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomReq

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		RoomId: req.ID,
		Name:   req.Name,
		Users:  make(map[string]*User),
	}
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomID := r.URL.Query().Get("roomId")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	user := &User{
		ID:       clientID,
		Username: username,
		RoomID:   roomID,
		Message:  make(chan *Message),
		Con:      ws,
	}

	message := &Message{
		Content:  fmt.Sprintf("%s has joined the group", username),
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- user
	h.hub.Broadcast <- message

	// TODO go routine writeMessage
	// TODO readMessage
}
