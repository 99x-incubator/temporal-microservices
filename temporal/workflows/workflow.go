package workflows

import (
	"os"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.temporal.io/sdk/workflow"
)

// NotificationMessage represents the structure of the notification to be sent
type NotificationMessage struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// DisableRobotWorkflow orchestrates disabling the robot and sending a notification
func DisableRobotWorkflow(ctx workflow.Context, robotID string, userID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Disabling robot", "robotID", robotID)

	// Activity options for disabling the robot
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Disable the robot (simulated activity)
	var result string
	err := workflow.ExecuteActivity(ctx, DisableRobotActivity, robotID).Get(ctx, &result)
	if err != nil {
		logger.Error("Failed to disable robot", "error", err)
		return err
	}

	logger.Info("Successfully disabled robot", "result", result)

	// Prepare notification message
	message := fmt.Sprintf("Robot %s has been disabled.", robotID)
	notification := NotificationMessage{
		UserID:  userID,
		Message: message,
	}

	// Activity options for sending notification
	aoNotify := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctxNotify := workflow.WithActivityOptions(ctx, aoNotify)

	// Send notification activity
	err = workflow.ExecuteActivity(ctxNotify, SendNotificationActivity, notification).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send notification", "error", err)
		return err
	}

	logger.Info("Notification sent successfully", "userID", userID)
	return nil
}

// DisableRobotActivity simulates disabling the robot
func DisableRobotActivity(ctx context.Context, robotID string) (string, error) {
	// Simulate the operation of disabling the robot (e.g., database update)
	time.Sleep(2 * time.Second)
	return "Robot " + robotID + " is now disabled.", nil
}

// SendNotificationActivity sends a notification to the notification service
func SendNotificationActivity(ctx context.Context, notification NotificationMessage) error {
	// Convert the notification message to JSON
	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Send the notification as an HTTP POST request

	// where to send the notification 
	var notify_server = "http://localhost:8082/notify"
	
	// read from env. location of the temporal server , client CORS
	val, ok := os.LookupEnv("NOTIFY_SERVER")
	if ok {
		notify_server = val
	}
	fmt.Printf("%s=%s\n", "NOTIFY_SERVER", notify_server)

	resp, err := http.Post(notify_server, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification service returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
