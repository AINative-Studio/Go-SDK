package ainative

import (
	"context"
	"fmt"
)

// EmbeddingsService handles embedding operations using Railway HuggingFace service
type EmbeddingsService struct {
	client *Client
}

// EmbeddingModel represents an available embedding model
type EmbeddingModel struct {
	ID          string  `json:"id"`
	Dimensions  int     `json:"dimensions"`
	Description string  `json:"description"`
	Speed       string  `json:"speed"`
	Loaded      bool    `json:"loaded"`
	CostPer1k   float64 `json:"cost_per_1k"`
}

// GenerateRequest represents a request to generate embeddings
type GenerateRequest struct {
	Texts     []string `json:"texts"`
	Model     string   `json:"model,omitempty"`
	Normalize bool     `json:"normalize"`
}

// GenerateResponse represents the response from generating embeddings
type GenerateResponse struct {
	Embeddings       [][]float64 `json:"embeddings"`
	Model            string      `json:"model"`
	Dimensions       int         `json:"dimensions"`
	Count            int         `json:"count"`
	ProcessingTimeMs float64     `json:"processing_time_ms"`
	CostUSD          float64     `json:"cost_usd"`
}

// EmbedAndStoreRequest represents a request to embed texts and store them
type EmbedAndStoreRequest struct {
	ProjectID    string                   `json:"project_id"`
	Texts        []string                 `json:"texts"`
	MetadataList []map[string]interface{} `json:"metadata_list,omitempty"`
	Namespace    string                   `json:"namespace,omitempty"`
	Model        string                   `json:"model,omitempty"`
}

// EmbedAndStoreResponse represents the response from embed and store operation
type EmbedAndStoreResponse struct {
	Success              bool    `json:"success"`
	VectorsStored        int     `json:"vectors_stored"`
	EmbeddingsGenerated  int     `json:"embeddings_generated"`
	Model                string  `json:"model"`
	Dimensions           int     `json:"dimensions"`
	Namespace            string  `json:"namespace"`
	ProcessingTimeMs     float64 `json:"processing_time_ms"`
}

// SemanticSearchRequest represents a request for semantic search
type SemanticSearchRequest struct {
	ProjectID      string                 `json:"project_id"`
	Query          string                 `json:"query"`
	Limit          int                    `json:"limit,omitempty"`
	Threshold      float64                `json:"threshold,omitempty"`
	Namespace      string                 `json:"namespace,omitempty"`
	FilterMetadata map[string]interface{} `json:"filter_metadata,omitempty"`
	Model          string                 `json:"model,omitempty"`
}

// SemanticSearchResult represents a single search result
type SemanticSearchResult struct {
	VectorID   string                 `json:"vector_id"`
	Similarity float64                `json:"similarity"`
	Document   string                 `json:"document"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Namespace  string                 `json:"namespace"`
}

// SemanticSearchResponse represents the response from semantic search
type SemanticSearchResponse struct {
	Results          []SemanticSearchResult `json:"results"`
	Query            string                 `json:"query"`
	TotalResults     int                    `json:"total_results"`
	Model            string                 `json:"model"`
	ProcessingTimeMs float64                `json:"processing_time_ms"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status           string                 `json:"status"`
	EmbeddingService map[string]interface{} `json:"embedding_service"`
	URL              string                 `json:"url"`
	CostPerEmbedding float64                `json:"cost_per_embedding"`
}

// UsageResponse represents embedding usage statistics
type UsageResponse struct {
	UserID                   string  `json:"user_id"`
	EmbeddingsGeneratedToday int     `json:"embeddings_generated_today"`
	EmbeddingsGeneratedMonth int     `json:"embeddings_generated_month"`
	CostTodayUSD             float64 `json:"cost_today_usd"`
	CostMonthUSD             float64 `json:"cost_month_usd"`
	Model                    string  `json:"model"`
	Service                  string  `json:"service"`
}

