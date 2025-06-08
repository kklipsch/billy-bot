# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Billy Bot is a Go CLI application that categorizes prompts into relevant Simpsons quotes and retrieves corresponding screen captures from Frinkiac. It replaces Billy's skill of selecting appropriate Simpsons references for any occasion.

## Development Commands

```bash
# Build and run
go build                                    # Build the billy-bot binary
go run main.go frinkiac "your prompt"      # Run frinkiac command directly
go run main.go smee                        # Run smee webhook listener

# Testing
go test ./...                              # Run all tests
go test ./pkg/frinkiac/client/             # Run specific package tests

# Code quality (used in CI)
gofmt -s -w .                              # Format code
golint ./...                               # Lint code
```

The following commands must be run and have a successful result, prior to any commit:
```
go test ./...
gofmt -s -w .
golint ./...
```

## Architecture

### Core Components
- **main.go**: CLI entry point using Kong framework
- **pkg/frinkiac**: Simpsons quote/screencap engine with OpenRouter AI integration
- **pkg/openrouter**: Complete OpenRouter API client for chat completions
- **pkg/smee**: Server-Sent Events client for webhook processing
- **pkg/jsonschema**: Type-safe JSON schema creation for structured AI responses
- **pkg/config**: Configuration utilities with env/flag precedence

### Data Flow
1. User provides prompt → 2. OpenRouter AI categorizes into Simpsons quotes → 3. Frinkiac API retrieves matching screen captures → 4. Results displayed with captions

### External Dependencies
- **OpenRouter API**: AI-powered quote categorization (requires API key)
- **Frinkiac website**: Quote search and screen capture retrieval (JSON API + HTML fallback)
- **Smee.io**: Webhook event streaming for real-time processing

## Key Patterns

### Error Handling
- Explicit error checking throughout codebase
- Graceful fallbacks (API failures fall back to HTML parsing)
- Structured error responses with context

### Configuration
- Environment variables take precedence over flags
- `.env` file support with `--env-file` flag
- Required: `OPENROUTER_API_KEY` environment variable

### Testing Strategy
- Unit tests use testify framework (require/assert patterns)
- HTTP client mocking for external API testing
- Test data in `/pkg/frinkiac/client/testdata/`
- Integration tests with real API responses

## Environment Setup

The project uses a dev container configuration. Required environment variable:
- `OPENROUTER_API_KEY`: Your OpenRouter API key for AI functionality

Optional:
- `SMEE_SOURCE`: Webhook URL (auto-created if not provided)