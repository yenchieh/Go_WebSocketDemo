package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"

	"log"

	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("view/*")

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/chat/ws", func(c *gin.Context) {
		handler := websocket.Handler(connect)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	go testMessage()

	router.Run(":8080")
}

var clients = make(map[*websocket.Conn]string)
var incomingMessage = make(chan string)

func testMessage() {
	for i := 0; ; i++ {

		broadcastMessage(fmt.Sprintf("Message|%d", i))
		time.Sleep(time.Second * 5)
	}
}

func connect(connection *websocket.Conn) {
	clients[connection] = fmt.Sprintf("Client-%d", len(clients))
	fmt.Println("Client Added")
	for {
		var message string

		if err := websocket.JSON.Receive(connection, &message); err != nil {
			log.Println(err)
			delete(clients, connection)
			break
		}
		clients[connection] = message
		fmt.Println("Received Message", message)

		incomingMessage <- message
	}
}

func broadcastMessage(message string) {
	for client, name := range clients {
		if err := websocket.Message.Send(client, fmt.Sprintf("%s : %s", name, message)); err != nil {
			log.Println(err)
			delete(clients, client)
			break
		}
	}
}
