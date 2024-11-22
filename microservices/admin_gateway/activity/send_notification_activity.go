package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"99x.io/admin_gateway/dto"
)

// SendNotificationActivity sends a notification to the notification service
func SendNotificationActivity(ctx context.Context, notification dto.NotificationMessage) error {
	// Convert the notification message to JSON
	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Send the notification as an HTTP POST request
	resp, err := http.Post("http://localhost:8082/notify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification service returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
