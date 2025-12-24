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

// TestGenerate_Success tests successful embedding generation
func TestGenerate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/embeddings/generate", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req GenerateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(req.Texts))
		assert.Equal(t, "BAAI/bge-small-en-v1.5", req.Model)
		assert.True(t, req.Normalize)

		response := GenerateResponse{
			Embeddings:       [][]float64{{0.1, 0.2, 0.3}, {0.4, 0.5, 0.6}},
			Model:            "BAAI/bge-small-en-v1.5",
			Dimensions:       384,
			Count:            2,
			ProcessingTimeMs: 45.2,
			CostUSD:          0.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	result, err := client.ZeroDB.Embeddings.Generate(ctx, []string{"test text 1", "test text 2"}, "BAAI/bge-small-en-v1.5", true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Count)
	assert.Equal(t, 384, result.Dimensions)
	assert.Equal(t, 2, len(result.Embeddings))
	assert.Equal(t, "BAAI/bge-small-en-v1.5", result.Model)
	assert.Equal(t, 0.0, result.CostUSD)
}

// TestGenerate_EmptyTexts tests error on empty texts array
func TestGenerate_EmptyTexts(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{}, "", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "texts list cannot be empty")
}

// TestGenerate_TooManyTexts tests error on exceeding max texts
func TestGenerate_TooManyTexts(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	// Create 101 texts (exceeds limit)
	texts := make([]string, 101)
	for i := range texts {
		texts[i] = "test text"
	}

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, texts, "", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum 100 texts")
}

// TestGenerate_DefaultModel tests default model is used when not specified
func TestGenerate_DefaultModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req GenerateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "BAAI/bge-small-en-v1.5", req.Model) // Default model

		response := GenerateResponse{
			Embeddings:       [][]float64{{0.1, 0.2}},
			Model:            "BAAI/bge-small-en-v1.5",
			Dimensions:       384,
			Count:            1,
			ProcessingTimeMs: 20.0,
			CostUSD:          0.0,
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
	result, err := client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.NoError(t, err)
	assert.Equal(t, "BAAI/bge-small-en-v1.5", result.Model)
}

// TestGenerate_Normalize tests normalization parameter
func TestGenerate_Normalize(t *testing.T) {
	tests := []struct {
		name      string
		normalize bool
	}{
		{"with normalization", true},
		{"without normalization", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req GenerateRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assert.NoError(t, err)
				assert.Equal(t, tt.normalize, req.Normalize)

				response := GenerateResponse{
					Embeddings:       [][]float64{{0.1}},
					Model:            "BAAI/bge-small-en-v1.5",
					Dimensions:       384,
					Count:            1,
					ProcessingTimeMs: 10.0,
					CostUSD:          0.0,
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
			_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", tt.normalize)
			assert.NoError(t, err)
		})
	}
}

// TestGenerate_MultipleTexts tests handling of multiple texts
func TestGenerate_MultipleTexts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req GenerateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(req.Texts))

		// Generate 5 embeddings
		embeddings := make([][]float64, 5)
		for i := range embeddings {
			embeddings[i] = []float64{float64(i) * 0.1, float64(i) * 0.2}
		}

		response := GenerateResponse{
			Embeddings:       embeddings,
			Model:            "BAAI/bge-small-en-v1.5",
			Dimensions:       384,
			Count:            5,
			ProcessingTimeMs: 100.0,
			CostUSD:          0.0,
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
	texts := []string{"text1", "text2", "text3", "text4", "text5"}
	result, err := client.ZeroDB.Embeddings.Generate(ctx, texts, "", true)

	assert.NoError(t, err)
	assert.Equal(t, 5, result.Count)
	assert.Equal(t, 5, len(result.Embeddings))
}

