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
go test ./pkg/frinkiac/http/               # Run specific package tests

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
- Test data in `/pkg/frinkiac/http/testdata/`
- Integration tests with real API responses

## Memory Bank Maintenance

The `memory-bank/` directory contains critical project context files that **must be kept up-to-date** when making changes to the codebase. These files serve as persistent knowledge for Claude Code and future development work.

### Memory Bank Files

- **activeContext.md**: Current work focus, recent changes, and immediate next steps
- **progress.md**: What works, what's left to build, current status, and known issues  
- **productContext.md**: Product vision, user stories, and feature requirements
- **projectbrief.md**: High-level project overview and goals
- **systemPatterns.md**: Architectural patterns, design decisions, and technical conventions
- **techContext.md**: Technical implementation details, frameworks, and dependencies

### Update Requirements

**When making changes, always update relevant memory-bank files:**

1. **After adding new features**: Update `progress.md` (move items from "What's Left" to "What Works")
2. **When changing architecture**: Update `systemPatterns.md` and `techContext.md`
3. **For new decisions/learnings**: Update `activeContext.md` with insights and approach changes
4. **When shifting focus**: Update `activeContext.md` "Current Work Focus" section
5. **For bugs/issues discovered**: Update `progress.md` "Known Issues" section

### Reading Context

**Before starting work, always read memory-bank files to understand:**
- Current project state and priorities (`activeContext.md`)
- What's already working vs. what needs building (`progress.md`)
- Established patterns and conventions (`systemPatterns.md`)
- Technical context and dependencies (`techContext.md`)

This ensures continuity across development sessions and maintains project knowledge.

## Environment Setup

The project uses a dev container configuration. Required environment variable:
- `OPENROUTER_API_KEY`: Your OpenRouter API key for AI functionality

Optional:
- `SMEE_SOURCE`: Webhook URL (auto-created if not provided)