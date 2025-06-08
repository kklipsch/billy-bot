# Technical Context: Billy Bot

## Technologies Used

### Programming Language
- **Go (version 1.22)**: The primary programming language used for the project. Go was chosen for its simplicity, strong typing, excellent concurrency support, and efficient performance.

### Key Libraries and Packages
- **Kong (github.com/alecthomas/kong)**: Used for CLI command parsing and routing. Kong provides a declarative, struct-based approach to defining CLI commands and options.
- **Zerolog (github.com/rs/zerolog)**: Employed for structured, leveled logging. Zerolog offers high-performance, JSON-structured logging with a clean API.
- **Godotenv (github.com/joho/godotenv)**: Used for loading environment variables from .env files, facilitating configuration management across different environments.
- **Standard Go Libraries**: Leveraging Go's rich standard library, including:
  - `context`: For managing cancellation, timeouts, and request-scoped values
  - `os/signal`: For handling system signals for graceful shutdown
  - `fmt`: For formatted I/O operations
  - `os`: For operating system functionality

### External Services
- **Frinkiac Website**: A website designed for humans (not an API) that allows searching for Simpsons scenes based on quotes. Queries are made in the form of `https://frinkiac.com/?q=garbage%20water` and the site returns HTML rendered for human consumption.
- **Smee**: Service for receiving webhook events, allowing the bot to respond to external triggers.

## Development Setup

### Environment Requirements
- **Go 1.22+**: Required for building and running the application.
- **Git**: For version control and source code management.
- **Dev Container**: The project attempts to maintain an up-to-date development container configuration for consistent development environments.

### Configuration
- **.env File**: Used for storing environment variables and configuration settings. The default path is `.env` in the project root, but can be specified with the `-e` flag.
- **Log Level Configuration**: Configurable logging levels (debug, info, warn, error, fatal, panic) set via the `-l` flag, defaulting to "warn".

### Build and Run Process
1. **Build**: Standard Go build process (`go build`)
2. **Run**: Execute the binary with appropriate command and flags:
   - For Smee functionality: `./billy-bot smee [options]`
   - For Frinkiac functionality: `./billy-bot frinkiac [options]`
   - Global options include `-e` for env file path and `-l` for log level

## Technical Constraints

### Performance Considerations
- **Response Time**: The bot should provide quick responses to maintain a good user experience.
- **Resource Usage**: As a CLI tool, the application should have minimal resource footprint.

### External Website Limitations
- **Frinkiac Website Rate Limits**: The bot's operation may be constrained by any rate limits or anti-scraping measures imposed by the Frinkiac website.
- **HTML Parsing Challenges**: Frinkiac has an undocumented api, it presents as a full page app so discovering the api is a challenge. 
- **Webhook Processing Capacity**: The number of webhook events that can be processed simultaneously may be limited by system resources.

### Security Considerations
- **Environment Variables**: Sensitive configuration (API keys, tokens) should be stored in environment variables, not hardcoded.
- **Input Validation**: All user inputs and webhook payloads should be properly validated to prevent security issues.

## Dependencies

### Direct Dependencies
- **github.com/alecthomas/kong v1.11.0**: CLI parsing library
- **github.com/joho/godotenv v1.5.1**: Environment variable loading from .env files
- **github.com/rs/zerolog v1.34.0**: Structured logging library
- **github.com/stretchr/testify v1.10.0**: Testing toolkit with rich assertion capabilities
- **github.com/kklipsch/billy-bot/pkg/frinkiac**: Internal package for Frinkiac functionality
- **github.com/kklipsch/billy-bot/pkg/smee**: Internal package for Smee client functionality

### Indirect Dependencies
- **github.com/mattn/go-colorable v0.1.13**: Terminal color support
- **github.com/mattn/go-isatty v0.0.19**: Terminal type detection
- **golang.org/x/sys v0.12.0**: System-level operations

## Tool Usage Patterns

### Command-Line Interface
The application follows a subcommand pattern for its CLI:
```
billy-bot [global options] <command> [command options]
```

Global options include:
- `-e, --env-file`: Path to the .env file (default: ".env")
- `-l, --log-level`: Set the log level (default: "warn")

Available commands:
- `smee`: Run the Smee client to receive webhook events
- `frinkiac`: Engage the Frinkiac tool to find Simpsons scenes by querying the Frinkiac website

### Logging Practices
- **Structured Logging**: All logs are structured in JSON format for better parsing and analysis.
- **Log Levels**: Different log levels are used appropriately:
  - `debug`: Detailed information for debugging
  - `info`: General information about application progress
  - `warn`: Warning conditions that don't affect operation
  - `error`: Error conditions that affect specific operations
  - `fatal`: Critical errors that prevent the application from running

### Error Handling
- **Explicit Error Checking**: Following Go's convention of explicit error checking.
- **Contextual Errors**: Errors include context about where and why they occurred.
- **User-Friendly Messages**: Error messages presented to users are clear and actionable.

### Context Management
- **Signal Handling**: The application sets up proper signal handling for graceful shutdown.
- **Context Propagation**: Contexts are propagated throughout the application for proper cancellation.

## Development Workflow

### Code Organization
- **Main Package**: Entry point and CLI setup
- **Pkg Directory**: Contains reusable packages:
  - `frinkiac`: Functionality for interacting with the Frinkiac API
  - `smee`: Functionality for the Smee client

### Development Standards
- **Go Development Standards**: See [go-development-standards.md](./go-development-standards.md) for detailed coding standards and patterns
- **Functional Style Preferred**: The project uses functional programming patterns over object-oriented approaches for API clients and business logic
- **Explicit Dependencies**: Functions take dependencies as parameters rather than embedding them in structs

### Testing Strategy
- **Unit Tests**: For testing individual components in isolation
- **Integration Tests**: For testing interactions between components
- **End-to-End Tests**: For testing the complete application flow
- **Testify Package**: The project uses the testify package for writing tests:
  - `require`: For assertions that should terminate the test immediately if they fail
  - `assert`: For assertions that should report failures but continue test execution
  - This approach provides clearer test code and more descriptive error messages

### Deployment Considerations
- **Binary Distribution**: As a Go application, it can be compiled into a single binary for easy distribution.
- **Environment Configuration**: Different environments (development, staging, production) can be configured via .env files.
- **Container Support**: The application can be containerized for consistent deployment across environments.
