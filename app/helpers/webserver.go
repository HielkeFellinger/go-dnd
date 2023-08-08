package helpers

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func ServerStaticWebContent() {
	fs := http.FileServer(http.Dir("../web/static/"))
	http.Handle("static/", http.StripPrefix("/static/", fs))
}

func ServeHttpServer(router *mux.Router) *http.Server {
	return &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
	}
}
