package services

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type ProfileSummaryResponse struct {
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

func (agent *AIAgent) Run(documents, userInput, template string) (*AIWorkflowResponse, error) {
	workflowURL := fmt.Sprintf("%s/workflows/%s/run", agent.BaseURL, agent.WorkflowID)

	payload := AIWorkflowRequest{
		Documents: documents,
		UserInput: userInput,
		Template:  template,
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var aiResponse AIWorkflowResponse
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &aiResponse, nil
}

func (agent *AIAgent) RunWithEmptyPayload() (*AIWorkflowResponse, error) {
	return agent.Run("", "", "")
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

func (agent *AIAgent) RunProfileSummary(documents, template, userEmail string) (*ProfileSummaryResponse, error) {
	workflowURL := fmt.Sprintf("%s/workflows/%s/run", agent.BaseURL, agent.WorkflowID)

	// Truncate documents if too long (let's limit to 8000 characters to be safe)
	if len(documents) > 10000 {
		documents = documents[:10000] + "... [truncated]"
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
	var profileResponse ProfileSummaryResponse
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

func ProfileSummaryAgent(apiKey string) *AIAgent {
	return NewAIAgent("profileSummaryAgent", "1989", apiKey)
}
