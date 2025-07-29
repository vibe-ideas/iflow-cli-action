package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Config holds all configuration options
type Config struct {
	Prompt       string
	APIKey       string
	SettingsJSON string
	BaseURL      string
	Model        string
	WorkingDir   string
	Timeout      int
	UseEnvVars   bool // Flag to indicate whether to use environment variables (GitHub Actions mode)
}

// IFlowSettings represents the iFlow configuration
type IFlowSettings struct {
	Theme            string `json:"theme"`
	SelectedAuthType string `json:"selectedAuthType"`
	APIKey           string `json:"apiKey"`
	BaseURL          string `json:"baseUrl"`
	ModelName        string `json:"modelName"`
	SearchAPIKey     string `json:"searchApiKey"`
}

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iflow-action",
	Short: "iFlow CLI Action wrapper",
	Long: `A GitHub Action wrapper for iFlow CLI that provides intelligent code assistance.
	
This tool can run in two modes:
1. GitHub Actions mode: Uses environment variables (INPUT_*) for configuration
2. CLI mode: Uses command-line flags for configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIFlowAction()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define flags
	rootCmd.Flags().StringVarP(&config.Prompt, "prompt", "p", "", "The prompt to send to iFlow CLI (required in CLI mode)")
	rootCmd.Flags().StringVar(&config.APIKey, "api-key", "", "API key for iFlow authentication")
	rootCmd.Flags().StringVar(&config.SettingsJSON, "settings-json", "", "Complete settings JSON configuration")
	rootCmd.Flags().StringVar(&config.BaseURL, "base-url", "https://apis.iflow.cn/v1", "Base URL for the iFlow API")
	rootCmd.Flags().StringVar(&config.Model, "model", "Qwen3-Coder", "Model name to use")
	rootCmd.Flags().StringVar(&config.WorkingDir, "working-directory", ".", "Working directory for execution")
	rootCmd.Flags().IntVar(&config.Timeout, "timeout", 300, "Timeout in seconds (1-3600)")
	rootCmd.Flags().BoolVar(&config.UseEnvVars, "use-env-vars", false, "Use environment variables for configuration (GitHub Actions mode)")

	// Mark required flags only if not in GitHub Actions mode - this will be validated later
}

func runIFlowAction() error {
	// If use-env-vars is set or we detect GitHub Actions environment, use environment variables
	if config.UseEnvVars || isGitHubActions() {
		if err := loadConfigFromEnv(); err != nil {
			return fmt.Errorf("failed to load config from environment: %w", err)
		}
	}

	// Validate configuration
	if err := validateConfig(); err != nil {
		return err
	}

	// Setup working directory
	if config.WorkingDir != "." && config.WorkingDir != "" {
		if err := os.Chdir(config.WorkingDir); err != nil {
			return fmt.Errorf("failed to change working directory: %w", err)
		}
	}

	// iFlow CLI is pre-installed in Docker image
	info("iFlow CLI is pre-installed and ready to use")

	// Configure iFlow settings
	info("Configuring iFlow settings...")
	if err := configureIFlow(); err != nil {
		return fmt.Errorf("failed to configure iFlow: %w", err)
	}

	// Execute iFlow CLI command with --prompt and --yolo flags
	info(fmt.Sprintf("Executing iFlow CLI prompt with --prompt and --yolo: %s", config.Prompt))
	result, exitCode, err := executeIFlow()
	if err != nil {
		return fmt.Errorf("failed to execute iFlow CLI: %w", err)
	}

	// Set outputs (GitHub Actions mode) or print results (CLI mode)
	if config.UseEnvVars || isGitHubActions() {
		setOutput("result", result)
		setOutput("exit_code", fmt.Sprintf("%d", exitCode))
	} else {
		fmt.Printf("Exit Code: %d\n", exitCode)
		fmt.Printf("Result:\n%s\n", result)
	}

	if exitCode != 0 {
		if config.UseEnvVars || isGitHubActions() {
			setFailed(fmt.Sprintf("iFlow CLI exited with code %d", exitCode))
		}
		return fmt.Errorf("iFlow CLI exited with code %d", exitCode)
	}

	info("iFlow CLI execution completed successfully")
	return nil
}

func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func loadConfigFromEnv() error {
	// Load configuration from environment variables (GitHub Actions convention)
	if prompt := getInput("prompt"); prompt != "" {
		config.Prompt = strings.TrimSpace(prompt)
	}
	if apiKey := getInput("api_key"); apiKey != "" {
		config.APIKey = apiKey
	}
	if settingsJSON := getInput("settings_json"); settingsJSON != "" {
		config.SettingsJSON = settingsJSON
	}
	if baseURL := getInput("base_url"); baseURL != "" {
		config.BaseURL = baseURL
	}
	if model := getInput("model"); model != "" {
		config.Model = model
	}
	if workingDir := getInput("working_directory"); workingDir != "" {
		config.WorkingDir = workingDir
	}
	if timeoutStr := getInput("timeout"); timeoutStr != "" {
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %s", timeoutStr)
		}
		config.Timeout = timeout
	}

	return nil
}

func validateConfig() error {
	// Validate required inputs
	if config.Prompt == "" {
		if config.UseEnvVars || isGitHubActions() {
			setFailed("prompt input is required and cannot be empty")
			return fmt.Errorf("prompt input is required and cannot be empty")
		}
		return fmt.Errorf("prompt is required and cannot be empty")
	}

	if config.APIKey == "" && config.SettingsJSON == "" {
		if config.UseEnvVars || isGitHubActions() {
			setFailed("api_key input is required and cannot be empty")
			return fmt.Errorf("api_key input is required and cannot be empty")
		}
		return fmt.Errorf("api-key is required when settings-json is not provided")
	}

	// Validate timeout range (1 second to 1 hour)
	if config.Timeout < 1 || config.Timeout > 3600 {
		if config.UseEnvVars || isGitHubActions() {
			setFailed(fmt.Sprintf("Timeout must be between 1 and 3600 seconds, got: %d", config.Timeout))
		}
		return fmt.Errorf("timeout must be between 1 and 3600 seconds, got: %d", config.Timeout)
	}

	return nil
}

// Helper functions for GitHub Actions
func getInput(name string) string {
	// GitHub Actions sets inputs as environment variables with INPUT_ prefix
	envName := "INPUT_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
	return os.Getenv(envName)
}

func setOutput(name, value string) {
	// GitHub Actions outputs can be set using the GITHUB_OUTPUT file
	if outputFile := os.Getenv("GITHUB_OUTPUT"); outputFile != "" {
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("::error::Failed to open output file: %v\n", err)
			return
		}
		defer f.Close()

		// Use proper GitHub Actions output format with multiline support
		delimiter := fmt.Sprintf("EOF_%d", os.Getpid())
		_, err = f.WriteString(fmt.Sprintf("%s<<%s\n%s\n%s\n", name, delimiter, value, delimiter))
		if err != nil {
			fmt.Printf("::error::Failed to write output: %v\n", err)
		}
	} else {
		// Fallback to legacy format if GITHUB_OUTPUT is not available
		fmt.Printf("::set-output name=%s::%s\n", name, value)
	}
}

func info(message string) {
	if isGitHubActions() {
		fmt.Printf("::notice::%s\n", message)
	} else {
		fmt.Printf("INFO: %s\n", message)
	}
}

func setFailed(message string) {
	fmt.Printf("::error::%s\n", message)
	os.Exit(1)
}

func configureIFlow() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	iflowDir := filepath.Join(homeDir, ".iflow")
	if err := os.MkdirAll(iflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create .iflow directory: %w", err)
	}

	settingsFile := filepath.Join(iflowDir, "settings.json")

	var settingsData []byte

	if config.SettingsJSON != "" {
		// Use provided settings JSON directly
		info("Using provided settings.json content")

		// Validate that it's valid JSON
		var testSettings map[string]interface{}
		if err := json.Unmarshal([]byte(config.SettingsJSON), &testSettings); err != nil {
			return fmt.Errorf("invalid settings_json provided: %w", err)
		}

		// Pretty format the JSON
		var prettyJSON json.RawMessage = []byte(config.SettingsJSON)
		settingsData, err = json.MarshalIndent(prettyJSON, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format settings JSON: %w", err)
		}
	} else {
		// Create settings from individual parameters
		info("Creating settings from individual parameters")
		settings := IFlowSettings{
			Theme:            "Default",
			SelectedAuthType: "iflow",
			APIKey:           config.APIKey,
			BaseURL:          config.BaseURL,
			ModelName:        config.Model,
			SearchAPIKey:     config.APIKey,
		}

		settingsData, err = json.MarshalIndent(settings, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal settings: %w", err)
		}
	}

	if err := os.WriteFile(settingsFile, settingsData, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	info(fmt.Sprintf("iFlow settings configured at %s", settingsFile))
	return nil
}

func executeIFlow() (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
	defer cancel()

	// Prepare the command with --prompt and --yolo flags by default
	// Use --prompt and --yolo flags for all commands
	cmd := exec.CommandContext(ctx, "iflow", "--yolo", "--prompt", config.Prompt)

	output, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// Non-exit error (e.g., timeout, command not found)
			return string(output), 1, err
		}
	}

	return string(output), exitCode, nil
}