// TestEmbedAndStore_Success tests successful embed and store operation
func TestEmbedAndStore_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/embeddings/embed-and-store", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req EmbedAndStoreRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "proj-123", req.ProjectID)
		assert.Equal(t, 2, len(req.Texts))
		assert.Equal(t, "tutorials", req.Namespace)

		response := EmbedAndStoreResponse{
			Success:              true,
			VectorsStored:        2,
			EmbeddingsGenerated:  2,
			Model:                "BAAI/bge-small-en-v1.5",
			Dimensions:           384,
			Namespace:            "tutorials",
			ProcessingTimeMs:     150.5,
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
	texts := []string{"text 1", "text 2"}
	result, err := client.ZeroDB.Embeddings.EmbedAndStore(ctx, "proj-123", texts, nil, "tutorials", "")

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.VectorsStored)
	assert.Equal(t, 2, result.EmbeddingsGenerated)
	assert.Equal(t, "tutorials", result.Namespace)
}

// TestEmbedAndStore_EmptyProjectID tests error on empty project ID
func TestEmbedAndStore_EmptyProjectID(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.EmbedAndStore(ctx, "", []string{"test"}, nil, "default", "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project ID is required")
}

// TestEmbedAndStore_WithMetadata tests embed and store with metadata
func TestEmbedAndStore_WithMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req EmbedAndStoreRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.NotNil(t, req.MetadataList)
		assert.Equal(t, 2, len(req.MetadataList))
		assert.Equal(t, "ml", req.MetadataList[0]["category"])
		assert.Equal(t, "dl", req.MetadataList[1]["category"])

		response := EmbedAndStoreResponse{
			Success:              true,
			VectorsStored:        2,
			EmbeddingsGenerated:  2,
			Model:                "BAAI/bge-small-en-v1.5",
			Dimensions:           384,
			Namespace:            "default",
			ProcessingTimeMs:     100.0,
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
	texts := []string{"ML tutorial", "DL guide"}
	metadata := []map[string]interface{}{
		{"category": "ml", "difficulty": "beginner"},
		{"category": "dl", "difficulty": "advanced"},
	}

	result, err := client.ZeroDB.Embeddings.EmbedAndStore(ctx, "proj-123", texts, metadata, "", "")

	assert.NoError(t, err)
	assert.True(t, result.Success)
}

// TestEmbedAndStore_MetadataMismatch tests error when metadata length doesn't match texts
func TestEmbedAndStore_MetadataMismatch(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()
	texts := []string{"text1", "text2"}
	metadata := []map[string]interface{}{{"key": "value"}} // Only 1 metadata for 2 texts

	_, err = client.ZeroDB.Embeddings.EmbedAndStore(ctx, "proj-123", texts, metadata, "", "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata_list length must match texts length")
}

// TestSemanticSearch_Success tests successful semantic search
func TestSemanticSearch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/embeddings/semantic-search", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req SemanticSearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "proj-123", req.ProjectID)
		assert.Equal(t, "machine learning tutorial", req.Query)
		assert.Equal(t, 5, req.Limit)
		assert.Equal(t, 0.8, req.Threshold)

		response := SemanticSearchResponse{
			Results: []SemanticSearchResult{
				{
					VectorID:   "vec-1",
					Similarity: 0.95,
					Document:   "Introduction to machine learning",
					Metadata:   map[string]interface{}{"category": "ml"},
					Namespace:  "default",
				},
				{
					VectorID:   "vec-2",
					Similarity: 0.87,
					Document:   "ML basics and fundamentals",
					Metadata:   map[string]interface{}{"category": "ml"},
					Namespace:  "default",
				},
			},
			Query:            "machine learning tutorial",
			TotalResults:     2,
			Model:            "BAAI/bge-small-en-v1.5",
			ProcessingTimeMs: 75.3,
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
	result, err := client.ZeroDB.Embeddings.SemanticSearch(
		ctx,
		"proj-123",
		"machine learning tutorial",
		5,
		0.8,
		"default",
		nil,
		"",
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Results))
	assert.Equal(t, 2, result.TotalResults)
	assert.Equal(t, "machine learning tutorial", result.Query)
	assert.Equal(t, 0.95, result.Results[0].Similarity)
}

