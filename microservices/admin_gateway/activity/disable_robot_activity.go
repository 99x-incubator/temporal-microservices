package activity

import (
	"context"
	"time"
)

// DisableRobotActivity simulates disabling the robot
func DisableRobotActivity(ctx context.Context, robotID string) (string, error) {
	// Simulate the operation of disabling the robot (e.g., database update)
	time.Sleep(2 * time.Second)
	return "Robot " + robotID + " is now disabled.", nil
}
