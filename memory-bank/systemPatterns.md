# System Patterns: Billy Bot

## System Architecture

Billy Bot follows a command-line application architecture built in Go, with a clear separation of concerns between different components. The architecture is designed to be modular, allowing for easy extension and maintenance.

```
flowchart TD
    CLI[CLI Interface] --> Commands[Command Handlers]
    Commands --> Frinkiac[Frinkiac Module]
    Commands --> Smee[Smee Client]
    Frinkiac --> Website[Frinkiac Website]
    Frinkiac --> Parser[HTML Parser]
    Smee --> Webhooks[Webhook Events]
```

### Key Components

1. **CLI Interface**: The main entry point for the application, handling command-line arguments and routing to appropriate command handlers.
2. **Command Handlers**: Specialized modules that implement specific functionality (Smee, Frinkiac).
3. **Frinkiac Module**: Handles interaction with the Frinkiac website to find and retrieve Simpsons scenes by querying the website and parsing the HTML responses.
4. **HTML Parser**: Component responsible for extracting relevant information from the HTML responses returned by the Frinkiac website.
5. **Smee Client**: Manages webhook event reception and processing.

## Key Technical Decisions

1. **Go Language**: The project is implemented in Go (version 1.22), leveraging its strong concurrency model, performance, and simplicity.

2. **Command Pattern**: The application uses the command pattern for CLI interactions, with each major function (Smee, Frinkiac) implemented as a separate command.

3. **Environment Configuration**: The application uses environment variables (via .env files) for configuration, allowing for flexible deployment across different environments.

4. **Structured Logging**: The project employs zerolog for structured logging, enabling better log analysis and monitoring.

5. **Context-Based Cancellation**: The application uses Go's context package for proper handling of cancellation signals, ensuring graceful shutdown.

6. **Web Scraping Integration**: Rather than implementing screen cap selection logic from scratch, the project interacts with the existing Frinkiac website, parsing its HTML responses to extract the necessary information.

## Design Patterns in Use

1. **Command Pattern**: Each major function is implemented as a command, with a consistent interface for execution.

2. **Dependency Injection**: Dependencies are injected into components rather than created internally, facilitating testing and flexibility.

3. **Context Propagation**: Go contexts are propagated throughout the application to manage cancellation and timeouts.

4. **Configuration Management**: Environment variables and configuration files are used to manage application settings.

5. **Structured Error Handling**: Errors are handled in a structured way, with appropriate logging and user feedback.

## Component Relationships

### CLI to Commands
The main CLI interface parses command-line arguments using the Kong package and routes execution to the appropriate command handler (Smee or Frinkiac).

```go
cli := CLI{}
k := kong.Parse(&cli, ...)
err = k.Run()
```

### Commands to Services
Each command handler interacts with its respective service (Smee client or Frinkiac service) to perform the requested operation.

### Module to External Websites
The Frinkiac module communicates with the Frinkiac website by constructing appropriate URLs (e.g., `https://frinkiac.com/?q=garbage%20water`), sending HTTP requests, and then parsing the returned HTML to extract the relevant Simpsons scenes based on quotes.

## Critical Implementation Paths

### Command Execution Path
1. User inputs a command via CLI
2. Kong parses the command and arguments
3. The appropriate command handler is invoked
4. The command handler performs its operation
5. Results are returned to the user

### Frinkiac Quote Matching Path
1. User provides a prompt or quote
2. The Frinkiac command processes the input
3. The system constructs a URL query to the Frinkiac website (e.g., `https://frinkiac.com/?q=prompt`)
4. The system sends an HTTP request to the Frinkiac website
5. The HTML response is parsed to extract matching Simpsons quotes with confidence levels
6. The system extracts screen captures from the parsed HTML
7. Results are returned to the user

### Webhook Event Handling Path
1. External service sends a webhook event
2. Smee client receives the event
3. Event is processed according to its type and content
4. Appropriate action is taken (e.g., finding a Simpsons reference)
5. Response is sent back if required

## Error Handling Strategy

The application employs a consistent error handling strategy:

1. Errors are logged with appropriate context using zerolog
2. User-facing errors are presented in a clear, actionable format
3. Critical errors result in a non-zero exit code
4. Context cancellation is properly handled for graceful shutdown

## Future Architectural Considerations

1. **Microservice Evolution**: As the application grows, consider splitting functionality into separate microservices.
2. **API Layer**: Add a REST or GraphQL API for programmatic access beyond CLI.
3. **Persistent Storage**: Implement caching or storage for frequently used quotes and results.
4. **Advanced NLP**: Integrate more sophisticated natural language processing for better prompt matching.
