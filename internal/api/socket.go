package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"blockbuffer/internal/io"
	store "blockbuffer/internal/store"
	types "blockbuffer/internal/types"
	"github.com/gorilla/websocket"
)

const (
	pollingInterval = 2 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan types.Message)
var outboundMessages = make(map[string]time.Time)

func HandleSocketConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		var msg = "Failed to upgrade connection"
		io.Log(msg, io.Error)
		io.ErrorJSON(w, msg, http.StatusInternalServerError)
		return
	}

	defer ws.Close()
	clients[ws] = true
	ws.WriteJSON(types.Message{MessageType: types.RefreshFiles, Data: store.FileList})
	io.Log("new socket connection established", io.Info)
	for {
		var messages []types.File
		err := ws.ReadJSON(&messages)
		if err != nil {
			delete(clients, ws)
			return
		}
	}
}

func HandleMessages() {
	for {
		// get next message
		message := <-broadcast

		// verify should send message
		hash := hashMessage(message)
		if message.MessageType != types.RefreshFiles && !message.MustSend {
			if lastSent, exists := outboundMessages[hash]; exists && time.Since(lastSent) < pollingInterval {
				continue
			}
		}

		// send message to all clients, close connection if error
		outboundMessages[hash] = time.Now()
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(clients, client)
				io.Log("socket connection closed", io.Info)
			}
		}
	}
}

func hashMessage(message types.Message) string {
	h := sha256.New()
	h.Write([]byte(message.MessageType))
	if dataMap, ok := message.Data.(map[string]interface{}); ok {
		for key := range dataMap {
			h.Write([]byte(key))
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

func BroadcastMessage(message types.Message) {
	broadcast <- message
}

func BroadcastFiles(fileList map[string]types.File) {
	broadcast <- types.Message{MessageType: types.UpdateFile, Data: fileList}
}
