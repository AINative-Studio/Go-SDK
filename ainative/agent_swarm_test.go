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

func TestAgentSwarmService_Start(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req StartSwarmRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "proj_123", req.ProjectID)
		assert.Equal(t, "Test Swarm", req.Name)
		assert.Equal(t, 2, len(req.Agents))

		response := AgentSwarm{
			ID:        "swarm_123",
			Name:      req.Name,
			ProjectID: req.ProjectID,
			Objective: req.Objective,
			Status:    SwarmStatusRunning,
			Agents: []Agent{
				{
					ID:           "agent_1",
					Type:         AgentTypeAnalyzer,
					Status:       AgentStatusActive,
					Capabilities: []string{"code_analysis"},
				},
				{
					ID:           "agent_2",
					Type:         AgentTypeOptimizer,
					Status:       AgentStatusActive,
					Capabilities: []string{"performance_analysis"},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &StartSwarmRequest{
		ProjectID: "proj_123",
		Name:      "Test Swarm",
		Objective: "Test objective",
		Agents: []AgentConfig{
			{
				Type:         AgentTypeAnalyzer,
				Count:        1,
				Capabilities: []string{"code_analysis"},
			},
			{
				Type:         AgentTypeOptimizer,
				Count:        1,
				Capabilities: []string{"performance_analysis"},
			},
		},
	}

	swarm, err := client.AgentSwarm.Start(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, swarm)
	assert.Equal(t, "swarm_123", swarm.ID)
	assert.Equal(t, "Test Swarm", swarm.Name)
	assert.Equal(t, SwarmStatusRunning, swarm.Status)
	assert.Equal(t, 2, len(swarm.Agents))
}

func TestAgentSwarmService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms/swarm_123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := AgentSwarm{
			ID:        "swarm_123",
			Name:      "Test Swarm",
			ProjectID: "proj_123",
			Status:    SwarmStatusRunning,
			Metrics: &SwarmMetrics{
				TasksCompleted:     5,
				TasksInProgress:    2,
				TasksFailed:        1,
				Efficiency:         0.85,
				AverageTaskTime:    30 * time.Second,
				TotalExecutionTime: 5 * time.Minute,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	swarm, err := client.AgentSwarm.Get(ctx, "swarm_123")

	assert.NoError(t, err)
	assert.NotNil(t, swarm)
	assert.Equal(t, "swarm_123", swarm.ID)
	assert.Equal(t, SwarmStatusRunning, swarm.Status)
	assert.NotNil(t, swarm.Metrics)
	assert.Equal(t, 5, swarm.Metrics.TasksCompleted)
}

func TestAgentSwarmService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "proj_123", r.URL.Query().Get("project_id"))

		response := ListSwarmsResponse{
			Swarms: []AgentSwarm{
				{
					ID:        "swarm_1",
					Name:      "Swarm 1",
					ProjectID: "proj_123",
					Status:    SwarmStatusRunning,
				},
				{
					ID:        "swarm_2",
					Name:      "Swarm 2",
					ProjectID: "proj_123",
					Status:    SwarmStatusPaused,
				},
			},
			TotalCount: 2,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &ListSwarmsRequest{
		ProjectID: "proj_123",
		Limit:     10,
	}

	response, err := client.AgentSwarm.List(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Swarms))
	assert.Equal(t, 2, response.TotalCount)
	assert.Equal(t, "swarm_1", response.Swarms[0].ID)
}

func TestAgentSwarmService_Orchestrate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/orchestrate", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req OrchestrationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "analyze_code", req.Task)
		assert.Equal(t, TaskPriorityHigh, req.Priority)

		response := OrchestrationResponse{
			TaskID:     "task_123",
			Status:     TaskStatusAssigned,
			AssignedTo: []string{"agent_1", "agent_2"},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &OrchestrationRequest{
		SwarmID:  "swarm_123",
		Task:     "analyze_code",
		Priority: TaskPriorityHigh,
		Context: map[string]interface{}{
			"repository": "test-repo",
		},
	}

	response, err := client.AgentSwarm.Orchestrate(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "task_123", response.TaskID)
	assert.Equal(t, TaskStatusAssigned, response.Status)
	assert.Equal(t, 2, len(response.AssignedTo))
}

func TestAgentSwarmService_GetTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/tasks/task_123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		now := time.Now()
		response := OrchestrationTask{
			ID:          "task_123",
			SwarmID:     "swarm_123",
			Status:      TaskStatusCompleted,
			CreatedAt:   now.Add(-10 * time.Minute),
			StartedAt:   &[]time.Time{now.Add(-8 * time.Minute)}[0],
			CompletedAt: &now,
			Result: map[string]interface{}{
				"status":  "success",
				"details": "Analysis completed",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	task, err := client.AgentSwarm.GetTask(ctx, "task_123")

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "task_123", task.ID)
	assert.Equal(t, TaskStatusCompleted, task.Status)
	assert.NotNil(t, task.StartedAt)
	assert.NotNil(t, task.CompletedAt)
}

// agentTypesBackendFixture mirrors, byte-for-byte in shape, what the real backend
// returns from GET /api/v1/admin/agent-swarm/agent-types
// (src/backend/app/api/admin/agent_swarm.py::get_available_agent_types, response_model
// AvailableAgentTypesResponse -> {"agent_types": [AgentTypeEntry, ...]}).
// Verified against the live deployment's /openapi.json. Refs #8478.
const agentTypesBackendFixture = `{
  "agent_types": [
    {
      "name": "architect",
      "description": "System architecture and database design",
      "capabilities": ["system_design", "database_schema", "api_specification"],
      "models": ["gpt-4", "claude-3", "gpt-3.5"],
      "avg_completion_time": 5.2
    },
    {
      "name": "developer",
      "description": "Code generation and implementation",
      "capabilities": ["frontend_development", "backend_api", "database_integration"],
      "models": ["claude-3", "gpt-4", "codex"],
      "avg_completion_time": 8.7
    },
    {
      "name": "tester",
      "description": "Test generation and quality assurance",
      "capabilities": ["unit_testing", "integration_testing", "test_automation"],
      "models": ["gpt-3.5", "claude-3", "gpt-4"],
      "avg_completion_time": 3.1
    }
  ]
}`

func TestAgentSwarmService_ListAgentTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The real route is mounted by app/api/routers/admin.py at prefix
		// "/agent-swarm" under admin_router, which main.py mounts at /admin and
		// api_v1/api.py also mounts at /api/v1/admin. The SDK targets the
		// versioned mount.
		assert.Equal(t, "/api/v1/admin/agent-swarm/agent-types", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(agentTypesBackendFixture))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	agentTypes, err := client.AgentSwarm.ListAgentTypes(ctx)

	assert.NoError(t, err)
	require.Len(t, agentTypes, 3)

	// Every field the backend sends must survive decoding -- the old
	// []AgentType (string) return type silently dropped all of them.
	assert.Equal(t, AgentType("architect"), agentTypes[0].Name)
	assert.Equal(t, "System architecture and database design", agentTypes[0].Description)
	assert.Equal(t, []string{"system_design", "database_schema", "api_specification"}, agentTypes[0].Capabilities)
	assert.Equal(t, []string{"gpt-4", "claude-3", "gpt-3.5"}, agentTypes[0].Models)
	assert.InDelta(t, 5.2, agentTypes[0].AvgCompletionTime, 1e-9)

	assert.Equal(t, AgentType("developer"), agentTypes[1].Name)
	assert.Equal(t, []string{"claude-3", "gpt-4", "codex"}, agentTypes[1].Models)
	assert.InDelta(t, 8.7, agentTypes[1].AvgCompletionTime, 1e-9)

	assert.Equal(t, AgentType("tester"), agentTypes[2].Name)
	assert.InDelta(t, 3.1, agentTypes[2].AvgCompletionTime, 1e-9)
}

