package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
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
	Token          string        `json:"token"`

	conn  *websocket.Conn
	isNew bool
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

func NewUser(conn *websocket.Conn, token, nickname string, addr string) *User {
	user := &User{
		NickName:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),
		Token:          token,
		conn:           conn,
	}

	if user.Token != "" {
		uid, err := parseTokenAndValidation(token, nickname)
		if err == nil {
			user.UID = uid
		}
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddInt32(&globalUID, 1))
		user.Token = genToken(user.UID, user.NickName)
		user.isNew = true
	}

	return user
}

func genToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)

	messageMAC := macSha256([]byte(message), []byte(secret))
	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func parseTokenAndValidation(token string, nickname string) (int, error) {
	pos := strings.LastIndex(token, "uid")
	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	ok := validateMAC([]byte(message), messageMAC, []byte(secret))
	if ok {
		return uid, nil
	}
	return 0, errors.New("token is illegal")

}

func validateMAC(message, messageMAC, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
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
