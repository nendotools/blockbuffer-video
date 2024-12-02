package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	RefreshFiles MessageType = "refresh_files"
	UpdateFile   MessageType = "update_file"
	CreateFile   MessageType = "create_file"
	DeleteFile   MessageType = "delete_file"
)

type Message struct {
	MessageType MessageType `json:"type"`
	Data        interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var outboundMessages = make(map[string]time.Time)

func HandleSocketConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection")
		ErrorJSON(w, fmt.Sprintf("Failed to upgrade connection: %v", err), http.StatusInternalServerError)
		return
	}

	defer ws.Close()
	clients[ws] = true
	ws.WriteJSON(Message{MessageType: RefreshFiles, Data: fileList})
	log.Println("Connected to server")
	for {
		var messages []File
		err := ws.ReadJSON(&messages)
		if err != nil {
			delete(clients, ws)
			return
		}
	}
}

func HandleMessages() {
	for {
		message := <-broadcast
		hash := hashMessage(message)
		if message.MessageType != RefreshFiles {
			if lastSent, exists := outboundMessages[hash]; exists && time.Since(lastSent) < 2*time.Second {
				continue
			}
		}

		outboundMessages[hash] = time.Now()
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func hashMessage(message Message) string {
	h := sha256.New()
	h.Write([]byte(message.MessageType))
	if dataMap, ok := message.Data.(map[string]interface{}); ok {
		for key := range dataMap {
			h.Write([]byte(key))
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

func BroadcastMessage(message Message) {
	broadcast <- message
}

func BroadcastFiles(fileList map[string]File) {
	broadcast <- Message{MessageType: UpdateFile, Data: fileList}
}
