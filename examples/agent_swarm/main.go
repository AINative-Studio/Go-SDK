package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ainative/go-sdk/ainative"
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

	// Create a test project for agent swarm operations
	fmt.Println("🏗️ Creating Test Project for Agent Swarm Operations...")
	project, err := client.ZeroDB.Projects.Create(ctx, &ainative.CreateProjectRequest{
		Name:        "Agent Swarm Example",
		Description: "Demonstrating agent swarm operations with the Go SDK",
		Metadata: map[string]interface{}{
			"example": "agent_swarm",
			"sdk":     "go",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}

	fmt.Printf("✅ Created project: %s (ID: %s)\n", project.Name, project.ID)
	projectID := project.ID

	// Example 1: List Available Agent Types
	fmt.Println("\n🤖 Listing Available Agent Types...")
	agentTypes, err := client.AgentSwarm.ListAgentTypes(ctx)
	if err != nil {
		log.Fatalf("Failed to list agent types: %v", err)
	}

	fmt.Printf("✅ Available agent types (%d):\n", len(agentTypes))
	for i, agentType := range agentTypes {
		fmt.Printf("   %d. %s\n", i+1, agentType)
	}

	// Example 2: Start a Code Analysis Swarm
	fmt.Println("\n🚀 Starting Code Analysis Agent Swarm...")
	
	swarm, err := client.AgentSwarm.Start(ctx, &ainative.StartSwarmRequest{
		ProjectID: projectID,
		Name:      "Code Analysis Swarm",
		Objective: "Analyze codebase for quality, security, and optimization opportunities",
		Agents: []ainative.AgentConfig{
			{
				Type:  ainative.AgentTypeAnalyzer,
				Count: 2,
				Capabilities: []string{"code_analysis", "pattern_detection"},
				Config: map[string]interface{}{
					"analysis_depth": "comprehensive",
					"language_focus": []string{"go", "python", "javascript"},
				},
			},
			{
				Type:  ainative.AgentTypeSecurityScanner,
				Count: 1,
				Capabilities: []string{"vulnerability_scan", "security_audit"},
				Config: map[string]interface{}{
					"scan_level": "thorough",
					"check_dependencies": true,
				},
			},
			{
				Type:  ainative.AgentTypeOptimizer,
				Count: 1,
				Capabilities: []string{"performance_analysis", "resource_optimization"},
				Config: map[string]interface{}{
					"optimization_level": "aggressive",
					"target_metrics": []string{"speed", "memory"},
				},
			},
		},
		Config: &ainative.SwarmConfig{
			MaxConcurrentTasks: 10,
			TaskTimeout:        30 * time.Minute,
			RetryCount:         3,
			ResourceLimits: &ainative.ResourceLimits{
				MaxCPU:    2.0,
				MaxMemory: 4 * 1024 * 1024 * 1024, // 4GB
				MaxAPICalls: 1000,
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to start swarm: %v", err)
	}

	fmt.Printf("✅ Started swarm: %s (ID: %s)\n", swarm.Name, swarm.ID)
	fmt.Printf("   Status: %s\n", swarm.Status)
	fmt.Printf("   Agents: %d\n", len(swarm.Agents))

	swarmID := swarm.ID

	// Example 3: Monitor Swarm Status
	fmt.Println("\n📊 Monitoring Swarm Status...")
	
	// Wait for swarm to initialize
	time.Sleep(3 * time.Second)
	
	swarmStatus, err := client.AgentSwarm.Get(ctx, swarmID)
	if err != nil {
		log.Printf("Failed to get swarm status: %v", err)
	} else {
		fmt.Printf("✅ Swarm Status: %s\n", swarmStatus.Status)
		fmt.Printf("   Active Agents:\n")
		for i, agent := range swarmStatus.Agents {
			fmt.Printf("     %d. %s (%s) - Status: %s\n", 
				i+1, agent.Type, agent.ID, agent.Status)
		}
		
		if swarmStatus.Metrics != nil {
			fmt.Printf("   Metrics:\n")
			fmt.Printf("     Tasks Completed: %d\n", swarmStatus.Metrics.TasksCompleted)
			fmt.Printf("     Tasks In Progress: %d\n", swarmStatus.Metrics.TasksInProgress)
			fmt.Printf("     Efficiency: %.2f%%\n", swarmStatus.Metrics.Efficiency*100)
		}
	}

	// Example 4: Orchestrate Tasks
	fmt.Println("\n🎯 Orchestrating Analysis Tasks...")
	
	// Task 1: Code Quality Analysis
	task1, err := client.AgentSwarm.Orchestrate(ctx, &ainative.OrchestrationRequest{
		SwarmID:  swarmID,
		Task:     "analyze_code_quality",
		Priority: ainative.TaskPriorityHigh,
		Context: map[string]interface{}{
			"repository_url": "https://github.com/example/project",
			"branch":         "main",
			"focus_areas":    []string{"maintainability", "complexity", "coverage"},
		},
	})
	if err != nil {
		log.Printf("Failed to orchestrate task 1: %v", err)
	} else {
		fmt.Printf("✅ Task 1 orchestrated: %s (ID: %s)\n", "Code Quality Analysis", task1.TaskID)
		fmt.Printf("   Status: %s, Assigned to: %v\n", task1.Status, task1.AssignedTo)
	}

	// Task 2: Security Vulnerability Scan
	task2, err := client.AgentSwarm.Orchestrate(ctx, &ainative.OrchestrationRequest{
		SwarmID:  swarmID,
		Task:     "security_vulnerability_scan",
		Priority: ainative.TaskPriorityCritical,
		Context: map[string]interface{}{
			"scan_type":    "comprehensive",
			"check_deps":   true,
			"check_config": true,
		},
	})
	if err != nil {
		log.Printf("Failed to orchestrate task 2: %v", err)
	} else {
		fmt.Printf("✅ Task 2 orchestrated: %s (ID: %s)\n", "Security Scan", task2.TaskID)
		fmt.Printf("   Status: %s, Assigned to: %v\n", task2.Status, task2.AssignedTo)
	}

	// Task 3: Performance Optimization
	task3, err := client.AgentSwarm.Orchestrate(ctx, &ainative.OrchestrationRequest{
		SwarmID:  swarmID,
		Task:     "performance_optimization_analysis",
		Priority: ainative.TaskPriorityMedium,
		Context: map[string]interface{}{
			"target_metrics": []string{"response_time", "memory_usage", "cpu_utilization"},
			"optimization_level": "moderate",
		},
	})
	if err != nil {
		log.Printf("Failed to orchestrate task 3: %v", err)
	} else {
		fmt.Printf("✅ Task 3 orchestrated: %s (ID: %s)\n", "Performance Analysis", task3.TaskID)
		fmt.Printf("   Status: %s, Assigned to: %v\n", task3.Status, task3.AssignedTo)
	}

	// Example 5: Monitor Task Progress
	fmt.Println("\n⏳ Monitoring Task Progress...")
	
	tasks := []struct {
		name   string
		taskID string
	}{
		{"Code Quality Analysis", task1.TaskID},
		{"Security Scan", task2.TaskID},
		{"Performance Analysis", task3.TaskID},
	}

	// Check task status multiple times
	for round := 1; round <= 3; round++ {
		fmt.Printf("\n--- Round %d ---\n", round)
		
		for _, task := range tasks {
			if task.taskID == "" {
				continue
			}
			
			taskStatus, err := client.AgentSwarm.GetTask(ctx, task.taskID)
			if err != nil {
				fmt.Printf("❌ %s: Failed to get status - %v\n", task.name, err)
				continue
			}
			
			fmt.Printf("📋 %s: %s", task.name, taskStatus.Status)
			if taskStatus.CompletedAt != nil {
				elapsed := taskStatus.CompletedAt.Sub(taskStatus.CreatedAt)
				fmt.Printf(" (completed in %v)", elapsed)
			} else if taskStatus.StartedAt != nil {
				elapsed := time.Since(*taskStatus.StartedAt)
				fmt.Printf(" (running for %v)", elapsed)
			}
			fmt.Println()
		}
		
		if round < 3 {
			time.Sleep(5 * time.Second)
		}
	}

	// Example 6: Get Detailed Swarm Metrics
	fmt.Println("\n📈 Getting Detailed Swarm Metrics...")
	
	metrics, err := client.AgentSwarm.GetSwarmMetrics(ctx, swarmID)
	if err != nil {
		log.Printf("Failed to get swarm metrics: %v", err)
	} else {
		fmt.Printf("✅ Swarm Performance Metrics:\n")
		fmt.Printf("   Tasks Completed: %d\n", metrics.TasksCompleted)
		fmt.Printf("   Tasks Failed: %d\n", metrics.TasksFailed)
		fmt.Printf("   Tasks In Progress: %d\n", metrics.TasksInProgress)
		fmt.Printf("   Average Task Time: %v\n", metrics.AverageTaskTime)
		fmt.Printf("   Total Execution Time: %v\n", metrics.TotalExecutionTime)
		fmt.Printf("   Efficiency: %.2f%%\n", metrics.Efficiency*100)
		
		if metrics.ResourceUsage != nil {
			fmt.Printf("   Resource Usage:\n")
			fmt.Printf("     CPU: %.2f%%\n", metrics.ResourceUsage.CPUUsage*100)
			fmt.Printf("     Memory: %.2f MB\n", float64(metrics.ResourceUsage.MemoryUsage)/(1024*1024))
			fmt.Printf("     API Calls: %d\n", metrics.ResourceUsage.APICallsCounter)
		}
	}

	// Example 7: List All Swarms
	fmt.Println("\n📝 Listing All Swarms for Project...")
	
	swarmList, err := client.AgentSwarm.List(ctx, &ainative.ListSwarmsRequest{
		ProjectID: projectID,
		Limit:     10,
	})
	if err != nil {
		log.Printf("Failed to list swarms: %v", err)
	} else {
		fmt.Printf("✅ Found %d swarms:\n", len(swarmList.Swarms))
		for i, s := range swarmList.Swarms {
			fmt.Printf("   %d. %s (ID: %s, Status: %s)\n", 
				i+1, s.Name, s.ID, s.Status)
			fmt.Printf("      Objective: %s\n", s.Objective)
			fmt.Printf("      Agents: %d\n", len(s.Agents))
		}
	}

	// Example 8: Pause and Resume Swarm
	fmt.Println("\n⏸️ Pausing Swarm...")
	err = client.AgentSwarm.Pause(ctx, swarmID)
	if err != nil {
		log.Printf("Failed to pause swarm: %v", err)
	} else {
		fmt.Printf("✅ Swarm paused successfully\n")
		
		time.Sleep(2 * time.Second)
		
		fmt.Println("\n▶️ Resuming Swarm...")
		err = client.AgentSwarm.Resume(ctx, swarmID)
		if err != nil {
			log.Printf("Failed to resume swarm: %v", err)
		} else {
			fmt.Printf("✅ Swarm resumed successfully\n")
		}
	}

	// Clean up
	fmt.Println("\n🛑 Stopping Swarm...")
	err = client.AgentSwarm.Stop(ctx, swarmID)
	if err != nil {
		log.Printf("Failed to stop swarm: %v", err)
	} else {
		fmt.Printf("✅ Swarm stopped successfully\n")
	}

	fmt.Println("\n🧹 Cleaning up project...")
	err = client.ZeroDB.Projects.Suspend(ctx, projectID, "Agent swarm example completed")
	if err != nil {
		log.Printf("Failed to suspend project: %v", err)
	} else {
		fmt.Printf("✅ Project suspended successfully\n")
	}

	fmt.Println("\n🎉 Agent swarm example completed successfully!")
}