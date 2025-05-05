package server

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}

	return string(b)
}

func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	roomID := randomString(8)
	r.Map[roomID] = []Participant{}

	return roomID
}

func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{host, conn}

	log.Printf("Inserting into room with ID: %s", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}
