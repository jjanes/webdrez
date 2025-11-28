package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./web/templates/*")

	// isKick, _ := IsKickStreamLive("drezdin")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"message":    "Hello, World!",
			"kickIsLive": "blah",
		})
	})

	r.Static("/static", "./web/static")
	r.Static("/share", "./web/share")

	r.GET("/dev", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"message": "Hello, World!",
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/dev/hogan", func(c *gin.Context) {
		log.Println(c)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	r.GET("/dev/ws", func(c *gin.Context) {
		log.Println("abc")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading:", err)
				break
			}

			// Echo back or handle commands here
			response := handleCommand(string(msg))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("Error writing:", err)
				break
			}
		}
	})

	r.Run("0.0.0.0:8420") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all connections (be careful in production)
	},
}

func handleCommand(cmd string) string {
	switch cmd {
	case "look":
		return "You see a dusty cave with old relics."
	case "north\n":
		return "You move north into the Misty Forest."
	default:
		return fmt.Sprintf("Unknown command: %s", cmd)
	}
}
