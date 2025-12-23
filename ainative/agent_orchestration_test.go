package ainative

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentOrchestrationService_CreateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-orchestration/tasks", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "agent_123", req.AgentID)
		assert.Equal(t, "code_analysis", req.TaskType)
		assert.Equal(t, "Analyze codebase for vulnerabilities", req.Description)
		assert.Equal(t, "high", req.Priority)

		now := time.Now()
		response := CreateTaskResponse{
			ID:     "task_456",
			Status: "created",
			Task: Task{
				ID:          "task_456",
				AgentID:     req.AgentID,
				TaskType:    req.TaskType,
				Description: req.Description,
				Status:      "pending",
				Priority:    req.Priority,
				Context:     req.Context,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &CreateTaskRequest{
		AgentID:     "agent_123",
		TaskType:    "code_analysis",
		Description: "Analyze codebase for vulnerabilities",
		Priority:    "high",
		Context: map[string]interface{}{
			"repository": "test-repo",
			"branch":     "main",
		},
	}

	response, err := client.AgentOrchestration.CreateTask(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "task_456", response.ID)
	assert.Equal(t, "created", response.Status)
	assert.Equal(t, "agent_123", response.Task.AgentID)
	assert.Equal(t, "pending", response.Task.Status)
}

func TestAgentOrchestrationService_ListTasks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-orchestration/tasks", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "0", r.URL.Query().Get("offset"))
		assert.Equal(t, "agent_123", r.URL.Query().Get("agent_id"))
		assert.Equal(t, "completed", r.URL.Query().Get("status"))

		now := time.Now()
		response := ListTasksResponse{
			Tasks: []Task{
				{
					ID:          "task_1",
					AgentID:     "agent_123",
					TaskType:    "code_analysis",
					Description: "Task 1",
					Status:      "completed",
					CreatedAt:   now.Add(-2 * time.Hour),
					UpdatedAt:   now.Add(-1 * time.Hour),
				},
				{
					ID:          "task_2",
					AgentID:     "agent_123",
					TaskType:    "optimization",
					Description: "Task 2",
					Status:      "completed",
					CreatedAt:   now.Add(-4 * time.Hour),
					UpdatedAt:   now.Add(-3 * time.Hour),
				},
			},
			Total: 2,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &ListTasksRequest{
		AgentID: "agent_123",
		Status:  "completed",
		Limit:   10,
		Offset:  0,
	}

	response, err := client.AgentOrchestration.ListTasks(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Tasks))
	assert.Equal(t, 2, response.Total)
	assert.Equal(t, "task_1", response.Tasks[0].ID)
	assert.Equal(t, "completed", response.Tasks[0].Status)
}

func TestAgentOrchestrationService_GetTaskStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-orchestration/tasks/task_123/status", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		progress := 75
		response := TaskStatusResponse{
			ID:       "task_123",
			Status:   "running",
			Progress: &progress,
			Message:  "Processing data",
			Result: map[string]interface{}{
				"items_processed": 750,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	response, err := client.AgentOrchestration.GetTaskStatus(ctx, "task_123")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "task_123", response.ID)
	assert.Equal(t, "running", response.Status)
	assert.NotNil(t, response.Progress)
	assert.Equal(t, 75, *response.Progress)
	assert.Equal(t, "Processing data", response.Message)
}

func TestAgentOrchestrationService_ExecuteTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-orchestration/tasks/task_123/execute", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)
		assert.NotNil(t, body["params"])

		response := ExecuteTaskResponse{
			ID:     "task_123",
			Status: "executing",
			Result: map[string]interface{}{
				"execution_id": "exec_789",
				"started_at":   time.Now(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &ExecuteTaskRequest{
		TaskID: "task_123",
		Params: map[string]interface{}{
			"timeout":    300,
			"max_retries": 3,
		},
	}

	response, err := client.AgentOrchestration.ExecuteTask(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "task_123", response.ID)
	assert.Equal(t, "executing", response.Status)
	assert.NotNil(t, response.Result)
}

func TestAgentOrchestrationService_CreateTaskSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-orchestration/sequences", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateTaskSequenceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "CI/CD Pipeline", req.Name)
		assert.Equal(t, 3, len(req.Tasks))

		now := time.Now()
		response := CreateTaskSequenceResponse{
			ID: "sequence_123",
			Sequence: TaskSequence{
				ID:        "sequence_123",
				Name:      req.Name,
				Tasks:     req.Tasks,
				Status:    "created",
				CreatedAt: now,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &CreateTaskSequenceRequest{
		Name:        "CI/CD Pipeline",
		Tasks:       []string{"task_1", "task_2", "task_3"},
		Description: "Run tests, build, and deploy",
	}

	response, err := client.AgentOrchestration.CreateTaskSequence(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "sequence_123", response.ID)
	assert.Equal(t, "CI/CD Pipeline", response.Sequence.Name)
	assert.Equal(t, 3, len(response.Sequence.Tasks))
	assert.Equal(t, "created", response.Sequence.Status)
}

func TestAgentOrchestrationService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test CreateTask with nil request
	_, err = client.AgentOrchestration.CreateTask(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test CreateTask with empty agent_id
	_, err = client.AgentOrchestration.CreateTask(ctx, &CreateTaskRequest{AgentID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent ID is required")

	// Test CreateTask with empty task_type
	_, err = client.AgentOrchestration.CreateTask(ctx, &CreateTaskRequest{
		AgentID:  "agent_123",
		TaskType: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task type is required")

	// Test CreateTask with empty description
	_, err = client.AgentOrchestration.CreateTask(ctx, &CreateTaskRequest{
		AgentID:     "agent_123",
		TaskType:    "analysis",
		Description: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "description is required")

	// Test GetTaskStatus with empty task_id
	_, err = client.AgentOrchestration.GetTaskStatus(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task ID is required")

	// Test ExecuteTask with nil request
	_, err = client.AgentOrchestration.ExecuteTask(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test ExecuteTask with empty task_id
	_, err = client.AgentOrchestration.ExecuteTask(ctx, &ExecuteTaskRequest{TaskID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task ID is required")

	// Test CreateTaskSequence with nil request
	_, err = client.AgentOrchestration.CreateTaskSequence(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test CreateTaskSequence with empty name
	_, err = client.AgentOrchestration.CreateTaskSequence(ctx, &CreateTaskSequenceRequest{Name: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test CreateTaskSequence with empty tasks
	_, err = client.AgentOrchestration.CreateTaskSequence(ctx, &CreateTaskSequenceRequest{
		Name:  "Test Sequence",
		Tasks: []string{},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one task is required")
}
