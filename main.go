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
	Theme           string `json:"theme"`
	SelectedAuthType string `json:"selectedAuthType"`
	APIKey          string `json:"apiKey"`
	BaseURL         string `json:"baseUrl"`
	ModelName       string `json:"modelName"`
	SearchAPIKey    string `json:"searchApiKey"`
}

func main() {
	// Get inputs from environment variables (GitHub Actions convention)
	command := strings.TrimSpace(getInput("command"))
	apiKey := getInput("api-key")
	settingsJSON := getInput("settings-json")
	baseURL := getInput("base-url")
	if baseURL == "" {
		baseURL = "https://apis.iflow.cn/v1"
	}
	model := getInput("model")
	if model == "" {
		model = "Qwen3-Coder"
	}
	workingDir := getInput("working-directory")
	if workingDir == "" {
		workingDir = "."
	}
	timeoutStr := getInput("timeout")
	if timeoutStr == "" {
		timeoutStr = "300"
	}
	nodeVersion := getInput("node-version")
	if nodeVersion == "" {
		nodeVersion = "20"
	}

	// Validate required inputs
	if command == "" {
		setFailed("command input is required and cannot be empty")
		return
	}
	
	if apiKey == "" && settingsJSON == "" {
		setFailed("api-key input is required and cannot be empty")
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

	// Install dependencies
	info("Installing Node.js and iFlow CLI...")
	if err := installDependencies(nodeVersion); err != nil {
		setFailed(fmt.Sprintf("Failed to install dependencies: %v", err))
		return
	}

	// Configure iFlow settings
	info("Configuring iFlow settings...")
	if err := configureIFlow(apiKey, baseURL, model, settingsJSON); err != nil {
		setFailed(fmt.Sprintf("Failed to configure iFlow: %v", err))
		return
	}

	// Execute iFlow CLI command
	info(fmt.Sprintf("Executing iFlow CLI command: %s", command))
	result, exitCode, err := executeIFlow(command, timeout)
	if err != nil {
		setFailed(fmt.Sprintf("Failed to execute iFlow CLI: %v", err))
		return
	}

	// Set outputs
	setOutput("result", result)
	setOutput("exit-code", fmt.Sprintf("%d", exitCode))

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

func installDependencies(nodeVersion string) error {
	// Check if Node.js is already installed
	if cmd := exec.Command("node", "--version"); cmd.Run() == nil {
		info("Node.js is already installed")
	} else {
		info(fmt.Sprintf("Installing Node.js %s...", nodeVersion))
		if err := installNodeJS(nodeVersion); err != nil {
			return fmt.Errorf("failed to install Node.js: %w", err)
		}
	}

	// Check if iFlow CLI is already installed
	if cmd := exec.Command("iflow", "--version"); cmd.Run() == nil {
		info("iFlow CLI is already installed")
		return nil
	}

	// Install iFlow CLI
	info("Installing iFlow CLI...")
	cmd := exec.Command("bash", "-c", "curl -fsSL https://cloud.iflow.cn/iflow-cli/install.sh | bash")
	cmd.Env = append(os.Environ(), "NVM_DIR="+filepath.Join(os.Getenv("HOME"), ".nvm"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install iFlow CLI: %w\nOutput: %s", err, output)
	}

	info("iFlow CLI installed successfully")
	return nil
}

func installNodeJS(version string) error {
	// Install Node.js using nvm
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	nvmDir := filepath.Join(homeDir, ".nvm")
	
	// Download and install nvm
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		export NVM_DIR="%s"
		curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
		source "$NVM_DIR/nvm.sh"
		nvm install %s
		nvm use %s
		nvm alias default %s
	`, nvmDir, version, version, version))
	
	cmd.Env = append(os.Environ(), "NVM_DIR="+nvmDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Node.js: %w\nOutput: %s", err, output)
	}

	return nil
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
			return fmt.Errorf("invalid settings-json provided: %w", err)
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
			Theme:           "Default",
			SelectedAuthType: "iflow",
			APIKey:          apiKey,
			BaseURL:         baseURL,
			ModelName:       model,
			SearchAPIKey:    apiKey,
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

func executeIFlow(command string, timeoutSeconds int) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Setup environment for Node.js and npm
	homeDir, _ := os.UserHomeDir()
	nvmDir := filepath.Join(homeDir, ".nvm")
	npmGlobalDir := filepath.Join(homeDir, ".npm-global")
	
	// Build PATH with Node.js locations
	currentPath := os.Getenv("PATH")
	nodePaths := []string{
		filepath.Join(nvmDir, "versions", "node"),
		filepath.Join(npmGlobalDir, "bin"),
		filepath.Join(homeDir, ".local", "bin"),
	}
	
	for _, path := range nodePaths {
		if strings.Contains(currentPath, path) {
			continue
		}
		currentPath = path + ":" + currentPath
	}

	// Prepare the command
	var cmd *exec.Cmd
	if strings.Contains(command, "\n") || strings.Contains(command, ";") {
		// Multi-line or complex command - use interactive mode
		cmd = exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf(`
			source ~/.nvm/nvm.sh 2>/dev/null || true
			source ~/.bashrc 2>/dev/null || true
			export PATH="%s"
			echo "%s" | iflow
		`, currentPath, strings.ReplaceAll(command, "\"", "\\\"")))
	} else {
		// Simple command
		cmd = exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf(`
			source ~/.nvm/nvm.sh 2>/dev/null || true
			source ~/.bashrc 2>/dev/null || true
			export PATH="%s"
			iflow "%s"
		`, currentPath, command))
	}

	cmd.Env = append(os.Environ(),
		"PATH="+currentPath,
		"NVM_DIR="+nvmDir,
		"HOME="+homeDir,
	)

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