package logic

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

var globalUID int32 = 0
var System = &User{}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}

		sondMsg := NewMessage(u, receiveMsg["content"])
		Broadcaster.Broadcast(sondMsg)
	}
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	user := &User{
		NickName:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),
		conn:           conn,
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddInt32(&globalUID, 1))
	}

	return user
}