// TestSemanticSearch_EmptyQuery tests error on empty query
func TestSemanticSearch_EmptyQuery(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "", 10, 0.7, "default", nil, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query cannot be empty")
}

// TestSemanticSearch_DefaultLimit tests default limit is used when not specified
func TestSemanticSearch_DefaultLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SemanticSearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 10, req.Limit) // Default limit

		response := SemanticSearchResponse{
			Results:          []SemanticSearchResult{},
			Query:            "test",
			TotalResults:     0,
			Model:            "BAAI/bge-small-en-v1.5",
			ProcessingTimeMs: 10.0,
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
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "test", 0, 0, "", nil, "")

	assert.NoError(t, err)
}

// TestSemanticSearch_CustomThreshold tests custom similarity threshold
func TestSemanticSearch_CustomThreshold(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SemanticSearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, 0.9, req.Threshold)

		response := SemanticSearchResponse{
			Results:          []SemanticSearchResult{},
			Query:            "test",
			TotalResults:     0,
			Model:            "BAAI/bge-small-en-v1.5",
			ProcessingTimeMs: 10.0,
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
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "test", 10, 0.9, "", nil, "")

	assert.NoError(t, err)
}

// TestSemanticSearch_InvalidLimit tests error on invalid limit
func TestSemanticSearch_InvalidLimit(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test limit > 100
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "test", 101, 0.7, "", nil, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be between 1 and 100")

	// Test negative limit
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "test", -1, 0.7, "", nil, "")
	assert.Error(t, err)
}

// TestSemanticSearch_InvalidThreshold tests error on invalid threshold
func TestSemanticSearch_InvalidThreshold(t *testing.T) {
	client, err := NewClient(&Config{APIKey: "test-key"})
	require.NoError(t, err)

	ctx := context.Background()

	// Test threshold > 1.0
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "test", 10, 1.5, "", nil, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "threshold must be between 0.0 and 1.0")

	// Test negative threshold
	_, err = client.ZeroDB.Embeddings.SemanticSearch(ctx, "proj-123", "test", 10, -0.1, "", nil, "")
	assert.Error(t, err)
}

// TestListModels_Success tests successful model listing
func TestListModels_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/embeddings/models", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		models := []EmbeddingModel{
			{
				ID:          "BAAI/bge-small-en-v1.5",
				Dimensions:  384,
				Description: "Fast & efficient - General-purpose embeddings",
				Speed:       "⚡⚡⚡ Very Fast",
				Loaded:      true,
				CostPer1k:   0.0,
			},
			{
				ID:          "BAAI/bge-base-en-v1.5",
				Dimensions:  768,
				Description: "Balanced quality - Better semantic understanding",
				Speed:       "⚡⚡ Fast",
				Loaded:      true,
				CostPer1k:   0.0,
			},
			{
				ID:          "BAAI/bge-large-en-v1.5",
				Dimensions:  1024,
				Description: "Best quality - Mission-critical applications",
				Speed:       "⚡ Medium",
				Loaded:      false,
				CostPer1k:   0.0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	models, err := client.ZeroDB.Embeddings.ListModels(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, models)
	assert.Equal(t, 3, len(models))
	assert.Equal(t, "BAAI/bge-small-en-v1.5", models[0].ID)
	assert.Equal(t, 384, models[0].Dimensions)
	assert.True(t, models[0].Loaded)
	assert.Equal(t, 0.0, models[0].CostPer1k)
}

// TestHealthCheck_Success tests successful health check
func TestHealthCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/embeddings/health", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := HealthCheckResponse{
			Status: "healthy",
			EmbeddingService: map[string]interface{}{
				"model_loaded": true,
				"model_name":   "BAAI/bge-small-en-v1.5",
				"dimensions":   384,
			},
			URL:              "https://embedding-service.railway.app",
			CostPerEmbedding: 0.0,
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
	health, err := client.ZeroDB.Embeddings.HealthCheck(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, 0.0, health.CostPerEmbedding)
	assert.NotNil(t, health.EmbeddingService)
	assert.Equal(t, true, health.EmbeddingService["model_loaded"])
}

