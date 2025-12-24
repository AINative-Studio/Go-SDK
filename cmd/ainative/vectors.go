package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ainative/go-sdk/ainative"
	"github.com/spf13/cobra"
)

var vectorsCmd = &cobra.Command{
	Use:   "vectors",
	Short: "Vector operations",
	Long:  `Manage vector data for semantic search and AI operations.`,
}

var vectorsSearchCmd = &cobra.Command{
	Use:   "search [project-id] [vector-values...]",
	Short: "Search for similar vectors",
	Long:  `Search for similar vectors using a query vector.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		projectID := args[0]
		vectorArgs := args[1:]

		// Convert string arguments to float64 vector
		vector := make([]float64, len(vectorArgs))
		for i, v := range vectorArgs {
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing vector value '%s': %v\n", v, err)
				os.Exit(1)
			}
			vector[i] = val
		}

		topK, _ := cmd.Flags().GetInt("top-k")
		namespace, _ := cmd.Flags().GetString("namespace")

		ctx := context.Background()
		results, err := client.ZeroDB.Vectors.Search(ctx, &ainative.VectorSearchRequest{
			ProjectID: projectID,
			Vector:    vector,
			TopK:      topK,
			Namespace: namespace,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching vectors: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Found %d results\n", len(results.Matches))
		printOutput(results)
	},
}

var vectorsStatsCmd = &cobra.Command{
	Use:   "stats [project-id]",
	Short: "Get vector index statistics",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		projectID := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")

		ctx := context.Background()
		stats, err := client.ZeroDB.Vectors.GetStats(ctx, projectID, namespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vector stats: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Vector Index Statistics:")
		fmt.Println("=======================")
		printOutput(stats)
	},
}

func init() {
	rootCmd.AddCommand(vectorsCmd)
	vectorsCmd.AddCommand(vectorsSearchCmd)
	vectorsCmd.AddCommand(vectorsStatsCmd)

	// Flags for search
	vectorsSearchCmd.Flags().IntP("top-k", "k", 5, "Number of results to return")
	vectorsSearchCmd.Flags().StringP("namespace", "n", "default", "Vector namespace")

	// Flags for stats
	vectorsStatsCmd.Flags().StringP("namespace", "n", "", "Vector namespace")
}
