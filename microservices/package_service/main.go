package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type User struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	CurrentPackage RobotPackage `json:"current_package"`
}

type RobotPackage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

var (
	SDStreamingPackage = RobotPackage{ID: "2", Name: "SD Streaming", Description: "Standard Definition Video", IsActive: true}
	HDStreamingPackage = RobotPackage{ID: "3", Name: "HD Streaming", Description: "High Definition Video", IsActive: false}

	// In-memory user store
	users = map[string]User{
		"123": {ID: "123", Name: "John Doe", CurrentPackage: SDStreamingPackage},
	}

	mu sync.Mutex
)

func GetCurrentPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("userID")
	mu.Lock()
	user, exists := users[userID]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user.CurrentPackage)
}

func UpdatePackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UpdatePackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := users[req.UserID]
	if !exists {
		mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if req.PackageID == HDStreamingPackage.ID {
		user.CurrentPackage = HDStreamingPackage
		users[req.UserID] = user
		mu.Unlock()
		w.Write([]byte("Package updated to HD Streaming."))
		return
	}

	mu.Unlock()
	http.Error(w, "Invalid package ID", http.StatusBadRequest)
}

type UpdatePackageRequest struct {
	UserID    string `json:"user_id"`
	PackageID string `json:"package_id"`
}

func main() {
	log.Println("Starting package service on :8091")
	http.HandleFunc("/get-package", GetCurrentPackage)
	http.HandleFunc("/update-package", UpdatePackage)

	if err := http.ListenAndServe(":8091", nil); err != nil {
		log.Fatalf("Failed to start package service: %v", err)
	}
}
