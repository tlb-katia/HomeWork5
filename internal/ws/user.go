package ws

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"roomId"`
	Message  chan *Message
	Con      *websocket.Conn
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}
