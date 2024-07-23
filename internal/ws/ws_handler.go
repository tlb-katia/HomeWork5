package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log/slog"
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
	Log *slog.Logger
	hub *Hub
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewHandler(log *slog.Logger, hub *Hub) *Handler {
	return &Handler{
		Log: log,
		hub: hub,
	}
}

func (h *Handler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	logResponseStatusError(h.Log, message, statusCode)
}

func logResponseStatusError(log *slog.Logger, message string, statusCode int) {
	log.Error("Request error", "status", statusCode, "error", message)
}

// CreateRoom godoc
// @Summary      create a room
// @Description  create a room with id and name
// @Tags         room
// @Accept       json
// @Produce      json
// @Param        user  body      CreateRoomReq  true  "User request body"
// @Success      201   {string}  string  "Room created successfully"
// @Failure      400   {object}  ErrorResponse    "Bad request"
// @Failure      409   {object}  ErrorResponse    "Room ID already exists"
// @Router       /rooms [post]
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomReq

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if _, exists := h.hub.Rooms[req.ID]; exists {
		h.Log.Warn("Room ID already exists", "room_id", req.ID)
		http.Error(w, `{"error": "Room ID already exists"}`, http.StatusConflict)
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		RoomId: req.ID,
		Name:   req.Name,
		Users:  make(map[string]*User),
	}

	h.Log.Info("Room created successfully", "room_id", req.ID, "room_name", req.Name)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Room created successfully"}`))
}

// JoinRoom godoc
// @Summary      Join a room
// @Description  Join an existing room using WebSocket connection with roomId, userId, and username as query parameters.
// @Tags         room
// @Accept       json
// @Produce      json
// @Param        roomId   query     string  true  "Room ID"
// @Param        userId   query     string  true  "User ID"
// @Param        username query     string  true  "Username"
// @Success      101      {string}  string  "Switching Protocols"
// @Failure      400      {object}  ErrorResponse  "Bad request"
// @Router       /rooms/join [get]
func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.sendErrorResponse(w, "Failed to join the room", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	roomID := r.URL.Query().Get("roomId")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	if roomID == "" || clientID == "" || username == "" {
		h.sendErrorResponse(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

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

	go user.writeMessage()
	user.readMessage(h.hub)

	h.Log.Info("User joined room successfully", "user_id", clientID, "room_id", roomID, "username", username)
}

func (h *Hub) GetRooms() []*RoomReq {
	allRooms := make([]*RoomReq, 0, len(h.Rooms))

	for _, room := range h.Rooms {
		allRooms = append(allRooms, &RoomReq{
			ID:   room.RoomId,
			Name: room.Name,
		})
	}

	return allRooms
}

func (h *Hub) GetUsers(roomID string) []*UserReq {
	allUsers := make([]*UserReq, 0)

	room, ok := h.Rooms[roomID]
	if !ok {
		return nil
	}

	for _, user := range room.Users {
		allUsers = append(allUsers, &UserReq{
			ID:   user.ID,
			Name: user.Username,
		})
	}

	return allUsers
}
