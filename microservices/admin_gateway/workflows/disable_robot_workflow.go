package workflows

import (
	"fmt"
	"time"

	"99x.io/admin_gateway/activity"
	"99x.io/admin_gateway/dto"
	"go.temporal.io/sdk/workflow"
)

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
	err := workflow.ExecuteActivity(ctx, activity.DisableRobotActivity, robotID).Get(ctx, &result)
	if err != nil {
		logger.Error("Failed to disable robot", "error", err)
		return err
	}

	logger.Info("Successfully disabled robot", "result", result)

	// Prepare notification message
	message := fmt.Sprintf("Robot %s has been disabled.", robotID)
	notification := dto.NotificationMessage{
		UserID:  userID,
		Message: message,
	}

	// Activity options for sending notification
	aoNotify := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctxNotify := workflow.WithActivityOptions(ctx, aoNotify)

	// Send notification activity
	err = workflow.ExecuteActivity(ctxNotify, activity.SendNotificationActivity, notification).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send notification", "error", err)
		return err
	}

	logger.Info("Notification sent successfully", "userID", userID)
	return nil
}
