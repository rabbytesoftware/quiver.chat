package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed static/*
var staticFiles embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Time     string `json:"time"`
}

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan Message
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message, 256), // Make broadcast channel buffered
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			
			// Send user joined message
			joinMsg := Message{
				Type:     "system",
				Username: "System",
				Content:  fmt.Sprintf("%s joined the chat", client.username),
				Time:     getCurrentTime(),
			}
			h.broadcast <- joinMsg

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				
				// Send user left message
				leftMsg := Message{
					Type:     "system",
					Username: "System",
					Content:  fmt.Sprintf("%s left the chat", client.username),
					Time:     getCurrentTime(),
				}
				h.broadcast <- leftMsg
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func getCurrentTime() string {
	return time.Now().Format("15:04")
}

func (c *Client) writePump() {
	defer c.conn.Close()
	
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			c.conn.WriteJSON(message)
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		msg.Username = c.username
		msg.Time = getCurrentTime()
		
		hub.broadcast <- msg
	}
}

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		username = "Anonymous"
	}

	client := &Client{
		conn:     conn,
		username: username,
		send:     make(chan Message, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

func main() {
	hub := newHub()
	go hub.run()

	// Get the static subdirectory from embedded filesystem
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to create static filesystem:", err)
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.FS(staticFS)))
	
	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	port := "8080"
	fmt.Printf("🚀 Quiver Chat Server running on http://localhost:%s\n", port)
	fmt.Printf("📱 Open your browser and navigate to the URL above\n")
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
} 