// TestAgentSwarmService_ListAgentTypeNames covers the name-only convenience
// helper, which is what the pre-#8478 []AgentType return value was trying
// (and failing) to be.
func TestAgentSwarmService_ListAgentTypeNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/admin/agent-swarm/agent-types", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(agentTypesBackendFixture))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	names, err := client.AgentSwarm.ListAgentTypeNames(ctx)

	assert.NoError(t, err)
	assert.Equal(t, []AgentType{"architect", "developer", "tester"}, names)
}

// TestAgentSwarmService_ListAgentTypes_RejectsLegacyStringArray pins the bug
// that #8478 fixed: the backend never returned a bare []string, and the SDK
// must not silently decode one into an empty result.
func TestAgentSwarmService_ListAgentTypes_RejectsLegacyStringArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]string{"analyzer", "optimizer"})
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	_, err = client.AgentSwarm.ListAgentTypes(context.Background())
	assert.Error(t, err, "a bare JSON array is not the documented envelope and must not decode silently")
}

func TestAgentSwarmService_GetSwarmMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms/swarm_123/metrics", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := SwarmMetrics{
			TasksCompleted:     10,
			TasksInProgress:    2,
			TasksFailed:        1,
			Efficiency:         0.92,
			AverageTaskTime:    45 * time.Second,
			TotalExecutionTime: 15 * time.Minute,
			ResourceUsage: &ResourceUsage{
				CPUUsage:        0.65,
				MemoryUsage:     2 * 1024 * 1024 * 1024, // 2GB
				APICallsCounter: 150,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	metrics, err := client.AgentSwarm.GetSwarmMetrics(ctx, "swarm_123")

	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, 10, metrics.TasksCompleted)
	assert.Equal(t, 0.92, metrics.Efficiency)
	assert.NotNil(t, metrics.ResourceUsage)
	assert.Equal(t, 0.65, metrics.ResourceUsage.CPUUsage)
}

func TestAgentSwarmService_Pause(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms/swarm_123/pause", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "paused"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.AgentSwarm.Pause(ctx, "swarm_123")

	assert.NoError(t, err)
}

func TestAgentSwarmService_Resume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms/swarm_123/resume", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "running"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.AgentSwarm.Resume(ctx, "swarm_123")

	assert.NoError(t, err)
}

func TestAgentSwarmService_Stop(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-swarm/swarms/swarm_123/stop", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "stopped"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.AgentSwarm.Stop(ctx, "swarm_123")

	assert.NoError(t, err)
}

func TestAgentSwarmService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test Start with nil request
	_, err = client.AgentSwarm.Start(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Start with empty project ID
	_, err = client.AgentSwarm.Start(ctx, &StartSwarmRequest{ProjectID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project ID is required")

	// Test Start with empty objective
	_, err = client.AgentSwarm.Start(ctx, &StartSwarmRequest{
		ProjectID: "proj_123",
		Name:      "Test Swarm",
		Objective: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "objective is required")

	// Test Start with no agents
	_, err = client.AgentSwarm.Start(ctx, &StartSwarmRequest{
		ProjectID: "proj_123",
		Name:      "Test Swarm",
		Objective: "Test objective",
		Agents:    []AgentConfig{},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one agent configuration is required")

	// Test Get with empty swarm ID
	_, err = client.AgentSwarm.Get(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "swarm ID is required")

	// Test Orchestrate with empty swarm ID
	_, err = client.AgentSwarm.Orchestrate(ctx, &OrchestrationRequest{SwarmID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "swarm ID is required")

	// Test GetTask with empty task ID
	_, err = client.AgentSwarm.GetTask(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task ID is required")
}
