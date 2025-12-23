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

func TestAgentLearningService_SubmitFeedback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-learning/feedback", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req SubmitFeedbackRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "agent_123", req.AgentID)
		assert.Equal(t, "interaction_456", req.InteractionID)
		assert.Equal(t, 4.5, req.Rating)
		assert.Equal(t, "Great performance", req.Comments)

		response := SubmitFeedbackResponse{
			FeedbackID: "feedback_789",
			Status:     "received",
			Message:    "Feedback submitted successfully",
			Timestamp:  time.Now(),
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
	req := &SubmitFeedbackRequest{
		AgentID:       "agent_123",
		InteractionID: "interaction_456",
		Rating:        4.5,
		Comments:      "Great performance",
		Tags:          []string{"helpful", "accurate"},
		Metadata: map[string]interface{}{
			"task_type": "analysis",
		},
	}

	response, err := client.AgentLearning.SubmitFeedback(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "feedback_789", response.FeedbackID)
	assert.Equal(t, "received", response.Status)
	assert.Equal(t, "Feedback submitted successfully", response.Message)
}

func TestAgentLearningService_GetPerformanceMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-learning/metrics", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "agent_123", r.URL.Query().Get("agent_id"))
		assert.Equal(t, "7d", r.URL.Query().Get("period"))

		response := GetPerformanceMetricsResponse{
			Metrics: PerformanceMetrics{
				AgentID:               "agent_123",
				AverageRating:         4.2,
				TotalInteractions:     150,
				SuccessfulTasks:       135,
				FailedTasks:           15,
				SuccessRate:           0.9,
				AverageResponseTime:   250.5,
				AverageCompletionTime: 1500.0,
				LastUpdated:           time.Now(),
			},
			Trends: map[string]interface{}{
				"rating_trend":  "increasing",
				"success_trend": "stable",
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
	req := &GetPerformanceMetricsRequest{
		AgentID: "agent_123",
		Period:  "7d",
	}

	response, err := client.AgentLearning.GetPerformanceMetrics(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "agent_123", response.Metrics.AgentID)
	assert.Equal(t, 4.2, response.Metrics.AverageRating)
	assert.Equal(t, 150, response.Metrics.TotalInteractions)
	assert.Equal(t, 0.9, response.Metrics.SuccessRate)
	assert.NotNil(t, response.Trends)
}

func TestAgentLearningService_CompareAgents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/agent-learning/compare", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CompareAgentsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(req.Agents))
		assert.Equal(t, "success_rate", req.Metric)
		assert.Equal(t, "30d", req.Period)

		response := CompareAgentsResponse{
			Comparisons: []AgentComparison{
				{
					AgentID: "agent_1",
					Metrics: PerformanceMetrics{
						AgentID:     "agent_1",
						SuccessRate: 0.95,
					},
					Ranking: 1,
					Score:   95.0,
				},
				{
					AgentID: "agent_2",
					Metrics: PerformanceMetrics{
						AgentID:     "agent_2",
						SuccessRate: 0.88,
					},
					Ranking: 2,
					Score:   88.0,
				},
				{
					AgentID: "agent_3",
					Metrics: PerformanceMetrics{
						AgentID:     "agent_3",
						SuccessRate: 0.82,
					},
					Ranking: 3,
					Score:   82.0,
				},
			},
			Metric: "success_rate",
			Summary: map[string]interface{}{
				"average": 0.883,
				"median":  0.88,
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
	req := &CompareAgentsRequest{
		Agents: []string{"agent_1", "agent_2", "agent_3"},
		Metric: "success_rate",
		Period: "30d",
	}

	response, err := client.AgentLearning.CompareAgents(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Comparisons))
	assert.Equal(t, "success_rate", response.Metric)
	assert.Equal(t, "agent_1", response.Comparisons[0].AgentID)
	assert.Equal(t, 1, response.Comparisons[0].Ranking)
	assert.Equal(t, 95.0, response.Comparisons[0].Score)
}

func TestAgentLearningService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test SubmitFeedback with nil request
	_, err = client.AgentLearning.SubmitFeedback(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test SubmitFeedback with empty agent_id
	_, err = client.AgentLearning.SubmitFeedback(ctx, &SubmitFeedbackRequest{AgentID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent ID is required")

	// Test SubmitFeedback with invalid rating (too low)
	_, err = client.AgentLearning.SubmitFeedback(ctx, &SubmitFeedbackRequest{
		AgentID: "agent_123",
		Rating:  -1.0,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 0.0 and 5.0")

	// Test SubmitFeedback with invalid rating (too high)
	_, err = client.AgentLearning.SubmitFeedback(ctx, &SubmitFeedbackRequest{
		AgentID: "agent_123",
		Rating:  6.0,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 0.0 and 5.0")

	// Test GetPerformanceMetrics with nil request
	_, err = client.AgentLearning.GetPerformanceMetrics(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test GetPerformanceMetrics with empty agent_id
	_, err = client.AgentLearning.GetPerformanceMetrics(ctx, &GetPerformanceMetricsRequest{AgentID: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent ID is required")

	// Test CompareAgents with nil request
	_, err = client.AgentLearning.CompareAgents(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test CompareAgents with too few agents
	_, err = client.AgentLearning.CompareAgents(ctx, &CompareAgentsRequest{
		Agents: []string{"agent_1"},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least two agents are required for comparison")
}
