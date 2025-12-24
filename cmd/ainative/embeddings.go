package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ainative/go-sdk/ainative"
	"github.com/spf13/cobra"
)

var embeddingsCmd = &cobra.Command{
	Use:   "embeddings",
	Short: "Embedding generation operations (FREE)",
	Long:  `Generate vector embeddings using FREE Railway HuggingFace service with BAAI BGE models.`,
}

var embeddingsGenerateCmd = &cobra.Command{
	Use:   "generate [texts...]",
	Short: "Generate embeddings for texts",
	Long:  `Generate vector embeddings for one or more texts using BAAI BGE models.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		model, _ := cmd.Flags().GetString("model")
		normalize, _ := cmd.Flags().GetBool("normalize")

		ctx := context.Background()
		result, err := client.ZeroDB.Embeddings.Generate(ctx, &ainative.GenerateRequest{
			Texts:     args,
			Model:     model,
			Normalize: normalize,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating embeddings: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Generated %d embeddings\n", len(result.Embeddings))
		fmt.Printf("Model: %s\n", result.Model)
		fmt.Printf("Dimensions: %d\n", result.Dimensions)

		if verbose {
			printOutput(result)
		}
	},
}

var embeddingsSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Semantic search with automatic embedding",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		projectID, _ := cmd.Flags().GetString("project-id")
		limit, _ := cmd.Flags().GetInt("limit")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		namespace, _ := cmd.Flags().GetString("namespace")

		if projectID == "" {
			fmt.Fprintln(os.Stderr, "Error: --project-id is required")
			os.Exit(1)
		}

		ctx := context.Background()
		result, err := client.ZeroDB.Embeddings.SemanticSearch(ctx, &ainative.SemanticSearchRequest{
			ProjectID: projectID,
			Query:     args[0],
			Limit:     limit,
			Threshold: threshold,
			Namespace: namespace,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing semantic search: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Found %d results\n", len(result.Results))
		printOutput(result)
	},
}

var embeddingsModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available embedding models",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		models, err := client.ZeroDB.Embeddings.ListModels(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing models: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available Embedding Models:")
		fmt.Println("===========================")
		for _, model := range models.Models {
			fmt.Printf("\n%s\n", model.Name)
			fmt.Printf("  Dimensions: %d\n", model.Dimensions)
			fmt.Printf("  Description: %s\n", model.Description)
			fmt.Printf("  Context: %d tokens\n", model.MaxContextLength)
		}
	},
}

var embeddingsHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check embedding service health",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		health, err := client.ZeroDB.Embeddings.HealthCheck(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking health: %v\n", err)
			os.Exit(1)
		}

		if health.Status == "healthy" {
			fmt.Println("✓ Embedding service is healthy")
		} else {
			fmt.Printf("⚠ Embedding service status: %s\n", health.Status)
		}

		if verbose {
			printOutput(health)
		}
	},
}

func init() {
	rootCmd.AddCommand(embeddingsCmd)
	embeddingsCmd.AddCommand(embeddingsGenerateCmd)
	embeddingsCmd.AddCommand(embeddingsSearchCmd)
	embeddingsCmd.AddCommand(embeddingsModelsCmd)
	embeddingsCmd.AddCommand(embeddingsHealthCmd)

	// Flags for generate
	embeddingsGenerateCmd.Flags().String("model", "BAAI/bge-small-en-v1.5", "Embedding model")
	embeddingsGenerateCmd.Flags().Bool("normalize", true, "Normalize vectors")

	// Flags for search
	embeddingsSearchCmd.Flags().String("project-id", "", "Project ID (required)")
	embeddingsSearchCmd.Flags().IntP("limit", "l", 10, "Maximum results")
	embeddingsSearchCmd.Flags().Float64P("threshold", "t", 0.7, "Similarity threshold (0.0-1.0)")
	embeddingsSearchCmd.Flags().StringP("namespace", "n", "default", "Vector namespace")
	embeddingsSearchCmd.MarkFlagRequired("project-id")
}
