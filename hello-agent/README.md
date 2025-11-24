# hello-agent

`hello-agent` is a Go application built using the Google Agent Development Kit (ADK) that demonstrates the creation of an LLM-powered agent. This agent, named `hello_time_agent`, is designed to tell the current time in a specified city by leveraging a large language model (defaulting to Gemini 2.5 Flash) and integrating with the Google Search tool.

## Features

*   **LLM Integration:** Utilizes the Gemini 2.5 Flash model for natural language understanding and generation.
*   **Google ADK:** Built upon the Google Agent Development Kit for agent orchestration and management.
*   **Tool Usage:** Integrates with `geminitool.GoogleSearch` to retrieve real-time information.
*   **Contextual Responses:** Provides current time information for specified cities.
*   **Flexible Authentication:** Supports authentication via `GOOGLE_API_KEY` or Google Cloud's Vertex AI (default credentials).

## Prerequisites

Before running this application, ensure you have the following:

*   **Go:** Version 1.24.4 or later.
*   **Google Cloud Project (Optional):** If using Vertex AI for authentication, ensure your Google Cloud project is set up with the necessary permissions and that `gcloud` is authenticated.
*   **Google API Key (Optional):** Alternatively, you can provide a `GOOGLE_API_KEY` with access to the Gemini API.

## Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-repo/hello-agent.git
    cd hello-agent
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Authentication:**
    Choose one of the following methods:
    *   **Using a Google API Key:**
        ```bash
        export GOOGLE_API_KEY="YOUR_API_KEY"
        ```
    *   **Using Vertex AI (default credentials):**
        Ensure your `gcloud` CLI is authenticated to your Google Cloud project. No `GOOGLE_API_KEY` export is needed in this case.

## Running the Agent

To run the `hello-agent` application:

```bash
go run .
```

You can also specify a different model name:

```bash
MODEL_NAME="gemini-pro" go run .
```

After launching, the agent will start, and you can interact with it through the ADK launcher.

## Project Structure

*   `agent.go`: The main application file containing the agent's definition and logic.
*   `go.mod`: Go module definition and dependencies.
*   `go.sum`: Checksum for module dependencies.
*   `Makefile`: (Optional) Contains build or run commands for convenience.
*   `GEMINI.md`: Internal guidelines for Go development within this project.

## Go Project Guidelines

The following guidelines are adhered to in this project to ensure clean, performant, and idiomatic Go code:

### 1. Code Style and Formatting

- All Go code **must** be formatted with `gofmt`.
- Follow standard Go naming conventions:
    - **Packages:** Use short, concise, all-lowercase names.
    - **Variables, Functions, and Methods:** Use `camelCase` for unexported identifiers and `PascalCase` for exported identifiers.
    - **Interfaces:** Name interfaces based on what they do (e.g., `io.Reader`), avoiding prefixes like `I`.

### 2. Error Handling

- Errors are values and should be handled explicitly using `if err != nil`.
- Provide context to errors using `fmt.Errorf` or a dedicated error handling package for richer error information.
- Do not discard errors silently.

### 3. Concurrency

- Use goroutines and channels for concurrent operations as appropriate.
- Ensure proper synchronization to prevent race conditions (e.g., using `sync.Mutex` or `sync.WaitGroup`).
- Avoid global mutable state where possible.

### 4. Testing

- Write comprehensive unit tests for all significant functions and packages.
- Use Go's built-in `testing` package.
- Ensure test coverage is high and tests are maintainable.

### 5. Project Context

- This project is a backend service for a distributed system.
- Performance and reliability are critical.
- We utilize `PostgreSQL` as our primary database.

### 6. Agent Interaction Protocol

- When suggesting code changes, provide clear explanations for the reasoning behind the changes.
- If asked to refactor, prioritize readability and maintainability while considering performance implications.
- When reviewing code, highlight potential issues related to the above guidelines and suggest improvements.
