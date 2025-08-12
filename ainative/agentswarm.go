package ainative

import "context"

// AgentSwarmService provides access to Agent Swarm operations
type AgentSwarmService struct {
	client *Client
}

// AgentType represents the type of agent
type AgentType string

const (
	AgentTypeResearcher AgentType = "researcher"
	AgentTypeCoder      AgentType = "coder"
	AgentTypeAnalyst    AgentType = "analyst"
)

// AgentConfig represents agent configuration
type AgentConfig struct {
	ID           string      `json:"id"`
	Type         AgentType   `json:"type"`
	Capabilities []string    `json:"capabilities"`
	Config       interface{} `json:"config,omitempty"`
}

// StartSwarmRequest represents a swarm start request
type StartSwarmRequest struct {
	ProjectID string        `json:"project_id"`
	Agents    []AgentConfig `json:"agents"`
	Objective string        `json:"objective"`
}

// Swarm represents an agent swarm
type Swarm struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// StartSwarm starts a new agent swarm
func (a *AgentSwarmService) StartSwarm(ctx context.Context, req *StartSwarmRequest) (*Swarm, error) {
	var swarm Swarm
	_, err := a.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&swarm).
		Post("/agent-swarm/swarms")
	return &swarm, err
}

// GetStatus gets the status of a swarm
func (a *AgentSwarmService) GetStatus(ctx context.Context, swarmID string) (*Swarm, error) {
	var swarm Swarm
	_, err := a.client.resty.R().
		SetContext(ctx).
		SetResult(&swarm).
		Get("/agent-swarm/swarms/" + swarmID)
	return &swarm, err
}