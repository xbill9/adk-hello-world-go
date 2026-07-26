package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/prod"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const (
	defaultModelName = "gemini-3-pro-preview"
	agentName        = "hello_time_agent"
)

func main() {
	// Handle signal interrupts (Ctrl+C) gracefully.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		// slog.Error does not exit, so we use os.Exit(1) after logging.
		// Since we are at the end of main, defers in run() have already executed.
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

// run executes the main application logic.
// It initializes the model, creates the agent, and starts the launcher.
func run(ctx context.Context) error {
	slog.Info("Starting application...")

	apiKey := os.Getenv("GOOGLE_API_KEY")
	modelName := os.Getenv("MODEL_NAME")
	if modelName == "" {
		modelName = defaultModelName
	}

	slog.Info("Initializing model", "model", modelName)

	var model model.LLM
	var err error
	// use API KEY if set but otherwise Vertex AI
	if apiKey != "" {
		slog.Info("Using Google API Key for authentication")
		model, err = gemini.NewModel(ctx, modelName, &genai.ClientConfig{
			APIKey: apiKey,
		})
	} else {
		slog.Info("Using Vertex AI (default credentials) for authentication")
		model, err = gemini.NewModel(ctx, modelName, &genai.ClientConfig{
			Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
			Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
			Backend:  genai.BackendVertexAI,
		})
	}
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	ag, err := llmagent.New(llmagent.Config{
		Name:        agentName,
		Model:       model,
		Description: "Tells the current time in a specified city by searching on Google.",
		Instruction: "You are a helpful assistant that tells the current time in a city using Google Search. You cannot set the time.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	slog.Info("Agent created successfully", "agent", agentName)

	config := &launcher.Config{
		AgentLoader:    &singleAgentLoader{agent: ag},
		SessionService: session.InMemoryService(),
	}

	// Allow PORT to be set by the environment
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8095"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Warn("Invalid PORT environment variable, defaulting to 8095", "error", err)
		port = 8095
	}

	// Check if the port is available, if not find the next available one
	finalPort := findAvailablePort(port)
	if finalPort != port {
		slog.Info("Port was busy, switching to available port", "old_port", port, "new_port", finalPort)
		portStr = strconv.Itoa(finalPort)
	}

	// Set PORT env var for the launcher to pick up
	os.Setenv("PORT", portStr)

	args := []string{
		"web", // Launch web server
		"--port", portStr,
		"a2a", // Sublauncher for a2a functionality
		"--a2a_agent_url", "http://0.0.0.0:" + portStr,
	}

	l := prod.NewLauncher()
	slog.Info("Starting launcher...")
	// Pass the signal-aware context to the launcher.
	if err := l.Execute(ctx, config, args); err != nil {
		return fmt.Errorf("launcher execution failed: %w\n\nUsage:\n%s", err, l.CommandLineSyntax())
	}

	return nil
}

// singleAgentLoader adapts a single agent instance to the AgentLoader interface.
// It is used to provide a simple loader that always returns the same agent.
type singleAgentLoader struct {
	agent agent.Agent
}

// ListAgents returns a list of available agent names.
func (s *singleAgentLoader) ListAgents() []string {
	return []string{s.agent.Name()}
}

// LoadAgent retrieves an agent by name.
// It returns the stored agent if the name matches, otherwise nil.
func (s *singleAgentLoader) LoadAgent(name string) (agent.Agent, error) {
	if name == s.agent.Name() {
		return s.agent, nil
	}
	return nil, nil
}

// RootAgent returns the default or root agent for this loader.
func (s *singleAgentLoader) RootAgent() agent.Agent {
	return s.agent
}

// findAvailablePort attempts to find an available port starting from the given port.
// It returns the first available port it finds, or the original port if checking fails.
func findAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port
		}
	}
	return startPort // Fallback to original if no port found in range
}
