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
		assert.Equal(t, "/api/v1/zerodb/projects", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var req CreateProjectRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "Test Project", req.Name)
		assert.Equal(t, "Test Description", req.Description)

		// Return mock response
		now := time.Now()
		response := Project{
			ID:          "proj_123",
			Name:        req.Name,
			Description: req.Description,
			Status:      ProjectStatusActive,
			CreatedAt:   now,
			UpdatedAt:   NullableTime{Time: now, Valid: true},
			Metadata:    req.Metadata,
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

		now := time.Now()
		response := Project{
			ID:          "proj_123",
			Name:        "Test Project",
			Description: "Test Description",
			Status:      ProjectStatusActive,
			CreatedAt:   now,
			UpdatedAt:   NullableTime{Time: now, Valid: true},
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
	project, err := client.ZeroDB.Projects.Get(ctx, "proj_123")

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "proj_123", project.ID)
	assert.Equal(t, "Test Project", project.Name)
}

func TestProjectsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/zerodb/projects", r.URL.Path)
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
			TotalCount: 2,
			Limit:      10,
			Offset:     0,
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
		// Live-confirmed path/shape against production, Refs #8393.
		assert.Equal(t, "/api/v1/projects/proj_123/database/vectors/upsert-batch", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req []VectorItem
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(req))
		assert.Equal(t, "default", req[0].Namespace)

		response := UpsertVectorsResponse{
			SuccessCount: 2,
			ErrorCount:   0,
			TotalTimeMs:  12.5,
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
	vectors := []VectorItem{
		{
			VectorEmbedding: []float64{0.1, 0.2, 0.3},
			Namespace:       "default",
			Metadata: map[string]interface{}{
				"category": "test",
			},
		},
		{
			VectorEmbedding: []float64{0.4, 0.5, 0.6},
			Namespace:       "default",
		},
	}

	response, err := client.ZeroDB.Vectors.Upsert(ctx, "proj_123", vectors)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 0, response.ErrorCount)
}

func TestVectorsService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Live-confirmed path/shape against production, Refs #8393.
		assert.Equal(t, "/api/v1/projects/proj_123/database/vectors/search", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req VectorSearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(req.QueryVector))
		assert.Equal(t, 5, req.Limit)

		response := VectorSearchResponse{
			Vectors: []VectorMatch{
				{
					VectorID: "vec_1",
					Metadata: map[string]interface{}{
						"category": "test",
					},
					VectorEmbedding: []float64{0.1, 0.2, 0.3},
					Similarity:      0.95,
				},
				{
					VectorID:   "vec_2",
					Similarity: 0.87,
				},
			},
			TotalCount:   2,
			SearchTimeMs: 3.2,
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
	req := &VectorSearchRequest{
		QueryVector: []float64{0.1, 0.2, 0.3},
		Limit:       5,
		Namespace:   "default",
	}

	response, err := client.ZeroDB.Vectors.Search(ctx, "proj_123", req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Vectors))
	assert.Equal(t, "vec_1", response.Vectors[0].VectorID)
	assert.Equal(t, 0.95, response.Vectors[0].Similarity)
}

func TestMemoryService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Live-confirmed path/shape against production, Refs #8393.
		assert.Equal(t, "/api/v1/public/memory/simple/add", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateMemoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "Test content", req.Content)
		assert.Equal(t, []string{"test", "memory"}, req.Tags)

		response := CreateMemoryResponse{
			MemoryID: "mem_123",
			Status:   "stored",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	req := &CreateMemoryRequest{
		Content: "Test content",
		Tags:    []string{"test", "memory"},
	}

	memory, err := client.ZeroDB.Memory.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, memory)
	assert.Equal(t, "mem_123", memory.MemoryID)
	assert.Equal(t, "stored", memory.Status)
}

func TestMemoryService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Live-confirmed path/shape against production, Refs #8393:
		// a GET with query params, not a POST with a JSON body.
		assert.Equal(t, "/api/v1/public/memory/simple/search", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "test query", r.URL.Query().Get("q"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))

		response := SearchMemoryResponse{
			Results: []MemoryItem{
				{ID: "mem_1", Content: "Test content 1", RelevanceScore: 0.95},
				{ID: "mem_2", Content: "Test content 2", RelevanceScore: 0.80},
			},
			Count: 2,
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
	response, err := client.ZeroDB.Memory.Search(ctx, "test query", 10)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Results))
	assert.Equal(t, 2, response.Count)
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
	_, err = client.ZeroDB.Vectors.Upsert(ctx, "", []VectorItem{{VectorEmbedding: []float64{0.1}}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project ID is required")

	// Test Upsert with empty vectors
	_, err = client.ZeroDB.Vectors.Upsert(ctx, "proj_123", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vectors cannot be empty")

	// Test Search with nil request
	_, err = client.ZeroDB.Vectors.Search(ctx, "proj_123", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test Search with empty query_vector
	_, err = client.ZeroDB.Vectors.Search(ctx, "proj_123", &VectorSearchRequest{QueryVector: []float64{}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query_vector cannot be empty")
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

	// Test Search with empty query
	_, err = client.ZeroDB.Memory.Search(ctx, "", 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query is required")
}
