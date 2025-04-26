package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// type which have a upgrader type to upgrade http to websocket
type websocketHandler struct {
	upgrader websocket.Upgrader
}

// method to upgrade http to websocket
func (wsh websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Socket upgrade error: %v", err)
		return
	}

	defer func() {
		log.Println("Closing connection...")
		c.Close()
	}()

	// code to handle messaging between server and client through websocket
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client:%v", err)
			return
		}

		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("Server doesn't allow binary message"))
			if err != nil {
				log.Printf("Error writing messgae: %v", err)
				return
			}
		}

		if strings.Trim(string(message), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You didn't send the magic word"))
			if err != nil {
				log.Printf("Error writing messgae: %v", err)
				return
			}
			continue
		}

		log.Println("Sending message")
		i := 1
		for {
			response := fmt.Sprintf("Notification %d", i)
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error writing messgae: %v", err)
				return
			}
		}
	}
}

func main() {
	mySocketHandler := websocketHandler{
		upgrader: websocket.Upgrader{},
	}

	http.Handle("/", mySocketHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
