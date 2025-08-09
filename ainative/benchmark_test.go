package ainative

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Benchmark tests to compare Go SDK performance with other language SDKs

func BenchmarkClientCreation(b *testing.B) {
	config := &Config{
		APIKey: "bench-key",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := NewClient(config)
		if err != nil {
			b.Fatal(err)
		}
		_ = client
	}
}

func BenchmarkProjectOperations(b *testing.B) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Project{
			ID:     "bench_proj_123",
			Name:   "Benchmark Project",
			Status: ProjectStatusActive,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	req := &CreateProjectRequest{
		Name:        "Benchmark Project",
		Description: "Project for benchmarking",
	}

	b.ResetTimer()
	b.Run("CreateProject", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Projects.Create(ctx, req)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	
	b.Run("GetProject", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Projects.Get(ctx, "bench_proj_123")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkVectorOperations(b *testing.B) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/zerodb/projects/proj_123/vectors":
			// Upsert response
			response := UpsertVectorsResponse{
				UpsertedCount: 1,
				Namespace:     "default",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case "/api/v1/zerodb/projects/proj_123/vectors/search":
			// Search response
			response := VectorSearchResponse{
				Matches: []VectorSearchMatch{
					{
						ID:    "vec_1",
						Score: 0.95,
						Vector: []float64{0.1, 0.2, 0.3},
					},
				},
				Namespace: "default",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	projectID := "proj_123"

	// Prepare test data
	upsertReq := &UpsertVectorsRequest{
		Vectors: []VectorItem{
			{
				ID:     "bench_vec_1",
				Vector: []float64{0.1, 0.2, 0.3},
				Metadata: map[string]interface{}{
					"category": "benchmark",
				},
			},
		},
		Namespace: "default",
	}

	searchReq := &VectorSearchRequest{
		Vector:          []float64{0.1, 0.2, 0.3},
		TopK:            5,
		Namespace:       "default",
		IncludeMetadata: true,
	}

	b.ResetTimer()
	b.Run("UpsertVectors", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Vectors.Upsert(ctx, projectID, upsertReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SearchVectors", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Vectors.Search(ctx, projectID, searchReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMemoryOperations(b *testing.B) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/memory":
			// Create response
			response := MemoryItem{
				ID:       "bench_mem_123",
				Title:    "Benchmark Memory",
				Content:  "Test content for benchmarking",
				Priority: MemoryPriorityMedium,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case "/api/v1/memory/search":
			// Search response
			response := SearchMemoryResponse{
				Results: []MemoryItem{
					{
						ID:       "bench_mem_1",
						Title:    "Memory 1",
						Content:  "Benchmark content 1",
						Priority: MemoryPriorityMedium,
					},
				},
				Total: 1,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	createReq := &CreateMemoryRequest{
		Title:    "Benchmark Memory",
		Content:  "Test content for benchmarking",
		Priority: MemoryPriorityMedium,
		Tags:     []string{"benchmark", "test"},
	}

	searchReq := &SearchMemoryRequest{
		Query:    "benchmark",
		Limit:    10,
		Semantic: true,
	}

	b.ResetTimer()
	b.Run("CreateMemory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Memory.Create(ctx, createReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SearchMemory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Memory.Search(ctx, searchReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAgentSwarmOperations(b *testing.B) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/agent-swarm/swarms":
			// Start swarm response
			response := AgentSwarm{
				ID:        "bench_swarm_123",
				Name:      "Benchmark Swarm",
				ProjectID: "proj_123",
				Status:    SwarmStatusRunning,
				Agents: []Agent{
					{
						ID:     "agent_1",
						Type:   AgentTypeAnalyzer,
						Status: AgentStatusActive,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case r.URL.Path == "/api/v1/agent-swarm/orchestrate":
			// Orchestrate response
			response := OrchestrationResponse{
				TaskID:     "bench_task_123",
				Status:     TaskStatusAssigned,
				AssignedTo: []string{"agent_1"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	startReq := &StartSwarmRequest{
		ProjectID: "proj_123",
		Name:      "Benchmark Swarm",
		Objective: "Benchmark agent swarm operations",
		Agents: []AgentConfig{
			{
				Type:         AgentTypeAnalyzer,
				Count:        1,
				Capabilities: []string{"analysis"},
			},
		},
	}

	orchestrateReq := &OrchestrationRequest{
		SwarmID:  "bench_swarm_123",
		Task:     "benchmark_task",
		Priority: TaskPriorityMedium,
	}

	b.ResetTimer()
	b.Run("StartSwarm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.AgentSwarm.Start(ctx, startReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("OrchestrateTasks", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.AgentSwarm.Orchestrate(ctx, orchestrateReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAuthOperations(b *testing.B) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/me":
			// User info response
			response := UserInfo{
				ID:    "bench_user_123",
				Email: "bench@example.com",
				Name:  "Benchmark User",
				Role:  "user",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case "/api/v1/auth/api-keys":
			if r.Method == "GET" {
				// List API keys response
				response := struct {
					APIKeys []APIKeyInfo `json:"api_keys"`
				}{
					APIKeys: []APIKeyInfo{
						{
							ID:       "key_1",
							Name:     "Benchmark Key",
							Prefix:   "ak_bench",
							IsActive: true,
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				// Create API key response
				response := CreateAPIKeyResponse{
					ID:     "new_key_123",
					Name:   "New Benchmark Key",
					Key:    "ak_bench_newkey123456789",
					Prefix: "ak_bench",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	createKeyReq := &CreateAPIKeyRequest{
		Name:        "Benchmark Key",
		Permissions: []string{"read"},
	}

	b.ResetTimer()
	b.Run("GetUserInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.Auth.GetUserInfo(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ListAPIKeys", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.Auth.ListAPIKeys(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("CreateAPIKey", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.Auth.CreateAPIKey(ctx, createKeyReq)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark concurrent operations to test connection pooling and rate limiting
func BenchmarkConcurrentOperations(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add small delay to simulate real API latency
		time.Sleep(10 * time.Millisecond)
		
		response := Project{
			ID:     "concurrent_proj_123",
			Name:   "Concurrent Project",
			Status: ProjectStatusActive,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:    "bench-key",
		BaseURL:   server.URL,
		RateLimit: 1000, // High rate limit for benchmarking
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.ResetTimer()
	b.Run("ConcurrentRequests", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := client.ZeroDB.Projects.Get(ctx, "concurrent_proj_123")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// Benchmark large payload operations
func BenchmarkLargePayloads(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := UpsertVectorsResponse{
			UpsertedCount: 100,
			Namespace:     "large_batch",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	projectID := "proj_123"

	// Generate large batch of vectors
	generateLargeBatch := func(size int) *UpsertVectorsRequest {
		vectors := make([]VectorItem, size)
		for i := 0; i < size; i++ {
			vector := make([]float64, 384) // Standard embedding size
			for j := range vector {
				vector[j] = float64(i*j) / 1000.0
			}
			
			vectors[i] = VectorItem{
				ID:     fmt.Sprintf("large_vec_%d", i),
				Vector: vector,
				Metadata: map[string]interface{}{
					"batch_id": "large_batch",
					"index":    i,
				},
			}
		}
		
		return &UpsertVectorsRequest{
			Vectors:   vectors,
			Namespace: "large_batch",
		}
	}

	b.ResetTimer()
	b.Run("SmallBatch_10", func(b *testing.B) {
		req := generateLargeBatch(10)
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Vectors.Upsert(ctx, projectID, req)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MediumBatch_100", func(b *testing.B) {
		req := generateLargeBatch(100)
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Vectors.Upsert(ctx, projectID, req)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LargeBatch_1000", func(b *testing.B) {
		req := generateLargeBatch(1000)
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Vectors.Upsert(ctx, projectID, req)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark memory allocations
func BenchmarkMemoryAllocations(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Project{
			ID:     "alloc_proj_123",
			Name:   "Allocation Project",
			Status: ProjectStatusActive,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		APIKey:  "bench-key",
		BaseURL: server.URL,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	b.Run("AllocationEfficiency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := client.ZeroDB.Projects.Get(ctx, "alloc_proj_123")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}