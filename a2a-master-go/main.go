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
	"log"
	"math/rand"
	"strconv"


	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/remoteagent"
	adkservices "google.golang.org/adk/server/restapi/services"

	"google.golang.org/adk/model/gemini"

	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

        "google.golang.org/adk/cmd/launcher/adk"
        "google.golang.org/adk/cmd/launcher/web"
        "google.golang.org/adk/cmd/launcher/web/a2a"
        

	"google.golang.org/genai"
)

// --- Local Roll Agent ---

type rollDieToolArgs struct {
	Sides int `json:"sides" jsonschema:"The number of sides on the die."`
}

func rollDieTool(tc tool.Context, args rollDieToolArgs) int {
	return rand.Intn(args.Sides) + 1
}

func newRollAgent(ctx context.Context) (agent.Agent, error) {
	rollTool, err := functiontool.New(functiontool.Config{
		Name:        "roll_die",
		Description: "Roll a die and return the rolled result.",
	}, rollDieTool)
	if err != nil {
		return nil, fmt.Errorf("failed to create roll_die tool: %w", err)
	}

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to create model for roll agent: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "roll_agent",
		Description: "Handles rolling dice of different sizes.",
		Instruction: "You are responsible for rolling dice based on the user's request. When asked to roll a die, you must call the roll_die tool with the number of sides as an integer.",
		Model:       model,
		Tools:       []tool.Tool{rollTool},
	})
}

// --- Remote Prime Agent ---

// --8<-- [start:new-prime-agent]
func newPrimeAgent() (agent.Agent, error) {
	remoteAgent, err := remoteagent.New(remoteagent.A2AConfig{
		Name:            "prime_agent",
		Description:     "Agent that handles checking if numbers are prime.",
		AgentCardSource: "http://localhost:8086",
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
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{})
	if err != nil {
		return nil, err
	}
	return llmagent.New(llmagent.Config{
		Name:  "root_agent",
		Model: model,
		Instruction: `
      You are a helpful assistant that can roll dice and check if numbers are prime.
      You delegate rolling dice tasks to the roll_agent and prime checking tasks to the prime_agent.
      Follow these steps:
      1. If the user asks to roll a die, delegate to the roll_agent.
      2. If the user asks to check primes, delegate to the prime_agent.
      3. If the user asks to roll a die and then check if the result is prime, call roll_agent first, then pass the result to prime_agent.
      Always clarify the results before proceeding.
    `,
		SubAgents: []agent.Agent{rollAgent, primeAgent},
		Tools:     []tool.Tool{},
	})
}



// --8<-- [end:new-root-agent]

// --- Main Function ---

func main() {
	ctx := context.Background()

	primeAgent, err := newPrimeAgent()
	if err != nil {
		log.Fatalf("Failed to create prime agent: %v", err)
	}

	rollAgent, err := newRollAgent(ctx)
	if err != nil {
		log.Fatalf("Failed to create roll agent: %v", err)
	}

	rootAgent, err := newRootAgent(ctx, rollAgent, primeAgent)
	if err != nil {
		log.Fatalf("Failed to create root agent: %v", err)
	}

	sessionService := session.InMemoryService()

	_, err = sessionService.Create(ctx, &session.CreateRequest{
		AppName:   rootAgent.Name(),
		UserID:    "user-123",
		SessionID: "session-abc",
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}



        port := 8092
        launcher := web.NewLauncher(a2a.NewLauncher())
        _, parseErr := launcher.Parse([]string{
                "--port", strconv.Itoa(port),
                "a2a", "--a2a_agent_url", "http://0.0.0.0:" + strconv.Itoa(port),
        })
        if parseErr != nil {
                log.Fatalf("launcher.Parse() error = %v", parseErr)
        }

        // Create ADK config
        config := &adk.Config{
                AgentLoader:    adkservices.NewSingleAgentLoader(rootAgent),
                SessionService: session.InMemoryService(),
        }

        log.Printf("Starting A2A prime checker server on port %d\n", port)
        // Run launcher
        if err := launcher.Run(context.Background(), config); err != nil {
                log.Fatalf("launcher.Run() error = %v", err)
        }


}
