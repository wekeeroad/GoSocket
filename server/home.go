package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/wekeeroad/GoSocket/logic"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles(rootDir + "/template/home.html"))
	tpl.Execute(w, nil)
}

func userListHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)

	if err != nil {
		fmt.Fprint(w, `[]`)
		fmt.Println("list json.Marshal err:", err)
	} else {
		fmt.Fprint(w, string(b))
	}
}
