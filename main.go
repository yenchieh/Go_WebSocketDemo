package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/websocket"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	if address == "" {
		address = "localhost"
	}
	if port == "" {
		port = ":8080"
	}

	router := gin.Default()

	router.LoadHTMLGlob("view/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Address": fmt.Sprintf("%s%s", address, port),
		})
	})

	// WebSocket
	router.GET("/chat/ws", func(c *gin.Context) {
		handler := websocket.Handler(connect)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	// Create goroutine to listen message channel and broadcast it to all of connection
	go broadcastMessage()

	router.Run(port)
}

/**
 * Chat Room
 */

type ChatUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const (
	UserMessage = "MESSAGE"
	UserName    = "NAME"
)

type ReceiveMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SendMessage struct {
	UserID   string     `json:"userId"`
	Type     string     `json:"type"`
	Message  string     `json:"message"`
	UserList []ChatUser `json:"userList"`
}

var incomingMessage = make(chan SendMessage)

var connectionNum = 0
var chatUsers = make(map[*websocket.Conn]*ChatUser)

func connect(connection *websocket.Conn) {
	connectionNum++
	chatUser := ChatUser{
		ID:   fmt.Sprintf("%d", connectionNum),
		Name: "",
	}
	chatUsers[connection] = &chatUser

	fmt.Println("Client Added")

	for {
		var receiveMessage ReceiveMessage

		if err := websocket.JSON.Receive(connection, &receiveMessage); err != nil {
			delete(chatUsers, connection)
			message := SendMessage{
				UserID: chatUser.ID,
				Type:   "offline",
			}
			incomingMessage <- message
			break
		}

		fmt.Printf("\n%#v\n", receiveMessage)
		if receiveMessage.Type == UserMessage {
			incomingMessage <- SendMessage{
				Type:    "message",
				Message: fmt.Sprintf("%s: %s", chatUser.Name, receiveMessage.Text),
			}
		} else if receiveMessage.Type == UserName {
			chatUser.Name = receiveMessage.Text

			var userList []ChatUser
			for _, g := range chatUsers {
				userList = append(userList, *g)
			}

			incomingMessage <- SendMessage{
				Type:     "online",
				UserID:   chatUser.ID,
				UserList: userList,
			}
		}

	}
}

func broadcastMessage() {
	for {
		message := <-incomingMessage
		fmt.Println(message)
		for connection := range chatUsers {
			if err := websocket.JSON.Send(connection, message); err != nil {
				log.Println(err)
				delete(chatUsers, connection)
				break
			}
		}
	}

}
