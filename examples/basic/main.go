package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ainative/go-sdk/ainative"
)

func main() {
	// Create client with API key from environment
	apiKey := os.Getenv("AINATIVE_API_KEY")
	if apiKey == "" {
		log.Fatal("AINATIVE_API_KEY environment variable is required")
	}

	client, err := ainative.NewClient(&ainative.Config{
		APIKey:  apiKey,
		BaseURL: "https://api.ainative.studio", // Optional: defaults to production
		Debug:   true,                          // Optional: enable debug logging
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Check API Health
	fmt.Println("🏥 Checking API Health...")
	health, err := client.Health(ctx)
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Printf("✅ API Status: %s (Version: %s)\n", health.Status, health.Version)
	}

	// Example 2: Get User Information
	fmt.Println("\n👤 Getting User Information...")
	userInfo, err := client.Auth.GetUserInfo(ctx)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
	} else {
		fmt.Printf("✅ User: %s (%s)\n", userInfo.Name, userInfo.Email)
		fmt.Printf("   Organization: %s\n", userInfo.Organization)
		fmt.Printf("   Role: %s\n", userInfo.Role)
	}

	// Example 3: List Projects
	fmt.Println("\n📁 Listing Projects...")
	projects, err := client.ZeroDB.Projects.List(ctx, &ainative.ListProjectsRequest{
		Limit: 5,
	})
	if err != nil {
		log.Printf("Failed to list projects: %v", err)
	} else {
		fmt.Printf("✅ Found %d projects:\n", len(projects.Projects))
		for i, project := range projects.Projects {
			fmt.Printf("   %d. %s (ID: %s, Status: %s)\n", 
				i+1, project.Name, project.ID, project.Status)
			if project.Description != "" {
				fmt.Printf("      Description: %s\n", project.Description)
			}
		}
	}

	// Example 4: Create a New Project
	fmt.Println("\n🆕 Creating New Project...")
	newProject, err := client.ZeroDB.Projects.Create(ctx, &ainative.CreateProjectRequest{
		Name:        "Go SDK Example Project",
		Description: "A test project created using the Go SDK",
		Metadata: map[string]interface{}{
			"sdk":     "go",
			"example": "basic",
			"version": "1.0.0",
		},
	})
	if err != nil {
		log.Printf("Failed to create project: %v", err)
	} else {
		fmt.Printf("✅ Created project: %s (ID: %s)\n", newProject.Name, newProject.ID)
		
		// Use this project for the rest of the examples
		projectID := newProject.ID
		
		// Example 5: Create Memory
		fmt.Println("\n🧠 Creating Memory...")
		memory, err := client.ZeroDB.Memory.Create(ctx, &ainative.CreateMemoryRequest{
			Content:  "This is a test memory created from the Go SDK example",
			Title:    "Go SDK Test Memory",
			Tags:     []string{"go", "sdk", "test"},
			Priority: ainative.MemoryPriorityMedium,
			Metadata: map[string]interface{}{
				"project_id": projectID,
				"example":    "basic",
			},
		})
		if err != nil {
			log.Printf("Failed to create memory: %v", err)
		} else {
			fmt.Printf("✅ Created memory: %s (ID: %s)\n", memory.Title, memory.ID)
		}

		// Example 6: Search Memory
		fmt.Println("\n🔍 Searching Memory...")
		searchResults, err := client.ZeroDB.Memory.Search(ctx, &ainative.SearchMemoryRequest{
			Query:    "test memory",
			Limit:    5,
			Semantic: true,
		})
		if err != nil {
			log.Printf("Failed to search memory: %v", err)
		} else {
			fmt.Printf("✅ Found %d memory items:\n", len(searchResults.Results))
			for i, item := range searchResults.Results {
				fmt.Printf("   %d. %s (Priority: %s)\n", i+1, item.Title, item.Priority)
				fmt.Printf("      Tags: %v\n", item.Tags)
			}
		}

		// Example 7: List Agent Types
		fmt.Println("\n🤖 Listing Agent Types...")
		agentTypes, err := client.AgentSwarm.ListAgentTypes(ctx)
		if err != nil {
			log.Printf("Failed to list agent types: %v", err)
		} else {
			fmt.Printf("✅ Available agent types: %v\n", agentTypes)
		}

		// Example 8: Clean up - Suspend the project
		fmt.Println("\n🧹 Cleaning up - Suspending Project...")
		err = client.ZeroDB.Projects.Suspend(ctx, projectID, "Example completed")
		if err != nil {
			log.Printf("Failed to suspend project: %v", err)
		} else {
			fmt.Printf("✅ Project suspended successfully\n")
		}
	}

	fmt.Println("\n🎉 Basic example completed successfully!")
}