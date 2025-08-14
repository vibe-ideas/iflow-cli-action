# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Real-time Output Streaming**: Implemented real-time output streaming for iFlow CLI execution, providing immediate feedback during long-running operations
- **MCP Server Support**: Added support for Model Context Protocol (MCP) servers, enabling integration with external tools and services
- **GitHub CLI Installation**: Added GitHub CLI (gh) to the Docker image for enhanced GitHub integration capabilities
- **Enhanced Timeout Support**: Extended timeout configuration to allow up to 24 hours (86400 seconds) for long-running operations
- **Go Runtime Environment**: Added Go installation to Docker image to support github-mcp-server and other Go-based tools
- **wget Dependency**: Added wget to Dockerfile dependencies for improved compatibility

### Changed

- **Docker Base Image**: Switched Docker base image from node:22-slim to ubuntu:22.04 for better stability and compatibility
- **Improved Docker Configuration**: Enhanced Go installation and environment configuration in Dockerfile
- **Enhanced GitHub Actions Workflows**: Added automated issue triage, PR review automation, and issue killer workflows
- **Updated GitHub Actions Integration**: Improved GitHub Actions step summaries with better formatting and more detailed information
- **Enhanced Error Handling**: Improved error handling and reporting in GitHub Actions workflows

### Fixed

- **Command Streaming Issues**: Fixed issues with output streaming and buffering during iFlow CLI execution
- **Go Runtime Availability**: Ensured Go command is properly available in the runtime environment
- **Regex Pattern Matching**: Corrected regex for extracting PR numbers in workflow scripts
- **Comment Parsing**: Improved comment parsing for review-killer scripts

## [1.2.0] - 2025-08-02

### Added

- **Extra Arguments Support**: New `extra_args` input parameter allows passing additional command-line arguments to iFlow CLI
- Enhanced argument parsing with support for quoted arguments containing spaces
- Dynamic argument support for workflow inputs
- Updated documentation with extra_args examples and use cases
- New example workflow demonstrating extra_args functionality
- iFlow CLI version display functionality
- GitHub Actions workflow for iFlow with MCP (Model Context Protocol) support

### Changed

- Updated command execution logic to support additional arguments
- Enhanced GitHub Actions step summary to display extra arguments
- Improved configuration display in workflow summaries

## [1.1.0] - 2025-07-30

### Added

- Chinese documentation (README_zh.md)
- IFLOW.md development guide
- Examples for code review and documentation workflows

### Changed

- Updated iFlow CLI version in Dockerfile
- Updated timeout format in deploy workflow
- Improved homepage generation instructions
- Updated deploy homepage workflow prompt
- Updated iflow-cli-action version in README.md

### Refactored

- Moved timeout check before error handling in executeIFlow function
- Updated timeout input range and added flag in config struct

## [1.0.0] - 2025-07-29

### Added

- Initial release of iFlow CLI GitHub Action
- Docker-based action with pre-installed Node.js 22 and iFlow CLI
- Support for configurable authentication with iFlow API
- Support for custom models and API endpoints
- Flexible command execution with timeout control
- GitHub Actions Summary integration for rich execution reports
- Action outputs for result and exit code
