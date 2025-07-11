package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hackathon-2025/pkg/models"
	"io"
	"net/http"
)

type AIWorkflowRequest struct {
	Documents string `json:"documents"`
	UserInput string `json:"user_input"`
	Template  string `json:"template"`
}

type AIWorkflowResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  string      `json:"status,omitempty"`
}

type ProfileSummaryRequest struct {
	Documents string `json:"documents"`
	Template  string `json:"template"`
	UserEmail string `json:"user_email"`
}

type AgentResponse struct {
	Data struct {
		Outputs interface{} `json:"outputs,omitempty"`
		Status  string      `json:"status,omitempty"`
	} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

type AIAgent struct {
	BaseURL    string
	APIKey     string
	WorkflowID string
	Name       string
}

func NewAIAgent(name, workflowID, apiKey string) *AIAgent {
	return &AIAgent{
		BaseURL:    "https://ai.insea.io/api",
		APIKey:     apiKey,
		WorkflowID: workflowID,
		Name:       name,
	}
}

func (agent *AIAgent) GetName() string {
	return agent.Name
}

func (agent *AIAgent) GetWorkflowID() string {
	return agent.WorkflowID
}

func (agent *AIAgent) SetAPIKey(apiKey string) {
	agent.APIKey = apiKey
}

func (agent *AIAgent) RunProfileSummary(documents, template, userEmail string) (*AgentResponse, error) {
	workflowURL := fmt.Sprintf("%s/workflows/%s/run", agent.BaseURL, agent.WorkflowID)

	// Truncate documents if too long (let's limit to 8000 characters to be safe)
	if len(documents) > 90000 {
		documents = documents[:90000] + "... [truncated]"
	}

	payload := ProfileSummaryRequest{
		Documents: documents,
		Template:  template,
		UserEmail: userEmail,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", workflowURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", agent.APIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Always try to parse the response, even if it's an error
	var profileResponse AgentResponse
	if err := json.Unmarshal(body, &profileResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// If there's an error in the response, return it but don't treat as HTTP error
	if profileResponse.Error != nil {
		return &profileResponse, nil // Return the response with error, but don't fail
	}

	// Only treat HTTP errors as failures
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return &profileResponse, nil
}

func (agent *AIAgent) RunContentFilter(userContent string) (*AgentResponse, error) {
	workflowURL := fmt.Sprintf("%s/workflows/%s/run", agent.BaseURL, agent.WorkflowID)

	payload := models.FilterRequest{
		UserContent: userContent,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", workflowURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", agent.APIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Always try to parse the response, even if it's an error
	var filterResp AgentResponse
	if err := json.Unmarshal(body, &filterResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// If there's an error in the response, return it but don't treat as HTTP error
	if filterResp.Error != nil {
		return &filterResp, nil // Return the response with error, but don't fail
	}

	// Only treat HTTP errors as failures
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return &filterResp, nil
}

func ProfileSummaryAgent(apiKey string) *AIAgent {
	return NewAIAgent("profileSummaryAgent", "1989", apiKey)
}

func ContentFilterAgent(apiKey string) *AIAgent {
	return NewAIAgent("ContentFilterAgent", "2014", apiKey)
}
