package ainative

import (
	"context"
	"fmt"
	"time"
)

// AgentStateService handles agent state management operations
type AgentStateService struct {
	client *Client
}

// NewAgentStateService creates a new agent state service
func NewAgentStateService(client *Client) *AgentStateService {
	return &AgentStateService{
		client: client,
	}
}

// AgentState represents the current state of an agent
type AgentState struct {
	AgentID   string                 `json:"agent_id"`
	Version   int                    `json:"version"`
	State     map[string]interface{} `json:"state"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// GetStateRequest represents a request to get agent state
type GetStateRequest struct {
	AgentID string `json:"agent_id"`
	Version int    `json:"version,omitempty"` // Optional: specific version
}

// GetStateResponse represents the response containing agent state
type GetStateResponse struct {
	AgentID   string                 `json:"agent_id"`
	Version   int                    `json:"version"`
	State     map[string]interface{} `json:"state"`
	Timestamp time.Time              `json:"timestamp"`
}

// Checkpoint represents a saved state checkpoint
type Checkpoint struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data"`
	Version     int                    `json:"version"`
	CreatedAt   time.Time              `json:"created_at"`
}

// CreateCheckpointRequest represents a request to create a state checkpoint
type CreateCheckpointRequest struct {
	AgentID     string                 `json:"agent_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data"`
}

// CreateCheckpointResponse represents the response from creating a checkpoint
type CreateCheckpointResponse struct {
	CheckpointID string    `json:"checkpoint_id"`
	Version      int       `json:"version"`
	Message      string    `json:"message,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// RestoreCheckpointRequest represents a request to restore from a checkpoint
type RestoreCheckpointRequest struct {
	CheckpointID string `json:"checkpoint_id"`
	AgentID      string `json:"agent_id"`
}

// RestoreCheckpointResponse represents the response from restoring a checkpoint
type RestoreCheckpointResponse struct {
	AgentID   string                 `json:"agent_id"`
	Version   int                    `json:"version"`
	State     map[string]interface{} `json:"state"`
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
}

// ListCheckpointsRequest represents a request to list checkpoints
type ListCheckpointsRequest struct {
	AgentID string `json:"agent_id,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}

// ListCheckpointsResponse represents the response containing checkpoints
type ListCheckpointsResponse struct {
	Checkpoints []Checkpoint `json:"checkpoints"`
	TotalCount  int          `json:"total_count"`
	Limit       int          `json:"limit"`
	Offset      int          `json:"offset"`
}

// GetState retrieves the current state of an agent
func (s *AgentStateService) GetState(ctx context.Context, req *GetStateRequest) (*GetStateResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.AgentID == "" {
		return nil, NewValidationError("agent_id", "agent ID is required", req.AgentID)
	}

	path := fmt.Sprintf("/api/v1/agent-state/state?agent_id=%s", req.AgentID)

	if req.Version > 0 {
		path += fmt.Sprintf("&version=%d", req.Version)
	}

	var result GetStateResponse

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateCheckpoint creates a state checkpoint for an agent
func (s *AgentStateService) CreateCheckpoint(ctx context.Context, req *CreateCheckpointRequest) (*CreateCheckpointResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.AgentID == "" {
		return nil, NewValidationError("agent_id", "agent ID is required", req.AgentID)
	}

	if req.Name == "" {
		return nil, NewValidationError("name", "checkpoint name is required", req.Name)
	}

	if req.Data == nil {
		return nil, NewValidationError("data", "checkpoint data is required", req.Data)
	}

	var result CreateCheckpointResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-state/checkpoints", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RestoreCheckpoint restores an agent's state from a checkpoint
func (s *AgentStateService) RestoreCheckpoint(ctx context.Context, req *RestoreCheckpointRequest) (*RestoreCheckpointResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.CheckpointID == "" {
		return nil, NewValidationError("checkpoint_id", "checkpoint ID is required", req.CheckpointID)
	}

	if req.AgentID == "" {
		return nil, NewValidationError("agent_id", "agent ID is required", req.AgentID)
	}

	var result RestoreCheckpointResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-state/restore", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListCheckpoints retrieves a list of available checkpoints
func (s *AgentStateService) ListCheckpoints(ctx context.Context, req *ListCheckpointsRequest) (*ListCheckpointsResponse, error) {
	if req == nil {
		req = &ListCheckpointsRequest{}
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	path := fmt.Sprintf("/api/v1/agent-state/checkpoints?limit=%d&offset=%d", req.Limit, req.Offset)

	if req.AgentID != "" {
		path += fmt.Sprintf("&agent_id=%s", req.AgentID)
	}

	var result ListCheckpointsResponse

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
