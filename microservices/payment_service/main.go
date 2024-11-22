package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func ProcessPayment(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := req.UserID

	if userID == "123" {
		time.Sleep(5 * time.Second) // Simulate processing time
		w.Write([]byte("Success"))
		return
	}

	http.Error(w, "Payment Failed", http.StatusBadRequest)
}

func main() {
	log.Println("Starting payment service on :8092")
	http.HandleFunc("/process-payment", ProcessPayment)
	http.ListenAndServe(":8092", nil)

	if err := http.ListenAndServe(":8092", nil); err != nil {
		log.Fatalf("Failed to start payment service: %v", err)
	}
}

type UpdatePaymentRequest struct {
	UserID    string `json:"user_id"`
	PackageID string `json:"package_id"`
}
