# 🤖 iFlow CLI GitHub Action

A GitHub Action that enables you to run [iFlow CLI](https://github.com/iflow-ai/iflow-cli) commands within your GitHub workflows. This Docker-based action comes with Node.js 22 and npm pre-installed for optimal performance, and executes your specified commands using the iFlow CLI.

- [中文文档](README_zh.md)

> Docs Site (generated with iFlow CLI GitHub Action): [https://vibe-ideas.github.io/iflow-cli-action/](https://vibe-ideas.github.io/iflow-cli-action/)

## Features

- ✅ Docker-based action with pre-installed Node.js 22 and npm
- ✅ Configurable authentication with iFlow API
- ✅ Support for custom models and API endpoints
- ✅ Flexible command execution with timeout control
- ✅ Works in any working directory
- ✅ Built with Go for fast, reliable execution
- ✅ **GitHub Actions Summary integration**: Rich execution reports in PR summaries

## Usage

### Basic Example

```yaml
name: iFlow CLI Example
on: [push]

jobs:
  analyze-code:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run iFlow CLI
        uses: vibe-ideas/iflow-cli-action@v1.2.0
        with:
          prompt: "Analyze this codebase and suggest improvements"
          api_key: ${{ secrets.IFLOW_API_KEY }}
```

### Advanced Example

```yaml
name: Advanced iFlow CLI Usage
on: 
  pull_request:
    types: [opened, synchronize]

jobs:
  code-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Initialize Project Analysis
        uses: vibe-ideas/iflow-cli-action@v1.2.0
        with:
          prompt: "/init"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Qwen3-Coder"
          timeout: "600"
          working_directory: "."
      
      - name: Generate Technical Documentation
        uses: vibe-ideas/iflow-cli-action@v1.2.0
        with:
          prompt: "Generate technical documentation based on the codebase analysis"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          base_url: "https://apis.iflow.cn/v1"
          model: "DeepSeek-V3"
        id: docs
      
      - name: Display Results
        run: |
          echo "Documentation generated:"
          echo "${{ steps.docs.outputs.result }}"
```

### Multiple Commands Example

```yaml
name: Multi-step iFlow Analysis
on: [workflow_dispatch]

jobs:
  comprehensive-analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Project Overview
        uses: vibe-ideas/iflow-cli-action@v1.2.0
        with:
          prompt: |
            Analyze the project structure and provide:
            1. Main architectural components
            2. Key dependencies and their purposes
            3. Potential security considerations
          api_key: ${{ secrets.IFLOW_API_KEY }}
          timeout: "900"
      
      - name: Code Quality Assessment
        uses: vibe-ideas/iflow-cli-action@v1.2.0
        with:
          prompt: "Review the code for best practices, potential bugs, and performance improvements"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Kimi-K2"
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `prompt` | The prompt to execute with iFlow CLI | ✅ Yes | - |
| `api_key` | iFlow API key for authentication | ✅ Yes | - |
| `settings_json` | Complete iFlow settings.json content (JSON string). If provided, this will override other configuration options. | ❌ No | - |
| `base_url` | Custom base URL for iFlow API | ❌ No | `https://apis.iflow.cn/v1` |
| `model` | Model name to use | ❌ No | `Qwen3-Coder` |
| `working_directory` | Working directory to run iFlow CLI from | ❌ No | `.` |
| `timeout` | Timeout for iFlow CLI execution in seconds (1-86400) | ❌ No | `86400` |
| `extra_args` | Additional command line arguments to pass to iFlow CLI (space-separated string) | ❌ No | `` |

## Outputs

| Output | Description |
|--------|-------------|
| `result` | Output from iFlow CLI execution |
| `exit_code` | Exit code from iFlow CLI execution |

## Authentication

### Getting an iFlow API Key

1. Register for an iFlow account at [iflow.cn](https://iflow.cn)
2. Go to your profile settings or [click here to get your API key](https://iflow.cn/?open=setting)
3. Click "Reset" in the pop-up dialog to generate a new API key
4. Add the API key to your GitHub repository secrets as `IFLOW_API_KEY`

### Available Models

- `Qwen3-Coder` (default) - Excellent for code analysis and generation
- `Kimi-K2` - Good for general AI tasks and longer contexts
- `DeepSeek-V3` - Advanced reasoning and problem-solving
- Custom models supported via OpenAI-compatible APIs

## Custom Configuration

### Using Extra Arguments

The `extra_args` input allows you to pass additional command-line arguments directly to the iFlow CLI. This provides flexibility to use advanced iFlow CLI features that aren't exposed as dedicated action inputs.

```yaml
- name: iFlow with Custom Arguments
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Analyze this codebase with debug output"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    extra_args: "--debug --max-tokens 3000"
```

#### Examples of Extra Arguments

- `--debug` - Enable debug mode  

#### Quoted Arguments

For arguments containing spaces, use quotes:

```yaml
extra_args: '--debug'
```

### Using Custom Settings

For advanced users who need complete control over the iFlow configuration, you can provide a custom `settings.json` directly:

```yaml
- name: Custom iFlow Configuration
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Analyze this codebase with custom configuration"
    api_key: ${{ secrets.IFLOW_API_KEY }}  # Still required for basic validation
    settings_json: |
      {
        "theme": "Dark",
        "selectedAuthType": "iflow",
        "apiKey": "${{ secrets.IFLOW_API_KEY }}",
        "baseUrl": "https://custom-api.example.com/v1",
        "modelName": "custom-model",
        "searchApiKey": "${{ secrets.SEARCH_API_KEY }}",
        "customField": "customValue"
      }
```

When `settings_json` is provided, it takes precedence over individual configuration inputs (`base_url`, `model`, etc.). This allows you to:

- Use custom authentication types
- Configure additional fields not available as inputs
- Maintain complex configurations across multiple workflow runs
- Support custom API endpoints and models

**Note:** The `api_key` input is still required for validation, but the actual API key used will be the one specified in your `settings_json`.

## Using MCP Servers

[MCP (Model Context Protocol)](https://modelcontextprotocol.io) allows iFlow CLI to connect to external tools and services, extending its capabilities beyond just AI model interactions. You can configure MCP servers in your workflow to enable features like code search, database querying, or custom tool integrations.

### Example: Using DeepWiki MCP Server

The following example demonstrates how to configure and use the DeepWiki MCP server for enhanced code search capabilities:

```yaml
- name: iFlow CLI with MCP Server
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "use @deepwiki to search how to use Skynet to build a game"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    settings_json: |
      {
        "selectedAuthType": "iflow",
        "apiKey": "${{ secrets.IFLOW_API_KEY }}",
        "baseUrl": "https://apis.iflow.cn/v1",
        "modelName": "Qwen3-Coder",
        "searchApiKey": "${{ secrets.IFLOW_API_KEY }}",
        "mcpServers": {
          "deepwiki": {
            "command": "npx",
            "args": ["-y", "mcp-deepwiki@latest"]
          }
        }
      }
    model: "Qwen3-Coder"
    timeout: "1800"
    extra_args: "--debug"
```

In this example:

- The `mcpServers` configuration defines a server named `deepwiki`
- The server is executed using `npx -y mcp-deepwiki@latest`
- The prompt references the server with `@deepwiki` to utilize its capabilities
- The `searchApiKey` is used for authentication with the DeepWiki service

### When to Use MCP Servers

MCP servers are particularly useful when you need:

- Enhanced code search and documentation lookup capabilities
- Integration with external tools and services
- Access to specialized knowledge bases or databases
- Custom tooling that extends the iFlow CLI functionality

## Common Use Cases

### Code Analysis and Review

```yaml
- name: Code Review
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Review this pull request for code quality, security issues, and best practices"
    api_key: ${{ secrets.IFLOW_API_KEY }}
```

### Documentation Generation

```yaml
- name: Generate Documentation
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "/init && Generate comprehensive API documentation"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "600"
```

### Automated Testing Suggestions

```yaml
- name: Test Strategy
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Analyze the codebase and suggest a comprehensive testing strategy with specific test cases"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    model: "DeepSeek-V3"
```

### Architecture Analysis

```yaml
- name: Architecture Review
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Analyze the system architecture and suggest improvements for scalability and maintainability"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "900"
```

### Using MCP Servers

```yaml
- name: iFlow CLI with MCP Server
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "use @deepwiki to search how to use Skynet to build a game"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    settings_json: |
      {
        "selectedAuthType": "iflow",
        "apiKey": "${{ secrets.IFLOW_API_KEY }}",
        "baseUrl": "https://apis.iflow.cn/v1",
        "modelName": "Qwen3-Coder",
        "searchApiKey": "${{ secrets.IFLOW_API_KEY }}",
        "mcpServers": {
          "deepwiki": {
            "command": "npx",
            "args": ["-y", "mcp-deepwiki@latest"]
          }
        }
      }
    model: "Qwen3-Coder"
    timeout: "1800"
    extra_args: "--debug"
```

### Using Extra Arguments

```yaml
- name: iFlow with Custom Arguments
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Analyze this codebase with debug output"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    extra_args: "--debug --max-tokens 3000"
```

### Using Custom Settings

```yaml
- name: Custom iFlow Configuration
  uses: vibe-ideas/iflow-cli-action@v1.2.0
  with:
    prompt: "Analyze this codebase with custom configuration"
    api_key: ${{ secrets.IFLOW_API_KEY }}  # Still required for basic validation
    settings_json: |
      {
        "theme": "Dark",
        "selectedAuthType": "iflow",
        "apiKey": "${{ secrets.IFLOW_API_KEY }}",
        "baseUrl": "https://custom-api.example.com/v1",
        "modelName": "custom-model",
        "searchApiKey": "${{ secrets.SEARCH_API_KEY }}",
        "customField": "customValue"
      }
```

### Complete Workflow Example

```yaml
name: iFlow CLI Code Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  iflow-review:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Review code with iFlow CLI
        uses: vibe-ideas/iflow-cli-action@v1.2.0
        with:
          prompt: "Review this pull request for code quality, security issues, and best practices. Provide specific suggestions for improvement."
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Qwen3-Coder"
          timeout: "600"
        id: review
      
      - name: Comment on PR
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '## 🤖 iFlow CLI Code Review\n\n' + '${{ steps.review.outputs.result }}'
            })
```

## Requirements

- **Runner**: Linux-based GitHub Actions runners (ubuntu-latest recommended)
- **Permissions**: Actions need internet access to download dependencies
- **Resources**: Sufficient timeout for command execution (adjust based on complexity)

## Troubleshooting

### Common Issues

**Command timeout:** Increase the `timeout` value for complex operations

```yaml
timeout: "900"  # 15 minutes
```

**API authentication failed:** Verify your API key is correctly set in repository secrets

**Working directory not found:** Ensure the path exists and checkout action is used

```yaml
- uses: actions/checkout@v4  # Required before using iFlow action
```

### Debug Mode

Enable verbose logging by setting environment variables:

```yaml
env:
  ACTIONS_STEP_DEBUG: true
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related

- [iFlow CLI](https://github.com/iflow-ai/iflow-cli) - The underlying CLI tool
- [iFlow Platform](https://docs.iflow.cn/en/docs) - Official documentation
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
