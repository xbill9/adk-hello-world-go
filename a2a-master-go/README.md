# a2a-master-go

`a2a-master-go` is a Go application demonstrating a multi-agent system built with the Google Agent Development Kit (ADK). This project showcases how to combine local and remote agents to handle different tasks, specifically rolling dice and checking for prime numbers. The application runs as a web service, exposing its functionalities via a configurable port.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Building](#building)
  - [Running](#running)
  - [Testing](#testing)
- [Deployment](#deployment)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Multi-Agent Architecture**: Utilizes Google ADK to orchestrate multiple agents.
- **Dice Roll Agent**: A local agent capable of simulating dice rolls with a specified number of sides.
- **Prime Check Agent**: Integrates a remote agent to determine if a given number is prime.
- **Root Agent**: Delegates tasks to the appropriate sub-agents based on user requests.
- **Web Service**: Exposes agent functionalities as a web service, accessible via HTTP.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (version 1.22 or higher recommended)
- Docker (for building and running containerized versions)
- Google Cloud SDK with `gcloud` CLI (for deployment to Google Cloud)

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/a2aproject/a2a-master-go.git
    cd a2a-master-go
    ```
2.  **Initialize Go modules and tidy dependencies:**
    ```bash
    go mod tidy
    ```

### Building

To build the executable for development:

```bash
make build
```

To build a release-ready executable (Linux, no CGO):

```bash
make release
```

### Running

To run the application locally:

```bash
make run
```

The application will start on port `8092` by default. You can specify a different port by setting the `PORT` environment variable:

```bash
PORT=8080 make run
```

### Testing

To run the unit tests:

```bash
make test
```

## Deployment

The project includes targets for Docker and Google Cloud Build for deployment.

### Docker

To build a Docker image:

```bash
make docker-build
```

### Google Cloud Build

To deploy to Google Cloud using Cloud Build (requires `cloudbuild.yaml` to be configured):

```bash
make deploy
```

## Usage

Once the application is running, you can interact with the agent system through its exposed web interface or API endpoints. Specific interaction details (e.g., API paths, request formats) would depend on the ADK's default web UI or further API documentation.

## Contributing

Contributions are welcome! Please feel free to open issues or pull requests.

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.
