# Running the reference application

This reference application demonstrates how to build a real-time, event-driven, microservices architecture using Temporal, Go, and Vue.js.

See the [Architecture](architecture.md) for technical information. Below are the steps to run the reference application on your local machine.

## Running Temporal Cluster

To run the Temporal cluster, navigate to the `temporal/cluster` directory and run the following commands:

[This folder contains a clone from the [temporalio/docker-compose](https://github.com/temporalio/docker-compose) GitHub repository]

```bash
docker-compose up
```

The Temporal Web UI will be available at `http://localhost:8233/namespaces/default/workflows`.

### Temporal CLI

Set up an alias for the Temporal CLI to interact with the cluster:

```bash
alias temporal_docker="docker exec temporal-admin-tools temporal"
```

Now you can use `temporal_docker` command to start workflows, list workflow executions, and more.

## Running Admin-Gateway Microservice

To run the Admin-Gateway microservice, navigate to the `microservices/admin_gateway` directory and run the following commands:

```bash
go run ./cmd
```

### Test the Admin-Gateway Microservice

To test the Admin-Gateway microservice, you can use the following curl command:

```bash
curl -X POST http://localhost:8081/disable_robot \
     -H "Content-Type: application/json" \
     -d '{"robot_id": "test_robot_101", "user_id": "user123"}'
```

This should initiate a workflow in the Temporal cluster. Visit the Temporal Web UI to see the workflow execution.

This workflow will not complete because there are no workers running to process the workflow tasks. To complete the workflow, you need to start a Temporal Worker.

## Running the Notification Microservice

To run the Notification microservice, navigate to the `microservices/admin_notifications` directory and run the following commands:

```bash
go run .
```

The notification microservice is a Go-based service that sends real-time, targeted notifications to specific users via WebSockets. It listens for HTTP requests containing a userID and a message, then forwards the message to the corresponding WebSocket connection for the user, enabling real-time updates to the connected client.

### Test the Notification Microservice

To test the Notification microservice, you can use the following curl command:

```bash
curl -X POST http://localhost:8082/notify \
-H "Content-Type: application/json" \
-d '{"user_id": "user123", "message": "Robot Alpha has been disabled."}'
```

You should see the `User not connected` message appear in the client application. This is because we haven't started the web client application to connect to the Notification microservice yet.

## Starting a Temporal Worker

To start a Temporal Worker, navigate to the `temporal/workers` directory and run the following commands:

```bash
go run .
```

This will start a worker that listens for tasks from the Temporal cluster and processes them. The worker will process the workflow task we created earlier and complete the workflow. You can see the workflow completion in the Temporal Web UI.

## Running Web Client Application

To run the client application, navigate to the `client` directory and run the following commands:

```bash
cd client/av1-admin 
npm install
npm run serve
```

The client application will be available at `http://localhost:8080/`.

### Test the Web Client Application

To test the client application, open the browser and navigate to `http://localhost:8083/`. You should see the client application with a WebSocket connection to the Notification microservice. You can test the full round trip of the architecture by clicking the `Disable` button in the client application.

- This will send a request to the Admin-Gateway microservice.
- Which will initiate a workflow in the Temporal cluster.
- The Temporal worker will process the workflow task, and inform the notification microservice.
- The Notification microservice will send a real-time notification to the client application via WebSocket.
