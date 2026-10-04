package ainative

import (
	"context"
	"fmt"
	"net/url"
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
	CreatedAt   time.Time              `json:"created_at"`
	// UpdatedAt uses NullableTime: the live API returns "" (not a valid
	// RFC3339 timestamp, not omitted/null) for a project that's never been
	// updated since creation, which fails to unmarshal into time.Time or
	// *time.Time via the standard decoder. Refs #8393.
	UpdatedAt NullableTime  `json:"updated_at"`
	Owner     *ProjectOwner `json:"owner,omitempty"`
	Stats     *ProjectStats `json:"stats,omitempty"`
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
	VectorCount int64      `json:"vector_count"`
	MemoryCount int64      `json:"memory_count"`
	StorageSize int64      `json:"storage_size_bytes"`
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

	// Single-project GET lives at /api/v1/projects/{id}, NOT under the
	// /api/v1/zerodb/projects prefix used for List/Create. Confirmed live
	// against production (the zerodb-prefixed path 405s). Refs #8393.
	path := fmt.Sprintf("/api/v1/projects/%s", projectID)

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

// Delete deletes a project. Lives at /api/v1/projects/{id}, same prefix as
// Get (not the /api/v1/zerodb/projects prefix used for List/Create) —
// confirmed live against production. Refs #8393.
func (s *ProjectsService) Delete(ctx context.Context, projectID string) error {
	if projectID == "" {
		return NewValidationError("project_id", "project ID is required", projectID)
	}

	path := fmt.Sprintf("/api/v1/projects/%s", projectID)

	return s.client.makeRequest(ctx, "DELETE", path, nil, nil)
}

// VectorsService handles vector operations
type VectorsService struct {
	client *Client
}

// VectorItem represents a vector to upsert, matching the live API's
// VectorCreate schema (app/zerodb/schemas/database.py). Refs #8393.
type VectorItem struct {
	VectorID        string                 `json:"vector_id,omitempty"`
	VectorEmbedding []float64              `json:"vector_embedding"`
	Namespace       string                 `json:"namespace,omitempty"`
	Metadata        map[string]interface{} `json:"vector_metadata,omitempty"`
	Document        string                 `json:"document,omitempty"`
	Source          string                 `json:"source,omitempty"`
}

// VectorSearchRequest represents a vector search request, matching the
// live API's VectorSearchRequest schema — the field is query_vector, not
// vector, and the live route is POST /database/vectors/search (not the
// zerodb-prefixed path this SDK used to call). Refs #8393.
type VectorSearchRequest struct {
	QueryVector    []float64              `json:"query_vector"`
	Namespace      string                 `json:"namespace,omitempty"`
	Limit          int                    `json:"limit,omitempty"`
	Threshold      *float64               `json:"threshold,omitempty"`
	MetadataFilter map[string]interface{} `json:"metadata_filter,omitempty"`
}

// VectorMatch represents a single vector in a search result, matching the
// live API's VectorResponse schema.
type VectorMatch struct {
	VectorID        string                 `json:"vector_id"`
	VectorEmbedding []float64              `json:"vector_embedding,omitempty"`
	Namespace       string                 `json:"namespace"`
	Metadata        map[string]interface{} `json:"vector_metadata,omitempty"`
	Document        string                 `json:"document,omitempty"`
	Source          string                 `json:"source,omitempty"`
	Similarity      float64                `json:"similarity,omitempty"`
}

// VectorSearchResponse represents a vector search response, matching the
// live API's VectorSearchResult schema. Refs #8393.
type VectorSearchResponse struct {
	Vectors      []VectorMatch `json:"vectors"`
	TotalCount   int           `json:"total_count"`
	SearchTimeMs float64       `json:"search_time_ms"`
}

// UpsertVectorsResponse represents the live API's BulkOperationResult
// response from POST /database/vectors/upsert-batch. Refs #8393.
type UpsertVectorsResponse struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors,omitempty"`
	TotalTimeMs  float64  `json:"total_time_ms"`
}