// TestGetUsage_Success tests successful usage retrieval
func TestGetUsage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/embeddings/usage", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := UsageResponse{
			UserID:                   "user-123",
			EmbeddingsGeneratedToday: 1500,
			EmbeddingsGeneratedMonth: 45000,
			CostTodayUSD:             0.0,
			CostMonthUSD:             0.0,
			Model:                    "BAAI/bge-small-en-v1.5",
			Service:                  "Railway HuggingFace",
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
	usage, err := client.ZeroDB.Embeddings.GetUsage(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, usage)
	assert.Equal(t, "user-123", usage.UserID)
	assert.Equal(t, 1500, usage.EmbeddingsGeneratedToday)
	assert.Equal(t, 45000, usage.EmbeddingsGeneratedMonth)
	assert.Equal(t, 0.0, usage.CostTodayUSD)
	assert.Equal(t, 0.0, usage.CostMonthUSD)
	assert.Equal(t, "Railway HuggingFace", usage.Service)
}

// TestHTTPError_401 tests handling of 401 Unauthorized errors
func TestHTTPError_401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid API key",
			"code":    "UNAUTHORIZED",
		})
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "invalid-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 401, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Invalid API key")
}

// TestHTTPError_429 tests handling of 429 Rate Limit errors
func TestHTTPError_429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", "req-456")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Rate limit exceeded",
			"code":    "RATE_LIMIT_EXCEEDED",
		})
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:   0, // Disable retries for this test
			InitialDelay: 1 * time.Millisecond,
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 429, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Rate limit exceeded")
}

// TestHTTPError_422 tests handling of 422 Validation errors
func TestHTTPError_422(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid model specified",
			"code":    "VALIDATION_ERROR",
			"details": map[string]interface{}{
				"field": "model",
				"error": "Model not found",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "invalid-model", true)

	assert.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 422, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Invalid model")
}

// TestHTTPError_500 tests handling of 500 Internal Server errors
func TestHTTPError_500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Internal server error",
			"code":    "INTERNAL_ERROR",
		})
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:   0, // Disable retries
			InitialDelay: 1 * time.Millisecond,
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 500, apiErr.StatusCode)
}

// TestNetworkTimeout tests handling of network timeout
func TestNetworkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Timeout: 100 * time.Millisecond, // Very short timeout
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

// TestMalformedJSONResponse tests handling of malformed JSON response
func TestMalformedJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.Error(t, err)
}

// TestContextCancellation tests handling of context cancellation
func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = client.ZeroDB.Embeddings.Generate(ctx, []string{"test"}, "", true)

	assert.Error(t, err)
}

// TestSemanticSearch_WithMetadataFilter tests semantic search with metadata filters
func TestSemanticSearch_WithMetadataFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SemanticSearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.NotNil(t, req.FilterMetadata)
		assert.Equal(t, "ml", req.FilterMetadata["category"])
		assert.Equal(t, "beginner", req.FilterMetadata["difficulty"])

		response := SemanticSearchResponse{
			Results: []SemanticSearchResult{
				{
					VectorID:   "vec-1",
					Similarity: 0.92,
					Document:   "ML for beginners",
					Metadata: map[string]interface{}{
						"category":   "ml",
						"difficulty": "beginner",
					},
					Namespace: "default",
				},
			},
			Query:            "machine learning basics",
			TotalResults:     1,
			Model:            "BAAI/bge-small-en-v1.5",
			ProcessingTimeMs: 50.0,
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
	filterMetadata := map[string]interface{}{
		"category":   "ml",
		"difficulty": "beginner",
	}

	result, err := client.ZeroDB.Embeddings.SemanticSearch(
		ctx,
		"proj-123",
		"machine learning basics",
		10,
		0.7,
		"default",
		filterMetadata,
		"",
	)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "ml", result.Results[0].Metadata["category"])
}
