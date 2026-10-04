package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/AINative-Studio/Go-SDK/ainative"
)

func main() {
	// Create client with API key from environment
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		log.Fatal("AINATIVE_API_KEY environment variable is required")
	}

	client, err := ainative.NewClient(&ainative.Config{
		APIKey: apiKey,
		Debug:  true,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Create a test project for vectors
	fmt.Println("🏗️ Creating Test Project for Vector Operations...")
	project, err := client.ZeroDB.Projects.Create(ctx, &ainative.CreateProjectRequest{
		Name:        "Vector Operations Example",
		Description: "Demonstrating vector operations with the Go SDK",
		Metadata: map[string]interface{}{
			"example": "vectors",
			"sdk":     "go",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}

	fmt.Printf("✅ Created project: %s (ID: %s)\n", project.Name, project.ID)
	projectID := project.ID

	// Example 1: Generate and Batch-Upsert Vectors
	fmt.Println("\n📊 Generating and Upserting Vectors...")

	// Generate sample vectors (simulating embeddings)
	vectors := generateSampleVectors(10, 384) // 384-dimensional vectors

	vectorItems := make([]ainative.VectorItem, len(vectors))
	for i, vector := range vectors {
		vectorItems[i] = ainative.VectorItem{
			VectorEmbedding: vector,
			Namespace:       "default",
			Document:        fmt.Sprintf("doc_%d", i),
			Metadata: map[string]interface{}{
				"category": fmt.Sprintf("category_%d", i%3),
			},
		}
	}

	upsertResp, err := client.ZeroDB.Vectors.Upsert(ctx, projectID, vectorItems)
	if err != nil {
		log.Fatalf("Failed to upsert vectors: %v", err)
	}

	fmt.Printf("✅ Upserted %d/%d vectors (errors: %d)\n",
		upsertResp.SuccessCount, len(vectorItems), upsertResp.ErrorCount)

	// Wait a moment for indexing
	fmt.Println("⏳ Waiting for vector indexing...")
	time.Sleep(2 * time.Second)

	// Example 2: Vector Similarity Search
	fmt.Println("\n🔍 Performing Vector Similarity Search...")

	// Use the first vector as query
	queryVector := vectors[0]

	searchResp, err := client.ZeroDB.Vectors.Search(ctx, projectID, &ainative.VectorSearchRequest{
		QueryVector: queryVector,
		Limit:       5,
		Namespace:   "default",
	})
	if err != nil {
		log.Fatalf("Failed to search vectors: %v", err)
	}

	fmt.Printf("✅ Found %d similar vectors:\n", len(searchResp.Vectors))
	for i, match := range searchResp.Vectors {
		fmt.Printf("   %d. ID: %s, Similarity: %.6f\n", i+1, match.VectorID, match.Similarity)
		if match.Metadata != nil {
			fmt.Printf("      Category: %v, Document: %v\n",
				match.Metadata["category"], match.Document)
		}
	}

	// Example 3: Filtered Vector Search
	fmt.Println("\n🎯 Performing Filtered Vector Search...")

	filteredSearchResp, err := client.ZeroDB.Vectors.Search(ctx, projectID, &ainative.VectorSearchRequest{
		QueryVector: queryVector,
		Limit:       3,
		Namespace:   "default",
		MetadataFilter: map[string]interface{}{
			"category": "category_0", // Filter by specific category
		},
	})
	if err != nil {
		log.Printf("Failed to perform filtered search: %v", err)
	} else {
		fmt.Printf("✅ Found %d vectors with category='category_0':\n", len(filteredSearchResp.Vectors))
		for i, match := range filteredSearchResp.Vectors {
			fmt.Printf("   %d. ID: %s, Similarity: %.6f, Category: %v\n",
				i+1, match.VectorID, match.Similarity, match.Metadata["category"])
		}
	}

	// Example 4: Batch Vector Operations
	fmt.Println("\n📦 Performing Batch Vector Operations...")

	// Generate more vectors for batch operations
	batchVectors := generateSampleVectors(20, 384)
	batchItems := make([]ainative.VectorItem, len(batchVectors))

	for i, vector := range batchVectors {
		batchItems[i] = ainative.VectorItem{
			VectorEmbedding: vector,
			Namespace:       "batch",
			Document:        fmt.Sprintf("batch_doc_%d", i),
			Metadata: map[string]interface{}{
				"batch": "batch_1",
				"index": i,
				"type":  "batch_vector",
			},
		}
	}

	batchUpsertResp, err := client.ZeroDB.Vectors.Upsert(ctx, projectID, batchItems)
	if err != nil {
		log.Printf("Failed to upsert batch vectors: %v", err)
	} else {
		fmt.Printf("✅ Upserted %d/%d batch vectors (errors: %d)\n",
			batchUpsertResp.SuccessCount, len(batchItems), batchUpsertResp.ErrorCount)
	}

	// Example 5: Multiple Vector Searches (Concurrent)
	fmt.Println("\n🚀 Performing Concurrent Vector Searches...")

	// Create multiple search queries
	searchQueries := make([][]float64, 5)
	for i := range searchQueries {
		searchQueries[i] = generateRandomVector(384)
	}

	// Perform concurrent searches
	results := make(chan searchResult, len(searchQueries))

	for i, query := range searchQueries {
		go func(index int, queryVec []float64) {
			resp, err := client.ZeroDB.Vectors.Search(ctx, projectID, &ainative.VectorSearchRequest{
				QueryVector: queryVec,
				Limit:       3,
				Namespace:   "default",
			})
			results <- searchResult{index: index, response: resp, err: err}
		}(i, query)
	}

	// Collect results
	for i := 0; i < len(searchQueries); i++ {
		result := <-results
		if result.err != nil {
			fmt.Printf("   Search %d failed: %v\n", result.index, result.err)
		} else {
			fmt.Printf("   Search %d: Found %d matches, best similarity: %.6f\n",
				result.index, len(result.response.Vectors),
				getBestSimilarity(result.response.Vectors))
		}
	}

	// Example 6: Performance Demonstration
	fmt.Println("\n⚡ Performance Test - Batch Searches...")

	startTime := time.Now()
	numSearches := 10

	for i := 0; i < numSearches; i++ {
		queryVec := generateRandomVector(384)
		_, err := client.ZeroDB.Vectors.Search(ctx, projectID, &ainative.VectorSearchRequest{
			QueryVector: queryVec,
			Limit:       5,
			Namespace:   "default",
		})
		if err != nil {
			fmt.Printf("   Search %d failed: %v\n", i, err)
		}
	}

	elapsed := time.Since(startTime)
	avgTime := elapsed / time.Duration(numSearches)

	fmt.Printf("✅ Completed %d searches in %v (avg: %v per search)\n",
		numSearches, elapsed, avgTime)

	// Clean up
	fmt.Println("\n🧹 Cleaning up...")
	err = client.ZeroDB.Projects.Delete(ctx, projectID)
	if err != nil {
		log.Printf("Failed to delete project: %v", err)
	} else {
		fmt.Printf("✅ Project deleted successfully\n")
	}

	fmt.Println("\n🎉 Vector operations example completed successfully!")
}

// Helper types and functions

type searchResult struct {
	index    int
	response *ainative.VectorSearchResponse
	err      error
}

func generateSampleVectors(count, dimensions int) [][]float64 {
	vectors := make([][]float64, count)

	for i := range vectors {
		vectors[i] = generateRandomVector(dimensions)
	}

	return vectors
}

func generateRandomVector(dimensions int) []float64 {
	vector := make([]float64, dimensions)

	for j := range vector {
		vector[j] = rand.NormFloat64() // Normal distribution
	}

	// Normalize the vector
	norm := 0.0
	for _, v := range vector {
		norm += v * v
	}
	norm = 1.0 / (norm + 1e-8) // Add small epsilon to avoid division by zero

	for j := range vector {
		vector[j] *= norm
	}

	return vector
}

func getBestSimilarity(matches []ainative.VectorMatch) float64 {
	if len(matches) == 0 {
		return 0.0
	}

	best := matches[0].Similarity
	for _, match := range matches[1:] {
		if match.Similarity > best {
			best = match.Similarity
		}
	}

	return best
}
