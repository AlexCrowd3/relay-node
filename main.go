package main

import (
	"log"
	"net/http"

	"relay-node/server"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(hub *server.Hub, w http.ResponseWriter, r *http.Request) {

	pubKey := r.URL.Query().Get("pubKey")

	if pubKey == "" {

		http.Error(w, "missing pubKey", 400)

		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {

		log.Println(err)

		return
	}

	client := server.NewClient(pubKey, hub, conn)

	hub.Register <- client

	go client.ReadPump()

	go client.WritePump()
}

func main() {

	hub := server.NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		serveWs(hub, w, r)
	})

	log.Println("Relay node started on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {

		log.Fatal(err)
	}
}
