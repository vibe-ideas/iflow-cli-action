package main

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
)

// IFlowSettings represents the iFlow configuration
type IFlowSettings struct {
	Theme            string `json:"theme"`
	SelectedAuthType string `json:"selectedAuthType"`
	APIKey           string `json:"apiKey"`
	BaseURL          string `json:"baseUrl"`
	ModelName        string `json:"modelName"`
	SearchAPIKey     string `json:"searchApiKey"`
}

func main() {
	// Get inputs from environment variables (GitHub Actions convention)
	prompt := strings.TrimSpace(getInput("prompt"))
	apiKey := getInput("api_key")
	settingsJSON := getInput("settings_json")
	baseURL := getInput("base_url")
	if baseURL == "" {
		baseURL = "https://apis.iflow.cn/v1"
	}
	model := getInput("model")
	if model == "" {
		model = "Qwen3-Coder"
	}
	workingDir := getInput("working_directory")
	if workingDir == "" {
		workingDir = "."
	}
	timeoutStr := getInput("timeout")
	if timeoutStr == "" {
		timeoutStr = "300"
	}

	// Validate required inputs
	if prompt == "" {
		setFailed("prompt input is required and cannot be empty")
		return
	}

	if apiKey == "" && settingsJSON == "" {
		setFailed("api_key input is required and cannot be empty")
		return
	}

	// Parse timeout
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		setFailed(fmt.Sprintf("Invalid timeout value: %s", timeoutStr))
		return
	}

	// Validate timeout range (1 second to 1 hour)
	if timeout < 1 || timeout > 3600 {
		setFailed(fmt.Sprintf("Timeout must be between 1 and 3600 seconds, got: %d", timeout))
		return
	}

	// Setup working directory
	if workingDir != "." && workingDir != "" {
		if err := os.Chdir(workingDir); err != nil {
			setFailed(fmt.Sprintf("Failed to change working directory: %v", err))
			return
		}
	}

	// iFlow CLI is pre-installed in Docker image
	info("iFlow CLI is pre-installed and ready to use")

	// Configure iFlow settings
	info("Configuring iFlow settings...")
	if err := configureIFlow(apiKey, baseURL, model, settingsJSON); err != nil {
		setFailed(fmt.Sprintf("Failed to configure iFlow: %v", err))
		return
	}

	// Execute iFlow CLI command with --prompt and --yolo flags
	info(fmt.Sprintf("Executing iFlow CLI prompt with --prompt and --yolo: %s", prompt))
	result, exitCode, err := executeIFlow(prompt, timeout)
	if err != nil {
		setFailed(fmt.Sprintf("Failed to execute iFlow CLI: %v", err))
		return
	}

	// Set outputs
	setOutput("result", result)
	setOutput("exit_code", fmt.Sprintf("%d", exitCode))

	if exitCode != 0 {
		setFailed(fmt.Sprintf("iFlow CLI exited with code %d", exitCode))
		return
	}

	info("iFlow CLI execution completed successfully")
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
		if err == nil {
			defer f.Close()
			f.WriteString(fmt.Sprintf("%s=%s\n", name, value))
		}
	}
}

func info(message string) {
	fmt.Printf("::notice::%s\n", message)
}

func setFailed(message string) {
	fmt.Printf("::error::%s\n", message)
	os.Exit(1)
}

func configureIFlow(apiKey, baseURL, model, settingsJSON string) error {
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

	if settingsJSON != "" {
		// Use provided settings JSON directly
		info("Using provided settings.json content")

		// Validate that it's valid JSON
		var testSettings map[string]interface{}
		if err := json.Unmarshal([]byte(settingsJSON), &testSettings); err != nil {
			return fmt.Errorf("invalid settings_json provided: %w", err)
		}

		// Pretty format the JSON
		var prettyJSON json.RawMessage = []byte(settingsJSON)
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
			APIKey:           apiKey,
			BaseURL:          baseURL,
			ModelName:        model,
			SearchAPIKey:     apiKey,
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

func executeIFlow(prompt string, timeoutSeconds int) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Prepare the command with --prompt and --yolo flags by default
	var cmd *exec.Cmd

	// Use --prompt and --yolo flags for all commands
	cmd = exec.CommandContext(ctx, "iflow", "--prompt", prompt, "--yolo")

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
