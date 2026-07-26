package main

import (
	"testing"

	"google.golang.org/adk/v2/agent"
)

// MockAgent embeds agent.Agent to satisfy the interface requirements
// for methods we don't need to implement for these specific tests.
type MockAgent struct {
	agent.Agent
}

// Name implements the agent.Agent Name method.
func (m *MockAgent) Name() string {
	return "test-agent"
}

func TestSingleAgentLoader_ListAgents(t *testing.T) {
	mockAg := &MockAgent{}
	loader := &singleAgentLoader{agent: mockAg}

	agents := loader.ListAgents()
	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}
	if agents[0] != "test-agent" {
		t.Errorf("expected agent name 'test-agent', got %q", agents[0])
	}
}

func TestSingleAgentLoader_LoadAgent(t *testing.T) {
	mockAg := &MockAgent{}
	loader := &singleAgentLoader{agent: mockAg}

	// Test matching name
	ag, err := loader.LoadAgent("test-agent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ag == nil {
		t.Fatal("expected agent, got nil")
	}
	if ag.Name() != "test-agent" {
		t.Errorf("expected agent name 'test-agent', got %q", ag.Name())
	}

	// Test non-matching name
	ag, err = loader.LoadAgent("other-agent")
	if err == nil {
		t.Fatal("expected error for non-matching name, got nil")
	}
	if ag != nil {
		t.Errorf("expected nil agent, got %v", ag)
	}
}

func TestSingleAgentLoader_RootAgent(t *testing.T) {
	mockAg := &MockAgent{}
	loader := &singleAgentLoader{agent: mockAg}

	ag := loader.RootAgent()
	if ag == nil {
		t.Fatal("expected agent, got nil")
	}
	if ag.Name() != "test-agent" {
		t.Errorf("expected agent name 'test-agent', got %q", ag.Name())
	}
}
