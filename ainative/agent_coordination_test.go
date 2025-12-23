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

func TestAgentCoordinationService_SendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-coordination/messages", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req SendMessageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "agent_1", req.FromAgent)
		assert.Equal(t, "agent_2", req.ToAgent)
		assert.Equal(t, "Hello from agent 1", req.Message)
		assert.Equal(t, "high", req.Priority)

		response := SendMessageResponse{
			MessageID: "msg_123",
			Status:    "delivered",
			Timestamp: time.Now(),
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
	req := &SendMessageRequest{
		FromAgent:   "agent_1",
		ToAgent:     "agent_2",
		Message:     "Hello from agent 1",
		MessageType: "request",
		Priority:    "high",
		Metadata: map[string]interface{}{
			"task_id": "task_123",
		},
	}

	response, err := client.AgentCoordination.SendMessage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "msg_123", response.MessageID)
	assert.Equal(t, "delivered", response.Status)
}

func TestAgentCoordinationService_DistributeTasks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-coordination/distribute", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req DistributeTasksRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(req.Tasks))
		assert.Equal(t, "load_balanced", req.Strategy)

		response := DistributeTasksResponse{
			Assignments: []TaskAssignment{
				{
					TaskID:  "task_1",
					AgentID: "agent_1",
					Status:  "assigned",
				},
				{
					TaskID:  "task_2",
					AgentID: "agent_2",
					Status:  "assigned",
				},
				{
					TaskID:  "task_3",
					AgentID: "agent_3",
					Status:  "assigned",
				},
			},
			Strategy:  "load_balanced",
			Timestamp: time.Now(),
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
	req := &DistributeTasksRequest{
		Tasks:    []string{"task_1", "task_2", "task_3"},
		Agents:   []string{"agent_1", "agent_2", "agent_3"},
		Strategy: "load_balanced",
		Constraints: map[string]interface{}{
			"max_tasks_per_agent": 5,
		},
	}

	response, err := client.AgentCoordination.DistributeTasks(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Assignments))
	assert.Equal(t, "load_balanced", response.Strategy)
	assert.Equal(t, "task_1", response.Assignments[0].TaskID)
	assert.Equal(t, "agent_1", response.Assignments[0].AgentID)
}

func TestAgentCoordinationService_GetWorkloadStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-coordination/workload", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := GetWorkloadStatsResponse{
			Workloads: []AgentWorkload{
				{
					AgentID:        "agent_1",
					ActiveTasks:    3,
					QueuedTasks:    2,
					CompletedTasks: 10,
					FailedTasks:    1,
					CPUUsage:       0.65,
					MemoryUsage:    0.45,
					Availability:   0.7,
				},
				{
					AgentID:        "agent_2",
					ActiveTasks:    2,
					QueuedTasks:    1,
					CompletedTasks: 15,
					FailedTasks:    0,
					CPUUsage:       0.45,
					MemoryUsage:    0.35,
					Availability:   0.85,
				},
			},
			TotalTasks:  5,
			AverageLoad: 0.775,
			Timestamp:   time.Now(),
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
	req := &GetWorkloadStatsRequest{
		AgentIDs: []string{"agent_1", "agent_2"},
	}

	response, err := client.AgentCoordination.GetWorkloadStats(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Workloads))
	assert.Equal(t, 5, response.TotalTasks)
	assert.Equal(t, 0.775, response.AverageLoad)
	assert.Equal(t, "agent_1", response.Workloads[0].AgentID)
	assert.Equal(t, 3, response.Workloads[0].ActiveTasks)
}

func TestAgentCoordinationService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test SendMessage with nil request
	_, err = client.AgentCoordination.SendMessage(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test SendMessage with empty from_agent
	_, err = client.AgentCoordination.SendMessage(ctx, &SendMessageRequest{FromAgent: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "from agent is required")

	// Test SendMessage with empty to_agent
	_, err = client.AgentCoordination.SendMessage(ctx, &SendMessageRequest{
		FromAgent: "agent_1",
		ToAgent:   "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "to agent is required")

	// Test SendMessage with empty message
	_, err = client.AgentCoordination.SendMessage(ctx, &SendMessageRequest{
		FromAgent: "agent_1",
		ToAgent:   "agent_2",
		Message:   "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message is required")

	// Test DistributeTasks with nil request
	_, err = client.AgentCoordination.DistributeTasks(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test DistributeTasks with empty tasks
	_, err = client.AgentCoordination.DistributeTasks(ctx, &DistributeTasksRequest{Tasks: []string{}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one task is required")
}
