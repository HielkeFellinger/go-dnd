package root

import (
	"html/template"
	"net/http"
)

func Frontpage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/index.html"))
	_ = tmpl.Execute(w, nil)
}

// resembling
