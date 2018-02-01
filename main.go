package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("view/*")

	// Chat room page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// WebSocket
	router.GET("/chat/ws", func(c *gin.Context) {
		handler := websocket.Handler(connect)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	// Create goroutine to listen message channel and broadcast it to all of connection
	go broadcastMessage()

	router.Run(":8080")
}

/**
 * Chat Room handler
 * Client map store client number (ID) using connection as key
 */
var clients = make(map[*websocket.Conn]string)
var incomingMessage = make(chan string)
var connectionNum = 0

func connect(connection *websocket.Conn) {
	connectionNum++
	clients[connection] = fmt.Sprintf("Client-%d", connectionNum)
	fmt.Printf("\n%s Added\n", clients[connection])
	for {
		var message string

		if err := websocket.Message.Receive(connection, &message); err != nil {
			log.Println(err)
			delete(clients, connection)
			break
		}
		fmt.Printf("Received Message From %s \n - Message: %s\n", clients[connection], message)

		message = fmt.Sprintf("%s: %s", clients[connection], message)

		incomingMessage <- message
	}
}

func broadcastMessage() {
	for {
		message := <-incomingMessage
		fmt.Println(message)
		for client := range clients {
			if err := websocket.Message.Send(client, message); err != nil {
				log.Println(err)
				delete(clients, client)
				break
			}
		}
	}

}
