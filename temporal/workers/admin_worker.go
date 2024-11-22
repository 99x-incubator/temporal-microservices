package main

import (
	"log"

	activity "99x.io/admin_gateway/activity"
	workflows "99x.io/admin_gateway/workflows"

	"99x.io/shared/vars"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Set up Temporal client
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a worker for the task queue
	w := worker.New(c, vars.TaskQueue, worker.Options{})

	// Register the workflow and activity
	w.RegisterWorkflow(workflows.DisableRobotWorkflow)
	w.RegisterWorkflow(workflows.PackageUpgradeWorkflow)
	w.RegisterActivity(activity.DisableRobotActivity)
	w.RegisterActivity(activity.SendNotificationActivity)
	w.RegisterActivity(activity.GetPackageActivity)
	w.RegisterActivity(activity.UpdatePackageActivity)

	// Start listening for workflow tasks
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
