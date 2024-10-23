package main

import (
	"log"
	"fmt"
	"os"
	admin_workflows "99x.io/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	var temporal_server = "localhost:7233"
	// read from env. location of the temporal server , client CORS
	val, ok := os.LookupEnv("TEMPORAL_SERVER")
	if ok {
		temporal_server = val
	}
	fmt.Printf("%s=%s\n", "TEMPORAL_SERVER", temporal_server)

	// Set up Temporal client
	c, err := client.NewLazyClient(client.Options{
		HostPort: temporal_server,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a worker for the task queue
	w := worker.New(c, "ADMIN_TASK_QUEUE", worker.Options{})

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
