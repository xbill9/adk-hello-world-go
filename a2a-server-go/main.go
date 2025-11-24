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
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/prod"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// isPrime checks if a number is prime.
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

type checkPrimeToolArgs struct {
	Num int `json:"num" jsonschema:"A number to check for primality."`
}

func checkPrimeTool(tc tool.Context, args checkPrimeToolArgs) (string, error) {
	if isPrime(args.Num) {
		return fmt.Sprintf("%d is a prime number.", args.Num), nil
	}
	return fmt.Sprintf("%d is not a prime number.", args.Num), nil
}

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

// --8<-- [start:a2a-launcher]
func main() {
	ctx := context.Background()
	// Initialize structured logging with JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	primeTool, err := functiontool.New(functiontool.Config{
		Name:        "prime_checking",
		Description: "Check if a number is prime using efficient mathematical algorithms",
	}, checkPrimeTool)
	if err != nil {
		slog.Error("Failed to create prime_checking tool", "error", err)
		os.Exit(1)
	}

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{})
	if err != nil {
		slog.Error("Failed to create model", "error", err)
		os.Exit(1)
	}

	primeAgent, err := llmagent.New(llmagent.Config{
		Name:        "check_prime_agent",
		Description: "check prime agent that can check whether a number is prime.",
		Instruction: `
			You check whether a number is prime.
			When checking if a number is prime, call the check_prime tool with a single integer.
			You should not rely on the previous history on prime results.
    `,
		Model: model,
		Tools: []tool.Tool{primeTool},
	})
	if err != nil {
		slog.Error("Failed to create agent", "error", err)
		os.Exit(1)
	}

	// Create launcher.
	l := prod.NewLauncher()

	// Allow PORT to be set by the environment (e.g., Cloud Run), default to 8086
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8086"
	}
	// Set PORT env var for the launcher to pick up
	os.Setenv("PORT", portStr)

	// Create ADK config
	config := &launcher.Config{
		AgentLoader:    &SingleAgentLoader{Agent: primeAgent},
		SessionService: session.InMemoryService(),
	}

	slog.Info("Starting A2A prime checker server", "port", portStr)

	// Arguments for the launcher.
	// Note: ParseAndRun usually expects the first argument to be the program name if it parses full os.Args,
	// but here we are constructing args manually.
	// If full launcher uses standard flag parsing, it might expect the command "a2a" as a subcommand.
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

// --8<-- [end:a2a-launcher]
