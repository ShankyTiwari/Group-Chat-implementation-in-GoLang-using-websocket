package handlers

import (
	"bytes"
	"encoding/json"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub will will give an instance of an Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run will execute Go Routines to check incoming Socket events
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			handleUserJoinEvent(hub, client)

		case client := <-hub.unregister:
			handleUserDisconnectEvent(hub, client)

		case message := <-hub.broadcast:
			broadcastSocketEventToAllClient(hub, message)
		}
	}
}

func handleUserJoinEvent(hub *Hub, client *Client) {
	hub.clients[client] = true

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(socketEventStruct{
		EventName:    "join",
		EventPayload: client.username,
	})
	broadcastSocketEventToAllClient(hub, reqBodyBytes.Bytes())
}

func handleUserDisconnectEvent(hub *Hub, client *Client) {
	_, ok := hub.clients[client]
	if ok {

		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(socketEventStruct{
			EventName:    "disconnect",
			EventPayload: client.username,
		})
		broadcastSocketEventToAllClient(hub, reqBodyBytes.Bytes())

		delete(hub.clients, client)
		close(client.send)
	}
}

func broadcastSocketEventToAllClient(hub *Hub, message []byte) {
	for client := range hub.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(hub.clients, client)
		}
	}
}
