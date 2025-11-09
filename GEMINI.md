# Gemini Workspace for `adk-hello-world-go`

You are a GO Developer working with Google Cloud.
This document outlines best practices for developing Go applications, especially when using the Go ADK and deploying to Google Cloud. Adhering to these guidelines will ensure your code is maintainable, performant, and secure.

## 1. General Go Best Practices

### 1.1 Error Handling
*   **Explicit Error Checks:** Always check errors explicitly. Don't ignore them.
*   **Error Wrapping:** Use `fmt.Errorf` with `%w` to wrap errors, preserving the original error context. This allows for programmatic inspection of error chains.
*   **Custom Error Types:** Define custom error types for specific error conditions to provide more context and allow for type-based error handling.
*   **Contextual Errors:** Add sufficient context to error messages to aid debugging.

### 1.2 Concurrency (Goroutines and Channels)
*   **Use Goroutines for Parallelism:** Leverage goroutines for tasks that can run concurrently.
*   **Communicate with Channels:** Use channels to safely communicate and synchronize data between goroutines, avoiding shared memory issues.
*   **Avoid Naked Goroutines:** Ensure goroutines are properly managed (e.g., using `sync.WaitGroup` or context cancellation) to prevent leaks or unexpected behavior.
*   **Context for Cancellation/Timeouts:** Use `context.Context` to manage cancellation signals and timeouts across goroutine hierarchies.

### 1.3 Code Structure and Modularity
*   **Clear Package Structure:** Organize code into logical packages. Each package should have a single, clear responsibility.
*   **Small Functions:** Keep functions small and focused on a single task.
*   **Modules for Dependency Management:** Use Go Modules for managing dependencies.
*   **Internal Packages:** Use `internal` packages for code that should not be imported by other modules.

### 1.4 Naming Conventions
*   **Descriptive Names:** Use clear, descriptive names for variables, functions, and types.
*   **CamelCase for Exported Names:** Exported (public) names start with a capital letter. Unexported (private) names start with a lowercase letter.
*   **Short Variable Names for Local Scope:** Use shorter, concise names for variables with small scopes (e.g., `i` for loop counters, `err` for errors).

### 1.5 Testing
*   **Write Unit Tests:** Create unit tests for all critical logic.
*   **Table-Driven Tests:** Use table-driven tests for multiple test cases with similar logic.
*   **Benchmarking:** Write benchmarks for performance-critical code.
*   **Test Coverage:** Aim for good test coverage, but prioritize meaningful tests over 100% coverage.

### 1.6 Documentation
*   **Godoc:** Document all exported functions, types, and variables using Godoc comments.
*   **README.md:** Maintain a comprehensive `README.md` for project overview, setup, and usage.

## 2. Go ADK Specific Best Practices

When working with the Go ADK, consider the following:

### 2.1 ADK Integration
*   **Utilize ADK Libraries:** Leverage the provided ADK libraries and utilities for common tasks (e.g., authentication, configuration, service interaction).
*   **Follow ADK Patterns:** Adhere to any specific patterns or interfaces defined by the ADK to ensure seamless integration and future compatibility.
*   **Configuration Management:** Use ADK's recommended methods for managing application configuration, especially for environment-specific settings.

### 2.2 Google Cloud Deployment Considerations

### 2.2.1 Logging and Monitoring
*   **Structured Logging:** Use structured logging (e.g., JSON format) with libraries like `zap` or `logrus` to make logs easily parsable by Google Cloud Logging.
*   **Contextual Logging:** Include relevant context (e.g., request IDs, user IDs) in log entries.
*   **Cloud Monitoring:** Integrate with Google Cloud Monitoring (formerly Stackdriver) for metrics and alerts. Use OpenCensus or OpenTelemetry for tracing and custom metrics.

### 2.2.2 Security
*   **Least Privilege:** Grant only the necessary IAM permissions to your service accounts.
*   **Secret Management:** Use Google Secret Manager for sensitive information (API keys, database credentials) instead of hardcoding them or storing them in environment variables directly.
*   **Vulnerability Scanning:** Regularly scan your dependencies for known vulnerabilities.

### 2.2.3 Cost Optimization
*   **Resource Sizing:** Right-size your Cloud Run instances or GKE pods to avoid over-provisioning.
*   **Autoscaling:** Configure autoscaling where appropriate to scale resources up and down based on demand.
*   **Serverless First:** Prefer serverless options like Cloud Run or Cloud Functions for event-driven or stateless workloads to minimize operational overhead and cost.

### 2.2.4 Deployment
*   **Containerization:** Use Docker to containerize your Go applications for consistent deployment across environments.
*   **Cloud Build:** Automate your build and deployment pipelines using Google Cloud Build.
*   **Health Checks:** Implement health checks (`/healthz`, `/readyz`) for services deployed on platforms like Cloud Run or GKE.

By following these best practices, you can develop robust, efficient, and scalable Go applications within the Google Cloud ecosystem using the Go ADK.