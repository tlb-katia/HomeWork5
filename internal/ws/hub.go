package ws

type Room struct {
	RoomId string           `json:"roomId"`
	Name   string           `json:"name"`
	Users  map[string]*User `json:"users"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *User
	Unregister chan *User
	Broadcast  chan *Message
}
