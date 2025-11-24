package main

import (
	"context"
	"testing"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"iter"
)

// MockLLM is a mock implementation of the model.LLM interface for testing purposes.
type MockLLM struct{}

func (m *MockLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		// For mocking, we can yield nothing, or a predefined response.
		// For now, an empty sequence is fine as the tests are not concerned with the LLM's output.
	}
}

func (m *MockLLM) Name() string {
	return "mock-model"
}

func TestSingleAgentLoader_ListAgents(t *testing.T) {
	mockAgent, _ := llmagent.New(llmagent.Config{
		Name:  "test_agent",
		Model: &MockLLM{},
	})
	loader := &singleAgentLoader{agent: mockAgent}

	agents := loader.ListAgents()
	if len(agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(agents))
	}
	if agents[0] != "test_agent" {
		t.Errorf("Expected agent name 'test_agent', got '%s'", agents[0])
	}
}

func TestSingleAgentLoader_LoadAgent(t *testing.T) {
	mockAgent, _ := llmagent.New(llmagent.Config{
		Name:  "test_agent",
		Model: &MockLLM{},
	})
	loader := &singleAgentLoader{agent: mockAgent}

	// Test loading the correct agent
	ag, err := loader.LoadAgent("test_agent")
	if err != nil {
		t.Fatalf("Unexpected error loading agent: %v", err)
	}
	if ag == nil {
		t.Fatal("Expected agent to be loaded, got nil")
	}
	if ag.Name() != "test_agent" {
		t.Errorf("Expected agent name 'test_agent', got '%s'", ag.Name())
	}

	// Test loading a non-existent agent
	ag, err = loader.LoadAgent("non_existent_agent")
	if err != nil {
		t.Fatalf("Unexpected error loading non-existent agent: %v", err)
	}
	if ag != nil {
		t.Errorf("Expected nil for non-existent agent, got %v", ag)
	}
}

func TestSingleAgentLoader_RootAgent(t *testing.T) {
	mockAgent, _ := llmagent.New(llmagent.Config{
		Name:  "test_agent",
		Model: &MockLLM{},
	})
	loader := &singleAgentLoader{agent: mockAgent}

	rootAgent := loader.RootAgent()
	if rootAgent == nil {
		t.Fatal("Expected root agent to be loaded, got nil")
	}
	if rootAgent.Name() != "test_agent" {
		t.Errorf("Expected root agent name 'test_agent', got '%s'", rootAgent.Name())
	}
}
