# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Licenser is a Go CLI tool that verifies and applies Apache 2.0 licenses to source code files across multiple programming languages. It supports 11 languages including Go, JavaScript, TypeScript, Python, C/C++, Shell, YAML, and others.

## Development Commands

### Make Targets
- Run all tests: `make test`
- Run linting: `make lint` (license verification + golangci-lint)
- Format code: `make format` (apply licenses + goimports + golangci-lint)
- Build binary: `make build`
- Install dependencies: `make install`
- Show help: `make help`

### Direct Commands
- Run specific test: `go test ./pkg/processor -v`
- Build: `go build -o licenser main.go`
- Run from source: `go run main.go [command]`
- Install locally: `go install`

### License Operations
- Verify licenses: `go run main.go verify -r`
- Apply licenses: `go run main.go apply -r "Copyright Owner"`

## Architecture

### Core Components

1. **Command Layer** (`pkg/command/`): Cobra-based CLI interface
   - `root.go`: Main command setup with recursive flag
   - `apply.go`: License application command
   - `verify.go`: License verification command

2. **License Handlers** (`pkg/license/`): License format and detection logic
   - `handler.go`: Interface for license implementations
   - `apache.go`: Apache 2.0 license implementation
   - `generic.go`: Generic license handler

3. **File Processing** (`pkg/file/`): File mutation and styling
   - `licenser.go`: Core interface for license operations
   - `mutator.go`: File content modification logic
   - `style.go`: Language-specific comment styling

4. **Processor** (`pkg/processor/`): Directory traversal and filtering
   - Handles .gitignore and .licenserignore patterns
   - Concurrent file processing with worker pools
   - Language detection and file filtering

### Key Patterns

- Uses `github.com/spf13/cobra` for CLI structure
- File operations are concurrent using goroutines and channels
- Language-specific commenting styles determined by file extension
- Gitignore-style pattern matching for file exclusion
- Dry-run support for testing changes before applying

### Language Support

Each supported language has specific comment style handling in `pkg/file/style.go`. When adding new language support, update both the style mappings and the processor's file extension recognition.