// Search searches for similar vectors via
// POST /api/v1/projects/{id}/database/vectors/search. Refs #8393.
func (s *VectorsService) Search(ctx context.Context, projectID string, req *VectorSearchRequest) (*VectorSearchResponse, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}

	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if len(req.QueryVector) == 0 {
		return nil, NewValidationError("query_vector", "query_vector cannot be empty", req.QueryVector)
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	var result VectorSearchResponse

	path := fmt.Sprintf("/api/v1/projects/%s/database/vectors/search", projectID)

	err := s.client.makeRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Upsert batch-upserts vectors into the project via
// POST /api/v1/projects/{id}/database/vectors/upsert-batch. Refs #8393.
func (s *VectorsService) Upsert(ctx context.Context, projectID string, vectors []VectorItem) (*UpsertVectorsResponse, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}

	if len(vectors) == 0 {
		return nil, NewValidationError("vectors", "vectors cannot be empty", vectors)
	}

	var result UpsertVectorsResponse

	path := fmt.Sprintf("/api/v1/projects/%s/database/vectors/upsert-batch", projectID)

	err := s.client.makeRequest(ctx, "POST", path, vectors, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// MemoryService handles memory operations
type MemoryService struct {
	client *Client
}

// MemoryItem represents a single memory returned by Search. The live
// /api/v1/public/memory/simple/search endpoint returns a loosely-typed
// dict per result (id, content, role, timestamp, metadata, relevance_score)
// rather than the fixed tags/priority shape CreateMemoryRequest aspired to
// (that shape is not what the API accepts or returns). Refs #8393.
type MemoryItem struct {
	ID             string                 `json:"id"`
	Content        string                 `json:"content"`
	Role           string                 `json:"role,omitempty"`
	Timestamp      string                 `json:"timestamp,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	RelevanceScore float64                `json:"relevance_score,omitempty"`
}

// MemoryPriority represents memory priority levels
type MemoryPriority string

const (
	MemoryPriorityLow      MemoryPriority = "low"
	MemoryPriorityMedium   MemoryPriority = "medium"
	MemoryPriorityHigh     MemoryPriority = "high"
	MemoryPriorityCritical MemoryPriority = "critical"
)

// CreateMemoryRequest represents a request to create a memory via the
// Simple Memory API. The live endpoint only accepts content+tags (no
// title/priority/metadata) — Refs #8393.
type CreateMemoryRequest struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
}

// CreateMemoryResponse is the live response shape from POST .../simple/add.
type CreateMemoryResponse struct {
	MemoryID string `json:"memory_id"`
	Status   string `json:"status"`
}

// SearchMemoryResponse represents a memory search response, matching the
// live GET .../simple/search response shape. Refs #8393.
type SearchMemoryResponse struct {
	Results []MemoryItem `json:"results"`
	Count   int          `json:"count"`
}

// Create stores a memory via POST /api/v1/public/memory/simple/add.
func (s *MemoryService) Create(ctx context.Context, req *CreateMemoryRequest) (*CreateMemoryResponse, error) {
	if req == nil {
		return nil, NewValidationError("request", "request cannot be nil", nil)
	}

	if req.Content == "" {
		return nil, NewValidationError("content", "content is required", req.Content)
	}

	var result CreateMemoryResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/public/memory/simple/add", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Search searches memories via GET /api/v1/public/memory/simple/search
// (a query-param GET on the live API, not a JSON POST body). Refs #8393.
func (s *MemoryService) Search(ctx context.Context, query string, limit int) (*SearchMemoryResponse, error) {
	if query == "" {
		return nil, NewValidationError("query", "query is required", query)
	}

	if limit == 0 {
		limit = 5
	}

	var result SearchMemoryResponse

	path := fmt.Sprintf("/api/v1/public/memory/simple/search?q=%s&limit=%d", url.QueryEscape(query), limit)
	err := s.client.makeRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
