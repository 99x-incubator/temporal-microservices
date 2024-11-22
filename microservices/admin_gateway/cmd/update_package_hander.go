package main

import (
	"encoding/json"
	"net/http"

	"99x.io/admin_gateway/dto"

	"context"

	"99x.io/admin_gateway/workflows"
	"99x.io/shared/vars"
	"go.temporal.io/sdk/client"
)

func UpdateCurrentPackageHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	workflowOptions := client.StartWorkflowOptions{
		ID:        vars.WorkflowPackageID + req.PackageID,
		TaskQueue: vars.TaskQueue,
	}

	// Start the workflow for disabling the robot
	we, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.PackageUpgradeWorkflow, req.PackageID, req.UserID)
	if err != nil {
		http.Error(w, "Failed to initiate workflow", http.StatusInternalServerError)
		return
	}

	// Respond with workflow execution details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Upgrade robot package workflow initiated",
		"workflowID":    we.GetID(),
		"workflowRunID": we.GetRunID(),
		"PackageID":     req.PackageID,
		"UserID":        req.UserID,
	})

}
