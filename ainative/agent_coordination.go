package ainative

import (
	"context"
	"fmt"
	"time"
)

// AgentCoordinationService handles agent coordination operations
type AgentCoordinationService struct {
	client *Client
}

// NewAgentCoordinationService creates a new agent coordination service
func NewAgentCoordinationService(client *Client) *AgentCoordinationService {
	return &AgentCoordinationService{
		client: client,
	}
}

// AgentMessage represents a message between agents
type AgentMessage struct {
	ID          string                 `json:"id"`
	FromAgent   string                 `json:"from_agent"`
	ToAgent     string                 `json:"to_agent"`
	Message     string                 `json:"message"`
	MessageType string                 `json:"message_type,omitempty"`
	Priority    string                 `json:"priority,omitempty"` // low, medium, high, urgent
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	DeliveredAt *time.Time             `json:"delivered_at,omitempty"`
}

// SendMessageRequest represents a request to send a message between agents
type SendMessageRequest struct {
	FromAgent   string                 `json:"from_agent"`
	ToAgent     string                 `json:"to_agent"`
	Message     string                 `json:"message"`
	MessageType string                 `json:"message_type,omitempty"`
	Priority    string                 `json:"priority,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SendMessageResponse represents the response from sending a message
type SendMessageResponse struct {
	MessageID string    `json:"message_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// DistributeTasksRequest represents a request to distribute tasks across agents
type DistributeTasksRequest struct {
	Tasks       []string               `json:"tasks"`
	Agents      []string               `json:"agents,omitempty"`
	Strategy    string                 `json:"strategy,omitempty"` // round_robin, load_balanced, capability_based
	Constraints map[string]interface{} `json:"constraints,omitempty"`
}

// TaskAssignment represents a task assigned to an agent
type TaskAssignment struct {
	TaskID  string `json:"task_id"`
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}

// DistributeTasksResponse represents the response from task distribution
type DistributeTasksResponse struct {
	Assignments []TaskAssignment `json:"assignments"`
	Strategy    string           `json:"strategy"`
	Timestamp   time.Time        `json:"timestamp"`
}

// AgentWorkload represents workload information for an agent
type AgentWorkload struct {
	AgentID        string  `json:"agent_id"`
	ActiveTasks    int     `json:"active_tasks"`
	QueuedTasks    int     `json:"queued_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	FailedTasks    int     `json:"failed_tasks"`
	CPUUsage       float64 `json:"cpu_usage,omitempty"`
	MemoryUsage    float64 `json:"memory_usage,omitempty"`
	Availability   float64 `json:"availability"` // 0.0 to 1.0
}

// GetWorkloadStatsRequest represents a request for workload statistics
type GetWorkloadStatsRequest struct {
	AgentIDs []string `json:"agent_ids,omitempty"`
}

// GetWorkloadStatsResponse represents workload statistics across agents
type GetWorkloadStatsResponse struct {
	Workloads     []AgentWorkload `json:"workloads"`
	TotalTasks    int             `json:"total_tasks"`
	AverageLoad   float64         `json:"average_load"`
	Timestamp     time.Time       `json:"timestamp"`
}

// SendMessage sends a message from one agent to another
func (s *AgentCoordinationService) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.FromAgent == "" {
		return nil, NewValidationError("from_agent", "from agent is required", req.FromAgent)
	}

	if req.ToAgent == "" {
		return nil, NewValidationError("to_agent", "to agent is required", req.ToAgent)
	}

	if req.Message == "" {
		return nil, NewValidationError("message", "message is required", req.Message)
	}

	var result SendMessageResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-coordination/messages", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DistributeTasks distributes tasks across multiple agents
func (s *AgentCoordinationService) DistributeTasks(ctx context.Context, req *DistributeTasksRequest) (*DistributeTasksResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if len(req.Tasks) == 0 {
		return nil, NewValidationError("tasks", "at least one task is required", req.Tasks)
	}

	// Set default strategy if not provided
	if req.Strategy == "" {
		req.Strategy = "load_balanced"
	}

	var result DistributeTasksResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-coordination/distribute", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWorkloadStats retrieves workload statistics for agents
func (s *AgentCoordinationService) GetWorkloadStats(ctx context.Context, req *GetWorkloadStatsRequest) (*GetWorkloadStatsResponse, error) {
	if req == nil {
		req = &GetWorkloadStatsRequest{}
	}

	path := "/api/v1/agent-coordination/workload"

	// Add agent IDs as query parameters if provided
	if len(req.AgentIDs) > 0 {
		path += "?agent_ids=" + fmt.Sprintf("%v", req.AgentIDs)
	}

	var result GetWorkloadStatsResponse

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
