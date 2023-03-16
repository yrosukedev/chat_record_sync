package business

import "time"

type ChatRecord struct {
	MsgId   string
	Action  string
	From    *User
	To      []*User
	RoomId  string
	MsgTime time.Time
	MsgType string
	Content string
}

type User struct {
	UserId string
	Name   string
}
