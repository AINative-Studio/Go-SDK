package ainative

import (
	"context"
	"fmt"
	"time"
)

// AgentOrchestrationService handles agent task orchestration operations
type AgentOrchestrationService struct {
	client *Client
}

// NewAgentOrchestrationService creates a new agent orchestration service
func NewAgentOrchestrationService(client *Client) *AgentOrchestrationService {
	return &AgentOrchestrationService{
		client: client,
	}
}

// Task represents an agent task
type Task struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	TaskType    string                 `json:"task_type"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"` // pending, running, completed, failed
	Priority    string                 `json:"priority,omitempty"` // low, medium, high, critical
	Context     map[string]interface{} `json:"context,omitempty"`
	Result      interface{}            `json:"result,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	AgentID     string                 `json:"agent_id"`
	TaskType    string                 `json:"task_type"`
	Description string                 `json:"description"`
	Priority    string                 `json:"priority,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// CreateTaskResponse represents a response from task creation
type CreateTaskResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Task   Task   `json:"task"`
}

// ListTasksRequest represents a request to list tasks
type ListTasksRequest struct {
	AgentID string `json:"agent_id,omitempty"`
	Status  string `json:"status,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}

// ListTasksResponse represents a response containing tasks
type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Total int    `json:"total"`
}

// TaskStatusResponse represents task status details
type TaskStatusResponse struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"`
	Progress *int        `json:"progress,omitempty"`
	Message  string      `json:"message,omitempty"`
	Result   interface{} `json:"result,omitempty"`
}

// ExecuteTaskRequest represents a request to execute a task
type ExecuteTaskRequest struct {
	TaskID string                 `json:"task_id"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// ExecuteTaskResponse represents a response from task execution
type ExecuteTaskResponse struct {
	ID     string      `json:"id"`
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}

// TaskSequence represents a sequence of tasks
type TaskSequence struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tasks     []string  `json:"tasks"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTaskSequenceRequest represents a request to create a task sequence
type CreateTaskSequenceRequest struct {
	Name        string   `json:"name"`
	Tasks       []string `json:"tasks"`
	Description string   `json:"description,omitempty"`
}

// CreateTaskSequenceResponse represents a response from task sequence creation
type CreateTaskSequenceResponse struct {
	ID       string       `json:"id"`
	Sequence TaskSequence `json:"sequence"`
}

// CreateTask creates a new agent task
func (s *AgentOrchestrationService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.AgentID == "" {
		return nil, NewValidationError("agent_id", "agent ID is required", req.AgentID)
	}

	if req.TaskType == "" {
		return nil, NewValidationError("task_type", "task type is required", req.TaskType)
	}

	if req.Description == "" {
		return nil, NewValidationError("description", "description is required", req.Description)
	}

	var result CreateTaskResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-orchestration/tasks", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListTasks lists tasks with optional filtering
func (s *AgentOrchestrationService) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
	if req == nil {
		req = &ListTasksRequest{}
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	path := fmt.Sprintf("/api/v1/agent-orchestration/tasks?limit=%d&offset=%d", req.Limit, req.Offset)

	if req.AgentID != "" {
		path += fmt.Sprintf("&agent_id=%s", req.AgentID)
	}

	if req.Status != "" {
		path += fmt.Sprintf("&status=%s", req.Status)
	}

	var result ListTasksResponse

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTaskStatus retrieves task status and progress
func (s *AgentOrchestrationService) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatusResponse, error) {
	if taskID == "" {
		return nil, NewValidationError("task_id", "task ID is required", taskID)
	}

	var result TaskStatusResponse

	path := fmt.Sprintf("/api/v1/agent-orchestration/tasks/%s/status", taskID)

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ExecuteTask executes a task
func (s *AgentOrchestrationService) ExecuteTask(ctx context.Context, req *ExecuteTaskRequest) (*ExecuteTaskResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.TaskID == "" {
		return nil, NewValidationError("task_id", "task ID is required", req.TaskID)
	}

	var result ExecuteTaskResponse

	path := fmt.Sprintf("/api/v1/agent-orchestration/tasks/%s/execute", req.TaskID)

	body := map[string]interface{}{
		"params": req.Params,
	}

	err := s.client.makeRequest(ctx, "POST", path, body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateTaskSequence creates a task sequence
func (s *AgentOrchestrationService) CreateTaskSequence(ctx context.Context, req *CreateTaskSequenceRequest) (*CreateTaskSequenceResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.Name == "" {
		return nil, NewValidationError("name", "name is required", req.Name)
	}

	if len(req.Tasks) == 0 {
		return nil, NewValidationError("tasks", "at least one task is required", req.Tasks)
	}

	var result CreateTaskSequenceResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-orchestration/sequences", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