// Generate generates embeddings for texts using Railway HuggingFace service.
//
// Uses FREE self-hosted HuggingFace service via Railway with
// BAAI/bge-small-en-v1.5 model (384 dimensions).
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - texts: List of texts to embed (max 100)
//   - model: Embedding model to use (optional, defaults to BAAI/bge-small-en-v1.5)
//   - normalize: Normalize embeddings to unit length (optional, defaults to true)
//
// Returns:
//   - GenerateResponse containing embeddings and metadata
//   - error if validation fails or API request fails
//
// Example:
//
//	resp, err := client.ZeroDB.Embeddings.Generate(ctx, []string{
//	    "Machine learning tutorial",
//	    "Python programming guide",
//	}, "BAAI/bge-small-en-v1.5", true)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Generated %d embeddings\n", resp.Count)
//	fmt.Printf("Dimensions: %d\n", resp.Dimensions)
//	fmt.Printf("Cost: $%.2f\n", resp.CostUSD) // Always $0.00
func (s *EmbeddingsService) Generate(ctx context.Context, texts []string, model string, normalize bool) (*GenerateResponse, error) {
	if len(texts) == 0 {
		return nil, NewValidationError("texts", "texts list cannot be empty", texts)
	}

	if len(texts) > 100 {
		return nil, NewValidationError("texts", "maximum 100 texts per request", len(texts))
	}

	// Set defaults
	if model == "" {
		model = "BAAI/bge-small-en-v1.5"
	}

	req := &GenerateRequest{
		Texts:     texts,
		Model:     model,
		Normalize: normalize,
	}

	var result GenerateResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/embeddings/generate", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// EmbedAndStore generates embeddings and stores them in ZeroDB (one-step workflow).
//
// Combines embedding generation with vector storage in a single operation.
// Automatically embeds texts and stores them in your ZeroDB vectors.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - projectID: Project ID (UUID string)
//   - texts: List of texts to embed and store (max 100)
//   - metadataList: Optional metadata for each text (must match texts length)
//   - namespace: Vector namespace (optional, defaults to "default")
//   - model: Embedding model to use (optional, defaults to BAAI/bge-small-en-v1.5)
//
// Returns:
//   - EmbedAndStoreResponse containing storage results
//   - error if validation fails or API request fails
//
// Example:
//
//	texts := []string{
//	    "Introduction to machine learning",
//	    "Advanced deep learning techniques",
//	    "Natural language processing basics",
//	}
//	metadata := []map[string]interface{}{
//	    {"category": "ml", "difficulty": "beginner"},
//	    {"category": "dl", "difficulty": "advanced"},
//	    {"category": "nlp", "difficulty": "beginner"},
//	}
//	resp, err := client.ZeroDB.Embeddings.EmbedAndStore(
//	    ctx,
//	    "550e8400-e29b-41d4-a716-446655440000",
//	    texts,
//	    metadata,
//	    "tutorials",
//	    "",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Stored %d vectors\n", resp.VectorsStored)
func (s *EmbeddingsService) EmbedAndStore(ctx context.Context, projectID string, texts []string, metadataList []map[string]interface{}, namespace, model string) (*EmbedAndStoreResponse, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}

	if len(texts) == 0 {
		return nil, NewValidationError("texts", "texts list cannot be empty", texts)
	}

	if len(texts) > 100 {
		return nil, NewValidationError("texts", "maximum 100 texts per request", len(texts))
	}

	if metadataList != nil && len(metadataList) != len(texts) {
		return nil, NewValidationError("metadata_list", "metadata_list length must match texts length", fmt.Sprintf("texts: %d, metadata: %d", len(texts), len(metadataList)))
	}

	// Set defaults
	if namespace == "" {
		namespace = "default"
	}

	if model == "" {
		model = "BAAI/bge-small-en-v1.5"
	}

	req := &EmbedAndStoreRequest{
		ProjectID:    projectID,
		Texts:        texts,
		MetadataList: metadataList,
		Namespace:    namespace,
		Model:        model,
	}

	var result EmbedAndStoreResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/embeddings/embed-and-store", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// SemanticSearch performs text-based semantic search (auto-embed query + search).
//
// Performs semantic search using natural language queries. Automatically
// generates embedding for the query text and searches your vector database.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - projectID: Project ID (UUID string)
//   - query: Natural language search query
//   - limit: Maximum results (1-100, optional, defaults to 10)
//   - threshold: Similarity threshold (0.0-1.0, optional, defaults to 0.7)
//   - namespace: Vector namespace to search (optional, defaults to "default")
//   - filterMetadata: Optional metadata filters (MongoDB-style)
//   - model: Embedding model to use (optional, defaults to BAAI/bge-small-en-v1.5)
//
// Returns:
//   - SemanticSearchResponse containing search results
//   - error if validation fails or API request fails
//
// Example:
//
//	resp, err := client.ZeroDB.Embeddings.SemanticSearch(
//	    ctx,
//	    "550e8400-e29b-41d4-a716-446655440000",
//	    "How do I train a neural network?",
//	    5,
//	    0.8,
//	    "default",
//	    map[string]interface{}{"category": "ml"},
//	    "",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, result := range resp.Results {
//	    fmt.Printf("%s: %.2f\n", result.Document, result.Similarity)
//	}
func (s *EmbeddingsService) SemanticSearch(ctx context.Context, projectID, query string, limit int, threshold float64, namespace string, filterMetadata map[string]interface{}, model string) (*SemanticSearchResponse, error) {
	if projectID == "" {
		return nil, NewValidationError("project_id", "project ID is required", projectID)
	}

	if query == "" {
		return nil, NewValidationError("query", "query cannot be empty", query)
	}

	// Set defaults
	if limit == 0 {
		limit = 10
	}

	if limit < 1 || limit > 100 {
		return nil, NewValidationError("limit", "limit must be between 1 and 100", limit)
	}

	if threshold == 0.0 {
		threshold = 0.7
	}

	if threshold < 0.0 || threshold > 1.0 {
		return nil, NewValidationError("threshold", "threshold must be between 0.0 and 1.0", threshold)
	}

	if namespace == "" {
		namespace = "default"
	}

	if model == "" {
		model = "BAAI/bge-small-en-v1.5"
	}

	req := &SemanticSearchRequest{
		ProjectID:      projectID,
		Query:          query,
		Limit:          limit,
		Threshold:      threshold,
		Namespace:      namespace,
		FilterMetadata: filterMetadata,
		Model:          model,
	}

	var result SemanticSearchResponse

	err := s.client.makeRequest(ctx, "POST", "/api/v1/embeddings/semantic-search", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListModels lists available embedding models.
//
// Returns information about available embedding models including
// dimensions, speed, and whether they're currently loaded.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//
// Returns:
//   - List of EmbeddingModel objects
//   - error if API request fails
//
// Example:
//
//	models, err := client.ZeroDB.Embeddings.ListModels(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, model := range models {
//	    fmt.Printf("%s: %d dims, %s\n", model.ID, model.Dimensions, model.Speed)
//	}
func (s *EmbeddingsService) ListModels(ctx context.Context) ([]EmbeddingModel, error) {
	var result []EmbeddingModel

	err := s.client.makeRequest(ctx, "GET", "/api/v1/embeddings/models", nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// HealthCheck checks embedding service health.
//
// Verifies that the Railway HuggingFace embedding service is running
// and responding correctly.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//
// Returns:
//   - HealthCheckResponse containing service health status
//   - error if API request fails
//
// Example:
//
//	health, err := client.ZeroDB.Embeddings.HealthCheck(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Service status: %s\n", health.Status)
//	fmt.Printf("Model loaded: %v\n", health.EmbeddingService["model_loaded"])
func (s *EmbeddingsService) HealthCheck(ctx context.Context) (*HealthCheckResponse, error) {
	var result HealthCheckResponse

	err := s.client.makeRequest(ctx, "GET", "/api/v1/embeddings/health", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUsage gets embedding usage statistics.
//
// Returns embedding generation statistics for the current user.
// Since embeddings are FREE (self-hosted), cost is always $0.00.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//
// Returns:
//   - UsageResponse containing usage statistics
//   - error if API request fails
//
// Example:
//
//	usage, err := client.ZeroDB.Embeddings.GetUsage(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Embeddings today: %d\n", usage.EmbeddingsGeneratedToday)
//	fmt.Printf("Embeddings this month: %d\n", usage.EmbeddingsGeneratedMonth)
//	fmt.Printf("Total cost: $%.2f\n", usage.CostMonthUSD) // Always $0.00
func (s *EmbeddingsService) GetUsage(ctx context.Context) (*UsageResponse, error) {
	var result UsageResponse

	err := s.client.makeRequest(ctx, "GET", "/api/v1/embeddings/usage", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
