package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	handlers "./handlers"
)

// AddApproutes will add the routes for the application
func AddApproutes(route *mux.Router) {

	log.Println("Loadeding Routes...")

	hub := handlers.NewHub()
	go hub.Run()

	route.HandleFunc("/", handlers.RenderHome)

	route.HandleFunc("/ws/{username}", func(responseWriter http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		// Reading username from request parameter
		username := mux.Vars(request)["username"]

		// Upgrading the HTTP connection socket connection
		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		handlers.CreateNewSocketUser(hub, connection, username)

	})

	log.Println("Routes are Loaded.")
}
