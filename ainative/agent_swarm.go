package ainative

import (
	"context"
	"fmt"
	"time"
)

// AgentSwarmService handles agent swarm operations
type AgentSwarmService struct {
	client *Client
}

// NewAgentSwarmService creates a new agent swarm service
func NewAgentSwarmService(client *Client) *AgentSwarmService {
	return &AgentSwarmService{
		client: client,
	}
}

// AgentSwarm represents an agent swarm
type AgentSwarm struct {
	ID          string                 `json:"id"`
	ProjectID   string                 `json:"project_id"`
	Name        string                 `json:"name"`
	Objective   string                 `json:"objective"`
	Status      SwarmStatus            `json:"status"`
	Agents      []Agent                `json:"agents"`
	Metrics     *SwarmMetrics          `json:"metrics,omitempty"`
	Config      *SwarmConfig           `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// SwarmStatus represents the status of a swarm
type SwarmStatus string

const (
	SwarmStatusInitializing SwarmStatus = "initializing"
	SwarmStatusRunning      SwarmStatus = "running"
	SwarmStatusPaused       SwarmStatus = "paused"
	SwarmStatusCompleted    SwarmStatus = "completed"
	SwarmStatusFailed       SwarmStatus = "failed"
	SwarmStatusStopped      SwarmStatus = "stopped"
)

// Agent represents an individual agent in a swarm
type Agent struct {
	ID          string                 `json:"id"`
	Type        AgentType              `json:"type"`
	Status      AgentStatus            `json:"status"`
	Capabilities []string              `json:"capabilities"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Metrics     *AgentMetrics          `json:"metrics,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// AgentType represents different types of agents
type AgentType string

const (
	AgentTypeAnalyzer        AgentType = "analyzer"
	AgentTypeGenerator       AgentType = "generator"
	AgentTypeOptimizer       AgentType = "optimizer"
	AgentTypeValidator       AgentType = "validator"
	AgentTypeCoordinator     AgentType = "coordinator"
	AgentTypeSecurityScanner AgentType = "security_scanner"
	AgentTypeCodeReviewer    AgentType = "code_reviewer"
	AgentTypeDocumentWriter  AgentType = "document_writer"
)

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	AgentStatusIdle       AgentStatus = "idle"
	AgentStatusActive     AgentStatus = "active"
	AgentStatusWorking    AgentStatus = "working"
	AgentStatusCompleted  AgentStatus = "completed"
	AgentStatusFailed     AgentStatus = "failed"
	AgentStatusTerminated AgentStatus = "terminated"
)

// SwarmMetrics represents performance metrics for a swarm
type SwarmMetrics struct {
	TasksCompleted    int           `json:"tasks_completed"`
	TasksFailed       int           `json:"tasks_failed"`
	TasksInProgress   int           `json:"tasks_in_progress"`
	AverageTaskTime   time.Duration `json:"average_task_time"`
	TotalExecutionTime time.Duration `json:"total_execution_time"`
	Efficiency        float64       `json:"efficiency"`
	ResourceUsage     *ResourceUsage `json:"resource_usage,omitempty"`
}

// AgentMetrics represents performance metrics for an individual agent
type AgentMetrics struct {
	TasksCompleted  int           `json:"tasks_completed"`
	TasksFailed     int           `json:"tasks_failed"`
	AverageTaskTime time.Duration `json:"average_task_time"`
	LastActiveAt    time.Time     `json:"last_active_at"`
}

// ResourceUsage represents resource consumption metrics
type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage int64   `json:"memory_usage_bytes"`
	NetworkIO   int64   `json:"network_io_bytes"`
	APICallsCounter int     `json:"api_calls_count"`
}

// SwarmConfig represents configuration for a swarm
type SwarmConfig struct {
	MaxConcurrentTasks int                    `json:"max_concurrent_tasks"`
	TaskTimeout        time.Duration          `json:"task_timeout"`
	RetryCount         int                    `json:"retry_count"`
	ResourceLimits     *ResourceLimits        `json:"resource_limits,omitempty"`
	CustomSettings     map[string]interface{} `json:"custom_settings,omitempty"`
}

// ResourceLimits represents resource limits for a swarm
type ResourceLimits struct {
	MaxCPU    float64 `json:"max_cpu"`
	MaxMemory int64   `json:"max_memory_bytes"`
	MaxAPICalls int     `json:"max_api_calls"`
}

// AgentConfig represents configuration for creating agents
type AgentConfig struct {
	Type         AgentType              `json:"type"`
	Count        int                    `json:"count"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
}

// StartSwarmRequest represents a request to start a swarm
type StartSwarmRequest struct {
	ProjectID string        `json:"project_id"`
	Name      string        `json:"name,omitempty"`
	Objective string        `json:"objective"`
	Agents    []AgentConfig `json:"agents"`
	Config    *SwarmConfig  `json:"config,omitempty"`
}

// ListSwarmsRequest represents a request to list swarms
type ListSwarmsRequest struct {
	ProjectID string      `json:"project_id,omitempty"`
	Status    SwarmStatus `json:"status,omitempty"`
	Limit     int         `json:"limit,omitempty"`
	Offset    int         `json:"offset,omitempty"`
}

