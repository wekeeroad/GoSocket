package logic

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/wekeeroad/GoSocket/global"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	filter "github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
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
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sondMsg.Ats = reg.FindAllString(sondMsg.Content, -1)
		sondMsg.Content, err = FilterSensitive(sondMsg.Content)
		if err != nil {
			return err
		}
		fmt.Println(sondMsg)
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

func MatchSensitive(content string) ([]string, error) {
	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: global.SensitiveWords,
	})
	if err != nil {
		return nil, err
	}
	fileterManage := filter.NewDirtyManager(memStore)
	result, err := fileterManage.Filter().Filter(content, '#')
	if err != nil {
		return nil, err
	}
	return result, nil
}

func FilterSensitive(content string) (string, error) {
	matchWords, err := MatchSensitive(content)
	if err != nil {
		return "", err
	}
	var matchSlice = make([]string, 10)
	for _, i := range matchWords {
		if len(i) > 1 {
			exp := []string{}
			for _, j := range i {
				exp = append(exp, string(j))
			}
			expre := strings.Join(exp, "#*")
			fmt.Println(expre)
			reg := regexp.MustCompile(expre)
			matchPart := reg.FindAllString(content, -1)
			matchSlice = append(matchSlice, matchPart...)
		} else {
			matchSlice = append(matchSlice, i)
		}
	}
	for _, i := range matchSlice {
		if i == "" {
			continue
		}
		content = strings.ReplaceAll(content, i, "* *")
	}

	return content, nil
}
