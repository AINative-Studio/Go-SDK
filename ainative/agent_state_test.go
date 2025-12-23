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

func TestAgentStateService_GetState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-state/state", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "agent_123", r.URL.Query().Get("agent_id"))
		assert.Equal(t, "5", r.URL.Query().Get("version"))

		response := GetStateResponse{
			AgentID: "agent_123",
			Version: 5,
			State: map[string]interface{}{
				"current_task": "analyzing_code",
				"progress":     0.75,
				"status":       "active",
			},
			Timestamp: time.Now(),
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
	req := &GetStateRequest{
		AgentID: "agent_123",
		Version: 5,
	}

	response, err := client.AgentState.GetState(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "agent_123", response.AgentID)
	assert.Equal(t, 5, response.Version)
	assert.NotNil(t, response.State)
	assert.Equal(t, "analyzing_code", response.State["current_task"])
	assert.Equal(t, 0.75, response.State["progress"])
}

func TestAgentStateService_CreateCheckpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-state/checkpoints", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateCheckpointRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "agent_123", req.AgentID)
		assert.Equal(t, "checkpoint_before_deploy", req.Name)
		assert.Equal(t, "State before production deployment", req.Description)
		assert.NotNil(t, req.Data)

		response := CreateCheckpointResponse{
			CheckpointID: "checkpoint_456",
			Version:      6,
			Message:      "Checkpoint created successfully",
			Timestamp:    time.Now(),
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
	req := &CreateCheckpointRequest{
		AgentID:     "agent_123",
		Name:        "checkpoint_before_deploy",
		Description: "State before production deployment",
		Data: map[string]interface{}{
			"task_queue":   []string{"task1", "task2"},
			"configuration": map[string]interface{}{
				"timeout": 300,
			},
		},
	}

	response, err := client.AgentState.CreateCheckpoint(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "checkpoint_456", response.CheckpointID)
	assert.Equal(t, 6, response.Version)
	assert.Equal(t, "Checkpoint created successfully", response.Message)
}

func TestAgentStateService_RestoreCheckpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-state/restore", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req RestoreCheckpointRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "checkpoint_456", req.CheckpointID)
		assert.Equal(t, "agent_123", req.AgentID)

		response := RestoreCheckpointResponse{
			AgentID: "agent_123",
			Version: 5,
			State: map[string]interface{}{
				"task_queue":   []string{"task1", "task2"},
				"configuration": map[string]interface{}{
					"timeout": 300,
				},
			},
			Status:    "restored",
			Timestamp: time.Now(),
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
	req := &RestoreCheckpointRequest{
		CheckpointID: "checkpoint_456",
		AgentID:      "agent_123",
	}

	response, err := client.AgentState.RestoreCheckpoint(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "agent_123", response.AgentID)
	assert.Equal(t, 5, response.Version)
	assert.Equal(t, "restored", response.Status)
	assert.NotNil(t, response.State)
}

func TestAgentStateService_ListCheckpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-state/checkpoints", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "0", r.URL.Query().Get("offset"))
		assert.Equal(t, "agent_123", r.URL.Query().Get("agent_id"))

		response := ListCheckpointsResponse{
			Checkpoints: []Checkpoint{
				{
					ID:          "checkpoint_1",
					AgentID:     "agent_123",
					Name:        "before_deploy",
					Description: "State before deployment",
					Data: map[string]interface{}{
						"tasks": 5,
					},
					Version:   3,
					CreatedAt: time.Now().Add(-2 * time.Hour),
				},
				{
					ID:          "checkpoint_2",
					AgentID:     "agent_123",
					Name:        "after_training",
					Description: "State after ML training",
					Data: map[string]interface{}{
						"model_accuracy": 0.95,
					},
					Version:   2,
					CreatedAt: time.Now().Add(-24 * time.Hour),
				},
			},
			TotalCount: 2,
			Limit:      10,
			Offset:     0,
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
	req := &ListCheckpointsRequest{
		AgentID: "agent_123",
		Limit:   10,
		Offset:  0,
	}

	response, err := client.AgentState.ListCheckpoints(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Checkpoints))
	assert.Equal(t, 2, response.TotalCount)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, "checkpoint_1", response.Checkpoints[0].ID)
	assert.Equal(t, "before_deploy", response.Checkpoints[0].Name)
}

func TestAgentStateService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test GetState with nil request
	_, err = client.AgentState.GetState(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test GetState with empty agent_id
	_, err = client.AgentState.GetState(ctx, &GetStateRequest{AgentID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent ID is required")

	// Test CreateCheckpoint with nil request
	_, err = client.AgentState.CreateCheckpoint(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test CreateCheckpoint with empty agent_id
	_, err = client.AgentState.CreateCheckpoint(ctx, &CreateCheckpointRequest{AgentID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent ID is required")

	// Test CreateCheckpoint with empty name
	_, err = client.AgentState.CreateCheckpoint(ctx, &CreateCheckpointRequest{
		AgentID: "agent_123",
		Name:    "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checkpoint name is required")

	// Test CreateCheckpoint with nil data
	_, err = client.AgentState.CreateCheckpoint(ctx, &CreateCheckpointRequest{
		AgentID: "agent_123",
		Name:    "test_checkpoint",
		Data:    nil,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checkpoint data is required")

	// Test RestoreCheckpoint with nil request
	_, err = client.AgentState.RestoreCheckpoint(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test RestoreCheckpoint with empty checkpoint_id
	_, err = client.AgentState.RestoreCheckpoint(ctx, &RestoreCheckpointRequest{CheckpointID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checkpoint ID is required")

	// Test RestoreCheckpoint with empty agent_id
	_, err = client.AgentState.RestoreCheckpoint(ctx, &RestoreCheckpointRequest{
		CheckpointID: "checkpoint_456",
		AgentID:      "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent ID is required")
}
