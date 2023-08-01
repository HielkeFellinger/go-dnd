package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Test")

	root := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		_ = tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
