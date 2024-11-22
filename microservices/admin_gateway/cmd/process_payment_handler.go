package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"99x.io/admin_gateway/dto"
	"99x.io/admin_gateway/shared"
	"99x.io/shared/vars"
)

func PaymentConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentConfirmationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	paymentWorkflowID := vars.WorkflowPackageID + req.PackageID

	reqBody := `{"user_id": "` + req.UserID + `", "package_id": "` + req.PackageID + `"}`
	resp, err := http.Post(shared.PAYMENT_SERVICE, "application/json", ioutil.NopCloser(strings.NewReader(reqBody)))
	if err != nil {
		http.Error(w, "Failed to process payment"+err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()
	status := resp.Status

	if status != "200 OK" {
		http.Error(w, "Failed to process payment status"+status, http.StatusInternalServerError)
		return
	}

	sglErr := temporalClient.SignalWorkflow(context.Background(), paymentWorkflowID, "", "paymentConfirmation", true)
	if sglErr != nil {
		http.Error(w, "Failed to confirm payment "+sglErr.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with workflow execution details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Payment confirmed successfully",
		"PackageID": req.PackageID,
		"UserID":    req.UserID,
	})
}
