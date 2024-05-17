package server

import (
	"net/http"
	"os"
	"path/filepath"
)

var rootDir string

func RegisterHandle() {
	inferRootDir()

	//go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
}

func inferRootDir() {
	cmd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/template") {
			return d
		}
		return infer(filepath.Dir(d))
	}
	rootDir = infer(cmd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
