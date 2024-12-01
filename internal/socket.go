package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []File)

func HandleSocketConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection")
		ErrorJSON(w, fmt.Sprintf("Failed to upgrade connection: %v", err), http.StatusInternalServerError)
		return
	}
	err = ws.WriteJSON(fileList)
	if err != nil {
		log.Println("Failed to write to client")
		ws.Close()
		return
	}

	log.Println("Connected to server")
	defer ws.Close()
	clients[ws] = true
	for {
		var files []File
		err := ws.ReadJSON(&files)
		if err != nil {
			delete(clients, ws)
			return
		}
		broadcast <- files
	}
}

func HandleMessages() {
	for {
		files := <-broadcast
		for client := range clients {
			err := client.WriteJSON(files)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func BroadcastFiles(fileList map[string]File) {
	files := make([]File, 0, len(fileList))
	for _, file := range fileList {
		files = append(files, file)
	}
	broadcast <- files
}
