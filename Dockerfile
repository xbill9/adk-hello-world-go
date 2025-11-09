# Build stage
FROM golang:1.24.4 AS build

WORKDIR /app

# Copy go.mod and go.sum files
COPY hello-agent/go.mod hello-agent/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY hello-agent/ ./


# Create vendor directory
RUN go mod vendor

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /hello-agent

# Deploy stage
FROM gcr.io/distroless/static-debian11

# Set the working directory for the deploy stage
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /hello-agent /app/hello-agent

# Expose the port the application will listen on
EXPOSE 8080

# Run the application
CMD ["/app/hello-agent", "web", "-port", "8080", "api", "-webui_address", "127.0.0.1:8081", "a2a", "--a2a_agent_url", "http://127.0.0.1:8081", "webui", "--api_server_address", "http://127.0.0.1:8081/api"]
