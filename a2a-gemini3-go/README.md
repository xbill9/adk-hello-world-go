# a2a-gemini3-go

`a2a-gemini3-go` is a sample Go application demonstrating how to build an AI agent using the [Agent Development Kit (ADK)](https://pkg.go.dev/google.golang.org/adk).

This project creates a simple agent named `hello_time_agent` powered by Google's Gemini models. The agent is designed to tell the current time in a specified city, utilizing Google Search via the ADK's `geminitool` to fetch accurate information.

## Features

*   **ADK Integration:** Showcases usage of `google.golang.org/adk` for agent creation and management.
*   **Gemini Powered:** Uses the `gemini-3-pro-preview` model by default (configurable).
*   **Tool Use:** Equip the agent with `geminitool.GoogleSearch` to perform real-world queries.
*   **Dual Authentication:** Supports both Google API Key and Google Cloud Vertex AI (ADC) authentication methods.
*   **Standard Launcher:** Utilizes the ADK's full launcher for robust CLI execution.

## Prerequisites

*   **Go:** Version 1.24.4 or higher.
*   **Google Cloud Project:** Required if using Vertex AI or for deployment.
*   **Google API Key:** Required if running locally without Google Cloud default credentials.

## Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd adk-hello-world-go/a2a-gemini3-go
    ```

2.  **Install dependencies:**
    ```bash
    make deps
    ```

## Configuration

The application is configured via environment variables:

| Variable         | Description                                                                 | Default                |
| ---------------- | --------------------------------------------------------------------------- | ---------------------- |
| `GOOGLE_API_KEY` | Your Google API Key for authentication. If not set, Vertex AI credentials are used. | `""`                   |
| `MODEL_NAME`     | The name of the Gemini model to use.                                        | `gemini-3-pro-preview` |

## Usage

This project includes a `Makefile` to simplify common tasks.

### Building locally

To build the project binary:

```bash
make build
```

### Running the Agent

To run the agent locally (ensure you have set `GOOGLE_API_KEY` if needed):

```bash
export GOOGLE_API_KEY="your-api-key"
make run -- --help
```

Because this uses the ADK launcher, you can interact with the agent via the CLI. For example:

```bash
make run -- chat "What time is it in Tokyo?"
```

### Running Tests

To execute the unit tests:

```bash
make test
```

### Other Make Commands

*   `make format`: Format the code using `gofmt`.
*   `make lint`: Run `go vet` to identify potential issues.
*   `make clean`: Clean up build artifacts.
*   `make docker-build`: Build a Docker image for the service.

## Deployment

The project includes a target for submitting a build to Google Cloud Build.

1.  Ensure you have the `gcloud` CLI installed and configured with your project.
2.  Run the deploy command:
    ```bash
    make deploy
    ```

This will submit the build using `cloudbuild.yaml`.

## Project Structure

*   `main.go`: The application entry point. It configures the Gemini model, initializes the `hello_time_agent` with the Google Search tool, and starts the ADK launcher.
*   `go.mod` / `go.sum`: Go module definition and checksums.
*   `Makefile`: Automation for build, test, and deployment tasks.
