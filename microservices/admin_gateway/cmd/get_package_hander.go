package main

import (
	"io/ioutil"
	"net/http"

	"99x.io/admin_gateway/shared"
)

func GetCurrentPackageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")

	// Make the request to the external service
	resp, err := http.Get(shared.GET_PACKAGE_SERVICE + "?userID=" + userID)

	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Read the body from the external service's response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Write the response body directly back to the client
	w.Write(body)
}
