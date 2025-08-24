# 🤖 iFlow CLI GitHub Action

A GitHub Action that enables you to run [iFlow CLI](https://github.com/iflow-ai/iflow-cli) commands within your GitHub workflows. This Docker-based action comes with Node.js 22, npm, and uv (ultra-fast Python package manager) pre-installed for optimal performance, and executes your specified commands using the iFlow CLI.

- [中文文档](README_zh.md)

> Docs Site (generated with iFlow CLI GitHub Action): [https://iflow-ai.github.io/iflow-cli-action/](https://iflow-ai.github.io/iflow-cli-action/)

## Features

- ✅ Docker-based action with pre-installed Node.js 22, npm, and uv
- ✅ Configurable authentication with iFlow API
- ✅ Support for custom models and API endpoints
- ✅ Flexible command execution with timeout control
- ✅ Works in any working directory
- ✅ Built with Go for fast, reliable execution
- ✅ **GitHub Actions Summary integration**: Rich execution reports in PR summaries
- ✅ PR/Issue Integration: Works seamlessly with GitHub comments and PR reviews

## Usage

### Basic Example

Issue triage with iFLOW CLI:

```yaml
name: '🏷️ iFLOW CLI Automated Issue Triage'

on:
  issues:
    types:
      - 'opened'
      - 'reopened'
  issue_comment:
    types:
      - 'created'
  workflow_dispatch:
    inputs:
      issue_number:
        description: 'issue number to triage'
        required: true
        type: 'number'

concurrency:
  group: '${{ github.workflow }}-${{ github.event.issue.number }}'
  cancel-in-progress: true

defaults:
  run:
    shell: 'bash'

permissions:
  contents: 'read'
  issues: 'write'
  statuses: 'write'

jobs:
  triage-issue:
    if: |-
      github.event_name == 'issues' ||
      github.event_name == 'workflow_dispatch' ||
      (
        github.event_name == 'issue_comment' &&
        contains(github.event.comment.body, '@iflow-cli /triage') &&
        contains(fromJSON('["OWNER", "MEMBER", "COLLABORATOR"]'), github.event.comment.author_association)
      )
    timeout-minutes: 5
    runs-on: 'ubuntu-latest'
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: 'Run iFlow CLI Issue Triage'
        uses: iflow-ai/iflow-cli-action@v1.3.0
        id: 'iflow_cli_issue_triage'
        env:
          GITHUB_TOKEN: '${{ secrets.GITHUB_TOKEN }}'
          ISSUE_TITLE: '${{ github.event.issue.title }}'
          ISSUE_BODY: '${{ github.event.issue.body }}'
          ISSUE_NUMBER: '${{ github.event.issue.number }}'
          REPOSITORY: '${{ github.repository }}'
        with:
          api_key: ${{ secrets.IFLOW_API_KEY }}
          timeout: "3600"
          extra_args: "--debug"
          prompt: |
            ## Role

            You are an issue triage assistant. Analyze the current GitHub issue
            and apply the most appropriate existing labels. Use the available
            tools to gather information; do not ask for information to be
            provided.

            ## Steps

            1. Run: `gh label list` to get all available labels.
            2. Review the issue title and body provided in the environment
               variables: "${ISSUE_TITLE}" and "${ISSUE_BODY}".
            3. Classify issues by their kind (bug, enhancement, documentation,
               cleanup, etc) and their priority (p0, p1, p2, p3). Set the
               labels according to the format `kind/*` and `priority/*` patterns.
            4. Apply the selected labels to this issue using:
               `gh issue edit "${ISSUE_NUMBER}" --add-label "label1,label2"`
            5. If the "status/needs-triage" label is present, remove it using:
               `gh issue edit "${ISSUE_NUMBER}" --remove-label "status/needs-triage"`

            ## Guidelines

            - Only use labels that already exist in the repository
            - Do not add comments or modify the issue content
            - Triage only the current issue
            - Assign all applicable labels based on the issue content
            - Reference all shell variables as "${VAR}" (with quotes and braces)

      - name: 'Post Issue Triage Failure Comment'
        if: |-
          ${{ failure() && steps.iflow_cli_issue_triage.outcome == 'failure' }}
        uses: 'actions/github-script@60a0d83039c74a4aee543508d2ffcb1c3799cdea'
        with:
          github-token: '${{ secrets.GITHUB_TOKEN }}'
          script: |-
            github.rest.issues.createComment({
              owner: '${{ github.repository }}'.split('/')[0],
              repo: '${{ github.repository }}'.split('/')[1],
              issue_number: '${{ github.event.issue.number }}',
              body: 'There is a problem with the iFlow CLI issue triaging. Please check the [action logs](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}) for details.'
            })
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `prompt` | The prompt to execute with iFlow CLI | ✅ Yes | - |
| `api_key` | iFlow API key for authentication | ✅ Yes | - |
| `settings_json` | Complete `~/.iflow/settings.json` content (JSON string). If provided, this will override other configuration options. | ❌ No | - |
| `base_url` | Custom base URL for iFlow API | ❌ No | `https://apis.iflow.cn/v1` |
| `model` | Model name to use | ❌ No | `Qwen3-Coder` |
| `working_directory` | Working directory to run iFlow CLI from | ❌ No | `.` |
| `timeout` | Timeout for iFlow CLI execution in seconds (1-86400) | ❌ No | `86400` |
| `extra_args` | Additional command line arguments to pass to iFlow CLI (space-separated string) | ❌ No | `` |
| `precmd` | Shell command(s) to execute before running iFlow CLI (e.g., "npm install", "git fetch") | ❌ No | `` |
| `gh_version` | Version of GitHub CLI to install (e.g., "2.76.2"). If not specified, uses the pre-installed version. | ❌ No | `` |
| `iflow_version` | Version of iFlow CLI to install (e.g., "0.2.4"). If not specified, uses the pre-installed version. | ❌ No | `` |

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
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze this codebase with debug output"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    extra_args: "--debug --max-tokens 3000"
```

#### Examples of Extra Arguments

- `--debug` - Enable iFLOW CLI debug mode  

### Using Pre-Execution Commands

The `precmd` input allows you to run shell commands before executing the iFlow CLI. This is useful for setting up the environment or installing dependencies needed for your iFlow commands.

```yaml
- name: iFlow with Pre-Execution Commands
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze this codebase after installing dependencies"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    precmd: |
      npm install
      git fetch origin main
```

#### Multi-line Commands

You can specify multiple commands by separating them with newlines:

```yaml
precmd: |
  npm ci
  npm run build
```

#### Quoted Arguments

For arguments containing spaces, use quotes:

```yaml
extra_args: '--debug'
```

### Using Custom Settings

For advanced users who need complete control over the iFlow configuration, you can provide a custom `settings.json` directly:

```yaml
- name: Custom iFlow Configuration
  uses: iflow-ai/iflow-cli-action@v1.3.0
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

### Using Custom Tool Versions

You can specify custom versions of GitHub CLI and iFlow CLI to use in your workflow:

```yaml
- name: iFlow CLI with Custom Versions
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze this codebase with specific tool versions"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    gh_version: "2.76.2"
    iflow_version: "0.2.4"
```

This is useful when you need to ensure compatibility with specific versions of these tools or when you want to use features available only in certain versions.

### Using MCP Servers

[MCP (Model Context Protocol)](https://modelcontextprotocol.io) allows iFlow CLI to connect to external tools and services, extending its capabilities beyond just AI model interactions. You can configure MCP servers in your workflow to enable features like code search, database querying, or custom tool integrations.

### Example: Using DeepWiki MCP Server

The following example demonstrates how to configure and use the DeepWiki MCP server for enhanced code search capabilities:

```yaml
- name: iFlow CLI with MCP Server
  uses: iflow-ai/iflow-cli-action@v1.3.0
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
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "Review this pull request for code quality, security issues, and best practices"
    api_key: ${{ secrets.IFLOW_API_KEY }}
```

### Documentation Generation

```yaml
- name: Generate Documentation
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "/init && Generate comprehensive API documentation"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "600"
```

### Automated Testing Suggestions

```yaml
- name: Test Strategy
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze the codebase and suggest a comprehensive testing strategy with specific test cases"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    model: "DeepSeek-V3"
```

### Architecture Analysis

```yaml
- name: Architecture Review
  uses: iflow-ai/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze the system architecture and suggest improvements for scalability and maintainability"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "900"
```

## Troubleshooting

### Common Issues

**Command timeout:** Increase the `timeout` value for complex operations

```yaml
timeout: "900"  # 15 minutes
```

**API authentication failed:** Verify your API key is correctly set in repository secrets

**Working directory not found:** Ensure the path exists and checkout action is used

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
- [Gemini CLI GitHub Action](https://github.com/google-github-actions/run-gemini-cli)
