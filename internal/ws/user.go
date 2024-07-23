package ws

import (
	"github.com/gorilla/websocket"
	"log"
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

func (u *User) writeMessage() {
	defer func() {
		u.Con.Close()
	}()

	for {
		if message, ok := <-u.Message; ok {
			u.Con.WriteJSON(message)
		} else {
			return
		}
	}
}

func (u *User) readMessage(h *Hub) {
	defer func() {
		h.Unregister <- u
		u.Con.Close()
	}()

	for {
		_, message, err := u.Con.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("readMessageError: %v", err)
			}
			break
		}
		h.Broadcast <- &Message{
			Content:  string(message),
			RoomID:   u.RoomID,
			Username: u.Username,
		}
	}
}
