package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Manage CLI configuration including API keys, base URL, and organization settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Println("=====================")

		if apiKey != "" {
			fmt.Printf("API Key:          %s\n", maskString(apiKey))
		} else {
			fmt.Println("API Key:          Not set")
		}

		if apiSecret != "" {
			fmt.Println("API Secret:       ***")
		} else {
			fmt.Println("API Secret:       Not set")
		}

		fmt.Printf("Base URL:         %s\n", baseURL)

		if orgID != "" {
			fmt.Printf("Organization ID:  %s\n", orgID)
		} else {
			fmt.Println("Organization ID:  Not set")
		}

		fmt.Printf("Output Format:    %s\n", outputFormat)
		fmt.Printf("Verbose:          %v\n", verbose)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set configuration value",
	Long:  `Set configuration values for API access and CLI behavior.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		switch key {
		case "api-key", "api_key":
			fmt.Printf("Set AINATIVE_API_KEY environment variable to: %s\n", value)
			fmt.Println("Run: export AINATIVE_API_KEY=" + value)
		case "api-secret", "api_secret":
			fmt.Println("Set AINATIVE_API_SECRET environment variable")
			fmt.Println("Run: export AINATIVE_API_SECRET=" + value)
		case "base-url", "base_url":
			fmt.Printf("Set AINATIVE_BASE_URL environment variable to: %s\n", value)
			fmt.Println("Run: export AINATIVE_BASE_URL=" + value)
		case "org-id", "org_id":
			fmt.Printf("Set AINATIVE_ORG_ID environment variable to: %s\n", value)
			fmt.Println("Run: export AINATIVE_ORG_ID=" + value)
		default:
			fmt.Printf("Unknown configuration key: %s\n", key)
			fmt.Println("\nAvailable keys:")
			fmt.Println("  api-key, api_key       - AINative API key")
			fmt.Println("  api-secret, api_secret - API secret")
			fmt.Println("  base-url, base_url     - API base URL")
			fmt.Println("  org-id, org_id         - Organization ID")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "..." + s[len(s)-4:]
}
