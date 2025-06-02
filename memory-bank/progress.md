# Progress: Billy Bot

## What Works

As the project is in its initial development phase, the following components are currently working:

1. **Basic CLI Framework**:
   - Command-line parsing using Kong
   - Subcommand structure for Frinkiac functionality
   - Early prototype of Smee client for webhook testing
   - Environment variable loading from .env files
   - Logging setup with zerolog
   - Signal handling for graceful shutdown

2. **Project Structure**:
   - Go module setup with dependencies
   - Basic package organization
   - Command pattern implementation

3. **Frinkiac Integration**:
   - API integration with Frinkiac's JSON endpoint (`/api/search`)
   - URL query construction and parameter handling
   - JSON parsing for extracting screen captures information
   - Unit tests with saved API responses

## What's Left to Build

The following components and features are still pending implementation:

1. **Frinkiac Integration**:
   - ✅ API integration with Frinkiac's JSON endpoint (`/api/search`)
   - ✅ URL query construction (e.g., `https://frinkiac.com/api/search?q=quote`)
   - ✅ JSON parsing for extracting screen captures information
   - Result formatting and presentation

2. **GitHub Integration** (Primary Goal):
   - GitHub API client for issues, tasks, and pull requests
   - Webhook handling for GitHub events
   - Comment formatting for GitHub responses
   - Authentication and permissions handling

3. **Prompt Categorization**:
   - Prompt analysis logic
   - Quote matching algorithms
   - Confidence scoring
   - Fuzzy matching and synonyms support

4. **Future Platform Integrations**:
   - Discord integration
   - Slack integration
   - Abstraction layer for supporting multiple platforms

5. **Testing Infrastructure**:
   - ✅ Unit tests for Frinkiac client API integration
   - Integration tests for API interactions
   - End-to-end tests for command execution
   - ✅ Test fixtures for API responses

6. **Documentation**:
   - Usage instructions
   - API documentation
   - Example use cases
   - Contribution guidelines

## Current Status

The project is in the **early development phase**, with the basic structure and CLI framework in place. The focus is currently on implementing Step 1 as outlined in the README: categorizing a prompt into a Simpson's quote.

The example.json file suggests that some progress has been made on the Frinkiac integration, with the ability to search for quotes with confidence levels. However, this functionality may not be fully implemented or integrated into the CLI yet.

The main.go file shows that the command structure is in place, but the actual implementation of the commands (in the pkg directory) is still pending or in progress.

## Known Issues

As the project is in early development, there are several known issues and limitations:

1. **Incomplete Functionality**:
   - The core functionality of categorizing prompts into Simpson's quotes is still under development.
   - The GitHub integration is not yet implemented.
   - The Smee client is only an early prototype for webhook testing.
   - ✅ The Frinkiac client now successfully retrieves and parses results from the Frinkiac API endpoint. We've switched from HTML parsing to using the JSON API at `/api/search`, which provides more reliable and structured data.

2. **Missing Tests**:
   - Test coverage is limited or non-existent at this stage.
   - No CI/CD pipeline is in place for automated testing.

3. **Documentation Gaps**:
   - Limited documentation on how to use the application.
   - No API documentation for the Frinkiac and Smee functionality.

4. **Development Environment**:
   - The dev container mentioned in the README may not be fully up to date.

## Evolution of Project Decisions

### Initial Concept
The project started with a humorous premise: replacing a human (Billy) who was unreliable in providing Simpsons screen captures with a bot that would consistently deliver appropriate references. This joke at Billy's expense is not just a footnote but a fundamental aspect of the project's identity - the bot exists specifically to mock Billy's unreliability and "replace" him with superior technology.

### Technical Approach
1. **Language Selection**: Go was chosen for its simplicity, performance, and strong concurrency support, which is well-suited for handling webhook events and API interactions.

2. **CLI Framework**: Kong was selected as the CLI parsing library due to its declarative, struct-based approach, which aligns well with Go's design philosophy.

3. **Logging Strategy**: Zerolog was chosen for structured logging, providing better log analysis capabilities compared to traditional logging approaches.

4. **Command Pattern**: The decision to use the command pattern for CLI interactions provides a clean separation of concerns and makes it easier to add new functionality in the future.

5. **Frinkiac API Integration**: Rather than implementing screen cap selection from scratch, the project leverages the existing Frinkiac website through its JSON API, focusing on the prompt categorization and structured data parsing. This approach was updated from the initial HTML scraping strategy to provide more reliable results.

### Ongoing Considerations

1. **Frinkiac Integration Strategy**: We've updated our approach for interacting with the Frinkiac website from HTML parsing to using the JSON API endpoint. This provides several advantages:
   - More reliable data extraction with structured JSON responses
   - Less susceptibility to website UI changes
   - Cleaner code with standard JSON parsing instead of complex HTML traversal
   - Better performance with smaller response payloads
   - Easier testing with predictable response formats

2. **Command Structure**: Decisions about the optimal command structure and options for the CLI are still being made, including parameters for the Frinkiac command and output format options.

3. **Error Handling**: The project is developing a consistent approach to error handling, focusing on user-friendly error messages and appropriate logging.

4. **Performance Optimization**: Considerations about response time targets, concurrency models, and resource usage constraints are ongoing.

### Humor as a Design Principle
Throughout development, the project has maintained its humorous foundation:
- Error messages and logs are designed with subtle jabs at Billy's expense
- The bot's superiority over Billy is emphasized in various aspects of the design
- User-facing elements maintain the comedic premise of Billy being replaced
- Technical decisions are made with consideration for how they support the joke

This humorous approach isn't merely decorative - it informs design decisions, communication style, and feature prioritization. The project team recognizes that while building a technically sound application is important, maintaining the comedic premise of making fun of Billy is equally essential to the project's success.

### Future Direction

The project roadmap suggests a phased approach:
1. First, focus on the core functionality of categorizing prompts into Simpson's quotes.
2. Then, implement the GitHub integration for issues, tasks, and pull requests.
3. If successful, expand to include Discord and Slack integrations.
4. Finally, consider additional features like more sophisticated matching algorithms, caching, and potentially a web interface.

Throughout all phases, the project will maintain its humorous tone and continue to emphasize how the bot outperforms Billy in every metric. This includes potentially developing a "Billy would have failed here" metric to highlight situations where the bot succeeds where Billy would have likely failed.

This incremental approach allows for early validation of the core concept before investing in more complex features and additional platform integrations, while ensuring the humorous premise remains central to the project's identity.
