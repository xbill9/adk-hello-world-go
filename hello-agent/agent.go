package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const (
	defaultModelName = "gemini-2.5-flash"
	agentName        = "hello_time_agent"
)

func main() {
	// Handle signal interrupts (Ctrl+C) gracefully.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		// log.Fatal calls os.Exit(1), preventing defers from running.
		// Since we are at the end of main, it's acceptable here,
		// but using 'run' allows earlier defers (like cancel) to execute if we returned an error.
		log.Fatalf("Application failed: %v", err)
	}
}

func run(ctx context.Context) error {
	log.Println("Starting application...")

	apiKey := os.Getenv("GOOGLE_API_KEY")
	modelName := os.Getenv("MODEL_NAME")
	if modelName == "" {
		modelName = defaultModelName
	}

	log.Printf("Initializing model %q...", modelName)

	var model model.LLM
	var err error
	// use API KEY if set but otherwise Vertex AI
	if apiKey != "" {
		log.Println("Using Google API Key for authentication")
		model, err = gemini.NewModel(ctx, modelName, &genai.ClientConfig{
			APIKey: apiKey,
		})
	} else {
		log.Println("Using Vertex AI (default credentials) for authentication")
		model, err = gemini.NewModel(ctx, modelName, &genai.ClientConfig{})
	}
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	ag, err := llmagent.New(llmagent.Config{
		Name:        agentName,
		Model:       model,
		Description: "Tells the current time in a specified city.",
		Instruction: "You are a helpful assistant that tells the current time in a city.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	log.Printf("Agent %q created successfully", agentName)

	config := &launcher.Config{
		AgentLoader: &singleAgentLoader{agent: ag},
	}

	l := full.NewLauncher()
	log.Println("Starting launcher...")
	// Pass the signal-aware context to the launcher.
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		return fmt.Errorf("launcher execution failed: %w\n\nUsage:\n%s", err, l.CommandLineSyntax())
	}

	return nil
}

// singleAgentLoader adapts a single agent instance to the AgentLoader interface.
type singleAgentLoader struct {
	agent agent.Agent
}

func (s *singleAgentLoader) ListAgents() []string {
	return []string{s.agent.Name()}
}

func (s *singleAgentLoader) LoadAgent(name string) (agent.Agent, error) {
	if name == s.agent.Name() {
		return s.agent, nil
	}
	return nil, nil
}

func (s *singleAgentLoader) RootAgent() agent.Agent {
	return s.agent
}
