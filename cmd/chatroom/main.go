package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wekeeroad/GoSocket/server"
)

var (
	addr   = ":2022"
	banner = `
	++++++++++++++++++++++++++++++++
	+    Golang program travel     +
	+    ChatRoom, start on: %s    +
	++++++++++++++++++++++++++++++++
`
)

func main() {
	fmt.Printf(banner+"\n", addr)
	server.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}
