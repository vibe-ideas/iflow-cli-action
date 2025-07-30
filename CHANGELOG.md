# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
