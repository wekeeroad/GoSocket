package server

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("illegal nick name!"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("nickname exists:", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("The nickname exists!"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, nickname, req.RemoteAddr)
	go user.SendMessage(req.Context())
	user.MessageChannel <- logic.NewWelcomeMessage(nickname)

	msg := logic.NewNoticeMessage(nickname + "joined ChatRoom")
	logic.Broadcaster.Broadcast(msg)

	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	err = user.ReceiveMessage(req.Context())

	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewNoticeMessage(user.NickName + "leaved ChatRoom")
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
