package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ainative/go-sdk/ainative"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		log.Fatal("AINATIVE_API_KEY environment variable is required")
	}

	// Create client
	config := &ainative.Config{
		APIKey:  apiKey,
		BaseURL: "https://api.ainative.studio",
	}

	client, err := ainative.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Generate embeddings
	fmt.Println("=== Example 1: Generate Embeddings ===")
	generateExample(ctx, client)

	// Example 2: List available models
	fmt.Println("\n=== Example 2: List Available Models ===")
	listModelsExample(ctx, client)

	// Example 3: Check service health
	fmt.Println("\n=== Example 3: Health Check ===")
	healthCheckExample(ctx, client)

	// Example 4: Get usage statistics
	fmt.Println("\n=== Example 4: Get Usage Statistics ===")
	getUsageExample(ctx, client)

	// Example 5: Embed and store (requires project ID)
	projectID := os.Getenv("ZERODB_PROJECT_ID")
	if projectID != "" {
		fmt.Println("\n=== Example 5: Embed and Store ===")
		embedAndStoreExample(ctx, client, projectID)

		fmt.Println("\n=== Example 6: Semantic Search ===")
		semanticSearchExample(ctx, client, projectID)
	} else {
		fmt.Println("\n=== Skipping Examples 5 & 6 ===")
		fmt.Println("Set ZERODB_PROJECT_ID environment variable to run embed-and-store and semantic search examples")
	}
}

func generateExample(ctx context.Context, client *ainative.Client) {
	texts := []string{
		"Machine learning is a subset of artificial intelligence",
		"Python is a popular programming language for data science",
		"Neural networks are inspired by biological brains",
	}

	resp, err := client.ZeroDB.Embeddings.Generate(ctx, texts, "BAAI/bge-small-en-v1.5", true)
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		return
	}

	fmt.Printf("Generated %d embeddings\n", resp.Count)
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Dimensions: %d\n", resp.Dimensions)
	fmt.Printf("Processing time: %.2f ms\n", resp.ProcessingTimeMs)
	fmt.Printf("Cost: $%.2f (FREE!)\n", resp.CostUSD)

	// Show first few dimensions of first embedding
	if len(resp.Embeddings) > 0 && len(resp.Embeddings[0]) > 0 {
		fmt.Printf("First embedding (first 5 dims): [")
		for i := 0; i < 5 && i < len(resp.Embeddings[0]); i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%.4f", resp.Embeddings[0][i])
		}
		fmt.Println("...]")
	}
}

func listModelsExample(ctx context.Context, client *ainative.Client) {
	models, err := client.ZeroDB.Embeddings.ListModels(ctx)
	if err != nil {
		log.Printf("Error listing models: %v", err)
		return
	}

	fmt.Printf("Available models: %d\n", len(models))
	for _, model := range models {
		fmt.Printf("  - %s\n", model.ID)
		fmt.Printf("    Dimensions: %d\n", model.Dimensions)
		fmt.Printf("    Description: %s\n", model.Description)
		fmt.Printf("    Speed: %s\n", model.Speed)
		fmt.Printf("    Loaded: %v\n", model.Loaded)
		fmt.Printf("    Cost per 1k: $%.4f\n", model.CostPer1k)
	}
}

func healthCheckExample(ctx context.Context, client *ainative.Client) {
	health, err := client.ZeroDB.Embeddings.HealthCheck(ctx)
	if err != nil {
		log.Printf("Error checking health: %v", err)
		return
	}

	fmt.Printf("Status: %s\n", health.Status)
	fmt.Printf("URL: %s\n", health.URL)
	fmt.Printf("Cost per embedding: $%.4f\n", health.CostPerEmbedding)

	if service, ok := health.EmbeddingService["status"].(string); ok {
		fmt.Printf("Embedding service status: %s\n", service)
	}

	if modelLoaded, ok := health.EmbeddingService["model_loaded"].(bool); ok {
		fmt.Printf("Model loaded: %v\n", modelLoaded)
	}
}

func getUsageExample(ctx context.Context, client *ainative.Client) {
	usage, err := client.ZeroDB.Embeddings.GetUsage(ctx)
	if err != nil {
		log.Printf("Error getting usage: %v", err)
		return
	}

	fmt.Printf("User ID: %s\n", usage.UserID)
	fmt.Printf("Embeddings generated today: %d\n", usage.EmbeddingsGeneratedToday)
	fmt.Printf("Embeddings generated this month: %d\n", usage.EmbeddingsGeneratedMonth)
	fmt.Printf("Cost today: $%.2f\n", usage.CostTodayUSD)
	fmt.Printf("Cost this month: $%.2f\n", usage.CostMonthUSD)
	fmt.Printf("Model: %s\n", usage.Model)
	fmt.Printf("Service: %s\n", usage.Service)
}

func embedAndStoreExample(ctx context.Context, client *ainative.Client, projectID string) {
	texts := []string{
		"Introduction to machine learning fundamentals",
		"Advanced deep learning architectures and techniques",
		"Natural language processing with transformers",
		"Computer vision and convolutional neural networks",
		"Reinforcement learning and Q-learning algorithms",
	}

	metadata := []map[string]interface{}{
		{"category": "ml", "difficulty": "beginner", "topic": "fundamentals"},
		{"category": "dl", "difficulty": "advanced", "topic": "architectures"},
		{"category": "nlp", "difficulty": "intermediate", "topic": "transformers"},
		{"category": "cv", "difficulty": "intermediate", "topic": "cnn"},
		{"category": "rl", "difficulty": "advanced", "topic": "q-learning"},
	}

	resp, err := client.ZeroDB.Embeddings.EmbedAndStore(
		ctx,
		projectID,
		texts,
		metadata,
		"tutorials",
		"BAAI/bge-small-en-v1.5",
	)
	if err != nil {
		log.Printf("Error embedding and storing: %v", err)
		return
	}

	fmt.Printf("Success: %v\n", resp.Success)
	fmt.Printf("Vectors stored: %d\n", resp.VectorsStored)
	fmt.Printf("Embeddings generated: %d\n", resp.EmbeddingsGenerated)
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Dimensions: %d\n", resp.Dimensions)
	fmt.Printf("Namespace: %s\n", resp.Namespace)
	fmt.Printf("Processing time: %.2f ms\n", resp.ProcessingTimeMs)
}

func semanticSearchExample(ctx context.Context, client *ainative.Client, projectID string) {
	query := "How do I build neural networks for image classification?"

	filterMetadata := map[string]interface{}{
		"category": "cv",
	}

	resp, err := client.ZeroDB.Embeddings.SemanticSearch(
		ctx,
		projectID,
		query,
		5,
		0.7,
		"tutorials",
		filterMetadata,
		"BAAI/bge-small-en-v1.5",
	)
	if err != nil {
		log.Printf("Error performing semantic search: %v", err)
		return
	}

	fmt.Printf("Query: %s\n", resp.Query)
	fmt.Printf("Total results: %d\n", resp.TotalResults)
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Processing time: %.2f ms\n", resp.ProcessingTimeMs)

	fmt.Println("\nResults:")
	for i, result := range resp.Results {
		fmt.Printf("%d. [Similarity: %.4f] %s\n", i+1, result.Similarity, result.Document)
		if len(result.Metadata) > 0 {
			fmt.Printf("   Metadata: %v\n", result.Metadata)
		}
	}
}
