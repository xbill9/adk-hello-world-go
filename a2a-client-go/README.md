# a2a-client-go

This is a Go application built using the Agent Development Kit (ADK) that demonstrates an agent delegating tasks to sub-agents.

## Project Structure

The application consists of three agents:

- **`roll_agent`**: A local agent responsible for rolling dice.
- **`prime_agent`**: A remote agent that checks if numbers are prime.
- **`root_agent`**: The main agent that orchestrates requests, delegating dice rolling to `roll_agent` and primality checking to `prime_agent`. It can also handle requests that involve both, such as "Roll a 6-sided die and check if it's prime."

## Getting Started

### Prerequisites

- Go (version 1.22 or higher recommended)
- Docker (for building and running the remote prime agent, if not already deployed)
- `gcloud` CLI (for deployment to Google Cloud)

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/GoogleCloudPlatform/genkit.git
    cd genkit/go/a2a-client-go
    ```
2.  **Install dependencies:**
    ```bash
    make deps
    ```

### Running the Application

To run the `a2a-client-go` application:

```bash
make run
```

This will execute the `main.go` file, which includes an example interaction: "Roll a 6-sided die and check if it's prime."

### Building the Application

To build the executable:

```bash
make build
```

This will create an executable in the current directory.

### Running Tests

To run the tests for the project:

```bash
make test
```

### Code Quality

- **Format:**
    ```bash
    make format
    ```
- **Check Format:**
    ```bash
    make check-fmt
    ```
- **Lint:**
    ```bash
    make lint
    ```

### Deployment

This project includes targets for building a Docker image and deploying to Google Cloud using Cloud Build.

- **Build Docker Image:**
    ```bash
    make docker-build
    ```
- **Deploy to Google Cloud (via Cloud Build):**
    ```bash
    make deploy
    ```
    _Note: Ensure your `gcloud` CLI is configured with the correct project ID._

## Local Prime Agent (Optional)

The `prime_agent` is configured as a remote agent that expects to be available at `http://localhost:8086`. If you wish to run a local version of this agent for testing or development, you would typically build and run the corresponding `prime-agent-go` service, which is usually found in a sibling directory (e.g., `../a2a-prime-go`). Please refer to its specific `README` for instructions on how to set it up and run it locally.
