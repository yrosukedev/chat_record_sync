package business

import "time"

type ChatRecord struct {
	MsgId   string
	Action  string
	From    *User
	To      []*User
	Room    *Room
	MsgTime time.Time
	MsgType string
	Content string
}

type User struct {
	UserId string
	Name   string
}

type Room struct {
	RoomId string
	Name   string
}
