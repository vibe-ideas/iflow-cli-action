# iFlow CLI Action - Cobra Migration

This document describes the updated iFlow CLI Action after migrating from direct environment variable parsing to using the Cobra CLI framework.

## Key Changes

### 1. Command Structure
The application now uses Cobra for command-line argument parsing, providing better structure and help documentation.

### 2. Dual Mode Operation
The tool now supports two operation modes:

#### GitHub Actions Mode (Default)
- Automatically detected when `GITHUB_ACTIONS=true` environment variable is set
- Uses `INPUT_*` environment variables as before
- Maintains backward compatibility with existing GitHub Actions workflows

#### CLI Mode
- Activated with command-line flags
- Useful for local development and testing
- Can be forced with `--use-env-vars=false`

## Usage Examples

### GitHub Actions Mode (Backward Compatible)
```yaml
# .github/workflows/example.yml
- name: Run iFlow CLI
  uses: ./
  with:
    prompt: "Review this code and suggest improvements"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    model: "Qwen3-Coder"
    timeout: 300
```

### CLI Mode
```bash
# Basic usage
./iflow-action --prompt "Review this code" --api-key "your-api-key"

# With custom settings
./iflow-action \
  --prompt "Analyze this function" \
  --api-key "your-api-key" \
  --model "Claude-3" \
  --base-url "https://custom.api.com/v1" \
  --timeout 600 \
  --working-directory "./src"

# Using settings JSON
./iflow-action \
  --prompt "Fix the bug in main.go" \
  --settings-json '{"apiKey":"key","baseUrl":"https://apis.iflow.cn/v1","modelName":"Qwen3-Coder"}'

# Force environment variable mode
./iflow-action --use-env-vars --prompt "test"
```

## Available Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--prompt` | `-p` | *required* | The prompt to send to iFlow CLI |
| `--api-key` | | | API key for iFlow authentication |
| `--settings-json` | | | Complete settings JSON configuration |
| `--base-url` | | `https://apis.iflow.cn/v1` | Base URL for the iFlow API |
| `--model` | | `Qwen3-Coder` | Model name to use |
| `--working-directory` | | `.` | Working directory for execution |
| `--timeout` | | `300` | Timeout in seconds (1-3600) |
| `--use-env-vars` | | `false` | Force environment variables mode |
| `--help` | `-h` | | Show help information |

## Environment Variables (GitHub Actions)

When running in GitHub Actions mode, the following environment variables are used:

- `INPUT_PROMPT` - The prompt text
- `INPUT_API_KEY` - API key for authentication
- `INPUT_SETTINGS_JSON` - Complete settings JSON
- `INPUT_BASE_URL` - Custom base URL
- `INPUT_MODEL` - Model name
- `INPUT_WORKING_DIRECTORY` - Working directory
- `INPUT_TIMEOUT` - Timeout in seconds

## Benefits of Cobra Migration

1. **Better CLI Experience**: Rich help documentation and flag validation
2. **Flexible Usage**: Can be used both as GitHub Action and standalone CLI tool
3. **Maintainability**: Better code organization with the cmd package
4. **Extensibility**: Easy to add new commands and flags in the future
5. **Error Handling**: Improved error messages and validation
6. **Documentation**: Built-in help system with `--help` flag

## Backward Compatibility

The migration maintains 100% backward compatibility with existing GitHub Actions workflows. No changes are required to existing `.github/workflows/*.yml` files.

## Development and Testing

```bash
# Build the application
go build -o iflow-action

# Test with CLI flags
./iflow-action --prompt "test prompt" --api-key "test-key"

# Test with environment variables
export INPUT_PROMPT="test prompt"
export INPUT_API_KEY="test-key"
export GITHUB_ACTIONS="true"
./iflow-action

# Run with help
./iflow-action --help
```

## Future Enhancements

With Cobra in place, future enhancements become easier:

1. **Subcommands**: Add commands like `iflow-action config`, `iflow-action validate`
2. **Configuration Files**: Support for config files with `viper`
3. **Auto-completion**: Shell completion for better developer experience
4. **Plugins**: Extensible plugin system
5. **Interactive Mode**: Interactive prompts for missing required parameters
