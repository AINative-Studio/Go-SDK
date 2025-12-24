/*
AINative GO SDK Command Line Interface

Provides comprehensive CLI for interacting with AINative Studio APIs using GO.
Built with Cobra framework for enterprise-grade CLI experience.
*/

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	apiKey      string
	apiSecret   string
	baseURL     string
	orgID       string
	verbose     bool
	outputFormat string

	// Version information
	version = "1.1.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ainative",
	Short: "AINative CLI - Unified database and AI operations",
	Long: `AINative Command Line Interface

A unified CLI for AINative Studio platform including:
  • ZeroDB operations (vectors, tables, memory)
  • Agent Swarm orchestration
  • Analytics and monitoring
  • Embeddings generation (FREE)
  • Quantum operations`,
	Version: version,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "AINative API key (or set AINATIVE_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&apiSecret, "api-secret", "", "API secret (or set AINATIVE_API_SECRET)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "https://api.ainative.studio", "API base URL")
	rootCmd.PersistentFlags().StringVar(&orgID, "org-id", "", "Organization ID")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|table|yaml)")

	// Read from environment variables
	if apiKey == "" {
		apiKey = os.Getenv("AINATIVE_API_KEY")
	}
	if apiSecret == "" {
		apiSecret = os.Getenv("AINATIVE_API_SECRET")
	}
	if os.Getenv("AINATIVE_BASE_URL") != "" {
		baseURL = os.Getenv("AINATIVE_BASE_URL")
	}
	if orgID == "" {
		orgID = os.Getenv("AINATIVE_ORG_ID")
	}
}
