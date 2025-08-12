# AINative Go SDK

Official Go SDK for AINative Studio APIs - unified database and AI operations platform.

## Features

- 🚀 **Simple Integration** - Get started with just an API key
- 🗄️ **ZeroDB Operations** - Full support for vector storage, search, and memory management
- 🤖 **Agent Swarm** - Orchestrate multiple AI agents for complex tasks
- 🔐 **Enterprise Security** - Multi-tenant authentication with API key management
- ⚡ **High Performance** - Concurrent request handling and connection pooling
- 📊 **Analytics** - Built-in usage tracking and performance metrics

## Installation

```bash
go get github.com/AINative-Studio/Go-SDK
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/AINative-Studio/Go-SDK/ainative"
)

func main() {
    // Initialize client
    client := ainative.NewClient("your-api-key")
    ctx := context.Background()
    
    // Create a project
    project, err := client.ZeroDB.Projects.Create(ctx, &ainative.CreateProjectRequest{
        Name:        "My First Project",
        Description: "Testing AINative Go SDK",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created project: %s\n", project.ID)
}
```
