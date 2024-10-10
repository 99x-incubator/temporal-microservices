package main

import (
	"log"

	admin_workflows "99x.io/admin_gateway/workflows"

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
	w.RegisterWorkflow(admin_workflows.DisableRobotWorkflow)
	w.RegisterActivity(admin_workflows.DisableRobotActivity)
	w.RegisterActivity(admin_workflows.SendNotificationActivity)

	// Start listening for workflow tasks
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
