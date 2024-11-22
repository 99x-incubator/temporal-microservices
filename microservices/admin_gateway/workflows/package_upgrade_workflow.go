package workflows

import (
	"fmt"
	"time"

	"99x.io/admin_gateway/activity"
	"99x.io/admin_gateway/dto"
	"go.temporal.io/sdk/workflow"
)

func PackageUpgradeWorkflow(ctx workflow.Context, packageID string, UserID string) error {

	logger := workflow.GetLogger(ctx)
	logger.Info("Upgrading package", "packageID", packageID)

	// Activity options for disabling the robot
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Get current package
	var currentPackage string
	pkgGetErr := workflow.ExecuteActivity(ctx, activity.GetPackageActivity, UserID).Get(ctx, &currentPackage)
	if pkgGetErr != nil {
		logger.Error("Failed to activate package", "error", pkgGetErr)
		return pkgGetErr
	}
	workflow.GetLogger(ctx).Info("Current package:", "package", currentPackage)

	//Step 2: Temporarily upgrade the package
	var upgradedPackageId string
	err := workflow.ExecuteActivity(ctx, activity.UpdatePackageActivity, packageID, UserID).Get(ctx, &upgradedPackageId)
	if err != nil {
		logger.Error("Failed to activate package", "error", err)
		return err
	}
	workflow.GetLogger(ctx).Info("Package temporarily upgraded to packageID:", "packageID", upgradedPackageId)

	// Step 3: Wait for external payment confirmation via API (with a 20-minute timeout)
	signalChan := workflow.GetSignalChannel(ctx, "paymentConfirmation")
	timerCtx, cancel := workflow.WithCancel(ctx)
	defer cancel()

	timerFuture := workflow.NewTimer(timerCtx, 30*time.Minute)
	var paymentConfirmed bool

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, _ bool) {
		c.Receive(ctx, &paymentConfirmed)
	})
	selector.AddFuture(timerFuture, func(f workflow.Future) {
		workflow.GetLogger(ctx).Info("Payment timeout")
		paymentConfirmed = false
	})

	selector.Select(ctx)
	var message string
	if !paymentConfirmed {

		err := workflow.ExecuteActivity(ctx, activity.UpdatePackageActivity, currentPackage, UserID).Get(ctx, &upgradedPackageId)

		if err != nil {
			logger.Error("Failed to revert back package", "error", err)
			return err
		}
		message = fmt.Sprintf("Payment not confirmed, package %s upgrade not persisted for user %s", packageID, UserID)
		workflow.GetLogger(ctx).Info("Package downgraded to packageID:", "packageID", upgradedPackageId)
	} else {
		message = fmt.Sprintf("Payment confirmed, package %s upgrade persisted for user %s", packageID, UserID)
	}

	workflow.GetLogger(ctx).Info("Payment confirmed, package upgrade persisted")

	// Prepare notification message

	notification := dto.NotificationMessage{
		UserID:  UserID,
		Message: message,
	}

	err = workflow.ExecuteActivity(ctx, activity.SendNotificationActivity, notification).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send notification", "error", err)
		return err
	}

	return nil
}
