package ainative

import (
	"context"
	"fmt"
	"time"
)

// AgentLearningService handles agent learning and feedback operations
type AgentLearningService struct {
	client *Client
}

// NewAgentLearningService creates a new agent learning service
func NewAgentLearningService(client *Client) *AgentLearningService {
	return &AgentLearningService{
		client: client,
	}
}

// Feedback represents feedback for an agent interaction
type Feedback struct {
	ID            string                 `json:"id"`
	AgentID       string                 `json:"agent_id"`
	InteractionID string                 `json:"interaction_id,omitempty"`
	Rating        float64                `json:"rating"` // 0.0 to 5.0
	Comments      string                 `json:"comments,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
}

// SubmitFeedbackRequest represents a request to submit feedback
type SubmitFeedbackRequest struct {
	AgentID       string                 `json:"agent_id"`
	InteractionID string                 `json:"interaction_id,omitempty"`
	Rating        float64                `json:"rating"`
	Comments      string                 `json:"comments,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// SubmitFeedbackResponse represents the response from submitting feedback
type SubmitFeedbackResponse struct {
	FeedbackID string    `json:"feedback_id"`
	Status     string    `json:"status"`
	Message    string    `json:"message,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// PerformanceMetrics represents performance metrics for an agent
type PerformanceMetrics struct {
	AgentID              string    `json:"agent_id"`
	AverageRating        float64   `json:"avg_rating"`
	TotalInteractions    int       `json:"total_interactions"`
	SuccessfulTasks      int       `json:"successful_tasks"`
	FailedTasks          int       `json:"failed_tasks"`
	SuccessRate          float64   `json:"success_rate"` // 0.0 to 1.0
	AverageResponseTime  float64   `json:"avg_response_time_ms"`
	AverageCompletionTime float64   `json:"avg_completion_time_ms"`
	LastUpdated          time.Time `json:"last_updated"`
}

// GetPerformanceMetricsRequest represents a request for performance metrics
type GetPerformanceMetricsRequest struct {
	AgentID string `json:"agent_id"`
	Period  string `json:"period,omitempty"` // 24h, 7d, 30d, 90d
}

// GetPerformanceMetricsResponse represents the response containing performance metrics
type GetPerformanceMetricsResponse struct {
	Metrics   PerformanceMetrics     `json:"metrics"`
	Trends    map[string]interface{} `json:"trends,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// AgentComparison represents comparison data between agents
type AgentComparison struct {
	AgentID  string             `json:"agent_id"`
	Metrics  PerformanceMetrics `json:"metrics"`
	Ranking  int                `json:"ranking"`
	Score    float64            `json:"score"`
}

// CompareAgentsRequest represents a request to compare multiple agents
type CompareAgentsRequest struct {
	Agents     []string `json:"agents"`
	Metric     string   `json:"metric,omitempty"` // success_rate, avg_rating, response_time
	Period     string   `json:"period,omitempty"` // 24h, 7d, 30d
}

// CompareAgentsResponse represents the response from agent comparison
type CompareAgentsResponse struct {
	Comparisons   []AgentComparison      `json:"comparisons"`
	Metric        string                 `json:"metric"`
	Summary       map[string]interface{} `json:"summary,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

// SubmitFeedback submits feedback for an agent interaction
func (s *AgentLearningService) SubmitFeedback(ctx context.Context, req *SubmitFeedbackRequest) (*SubmitFeedbackResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.AgentID == "" {
		return nil, NewValidationError("agent_id", "agent ID is required", req.AgentID)
	}

	if req.Rating < 0.0 || req.Rating > 5.0 {
		return nil, NewValidationError("rating", "rating must be between 0.0 and 5.0", req.Rating)
	}

	var result SubmitFeedbackResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-learning/feedback", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPerformanceMetrics retrieves performance metrics for an agent
func (s *AgentLearningService) GetPerformanceMetrics(ctx context.Context, req *GetPerformanceMetricsRequest) (*GetPerformanceMetricsResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.AgentID == "" {
		return nil, NewValidationError("agent_id", "agent ID is required", req.AgentID)
	}

	path := fmt.Sprintf("/api/v1/agent-learning/metrics?agent_id=%s", req.AgentID)

	if req.Period != "" {
		path += fmt.Sprintf("&period=%s", req.Period)
	}

	var result GetPerformanceMetricsResponse

	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CompareAgents compares performance metrics across multiple agents
func (s *AgentLearningService) CompareAgents(ctx context.Context, req *CompareAgentsRequest) (*CompareAgentsResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if len(req.Agents) < 2 {
		return nil, NewValidationError("agents", "at least two agents are required for comparison", req.Agents)
	}

	// Set default metric if not provided
	if req.Metric == "" {
		req.Metric = "success_rate"
	}

	var result CompareAgentsResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/agent-learning/compare", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
