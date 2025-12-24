# AINative GO CLI

Official command-line interface for AINative Studio built with GO and Cobra.

## Features

- ✅ **Fast** - Single binary, instant startup
- ✅ **Cross-platform** - Works on macOS, Linux, Windows
- ✅ **Zero dependencies** - Standalone binary
- ✅ **Complete** - All AINative operations accessible
- ✅ **Enterprise-grade** - Built with Cobra (used by Kubernetes, Docker, GitHub)

## Installation

### From Source

```bash
cd /Users/aideveloper/core/developer-tools/sdks/go/cmd/ainative
go build -o ainative .
sudo mv ainative /usr/local/bin/
```

### Using go install

```bash
go install github.com/ainative/go-sdk/cmd/ainative@latest
```

### Binary Release (Coming Soon)

Download from GitHub releases for your platform.

## Quick Start

### 1. Configure API Key

```bash
export AINATIVE_API_KEY="your-api-key"

# Or set it inline
ainative --api-key your-api-key config show
```

### 2. Basic Commands

```bash
# Show configuration
ainative config show

# List projects
ainative projects list

# Create a project
ainative projects create "My Project" --description "Test project"

# Generate embeddings (FREE)
ainative embeddings generate "Hello world" "How are you?"

# Semantic search
ainative embeddings search "machine learning" --project-id proj_123

# Vector search
ainative vectors search proj_123 0.1 0.2 0.3 --top-k 5
```

## Available Commands

### Configuration

```bash
ainative config show              # Show current configuration
ainative config set api-key KEY   # Set API key
ainative config set base-url URL  # Set custom API URL
```

### Projects

```bash
ainative projects list                  # List all projects
ainative projects create NAME           # Create new project
ainative projects get PROJECT_ID        # Get project details
ainative projects delete PROJECT_ID     # Delete project
```

### Embeddings (FREE)

```bash
ainative embeddings generate TEXT...            # Generate embeddings
ainative embeddings generate --model MODEL      # Use specific model
ainative embeddings search QUERY                # Semantic search
ainative embeddings models                      # List available models
ainative embeddings health                      # Check service health
```

**Available Models:**
- `BAAI/bge-small-en-v1.5` (default, 384 dimensions)
- `BAAI/bge-base-en-v1.5` (768 dimensions)
- `BAAI/bge-large-en-v1.5` (1024 dimensions)

### Vectors

```bash
ainative vectors search PROJECT_ID VECTOR...    # Search vectors
ainative vectors stats PROJECT_ID               # Get vector statistics
```

## Global Flags

```bash
--api-key string      AINative API key
--api-secret string   API secret (optional)
--base-url string     API base URL (default: https://api.ainative.studio)
--org-id string       Organization ID
-v, --verbose         Verbose output
-o, --output string   Output format (json|yaml|table) (default: json)
```

## Examples

### Generate Embeddings and Search

```bash
# Generate embeddings for multiple texts
ainative embeddings generate \
  "Machine learning is amazing" \
  "AI is transforming industries" \
  "Deep learning powers modern AI" \
  --model BAAI/bge-base-en-v1.5

# Perform semantic search
ainative embeddings search "artificial intelligence" \
  --project-id proj_abc123 \
  --limit 10 \
  --threshold 0.8 \
  -o json
```

### Project Management Workflow

```bash
# Create project
ainative projects create "Production App" \
  --description "Main production environment"

# Get project ID from output, then use it
PROJECT_ID="proj_xyz789"

# Generate and store embeddings
ainative embeddings generate "Important data" | \
  jq -r '.embeddings[0]' > embedding.json

# Search similar content
ainative embeddings search "related content" \
  --project-id $PROJECT_ID \
  --limit 5
```

### Verbose Output for Debugging

```bash
# Enable verbose mode
ainative -v embeddings health

# Combine with output format
ainative -v -o yaml projects list
```

## Environment Variables

```bash
export AINATIVE_API_KEY="your-api-key"
export AINATIVE_API_SECRET="your-api-secret"  # Optional
export AINATIVE_BASE_URL="https://api.ainative.studio"
export AINATIVE_ORG_ID="your-org-id"  # For multi-tenant
```

## Output Formats

### JSON (default)

```bash
ainative projects list -o json
```

```json
{
  "projects": [
    {
      "id": "proj_123",
      "name": "My Project",
      "status": "active"
    }
  ]
}
```

### YAML

```bash
ainative projects list -o yaml
```

```yaml
projects:
  - id: proj_123
    name: My Project
    status: active
```

### Table (Coming Soon)

```bash
ainative projects list -o table
```

## Building from Source

### Prerequisites

- GO 1.21 or higher
- Git

### Build Steps

```bash
# Clone repository
git clone https://github.com/ainative/go-sdk.git
cd go-sdk/cmd/ainative

# Install dependencies
go mod download

# Build
go build -o ainative .

# Install globally
sudo mv ainative /usr/local/bin/

# Verify installation
ainative --version
```

### Build for Multiple Platforms

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ainative-darwin-amd64 .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o ainative-darwin-arm64 .

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o ainative-linux-amd64 .

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o ainative-windows-amd64.exe .
```

## Comparison with Python CLI

| Feature | GO CLI | Python CLI |
|---------|--------|------------|
| **Startup Time** | <10ms | ~200ms |
| **Binary Size** | ~15MB | N/A (requires Python) |
| **Dependencies** | None (standalone) | Python + packages |
| **Installation** | Single binary | pip install |
| **Cross-platform** | Easy (compile once) | Requires Python runtime |
| **Performance** | Native speed | Interpreted |

## Troubleshooting

### Command not found

```bash
# Check if installed
which ainative

# Add to PATH
export PATH=$PATH:/usr/local/bin
```

### API Key Issues

```bash
# Verify API key is set
ainative config show

# Set API key
export AINATIVE_API_KEY="your-key"
```

### Connection Errors

```bash
# Check API health
ainative embeddings health -v

# Use custom base URL
ainative --base-url https://custom.api.com embeddings health
```

## Development

### Adding New Commands

1. Create new file: `cmd/ainative/newfeature.go`
2. Define command using Cobra:

```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "New feature",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

3. Rebuild: `go build`

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
AINATIVE_API_KEY=test go test -tags=integration ./...
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT License - see LICENSE file for details

## Support

- **Documentation**: https://docs.ainative.studio
- **Issues**: https://github.com/ainative/go-sdk/issues
- **Discord**: https://discord.gg/ainative
- **Email**: support@ainative.studio

## Version

Current version: **1.1.0**

```bash
ainative --version
```
