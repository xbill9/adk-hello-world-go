// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strconv"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/agent/remoteagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/prod"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"
)

const (
	RollAgentName  = "roll_agent"
	PrimeAgentName = "prime_agent"
	RootAgentName  = "root_agent"
	GeminiModel    = "gemini-2.5-flash"
)

// --- Local Roll Agent ---

type rollDieToolArgs struct {
	Sides int `json:"sides" jsonschema:"The number of sides on the die."`
}

func rollDieTool(ctx agent.Context, args rollDieToolArgs) (string, error) {
	if args.Sides <= 0 {
		return "", fmt.Errorf("number of sides must be greater than 0, got %d", args.Sides)
	}
	result := rand.Intn(args.Sides) + 1
	return strconv.Itoa(result), nil
}

func newRollAgent(ctx context.Context) (agent.Agent, error) {
	rollTool, err := functiontool.New(functiontool.Config{
		Name:        "roll_die",
		Description: "Roll a die and return the rolled result.",
	}, rollDieTool)
	if err != nil {
		return nil, fmt.Errorf("failed to create roll_die tool: %w", err)
	}

	model, err := gemini.NewModel(ctx, GeminiModel, &genai.ClientConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to create model for roll agent: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        RollAgentName,
		Description: "Handles rolling dice of different sizes.",
		Instruction: "You are responsible for rolling dice based on the user's request. When asked to roll a die, you must call the roll_die tool with the number of sides as an integer.",
		Model:       model,
		Tools:       []tool.Tool{rollTool},
	})
}

// --- Remote Prime Agent ---

// --8<-- [start:new-prime-agent]
func newPrimeAgent() (agent.Agent, error) {
	primeAgentURL := os.Getenv("ADK_PRIME_AGENT_URL")
	if primeAgentURL == "" {
		primeAgentURL = "http://localhost:8086"
	}
	remoteAgent, err := remoteagent.NewA2A(remoteagent.A2AConfig{
		Name:            PrimeAgentName,
		Description:     "Agent that handles checking if a single integer is prime.",
		AgentCardSource: primeAgentURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create remote prime agent: %w", err)
	}
	return remoteAgent, nil
}

// --8<-- [end:new-prime-agent]

// --- Root Agent ---

// --8<-- [start:new-root-agent]
func newRootAgent(ctx context.Context, rollAgent, primeAgent agent.Agent) (agent.Agent, error) {
	model, err := gemini.NewModel(ctx, GeminiModel, &genai.ClientConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to create model for root agent: %w", err)
	}
	return llmagent.New(llmagent.Config{
		Name:  RootAgentName,
		Model: model,
		Instruction: `
      You are a helpful assistant that can roll dice and check if numbers are prime.
      You delegate rolling dice tasks to the roll_agent and prime checking tasks to the prime_agent.
      Follow these steps:
      1. If the user asks to roll a die, delegate to the roll_agent.
      2. If the user asks to check primes, delegate to the prime_agent with a single integer.
      3. If the user asks to roll a die and then check if the result is prime, call roll_agent first, then pass the result to prime_agent.
      Always clarify the results before proceeding.
    `,
		SubAgents: []agent.Agent{rollAgent, primeAgent},
		Tools:     []tool.Tool{},
	})
}

// --8<-- [end:new-root-agent]

// SingleAgentLoader is a simple implementation of agent.Loader for a single agent.
type SingleAgentLoader struct {
	Agent agent.Agent
}

func (l *SingleAgentLoader) LoadAgent(name string) (agent.Agent, error) {
	if name == l.Agent.Name() {
		return l.Agent, nil
	}
	return nil, fmt.Errorf("agent not found: %s", name)
}

func (l *SingleAgentLoader) ListAgents() []string {
	return []string{l.Agent.Name()}
}

func (l *SingleAgentLoader) RootAgent() agent.Agent {
	return l.Agent
}

// --- Main Function ---

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	primeAgent, err := newPrimeAgent()
	if err != nil {
		slog.Error("Failed to create prime agent", "error", err)
		os.Exit(1)
	}

	rollAgent, err := newRollAgent(ctx)
	if err != nil {
		slog.Error("Failed to create roll agent", "error", err)
		os.Exit(1)
	}

	rootAgent, err := newRootAgent(ctx, rollAgent, primeAgent)
	if err != nil {
		slog.Error("Failed to create root agent", "error", err)
		os.Exit(1)
	}

	// Create launcher.
	l := prod.NewLauncher()

	// Allow PORT to be set by the environment (e.g., Cloud Run), default to 8092
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8092"
	}

	// Create ADK config
	config := &launcher.Config{
		AgentLoader:    &SingleAgentLoader{Agent: rootAgent},
		SessionService: session.InMemoryService(),
	}

	slog.Info("Starting A2A root agent server", "port", portStr)

	// Arguments for the launcher.
	args := []string{
		"--port", portStr,
		"a2a",
		"--a2a_agent_url", "http://0.0.0.0:" + portStr,
	}

	// Run launcher
	if err := l.Execute(ctx, config, args); err != nil {
		slog.Error("launcher.Run() error", "error", err)
		os.Exit(1)
	}
}
