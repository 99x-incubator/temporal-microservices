package main

import (
	"context"
	"encoding/json"
	"net/http"

	"99x.io/admin_gateway/dto"
	"99x.io/admin_gateway/workflows"
	"99x.io/shared/vars"
	"go.temporal.io/sdk/client"
)

// REST endpoint to disable a robot
func DisableRobotHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.DisableRequest
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
