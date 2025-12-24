//go:build integration
// +build integration

package ainative

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests for embeddings operations
// Run with: go test -tags=integration ./ainative -v

// TestIntegration_Generate tests actual embedding generation against live API
func TestIntegration_Generate(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		t.Skip("AINATIVE_API_KEY not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	ctx := context.Background()
	texts := []string{
		"Machine learning is a subset of artificial intelligence",
		"Deep learning uses neural networks with multiple layers",
	}

	result, err := client.ZeroDB.Embeddings.Generate(ctx, texts, "BAAI/bge-small-en-v1.5", true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Count)
	assert.Equal(t, 384, result.Dimensions)
	assert.Equal(t, 2, len(result.Embeddings))
	assert.Equal(t, "BAAI/bge-small-en-v1.5", result.Model)
	assert.Equal(t, 0.0, result.CostUSD) // Should be free

	// Verify embedding dimensions
	for _, embedding := range result.Embeddings {
		assert.Equal(t, 384, len(embedding))
	}

	t.Logf("Successfully generated %d embeddings in %.2fms", result.Count, result.ProcessingTimeMs)
}

// TestIntegration_EmbedAndStore tests embedding and storing vectors
func TestIntegration_EmbedAndStore(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	projectID := os.Getenv("AINATIVE_PROJECT_ID")
	if apiKey == "" || projectID == "" {
		t.Skip("AINATIVE_API_KEY or AINATIVE_PROJECT_ID not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	ctx := context.Background()
	texts := []string{
		"Go SDK integration test document 1",
		"Go SDK integration test document 2",
	}
	metadata := []map[string]interface{}{
		{"source": "integration_test", "doc_id": 1},
		{"source": "integration_test", "doc_id": 2},
	}

	result, err := client.ZeroDB.Embeddings.EmbedAndStore(
		ctx,
		projectID,
		texts,
		metadata,
		"integration_tests",
		"BAAI/bge-small-en-v1.5",
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.VectorsStored)
	assert.Equal(t, 2, result.EmbeddingsGenerated)
	assert.Equal(t, 384, result.Dimensions)
	assert.Equal(t, "integration_tests", result.Namespace)

	t.Logf("Successfully stored %d vectors in %.2fms", result.VectorsStored, result.ProcessingTimeMs)
}

// TestIntegration_SemanticSearch tests semantic search functionality
func TestIntegration_SemanticSearch(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	projectID := os.Getenv("AINATIVE_PROJECT_ID")
	if apiKey == "" || projectID == "" {
		t.Skip("AINATIVE_API_KEY or AINATIVE_PROJECT_ID not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	// First, store some test data
	ctx := context.Background()
	texts := []string{
		"Python programming language tutorial",
		"JavaScript web development guide",
		"Go concurrent programming",
	}

	_, err = client.ZeroDB.Embeddings.EmbedAndStore(
		ctx,
		projectID,
		texts,
		nil,
		"search_integration_test",
		"BAAI/bge-small-en-v1.5",
	)
	require.NoError(t, err)

	// Now search
	result, err := client.ZeroDB.Embeddings.SemanticSearch(
		ctx,
		projectID,
		"programming languages",
		5,
		0.5,
		"search_integration_test",
		nil,
		"BAAI/bge-small-en-v1.5",
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.TotalResults, 0)
	assert.Greater(t, len(result.Results), 0)

	// Verify results are sorted by similarity
	for i := 0; i < len(result.Results)-1; i++ {
		assert.GreaterOrEqual(t, result.Results[i].Similarity, result.Results[i+1].Similarity)
	}

	t.Logf("Found %d results in %.2fms", result.TotalResults, result.ProcessingTimeMs)
	t.Logf("Top result: %s (similarity: %.4f)", result.Results[0].Document, result.Results[0].Similarity)
}

// TestIntegration_ListModels tests listing available embedding models
func TestIntegration_ListModels(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		t.Skip("AINATIVE_API_KEY not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	ctx := context.Background()
	models, err := client.ZeroDB.Embeddings.ListModels(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, models)
	assert.Greater(t, len(models), 0)

	// Verify expected models are present
	modelIDs := make(map[string]bool)
	for _, model := range models {
		modelIDs[model.ID] = true
		assert.Greater(t, model.Dimensions, 0)
		assert.NotEmpty(t, model.Description)
		assert.Equal(t, 0.0, model.CostPer1k) // All models should be free
		t.Logf("Model: %s (%d dims) - %s", model.ID, model.Dimensions, model.Speed)
	}

	// Check for expected models
	expectedModels := []string{
		"BAAI/bge-small-en-v1.5",
		"BAAI/bge-base-en-v1.5",
		"BAAI/bge-large-en-v1.5",
	}

	for _, expectedModel := range expectedModels {
		if !modelIDs[expectedModel] {
			t.Logf("Warning: Expected model %s not found", expectedModel)
		}
	}
}

// TestIntegration_HealthCheck tests embedding service health check
func TestIntegration_HealthCheck(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		t.Skip("AINATIVE_API_KEY not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	ctx := context.Background()
	health, err := client.ZeroDB.Embeddings.HealthCheck(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.NotNil(t, health.EmbeddingService)
	assert.Equal(t, 0.0, health.CostPerEmbedding)

	t.Logf("Health status: %s", health.Status)
	t.Logf("Service URL: %s", health.URL)
}

// TestIntegration_GetUsage tests usage statistics retrieval
func TestIntegration_GetUsage(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		t.Skip("AINATIVE_API_KEY not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	ctx := context.Background()
	usage, err := client.ZeroDB.Embeddings.GetUsage(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, usage)
	assert.NotEmpty(t, usage.UserID)
	assert.GreaterOrEqual(t, usage.EmbeddingsGeneratedToday, 0)
	assert.GreaterOrEqual(t, usage.EmbeddingsGeneratedMonth, 0)
	assert.Equal(t, 0.0, usage.CostTodayUSD)   // Should be free
	assert.Equal(t, 0.0, usage.CostMonthUSD)   // Should be free

	t.Logf("User: %s", usage.UserID)
	t.Logf("Embeddings today: %d", usage.EmbeddingsGeneratedToday)
	t.Logf("Embeddings this month: %d", usage.EmbeddingsGeneratedMonth)
	t.Logf("Service: %s", usage.Service)
}

// TestIntegration_FullWorkflow tests complete embedding workflow
func TestIntegration_FullWorkflow(t *testing.T) {
	apiKey := os.Getenv("AINATIVE_API_KEY")
	projectID := os.Getenv("AINATIVE_PROJECT_ID")
	if apiKey == "" || projectID == "" {
		t.Skip("AINATIVE_API_KEY or AINATIVE_PROJECT_ID not set, skipping integration test")
	}

	client, err := NewClient(&Config{
		APIKey: apiKey,
	})
	require.NoError(t, err)

	ctx := context.Background()

	// Step 1: List available models
	t.Log("Step 1: Listing models...")
	models, err := client.ZeroDB.Embeddings.ListModels(ctx)
	require.NoError(t, err)
	require.Greater(t, len(models), 0)
	t.Logf("Found %d models", len(models))

	// Step 2: Generate embeddings
	t.Log("Step 2: Generating embeddings...")
	texts := []string{
		"Artificial intelligence and machine learning",
		"Natural language processing",
		"Computer vision and image recognition",
	}
	genResult, err := client.ZeroDB.Embeddings.Generate(ctx, texts, "BAAI/bge-small-en-v1.5", true)
	require.NoError(t, err)
	require.Equal(t, 3, genResult.Count)
	t.Logf("Generated %d embeddings", genResult.Count)

	// Step 3: Embed and store
	t.Log("Step 3: Embedding and storing vectors...")
	metadata := []map[string]interface{}{
		{"topic": "AI/ML", "test": "workflow"},
		{"topic": "NLP", "test": "workflow"},
		{"topic": "CV", "test": "workflow"},
	}
	storeResult, err := client.ZeroDB.Embeddings.EmbedAndStore(
		ctx,
		projectID,
		texts,
		metadata,
		"workflow_test",
		"BAAI/bge-small-en-v1.5",
	)
	require.NoError(t, err)
	require.True(t, storeResult.Success)
	require.Equal(t, 3, storeResult.VectorsStored)
	t.Logf("Stored %d vectors", storeResult.VectorsStored)

	// Step 4: Semantic search
	t.Log("Step 4: Performing semantic search...")
	searchResult, err := client.ZeroDB.Embeddings.SemanticSearch(
		ctx,
		projectID,
		"machine learning and AI",
		5,
		0.7,
		"workflow_test",
		map[string]interface{}{"test": "workflow"},
		"BAAI/bge-small-en-v1.5",
	)
	require.NoError(t, err)
	require.Greater(t, searchResult.TotalResults, 0)
	t.Logf("Found %d results", searchResult.TotalResults)

	// Step 5: Check usage
	t.Log("Step 5: Checking usage statistics...")
	usage, err := client.ZeroDB.Embeddings.GetUsage(ctx)
	require.NoError(t, err)
	t.Logf("Total embeddings this month: %d", usage.EmbeddingsGeneratedMonth)

	// Step 6: Health check
	t.Log("Step 6: Checking service health...")
	health, err := client.ZeroDB.Embeddings.HealthCheck(ctx)
	require.NoError(t, err)
	require.Equal(t, "healthy", health.Status)

	t.Log("Full workflow completed successfully!")
}