// ListSwarmsResponse represents a response containing swarms
type ListSwarmsResponse struct {
	Swarms     []AgentSwarm `json:"swarms"`
	TotalCount int          `json:"total_count"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
}

// OrchestrationTask represents a task for agent orchestration
type OrchestrationTask struct {
	ID          string                 `json:"id"`
	SwarmID     string                 `json:"swarm_id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Input       map[string]interface{} `json:"input"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Priority    TaskPriority           `json:"priority"`
	Status      TaskStatus             `json:"status"`
	AssignedTo  []string               `json:"assigned_to,omitempty"`
	Result      map[string]interface{} `json:"result,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// TaskPriority represents task priority levels
type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

// TaskStatus represents task execution status
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// OrchestrationRequest represents a request to orchestrate agents
type OrchestrationRequest struct {
	SwarmID     string                 `json:"swarm_id"`
	Task        string                 `json:"task"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Priority    TaskPriority           `json:"priority,omitempty"`
	AgentIDs    []string               `json:"agent_ids,omitempty"`
}

// OrchestrationResponse represents a response from agent orchestration
type OrchestrationResponse struct {
	TaskID      string                 `json:"task_id"`
	Status      TaskStatus             `json:"status"`
	AssignedTo  []string               `json:"assigned_to"`
	EstimatedDuration time.Duration     `json:"estimated_duration,omitempty"`
}

// Start starts a new agent swarm
func (s *AgentSwarmService) Start(ctx context.Context, req *StartSwarmRequest) (*AgentSwarm, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if req.ProjectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", req.ProjectID)
	}
	
	if req.Objective == "" {
		return nil, NewValidationError("objective", "objective is required", req.Objective)
	}
	
	if len(req.Agents) == 0 {
		return nil, NewValidationError("agents", "at least one agent configuration is required", req.Agents)
	}
	
	var result AgentSwarm
	
	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-swarm/swarms", req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Get retrieves a specific swarm
func (s *AgentSwarmService) Get(ctx context.Context, swarmID string) (*AgentSwarm, error) {
	if swarmID == "" {
		return nil, NewValidationError("swarm_id", "swarm ID is required", swarmID)
	}
	
	var result AgentSwarm
	
	path := fmt.Sprintf("/api/v1/agent-swarm/swarms/%s", swarmID)
	
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// List lists swarms with optional filtering
func (s *AgentSwarmService) List(ctx context.Context, req *ListSwarmsRequest) (*ListSwarmsResponse, error) {
	if req == nil {
		req = &ListSwarmsRequest{}
	}
	
	if req.Limit == 0 {
		req.Limit = 10
	}
	
	path := fmt.Sprintf("/api/v1/agent-swarm/swarms?limit=%d&offset=%d", req.Limit, req.Offset)
	
	if req.ProjectID != "" {
		path += fmt.Sprintf("&project_id=%s", req.ProjectID)
	}
	
	if req.Status != "" {
		path += fmt.Sprintf("&status=%s", req.Status)
	}
	
	var result ListSwarmsResponse
	
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Stop stops a running swarm
func (s *AgentSwarmService) Stop(ctx context.Context, swarmID string) error {
	if swarmID == "" {
		return NewValidationError("swarm_id", "swarm ID is required", swarmID)
	}
	
	path := fmt.Sprintf("/api/v1/agent-swarm/swarms/%s/stop", swarmID)
	
	return s.client.makeRequest(ctx, "POST", path, nil, nil)
}

// Pause pauses a running swarm
func (s *AgentSwarmService) Pause(ctx context.Context, swarmID string) error {
	if swarmID == "" {
		return NewValidationError("swarm_id", "swarm ID is required", swarmID)
	}
	
	path := fmt.Sprintf("/api/v1/agent-swarm/swarms/%s/pause", swarmID)
	
	return s.client.makeRequest(ctx, "POST", path, nil, nil)
}

// Resume resumes a paused swarm
func (s *AgentSwarmService) Resume(ctx context.Context, swarmID string) error {
	if swarmID == "" {
		return NewValidationError("swarm_id", "swarm ID is required", swarmID)
	}
	
	path := fmt.Sprintf("/api/v1/agent-swarm/swarms/%s/resume", swarmID)
	
	return s.client.makeRequest(ctx, "POST", path, nil, nil)
}

// Orchestrate orchestrates agents to perform a task
func (s *AgentSwarmService) Orchestrate(ctx context.Context, req *OrchestrationRequest) (*OrchestrationResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if req.SwarmID == "" {
		return nil, NewValidationError("swarm_id", "swarm ID is required", req.SwarmID)
	}
	
	if req.Task == "" {
		return nil, NewValidationError("task", "task is required", req.Task)
	}
	
	if req.Priority == "" {
		req.Priority = TaskPriorityMedium
	}
	
	var result OrchestrationResponse
	
	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-swarm/orchestrate", req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// GetTask retrieves a specific orchestration task
func (s *AgentSwarmService) GetTask(ctx context.Context, taskID string) (*OrchestrationTask, error) {
	if taskID == "" {
		return nil, NewValidationError("task_id", "task ID is required", taskID)
	}
	
	var result OrchestrationTask
	
	path := fmt.Sprintf("/api/v1/agent-swarm/tasks/%s", taskID)
	
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// ListAgentTypes lists available agent types
func (s *AgentSwarmService) ListAgentTypes(ctx context.Context) ([]AgentType, error) {
	var result struct {
		AgentTypes []AgentType `json:"agent_types"`
	}
	
	err := s.client.makeRequest(ctx, "GET", "/api/v1/agent-swarm/agent-types", nil, &result)
	if err != nil {
		return nil, err
	}
	
	return result.AgentTypes, nil
}

// GetSwarmMetrics retrieves detailed metrics for a swarm
func (s *AgentSwarmService) GetSwarmMetrics(ctx context.Context, swarmID string) (*SwarmMetrics, error) {
	if swarmID == "" {
		return nil, NewValidationError("swarm_id", "swarm ID is required", swarmID)
	}
	
	var result SwarmMetrics
	
	path := fmt.Sprintf("/api/v1/agent-swarm/swarms/%s/metrics", swarmID)
	
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}