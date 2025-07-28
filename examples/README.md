# Example Workflows

This directory contains example GitHub Actions workflows demonstrating how to use the iFlow CLI Action.

**Note:** All iFlow CLI commands are automatically executed with `--prompt` and `--yolo` flags for non-interactive, streamlined operation.

## Basic Examples

### Code Review on Pull Request
```yaml
# .github/workflows/iflow-review.yml
name: iFlow Code Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Review Code
        uses: vibe-ideas/iflow-cli-action@v1
        with:
          prompt: "Review this pull request for code quality and suggest improvements"
          api_key: ${{ secrets.IFLOW_API_KEY }}
```

### Documentation Generation
```yaml
# .github/workflows/generate-docs.yml
name: Generate Documentation
on:
  push:
    branches: [main]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Generate Docs
        uses: vibe-ideas/iflow-cli-action@v1
        with:
          prompt: "/init && Generate comprehensive documentation for this project"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          timeout: "600"
```

### Security Analysis
```yaml
# .github/workflows/security-scan.yml
name: Security Analysis
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly on Monday

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Security Scan
        uses: vibe-ideas/iflow-cli-action@v1
        with:
          prompt: "Analyze this codebase for security vulnerabilities and provide recommendations"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "DeepSeek-V3"
          timeout: "900"
```

## Advanced Examples

### Multi-step Analysis
```yaml
# .github/workflows/comprehensive-analysis.yml
name: Comprehensive Analysis
on: [workflow_dispatch]

jobs:
  analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Project Overview
        uses: vibe-ideas/iflow-cli-action@v1
        with:
          prompt: "/init"
          api_key: ${{ secrets.IFLOW_API_KEY }}
        id: init
      
      - name: Architecture Analysis
        uses: vibe-ideas/iflow-cli-action@v1
        with:
          prompt: "Based on the project analysis, provide detailed architecture recommendations"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Qwen3-Coder"
        id: arch
      
      - name: Performance Review
        uses: vibe-ideas/iflow-cli-action@v1
        with:
          prompt: "Analyze the code for performance bottlenecks and optimization opportunities"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "DeepSeek-V3"
        id: perf
      
      - name: Create Summary Report
        run: |
          echo "# Comprehensive Analysis Report" > analysis-report.md
          echo "## Architecture Analysis" >> analysis-report.md
          echo "${{ steps.arch.outputs.result }}" >> analysis-report.md
          echo "## Performance Analysis" >> analysis-report.md
          echo "${{ steps.perf.outputs.result }}" >> analysis-report.md
      
      - name: Upload Report
        uses: actions/upload-artifact@v4
        with:
          name: analysis-report
          path: analysis-report.md
```

## Configuration Examples

### Custom Model Configuration
```yaml
- name: Use Custom Model
  uses: vibe-ideas/iflow-cli-action@v1
  with:
    prompt: "Analyze this code"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    model: "Kimi-K2"
    base_url: "https://apis.iflow.cn/v1"
```

### Extended Timeout for Complex Tasks
```yaml
- name: Complex Analysis
  uses: vibe-ideas/iflow-cli-action@v1
  with:
    prompt: "Perform comprehensive code analysis and refactoring suggestions"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "1800"  # 30 minutes
```

### Different Working Directory
```yaml
- name: Analyze Specific Module
  uses: vibe-ideas/iflow-cli-action@v1
  with:
    prompt: "Analyze this module for improvement opportunities"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    working_directory: "./src/core"
```

## Setup Instructions

1. **Add API Key to Secrets:**
   - Go to your repository settings
   - Navigate to Secrets and Variables > Actions
   - Click "New repository secret"
   - Name: `IFLOW_API_KEY`
   - Value: Your iFlow API key

2. **Create Workflow File:**
   - Create `.github/workflows/` directory in your repository
   - Add one of the example workflows above
   - Customize the `prompt` and other parameters as needed

3. **Test the Workflow:**
   - Commit and push the workflow file
   - Trigger the workflow based on its trigger conditions
   - Check the Actions tab for execution results