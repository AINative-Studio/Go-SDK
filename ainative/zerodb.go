package ainative

import (
	"context"
	"fmt"
)

// ZeroDBService provides access to ZeroDB operations
type ZeroDBService struct {
	client   *Client
	Projects *ProjectsService
	Vectors  *VectorsService
	Memory   *MemoryService
}

// ProjectsService handles project operations
type ProjectsService struct {
	client *Client
}

// VectorsService handles vector operations  
type VectorsService struct {
	client *Client
}

// MemoryService handles memory operations
type MemoryService struct {
	client *Client
}

// CreateProjectRequest represents a project creation request
type CreateProjectRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Project represents a ZeroDB project
type Project struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	CreatedAt   string                 `json:"created_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Create creates a new project
func (p *ProjectsService) Create(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	var project Project
	_, err := p.client.resty.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&project).
		Post("/zerodb/projects")
	return &project, err
}

// Get retrieves a project by ID
func (p *ProjectsService) Get(ctx context.Context, id string) (*Project, error) {
	var project Project
	_, err := p.client.resty.R().
		SetContext(ctx).
		SetResult(&project).
		Get(fmt.Sprintf("/zerodb/projects/%s", id))
	return &project, err
}