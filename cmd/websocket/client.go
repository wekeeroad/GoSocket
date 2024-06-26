package main

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil)
	if err != nil {
		panic(err)
	}

	defer c.Close(websocket.StatusInternalError, "internal err")

	err = wsjson.Write(ctx, c, "Hello WebSocker Server")
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Recived response from server: %v\n", v)
	c.Close(websocket.StatusNormalClosure, "")
}
