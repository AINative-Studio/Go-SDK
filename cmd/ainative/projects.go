package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ainative/go-sdk/ainative"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "ZeroDB project management",
	Long:  `Manage ZeroDB projects for organizing data, vectors, and AI operations.`,
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		projects, err := client.ZeroDB.Projects.List(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing projects: %v\n", err)
			os.Exit(1)
		}

		printOutput(projects)
	},
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		name := args[0]
		description, _ := cmd.Flags().GetString("description")

		ctx := context.Background()
		project, err := client.ZeroDB.Projects.Create(ctx, &ainative.ProjectCreateRequest{
			Name:        name,
			Description: description,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Created project: %s\n", project.ID)
		printOutput(project)
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get [project-id]",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		projectID := args[0]
		ctx := context.Background()

		project, err := client.ZeroDB.Projects.Get(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting project: %v\n", err)
			os.Exit(1)
		}

		printOutput(project)
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete [project-id]",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]

		confirm, _ := cmd.Flags().GetBool("yes")
		if !confirm {
			fmt.Printf("Are you sure you want to delete project %s? (y/N): ", projectID)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled.")
				return
			}
		}

		client, err := getClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		err = client.ZeroDB.Projects.Delete(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Deleted project: %s\n", projectID)
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)

	// Flags for create
	projectsCreateCmd.Flags().StringP("description", "d", "", "Project description")

	// Flags for delete
	projectsDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation")
}
