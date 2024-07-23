package ws

import "fmt"

type Room struct {
	RoomId string           `json:"roomId"`
	Name   string           `json:"name"`
	Users  map[string]*User `json:"users"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *User
	Unregister chan *User
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *User),
		Unregister: make(chan *User),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case user := <-h.Register:
			if r, ok := h.Rooms[user.RoomID]; ok {
				r.registerUserInRoom(user)
			}
		case user := <-h.Unregister:
			if r, ok := h.Rooms[user.RoomID]; ok {
				if msg := r.unregisterUserInRoom(user); msg != nil {
					h.Broadcast <- msg
				}
			}
		case message := <-h.Broadcast:
			if r, ok := h.Rooms[message.RoomID]; ok {
				r.broadcastToUserRoom(message)
			}
		}
	}
}

func (r *Room) registerUserInRoom(u *User) {
	if _, ok := r.Users[u.ID]; !ok {
		r.Users[u.ID] = u
	}
}

func (r *Room) unregisterUserInRoom(u *User) *Message {
	if _, ok := r.Users[u.ID]; !ok {
		delete(r.Users, u.ID)
	}
	close(u.Message)

	if len(r.Users) != 0 {
		return &Message{
			Content:  fmt.Sprintf("%s has left the group"),
			RoomID:   r.RoomId,
			Username: u.Username,
		}
	}

	return nil
}

func (r *Room) broadcastToUserRoom(message *Message) {
	for _, u := range r.Users {
		u.Message <- message
	}
}
