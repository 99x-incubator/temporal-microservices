package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection with a specific userID
type Client struct {
	UserID string
	Conn   *websocket.Conn
}

// NotificationHub manages all WebSocket connections
type NotificationHub struct {
	Clients map[string]*Client
	Mutex   sync.Mutex
}

// Global notification hub
var hub = NotificationHub{
	Clients: make(map[string]*Client),
}

// NotificationMessage represents the structure of a message to be sent
type NotificationMessage struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// handleWebSocket handles incoming WebSocket requests and registers a client with its userID
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get userID from query parameters
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	client := &Client{UserID: userID, Conn: conn}

	// Add the client to the hub
	hub.Mutex.Lock()
	hub.Clients[userID] = client
	hub.Mutex.Unlock()

	log.Println("Client connected:", userID)

	// Listen for messages from the client (optional)
	defer func() {
		// Remove the client from the hub upon disconnect
		hub.Mutex.Lock()
		delete(hub.Clients, userID)
		hub.Mutex.Unlock()

		conn.Close()
		log.Println("Client disconnected:", userID)
	}()

	// Keep the connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket error:", err)
			break
		}
	}
}

// sendNotification sends a message to a specific user
func sendNotification(w http.ResponseWriter, r *http.Request) {
	var notification NotificationMessage
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil || notification.UserID == "" || notification.Message == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the client by userID
	hub.Mutex.Lock()
	client, exists := hub.Clients[notification.UserID]
	hub.Mutex.Unlock()

	if !exists {
		http.Error(w, "User not connected", http.StatusNotFound)
		return
	}

	// Send the notification to the user's WebSocket connection
	err = client.Conn.WriteJSON(notification)
	if err != nil {
		log.Println("Failed to send notification:", err)
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification sent to user: " + notification.UserID))
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// WebSocket connection endpoint
	r.HandleFunc("/ws", handleWebSocket).Methods("GET")

	// HTTP endpoint to send a notification
	r.HandleFunc("/notify", sendNotification).Methods("POST")

	// Start the server
	log.Println("Notification service started on :8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
