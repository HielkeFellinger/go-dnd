package main

import (
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/controller"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"log"
)

func main() {
	fmt.Println("Test")

	r := controller.LoadRoutes()

	helpers.ServerStaticWebContent()
	httpServer := helpers.ServeHttpServer(r)

	log.Fatal(httpServer.ListenAndServe())
}
