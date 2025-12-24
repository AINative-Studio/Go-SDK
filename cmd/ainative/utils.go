package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ainative/go-sdk/ainative"
	"gopkg.in/yaml.v3"
)

// getClient creates and returns an AINative client instance
func getClient() (*ainative.Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required. Set AINATIVE_API_KEY environment variable or use --api-key flag")
	}

	config := &ainative.ClientConfig{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
		OrgID:     orgID,
	}

	return ainative.NewClient(config)
}

// printOutput prints the output in the specified format
func printOutput(v interface{}) {
	switch outputFormat {
	case "json":
		printJSON(v)
	case "yaml":
		printYAML(v)
	case "table":
		printTable(v)
	default:
		printJSON(v)
	}
}

// printJSON prints output in JSON format
func printJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
	}
}

// printYAML prints output in YAML format
func printYAML(v interface{}) {
	encoder := yaml.NewEncoder(os.Stdout)
	if err := encoder.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding YAML: %v\n", err)
	}
}

// printTable prints output in simple table format
func printTable(v interface{}) {
	// For now, fall back to JSON
	// TODO: Implement proper table formatting using tablewriter
	printJSON(v)
}
