package controller

import (
	"github.com/gorilla/mux"
	"github.com/hielkefellinger/go-dnd/app/root"
)

func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", root.Frontpage).Methods("GET")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	return r
}
