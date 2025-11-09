package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

// Launcher is an interface for executing the ADK launcher.
type Launcher interface {
	Execute(ctx context.Context, config *adk.Config, args []string) error
	CommandLineSyntax() string
}

func main() {
	if err := run(full.NewLauncher(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(l Launcher, args []string) error {
	ctx := context.Background()

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	agentName := os.Getenv("ADK_AGENT_NAME")
	if agentName == "" {
		agentName = "hello_time_agent"
	}

	model, err := gemini.NewModel(ctx, modelName, &genai.ClientConfig{})
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	agent, err := llmagent.New(llmagent.Config{
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

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}

	if err := l.Execute(ctx, config, args); err != nil {
		return fmt.Errorf("run failed: %w\n\n%s", err, l.CommandLineSyntax())
	}
	return nil
}
