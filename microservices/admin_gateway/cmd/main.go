package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.temporal.io/sdk/client"
)

var temporalClient client.Client

func main() {
	// Set up Temporal client
	var err error
	temporalClient, err = client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/disable_robot", DisableRobotHandler).Methods("POST")
	r.HandleFunc("/get_package", GetCurrentPackageHandler).Methods("GET")
	r.HandleFunc("/update_package", UpdateCurrentPackageHandler).Methods("POST")
	r.HandleFunc("/process_payment", PaymentConfirmationHandler).Methods("POST")

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8083"}, // Replace with your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Start the HTTP server with CORS middleware
	handler := c.Handler(jsonMiddleware(r))
	log.Println("Starting gateway on :8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Failed to gateway server: %v", err)
	}
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Type for the response
		w.Header().Set("Content-Type", "application/json")
		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
