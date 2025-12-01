package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"webdrez/pkg/config"
	"webdrez/pkg/kick"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./web/templates/**")

	live, err := kick.IsKickLive("drezdin")

	go func() {
		for true {
			time.Sleep(2 * time.Minute)
			live, err = kick.IsKickLive("drezdin")
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			if live {
				fmt.Println("KICK: LIVE")
			} else {
				fmt.Println("KICK: OFFLINE")
			}
		}
	}()

	if err != nil {
		fmt.Println("error:", err)
	} else {
		if live {
			fmt.Println("KICK: LIVE")
		} else {
			fmt.Println("KICK: OFFLINE")
		}
	}

	config, err := config.Load("./config/main.json")

	if err != nil {

		panic(err)
	}

	data, err := os.ReadFile("data/icons.json")
	if err != nil {
		panic(err)
	}

	items := make(map[string]string)
	if err := json.Unmarshal(data, &items); err != nil {
		panic(err)
	}

	// Fix: Use template.HTML as value type from the start
	icons := map[string]template.HTML{}

	for key, value := range items {
		icons[key] = template.HTML(value) // This now compiles!
	}
	// isKick, _ := IsKickStreamLive("drezdin")

	for i, social := range config.Socials {
		key := fmt.Sprintf("brand-%s", social.Name)
		if icons[key] != "" {
			config.Socials[i].Icon = icons[key]
		}
		switch social.Name {
		case "twitch":
		}
	}

	entries, err := os.ReadDir("./web/static/")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			r.Static(e.Name(), filepath.Join("./web/static", e.Name()))
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"socials":    config.Socials,
			"icons":      icons,
			"message":    "Hello, World!",
			"kickIsLive": "blah",
		})
	})

	// 	r.Static("/static", "./web/static")
	// 	r.Static("/share", "./web/share")

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
