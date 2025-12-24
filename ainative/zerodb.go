package ainative

import (
	"context"
	"fmt"
	"time"
)

// ZeroDBService handles ZeroDB operations
type ZeroDBService struct {
	client     *Client
	Projects   *ProjectsService
	Vectors    *VectorsService
	Memory     *MemoryService
	Embeddings *EmbeddingsService
}

// NewZeroDBService creates a new ZeroDB service
func NewZeroDBService(client *Client) *ZeroDBService {
	service := &ZeroDBService{
		client: client,
	}

	service.Projects = &ProjectsService{client: client}
	service.Vectors = &VectorsService{client: client}
	service.Memory = &MemoryService{client: client}
	service.Embeddings = &EmbeddingsService{client: client}

	return service
}

// ProjectsService handles project operations
type ProjectsService struct {
	client *Client
}

// Project represents a ZeroDB project
type Project struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      ProjectStatus          `json:"status"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Owner       *ProjectOwner         `json:"owner,omitempty"`
	Stats       *ProjectStats         `json:"stats,omitempty"`
}

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusSuspended ProjectStatus = "suspended"
	ProjectStatusDeleted   ProjectStatus = "deleted"
)

// ProjectOwner represents project ownership information
type ProjectOwner struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Organization string `json:"organization,omitempty"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
	VectorCount int64  `json:"vector_count"`
	MemoryCount int64  `json:"memory_count"`
	StorageSize int64  `json:"storage_size_bytes"`
	LastAccess  *time.Time `json:"last_access,omitempty"`
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ListProjectsRequest represents a request to list projects
type ListProjectsRequest struct {
	Limit  int           `json:"limit,omitempty"`
	Offset int           `json:"offset,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// ListProjectsResponse represents a response containing projects
type ListProjectsResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int       `json:"total_count"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

// Create creates a new project
func (s *ProjectsService) Create(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if req.Name == "" {
		return nil, NewValidationError("name", "project name is required", req.Name)
	}
	
	var result Project
	
	err := s.client.makeRequest(ctx, "POST", "/api/v1/zerodb/projects", req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// List lists projects with optional filtering
func (s *ProjectsService) List(ctx context.Context, req *ListProjectsRequest) (*ListProjectsResponse, error) {
	if req == nil {
		req = &ListProjectsRequest{}
	}
	
	// Set defaults
	if req.Limit == 0 {
		req.Limit = 10
	}
	
	var result ListProjectsResponse
	
	path := fmt.Sprintf("/api/v1/zerodb/projects?limit=%d&offset=%d", req.Limit, req.Offset)
	if req.Status != "" {
		path += fmt.Sprintf("&status=%s", req.Status)
	}
	
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Get retrieves a specific project
func (s *ProjectsService) Get(ctx context.Context, projectID string) (*Project, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}
	
	var result Project
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s", projectID)
	
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Update updates a project
func (s *ProjectsService) Update(ctx context.Context, projectID string, req *UpdateProjectRequest) (*Project, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}
	
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	var result Project
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s", projectID)
	
	err := s.client.makeRequest(ctx, "PUT", path, req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Suspend suspends a project
func (s *ProjectsService) Suspend(ctx context.Context, projectID string, reason string) error {
	if projectID == "" {
		return NewValidationError("project_id", "project ID is required", projectID)
	}
	
	req := map[string]interface{}{
		"reason": reason,
	}
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s/suspend", projectID)
	
	return s.client.makeRequest(ctx, "POST", path, req, nil)
}

// Activate activates a suspended project
func (s *ProjectsService) Activate(ctx context.Context, projectID string) error {
	if projectID == "" {
		return NewValidationError("project_id", "project ID is required", projectID)
	}
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s/activate", projectID)
	
	return s.client.makeRequest(ctx, "POST", path, nil, nil)
}

// Delete deletes a project
func (s *ProjectsService) Delete(ctx context.Context, projectID string) error {
	if projectID == "" {
		return NewValidationError("project_id", "project ID is required", projectID)
	}
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s", projectID)
	
	return s.client.makeRequest(ctx, "DELETE", path, nil, nil)
}

// VectorsService handles vector operations
type VectorsService struct {
	client *Client
}

// VectorItem represents a vector with metadata
type VectorItem struct {
	ID       string                 `json:"id"`
	Vector   []float64              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// VectorSearchRequest represents a vector search request
type VectorSearchRequest struct {
	Vector    []float64              `json:"vector"`
	TopK      int                    `json:"top_k"`
	Namespace string                 `json:"namespace,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	IncludeMetadata bool              `json:"include_metadata"`
	IncludeValues   bool              `json:"include_values"`
}

// VectorSearchMatch represents a search result match
type VectorSearchMatch struct {
	ID       string                 `json:"id"`
	Score    float64               `json:"score"`
	Vector   []float64              `json:"vector,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// VectorSearchResponse represents a vector search response
type VectorSearchResponse struct {
	Matches   []VectorSearchMatch `json:"matches"`
	Namespace string              `json:"namespace"`
}

// UpsertVectorsRequest represents a request to upsert vectors
type UpsertVectorsRequest struct {
	Vectors   []VectorItem `json:"vectors"`
	Namespace string       `json:"namespace,omitempty"`
}

// UpsertVectorsResponse represents a response from upserting vectors
type UpsertVectorsResponse struct {
	UpsertedCount int    `json:"upserted_count"`
	Namespace     string `json:"namespace"`
}

// Search searches for similar vectors
func (s *VectorsService) Search(ctx context.Context, projectID string, req *VectorSearchRequest) (*VectorSearchResponse, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}
	
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if len(req.Vector) == 0 {
		return nil, NewValidationError("vector", "vector cannot be empty", req.Vector)
	}
	
	if req.TopK <= 0 {
		req.TopK = 5
	}
	
	var result VectorSearchResponse
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s/vectors/search", projectID)
	
	err := s.client.makeRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Upsert upserts vectors into the project
func (s *VectorsService) Upsert(ctx context.Context, projectID string, req *UpsertVectorsRequest) (*UpsertVectorsResponse, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}
	
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if len(req.Vectors) == 0 {
		return nil, NewValidationError("vectors", "vectors cannot be empty", req.Vectors)
	}
	
	var result UpsertVectorsResponse
	
	path := fmt.Sprintf("/api/v1/zerodb/projects/%s/vectors", projectID)
	
	err := s.client.makeRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// MemoryService handles memory operations
type MemoryService struct {
	client *Client
}

// MemoryItem represents a memory item
type MemoryItem struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Title     string                 `json:"title,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Priority  MemoryPriority         `json:"priority"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// MemoryPriority represents memory priority levels
type MemoryPriority string

const (
	MemoryPriorityLow      MemoryPriority = "low"
	MemoryPriorityMedium   MemoryPriority = "medium"
	MemoryPriorityHigh     MemoryPriority = "high"
	MemoryPriorityCritical MemoryPriority = "critical"
)

// CreateMemoryRequest represents a request to create a memory
type CreateMemoryRequest struct {
	Content  string                 `json:"content"`
	Title    string                 `json:"title,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Priority MemoryPriority         `json:"priority,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SearchMemoryRequest represents a memory search request
type SearchMemoryRequest struct {
	Query     string         `json:"query"`
	Limit     int           `json:"limit,omitempty"`
	Tags      []string      `json:"tags,omitempty"`
	Priority  MemoryPriority `json:"priority,omitempty"`
	Semantic  bool          `json:"semantic"`
}

// SearchMemoryResponse represents a memory search response
type SearchMemoryResponse struct {
	Results []MemoryItem `json:"results"`
	Total   int         `json:"total"`
}

// Create creates a new memory
func (s *MemoryService) Create(ctx context.Context, req *CreateMemoryRequest) (*MemoryItem, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if req.Content == "" {
		return nil, NewValidationError("content", "content is required", req.Content)
	}
	
	if req.Priority == "" {
		req.Priority = MemoryPriorityMedium
	}
	
	var result MemoryItem
	
	err := s.client.makeRequest(ctx, "POST", "/api/v1/memory", req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Search searches for memories
func (s *MemoryService) Search(ctx context.Context, req *SearchMemoryRequest) (*SearchMemoryResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}
	
	if req.Query == "" {
		return nil, NewValidationError("query", "query is required", req.Query)
	}
	
	if req.Limit == 0 {
		req.Limit = 10
	}
	
	var result SearchMemoryResponse
	
	err := s.client.makeRequest(ctx, "POST", "/api/v1/memory/search", req, &result)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}