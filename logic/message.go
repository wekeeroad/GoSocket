package logic

import (
	"time"
)

const (
	MsgTypeNormal    = iota // 普通 用户消息
	MsgTypeWelcome          // 当前用户欢迎消息
	MsgTypeUserEnter        // 用户进入
	MsgTypeUserLeave        // 用户退出
	MsgTypeError            // 错误消息
)

type Message struct {
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`
	Ats     []string  `json:"ats"`
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: user.NickName + " 您好，欢迎加入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.NickName + " 离开了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
}
