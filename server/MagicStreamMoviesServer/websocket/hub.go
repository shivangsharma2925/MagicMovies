package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var mutex sync.Mutex

// creates the parameters for the http Upgarder
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(c *gin.Context) {
	// Upgrade is simply the http on steroids where it don't close the connection and keep it alive
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	for {
		// keep connection alive
		if _, _, err := conn.NextReader(); err != nil {
			// log.Printf("error: %v", err)
			break
		}
	}
}

func BroadcastJobUpdate(data any) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
