package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"99x.io/admin_gateway/workflows"
	"99x.io/shared/vars"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.temporal.io/sdk/client"
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
		ID:        vars.WorkflowID + req.RobotID,
		TaskQueue: vars.TaskQueue,
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
	// Set up Temporal client
	var err error
	temporalClient, err = client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/disable_robot", disableRobotHandler).Methods("POST")

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8083"}, // Replace with your frontend URL
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
