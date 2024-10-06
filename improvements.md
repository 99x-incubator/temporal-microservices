# Design Improvements Proposed

## 1.Client Layer (Vue.js Frontend)

- **Synchronous API Calls**: We need to implement a synchronous API call example to the Gateway Service to demonstrate real-time interactions and immediate feedback to the user. For example showing the list of Robots.
- **Error Handling**: We need to implement error handling in the frontend to display error messages to the user in case of a failed API call. For example showing an error message if one of the dependent microservices is down, and workflow fails to complete (after retries).
- **Loading Indicators**: We need to implement loading indicators to show the user that the system is processing their request. For example showing a spinner when the system is processing a long-running workflow.

## 2.Gateway Layer (Custom API Gateway)

- **Command Pattern**: We need to implement the Command pattern in the Gateway to decouple the request processing/validation logic from the API endpoint and make it easier to add new features. For example, creating a Command class for each type of request (e.g., GetRobotListCommand, CreateRobotCommand).

## 3.Temporal.io Workflows

- **Abstract Activity Logic**: We need to abstract the activity logic into separate objects/methods to improve code readability and maintainability. For example, creating a separate class for each activity (e.g., OrderActivity, InventoryActivity) and injecting them into the workflow.
- **Long-Running Processes**: We need to demonstrate a long-running processes in Temporal workflows to handle complex business logic that spans multiple microservices. For example, implementing an order processing workflow that involves inventory management, payment processing, and notification delivery.
- **Retry Policies**: We need to demonstrate different retry policies in Temporal workflows to automatically retry failed activities based on predefined rules. For example, retrying an activity with an exponential backoff strategy if it fails due to a temporary network issue.
- **Timeouts**: We need to demonstrate timeouts in Temporal workflows to prevent long-running processes from blocking resources indefinitely. For example, setting a timeout for each activity to ensure that it completes within a certain time frame.
- **Saga Pattern & Compensation**: We need to implement the Saga pattern and compensation logic in Temporal workflows to handle processes that involve multiple microservices. For example, implementing a compensation workflow to rollback changes if an activity of the workflow fails.
- **Testing**: We need to implement automated testing of workflows to ensure the correctness and reliability of the system. For example, writing unit tests for Activities and integration tests for Workflows as a whole.

## 4.(Priority 2) Advanced System Architecture

- **Rate Limiting**: We need to implement rate limiting in the Gateway to prevent abuse and ensure fair usage of the system. For example, limiting the number of requests a user can make in a given time period.
- **Circuit Breaker**: We need to implement a circuit breaker pattern in the Gateway to handle failures gracefully and prevent cascading failures. For example, if a downstream service is unavailable, the Gateway should stop sending requests to it and return an error to the client.
- **Logging and Monitoring**: We need to implement logging and monitoring in the Gateway to track the performance and health of the system. For example, logging the request and response data, and monitoring the latency and error rates.
- **Response Caching**: We need to implement response caching in the Gateway to improve performance and reduce the load on downstream services. For example, caching the response of a frequently requested endpoint for a certain period of time.
- **API Versioning**: We need to implement API versioning in the Gateway to support backward compatibility and allow for changes to the API without breaking existing clients. For example, adding a version number to the API endpoint (e.g., /v1/robots).
- **Security**: We need to implement security measures in the Gateway to protect against common web vulnerabilities (e.g., SQL injection, cross-site scripting). For example, validating and sanitizing user input, and using HTTPS for secure communication.
- **Authentication and Authorization**: We need to implement authentication and authorization in the Gateway to control access to the system and protect sensitive data. For example, using JSON Web Tokens (JWT) for authentication and attribute-based access control (ABAC) for authorization.
- **Service Discovery**: We need to implement service discovery in the Gateway to distribute traffic  across multiple instances of a service. For example, using a service registry like Consul or Eureka to register and discover services.
- **API Documentation**: We need to implement API documentation in the Gateway to provide a reference for developers on how to use the API. For example, using Swagger or OpenAPI to generate API documentation automatically from the code.
- **Middleware**: We need to implement middleware in the Gateway to add cross-cutting concerns (e.g., logging, rate limiting, authentication) to the request processing pipeline. For example, adding a middleware to log the request and response data for each API call.
- **Workflow Versioning**: We need to implement workflow versioning in Temporal to support backward compatibility and allow for changes to the workflow without breaking existing instances. For example, adding a version number to the workflow definition and handling different versions of the workflow in the worker.
