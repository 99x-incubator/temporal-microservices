package main

import (
	"fmt"
	"os"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"99x.io/workflows"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.temporal.io/sdk/client"
	"github.com/hellofresh/health-go/v5"
)

type DisableRequest struct {
	RobotID string `json:"robot_id"`
	UserID  string `json:"user_id"`
}

var temporalClient client.Client   

// REST endpoint to disable a robot
func disableRobotHandler(w http.ResponseWriter, r *http.Request) {
	var req DisableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RobotID == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "disable_robot_workflow_" + req.RobotID,
		TaskQueue: "ADMIN_TASK_QUEUE",
	}

	// Start the workflow for disabling the robot
	we, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.DisableRobotWorkflow, req.RobotID, req.UserID)
	if err != nil {
		http.Error(w, "Failed to initiate workflow", http.StatusInternalServerError)
		return
	}

	// Respond with workflow execution details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Disable robot workflow initiated",
		"workflowID":    we.GetID(),
		"workflowRunID": we.GetRunID(),
		"robotID":       req.RobotID,
	})
}

func main() {

	var temporal_server = "localhost:7233"
	var client_cors = "http://localhost:8083"
	
	// read from env. location of the temporal server , client CORS
	val, ok := os.LookupEnv("TEMPORAL_SERVER")
	if ok {
		temporal_server = val
	}
	fmt.Printf("%s=%s\n", "TEMPORAL_SERVER", temporal_server)
	val, ok = os.LookupEnv("CLIENT_CORS")
	if ok {
		client_cors = val
	}
	fmt.Printf("%s=%s\n", "CLIENT_CORS", client_cors)
	

    //health
	//probe need to add checks which looks at visibility of the temporal 
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "admin_gateway",
		Version: "v1.0",
	}))


	// Set up Temporal client
	var err error
	temporalClient, err = client.NewLazyClient(client.Options{
			HostPort: temporal_server,
		})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/disable_robot", disableRobotHandler).Methods("POST")

	r.HandleFunc("/status", h.HandlerFunc).Methods("GET")

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{client_cors}, // Replace with your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Start the HTTP server with CORS middleware
	handler := c.Handler(r)
	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
