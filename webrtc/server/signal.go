package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	AllRooms = RoomMap{Map: make(map[string][]Participant)}
	Mutex    sync.RWMutex
)

type response struct {
	RoomID string `json:"room_id"`
}

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	log.Printf("New Room created: %s", roomID)

	if err := json.NewEncoder(w).Encode(response{RoomID: roomID}); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // real origin check in production
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func StartBroadcaster() {
	go func() {
		for msg := range broadcast {
			Mutex.RLock()
			defer Mutex.RUnlock()

			clients, ok := AllRooms.Map[msg.RoomID]
			if !ok {
				log.Printf("Room %s not found", msg.RoomID)
				continue
			}

			for _, client := range clients {
				if client.Conn != msg.Client {
					if err := client.Conn.WriteJSON(msg.Message); err != nil {
						log.Printf("Write error: %v", err)
						client.Conn.Close()
					}
				}
			}
		}
	}()
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomIDs, ok := r.URL.Query()["roomID"]
	if !ok || len(roomIDs[0]) < 1 {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}

	roomID := roomIDs[0]

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error: %v", err)
		http.Error(w, "Websocket upgrade failed", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	AllRooms.InsertIntoRoom(roomID, false, ws)

	for {
		var msg broadcastMsg
		if err := ws.ReadJSON(&msg.Message); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		msg.Client = ws
		msg.RoomID = roomID
		log.Printf("Received message in room %s: %+v", roomID, msg.Message)
		broadcast <- msg
	}
}
