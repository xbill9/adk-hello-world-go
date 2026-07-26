package main

import (
	"context"
	"iter"
	"testing"

	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

// MockLLM is a mock implementation of the model.LLM interface for testing purposes.
type MockLLM struct {
	ResponseText string
}

func (m *MockLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		if m.ResponseText != "" {
			resp := &model.LLMResponse{
				Content: &genai.Content{
					Role: "model",
					Parts: []*genai.Part{
						{Text: m.ResponseText},
					},
				},
			}
			yield(resp, nil)
		}
	}
}

func (m *MockLLM) Name() string {
	return "mock-model"
}

func TestSingleAgentLoader_ListAgents(t *testing.T) {
	mockAgent, err := llmagent.New(llmagent.Config{
		Name:  "test_agent",
		Model: &MockLLM{},
	})
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
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
	mockAgent, err := llmagent.New(llmagent.Config{
		Name:  "test_agent",
		Model: &MockLLM{ResponseText: "Hello from Mock!"},
	})
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
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
	if err == nil {
		t.Fatal("Expected error loading non-existent agent, got nil")
	}
	if ag != nil {
		t.Errorf("Expected nil for non-existent agent, got %v", ag)
	}

	// Test loading with empty string
	ag, err = loader.LoadAgent("")
	if err == nil {
		t.Fatal("Expected error loading empty string agent name, got nil")
	}
	if ag != nil {
		t.Errorf("Expected nil for empty string agent name, got %v", ag)
	}
}

func TestSingleAgentLoader_RootAgent(t *testing.T) {
	mockAgent, err := llmagent.New(llmagent.Config{
		Name:  "test_agent",
		Model: &MockLLM{},
	})
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	loader := &singleAgentLoader{agent: mockAgent}

	rootAgent := loader.RootAgent()
	if rootAgent == nil {
		t.Fatal("Expected root agent to be loaded, got nil")
	}
	if rootAgent.Name() != "test_agent" {
		t.Errorf("Expected root agent name 'test_agent', got '%s'", rootAgent.Name())
	}
}

func TestMockLLM_GenerateContent(t *testing.T) {
	mock := &MockLLM{ResponseText: "Mock Response"}
	if mock.Name() != "mock-model" {
		t.Errorf("Expected model name 'mock-model', got '%s'", mock.Name())
	}

	ctx := context.Background()
	req := &model.LLMRequest{}
	var resultText string
	for resp, err := range mock.GenerateContent(ctx, req, false) {
		if err != nil {
			t.Fatalf("Unexpected error during GenerateContent: %v", err)
		}
		if resp.Content != nil && len(resp.Content.Parts) > 0 {
			resultText = resp.Content.Parts[0].Text
		}
	}

	if resultText != "Mock Response" {
		t.Errorf("Expected 'Mock Response', got '%s'", resultText)
	}
}
