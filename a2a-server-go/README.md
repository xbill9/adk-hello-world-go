# A2A Prime Checker Server

This project implements an Agent-to-Agent (A2A) server in Go that hosts an agent capable of checking if a list of numbers are prime. It leverages the Google ADK (Agent Development Kit) for building and deploying the agent.

## Features

- **Prime Number Checking**: An integrated agent (`check_prime_agent`) that can determine if a given list of integers contains prime numbers.
- **Google ADK Integration**: Built using the Google ADK, demonstrating agent capabilities and interaction patterns.
- **LLM Agent**: Utilizes a Large Language Model (LLM) agent (`gemini-2.5-flash`) for natural language understanding and tool invocation.

## Prerequisites

Before running this project, ensure you have the following installed:

- **Go**: Version 1.24.4 or later.
- **Google Cloud SDK**: For authentication and access to Google Cloud services, including the Generative AI API. Ensure you are authenticated (`gcloud auth application-default login`).
- **Google ADK Client**: To interact with this A2A server, you will need an ADK client or another compatible A2A client.

## Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-repo/a2a-server-go.git
    cd a2a-server-go
    ```
2.  **Download dependencies**:
    ```bash
    go mod tidy
    ```

## Usage

### Running the Server

To start the A2A Prime Checker Server, execute the `main.go` file:

```bash
go run main.go
```

The server will start on port `8086` by default. You can change this by setting the `PORT` environment variable:

```bash
PORT=8000 go run main.go
```

You should see output similar to:

```
2025/11/23 10:00:00 Starting A2A prime checker server on port 8086
```

### Interacting with the Agent

Once the server is running, you can interact with the `check_prime_agent` using an ADK-compatible client. The agent exposes a `prime_checking` tool that takes a list of integers and returns which of them are prime.

**Example (Conceptual Client Interaction):**

```
User: Are 7, 10, 13, 22 prime numbers?
Agent: 7, 13 are prime numbers.
```

## Project Structure

-   `main.go`: Contains the main application logic, including the definition of the `isPrime` function, the `checkPrimeTool`, and the A2A server setup.
-   `go.mod`: Go module file listing project dependencies.
-   `go.sum`: Go module checksums.
-   `main_test.go`: Unit tests for the `isPrime` function and potentially other components.
-   `Makefile`: Contains common commands for building, testing, or deploying the application (if present).

## License

This project is licensed under the Apache 2.0 License. See the header in `main.go` for details.
