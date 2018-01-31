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

	router := gin.Default()

	router.LoadHTMLGlob("view/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Address": fmt.Sprintf("%s%s", address, port),
		})
	})

	router.GET("/chat/ws", func(c *gin.Context) {
		handler := websocket.Handler(connect)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	go broadcastMessage()

	router.Run(port)
}

var incomingMessage = make(chan string)
var connectionNum = 0

type ChatUser struct {
	Name string
}

const (
	UserMessage = "MESSAGE"
	UserName    = "NAME"
)

type ReceiveMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var chatUsers = make(map[*websocket.Conn]ChatUser)

func connect(connection *websocket.Conn) {
	connectionNum++
	chatUser := ChatUser{
		Name: "",
	}
	chatUsers[connection] = chatUser
	fmt.Println("Client Added")
	for {
		var receiveMessage ReceiveMessage

		if err := websocket.JSON.Receive(connection, &receiveMessage); err != nil {
			log.Println(err)
			delete(chatUsers, connection)
			break
		}

		fmt.Printf("\n%#v\n", receiveMessage)
		if receiveMessage.Type == UserMessage {
			incomingMessage <- fmt.Sprintf("%s: %s", chatUser.Name, receiveMessage.Text)
		} else if receiveMessage.Type == UserName {
			chatUser.Name = receiveMessage.Text
		}

	}
}

func broadcastMessage() {
	for {
		message := <-incomingMessage
		fmt.Println(message)
		for connection := range chatUsers {
			if err := websocket.Message.Send(connection, message); err != nil {
				log.Println(err)
				delete(chatUsers, connection)
				break
			}
		}
	}

}
