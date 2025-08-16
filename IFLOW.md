# IFLOW.md

This file provides guidance to iFlow Cli when working with code in this repository.

## Project Overview

This repository is a GitHub Action that wraps the iFlow CLI, enabling users to run iFlow commands within GitHub workflows. The action is built with Go and packaged in a Docker container that includes Node.js 22 and the iFlow CLI.

## Architecture

1. **Main Components**:
   - `main.go`: Entry point that executes the root command
   - `cmd/root.go`: Core logic implementing the Cobra command structure
   - `Dockerfile`: Multi-stage build creating a runtime with Node.js 22, npm, and uv
   - `action.yml`: GitHub Action definition with inputs/outputs

2. **Key Features**:
   - Docker-based action with pre-installed Node.js 22 and npm
   - Configurable authentication with iFlow API
   - Support for custom models and API endpoints
   - Flexible command execution with timeout control
   - GitHub Actions Summary integration for rich execution reports

## Development Commands

### Building

```bash
# Build the Go binary
go build -o iflow-action .

# Build Docker image
docker build -t iflow-cli-action .
```

The Docker image includes Node.js 22, npm, and uv (ultra-fast Python package manager).

### Testing

```bash
# Run Go tests (if any exist)
go test ./...
```

### Running

```bash
# Run with CLI flags
./iflow-action --prompt "Analyze this code" --api-key <key> --use-env-vars=false

# Run in GitHub Actions mode (uses INPUT_* environment variables)
INPUT_PROMPT="Analyze this code" INPUT_API_KEY=<key> ./iflow-action --use-env-vars=true
```

## Code Structure

1. **Configuration Handling**:
   - Supports two modes: CLI flags and GitHub Actions environment variables
   - Configuration loaded from environment variables when `INPUT_*` prefixed variables are detected
   - Validates required inputs (prompt, API key) and timeout range

2. **iFlow CLI Integration**:
   - Pre-installs iFlow CLI in Docker image via npm
   - Automatically configures iFlow settings in `~/.iflow/settings.json`
   - Executes iFlow commands with `--prompt` and `--yolo` flags for non-interactive operation

3. **GitHub Actions Features**:
   - Sets action outputs (`result`, `exit_code`) via `GITHUB_OUTPUT`
   - Generates rich execution summaries via `GITHUB_STEP_SUMMARY`
   - Provides proper error handling with `::error::` and `::notice::` annotations

## Common Development Tasks

1. **Adding New Configuration Options**:
   - Add flag in `cmd/root.go` init function
   - Update `Config` struct
   - Add to `loadConfigFromEnv` function for GitHub Actions support
   - Update `action.yml` with new input definition

2. **Modifying iFlow Execution**:
   - Update `executeIFlow` function in `cmd/root.go`
   - Modify command arguments as needed

3. **Enhancing GitHub Summary Output**:
   - Update `generateSummaryMarkdown` function in `cmd/root.go`
   - Add new sections or metrics as needed

## Testing Changes

1. **Local Testing**:

   ```bash
   # Build and test with CLI flags
   go build -o iflow-action .
   ./iflow-action --prompt "Test prompt" --api-key test-key --use-env-vars=false
   
   # Test with environment variables
   INPUT_PROMPT="Test prompt" INPUT_API_KEY=test-key ./iflow-action --use-env-vars=true
   
   # Test with precmd
   INPUT_PROMPT="Test prompt" INPUT_API_KEY=test-key INPUT_PRECMD="echo 'Running pre-command'" ./iflow-action --use-env-vars=true
   ```

2. **Docker Testing**:

   ```bash
   # Build Docker image
   docker build -t iflow-cli-action .
   
   # Run with environment variables
   docker run -e INPUT_PROMPT="Test prompt" -e INPUT_API_KEY=test-key iflow-cli-action
   
   # Run with precmd
   docker run -e INPUT_PROMPT="Test prompt" -e INPUT_API_KEY=test-key -e INPUT_PRECMD="echo 'Running pre-command'" iflow-cli-action
   ```