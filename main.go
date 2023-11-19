package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Msg struct {
	Type     string `json:"type"`
	NewCount int    `json:"newCount"`
}

var (
	pongCount   = 0
	pingCount   = 0
	clients     = make(map[*websocket.Conn]bool)
	clientMutex = sync.Mutex{}
	broadcast   = make(chan Msg)
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length"},
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"count": pingCount,
		})
	})
	r.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"count": pongCount,
		})
	})

	r.POST("/ping", func(c *gin.Context) {
		pingCount++

		broadcast <- Msg{
			Type:     "ping",
			NewCount: pingCount,
		}

		c.JSON(200, gin.H{
			"count": pingCount,
		})
	})
	r.POST("/pong", func(c *gin.Context) {
		pongCount++

		broadcast <- Msg{
			Type:     "pong",
			NewCount: pongCount,
		}

		c.JSON(200, gin.H{
			"count": pongCount,
		})
	})

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			println(err.Error())
			return
		}

		clientMutex.Lock()
		clients[conn] = true
		clientMutex.Unlock()

		go handleWebsocket(conn)
	})

	go broadcaster()

	r.Run() // listen and serve on 0.0.0.0:8080
}

func handleWebsocket(conn *websocket.Conn) {
	defer func() {
		clientMutex.Lock()
		defer clientMutex.Unlock()
		delete(clients, conn)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			println(err.Error())
			return
		}
	}
}

func broadcaster() {
	for {
		msg := <-broadcast
		clientMutex.Lock()
		for client := range clients {
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				println(err.Error())
				continue
			}

			err = client.WriteMessage(websocket.TextMessage, jsonMsg)
			if err != nil {
				println(err.Error())
				client.Close()
				delete(clients, client)
			}
		}
		clientMutex.Unlock()
	}
}
