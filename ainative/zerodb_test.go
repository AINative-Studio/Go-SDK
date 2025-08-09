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

func TestProjectsService_Create(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var req CreateProjectRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "Test Project", req.Name)
		assert.Equal(t, "Test Description", req.Description)

		// Return mock response
		response := Project{
			ID:          "proj_123",
			Name:        req.Name,
			Description: req.Description,
			Status:      ProjectStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Metadata:    req.Metadata,
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
	req := &CreateProjectRequest{
		Name:        "Test Project",
		Description: "Test Description",
		Metadata: map[string]interface{}{
			"test": true,
		},
	}

	project, err := client.ZeroDB.Projects.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "proj_123", project.ID)
	assert.Equal(t, "Test Project", project.Name)
	assert.Equal(t, ProjectStatusActive, project.Status)
}

func TestProjectsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/proj_123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := Project{
			ID:          "proj_123",
			Name:        "Test Project",
			Description: "Test Description",
			Status:      ProjectStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
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
	project, err := client.ZeroDB.Projects.Get(ctx, "proj_123")

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "proj_123", project.ID)
	assert.Equal(t, "Test Project", project.Name)
}

func TestProjectsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check query parameters
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "0", r.URL.Query().Get("offset"))

		response := ListProjectsResponse{
			Projects: []Project{
				{
					ID:     "proj_1",
					Name:   "Project 1",
					Status: ProjectStatusActive,
				},
				{
					ID:     "proj_2",
					Name:   "Project 2",
					Status: ProjectStatusSuspended,
				},
			},
			TotalCount:  2,
			Limit:  10,
			Offset: 0,
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
	req := &ListProjectsRequest{
		Limit:  10,
		Offset: 0,
	}

	response, err := client.ZeroDB.Projects.List(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Projects))
	assert.Equal(t, 2, response.TotalCount)
	assert.Equal(t, "proj_1", response.Projects[0].ID)
}

func TestVectorsService_Upsert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/proj_123/vectors", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		var req UpsertVectorsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(req.Vectors))
		assert.Equal(t, "default", req.Namespace)

		response := UpsertVectorsResponse{
			UpsertedCount: 2,
			Namespace:     "default",
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
	req := &UpsertVectorsRequest{
		Vectors: []VectorItem{
			{
				ID:     "vec_1",
				Vector: []float64{0.1, 0.2, 0.3},
				Metadata: map[string]interface{}{
					"category": "test",
				},
			},
			{
				ID:     "vec_2",
				Vector: []float64{0.4, 0.5, 0.6},
			},
		},
		Namespace: "default",
	}

	response, err := client.ZeroDB.Vectors.Upsert(ctx, "proj_123", req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.UpsertedCount)
	assert.Equal(t, "default", response.Namespace)
}

func TestVectorsService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/projects/proj_123/vectors/search", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req VectorSearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(req.Vector))
		assert.Equal(t, 5, req.TopK)

		response := VectorSearchResponse{
			Matches: []VectorSearchMatch{
				{
					ID:    "vec_1",
					Score: 0.95,
					Metadata: map[string]interface{}{
						"category": "test",
					},
					Vector: []float64{0.1, 0.2, 0.3},
				},
				{
					ID:    "vec_2",
					Score: 0.87,
				},
			},
			Namespace: "default",
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
	req := &VectorSearchRequest{
		Vector:          []float64{0.1, 0.2, 0.3},
		TopK:            5,
		Namespace:       "default",
		IncludeMetadata: true,
		IncludeValues:   true,
	}

	response, err := client.ZeroDB.Vectors.Search(ctx, "proj_123", req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Matches))
	assert.Equal(t, "vec_1", response.Matches[0].ID)
	assert.Equal(t, 0.95, response.Matches[0].Score)
}

func TestMemoryService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/memories", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateMemoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "Test Memory", req.Title)
		assert.Equal(t, "Test content", req.Content)

		response := MemoryItem{
			ID:        "mem_123",
			Title:     req.Title,
			Content:   req.Content,
			Tags:      req.Tags,
			Priority:  req.Priority,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
	req := &CreateMemoryRequest{
		Title:    "Test Memory",
		Content:  "Test content",
		Tags:     []string{"test", "memory"},
		Priority: MemoryPriorityHigh,
	}

	memory, err := client.ZeroDB.Memory.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, memory)
	assert.Equal(t, "mem_123", memory.ID)
	assert.Equal(t, "Test Memory", memory.Title)
	assert.Equal(t, MemoryPriorityHigh, memory.Priority)
}

func TestMemoryService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/memories/search", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req SearchMemoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "test query", req.Query)
		assert.Equal(t, 10, req.Limit)

		response := SearchMemoryResponse{
			Results: []MemoryItem{
				{
					ID:       "mem_1",
					Title:    "Memory 1",
					Content:  "Test content 1",
					Priority: MemoryPriorityMedium,
				},
				{
					ID:       "mem_2",
					Title:    "Memory 2",
					Content:  "Test content 2",
					Priority: MemoryPriorityLow,
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
	req := &SearchMemoryRequest{
		Query:    "test query",
		Limit:    10,
		Semantic: true,
	}

	response, err := client.ZeroDB.Memory.Search(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Results))
	assert.Equal(t, 2, response.Total)
	assert.Equal(t, "mem_1", response.Results[0].ID)
}

func TestProjectsService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test Create with nil request
	_, err = client.ZeroDB.Projects.Create(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Create with empty name
	_, err = client.ZeroDB.Projects.Create(ctx, &CreateProjectRequest{Name: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test Get with empty project ID
	_, err = client.ZeroDB.Projects.Get(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project ID is required")
}

func TestVectorsService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test Upsert with empty project ID
	_, err = client.ZeroDB.Vectors.Upsert(ctx, "", &UpsertVectorsRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project ID is required")

	// Test Upsert with nil request
	_, err = client.ZeroDB.Vectors.Upsert(ctx, "proj_123", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Search with empty vector
	_, err = client.ZeroDB.Vectors.Search(ctx, "proj_123", &VectorSearchRequest{Vector: []float64{}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vector cannot be empty")

	// Test Search with invalid TopK
	_, err = client.ZeroDB.Vectors.Search(ctx, "proj_123", &VectorSearchRequest{
		Vector: []float64{0.1, 0.2, 0.3},
		TopK:   0,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "topK must be greater than 0")
}

func TestMemoryService_Validation(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test Create with nil request
	_, err = client.ZeroDB.Memory.Create(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Create with empty content
	_, err = client.ZeroDB.Memory.Create(ctx, &CreateMemoryRequest{Content: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "content is required")

	// Test Search with nil request
	_, err = client.ZeroDB.Memory.Search(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Search with empty query
	_, err = client.ZeroDB.Memory.Search(ctx, &SearchMemoryRequest{Query: ""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query is required")
}