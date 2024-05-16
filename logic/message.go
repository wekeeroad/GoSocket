package logic

import "time"

const (
	MsgTypeNormal = iota
	MsgTypeSystem
	MsgTypeError
	MsgTypeUserList
)

type Message struct {
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	Users map[string]*User `json:"users"`
}
