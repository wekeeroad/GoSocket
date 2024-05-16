package server

import (
	"html/template"
	"net/http"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles(rootDir + "/template/home.html"))
	tpl.Execute(w, nil)
}
