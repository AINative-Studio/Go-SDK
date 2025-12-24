# ZeroDB Embeddings Example

This example demonstrates how to use the ZeroDB Embeddings API with the Go SDK.

## Prerequisites

1. Go 1.21 or higher
2. AINative API key
3. (Optional) ZeroDB project ID for embed-and-store and semantic search examples

## Setup

1. Set your API key as an environment variable:

```bash
export AINATIVE_API_KEY="your-api-key-here"
```

2. (Optional) Set your ZeroDB project ID to run all examples:

```bash
export ZERODB_PROJECT_ID="550e8400-e29b-41d4-a716-446655440000"
```

## Running the Example

```bash
cd /Users/aideveloper/core/developer-tools/sdks/go/examples/embeddings
go run main.go
```

## What This Example Demonstrates

### 1. Generate Embeddings
Generates vector embeddings for a list of texts using the BAAI/bge-small-en-v1.5 model (384 dimensions).

### 2. List Available Models
Lists all available embedding models with their specifications (dimensions, speed, cost).

### 3. Health Check
Checks the health status of the Railway HuggingFace embedding service.

### 4. Get Usage Statistics
Retrieves embedding generation statistics (always FREE - $0.00).

### 5. Embed and Store (requires project ID)
Generates embeddings and automatically stores them in your ZeroDB vector database in a single operation.

### 6. Semantic Search (requires project ID)
Performs natural language semantic search by automatically embedding the query and searching your vector database.

## Features

- **FREE Embeddings**: Uses Railway self-hosted HuggingFace service (no costs)
- **High Quality**: BAAI/bge-small-en-v1.5 model (384 dimensions)
- **Fast**: Optimized for speed and efficiency
- **One-Step Workflows**: Combine embedding + storage in single API calls
- **Natural Language Search**: Text-based semantic search without manual embedding

## API Endpoints

All operations use the Railway HuggingFace service via these endpoints:

- `POST /api/v1/embeddings/generate` - Generate embeddings
- `POST /api/v1/embeddings/embed-and-store` - Embed and store in one step
- `POST /api/v1/embeddings/semantic-search` - Text-based semantic search
- `GET /api/v1/embeddings/models` - List available models
- `GET /api/v1/embeddings/health` - Check service health
- `GET /api/v1/embeddings/usage` - Get usage statistics

## Example Output

```
=== Example 1: Generate Embeddings ===
Generated 3 embeddings
Model: BAAI/bge-small-en-v1.5
Dimensions: 384
Processing time: 125.50 ms
Cost: $0.00 (FREE!)
First embedding (first 5 dims): [0.1234, -0.5678, 0.9012, -0.3456, 0.7890...]

=== Example 2: List Available Models ===
Available models: 1
  - BAAI/bge-small-en-v1.5
    Dimensions: 384
    Description: Modern, efficient, high quality (DEFAULT)
    Speed: ⚡⚡⚡ Very Fast
    Loaded: true
    Cost per 1k: $0.0000

=== Example 3: Health Check ===
Status: healthy
URL: http://embedding-service.railway.internal:8000
Cost per embedding: $0.0000
Embedding service status: running
Model loaded: true

=== Example 4: Get Usage Statistics ===
User ID: 550e8400-e29b-41d4-a716-446655440000
Embeddings generated today: 1000
Embeddings generated this month: 50000
Cost today: $0.00
Cost this month: $0.00
Model: BAAI/bge-small-en-v1.5
Service: Railway Self-Hosted (FREE)

=== Example 5: Embed and Store ===
Success: true
Vectors stored: 5
Embeddings generated: 5
Model: BAAI/bge-small-en-v1.5
Dimensions: 384
Namespace: tutorials
Processing time: 350.20 ms

=== Example 6: Semantic Search ===
Query: How do I build neural networks for image classification?
Total results: 1
Model: BAAI/bge-small-en-v1.5
Processing time: 200.50 ms

Results:
1. [Similarity: 0.8756] Computer vision and convolutional neural networks
   Metadata: map[category:cv difficulty:intermediate topic:cnn]
```

## Error Handling

The SDK includes comprehensive error handling with validation errors and API errors:

```go
resp, err := client.ZeroDB.Embeddings.Generate(ctx, texts, "", true)
if err != nil {
    if validationErr, ok := err.(*ainative.ValidationError); ok {
        fmt.Printf("Validation error: %s\n", validationErr.Message)
    } else if apiErr, ok := err.(*ainative.APIError); ok {
        fmt.Printf("API error: %s (status: %d)\n", apiErr.Message, apiErr.StatusCode)
    } else {
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}
```

## Learn More

- [GO SDK Documentation](https://pkg.go.dev/github.com/ainative/go-sdk)
- [ZeroDB Documentation](https://docs.ainative.studio/zerodb)
- [Embeddings API Guide](https://docs.ainative.studio/embeddings)
- [BAAI/bge-small-en-v1.5 Model](https://huggingface.co/BAAI/bge-small-en-v1.5)
