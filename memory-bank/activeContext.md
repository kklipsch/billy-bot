# Active Context: Billy Bot

## Current Work Focus

The current focus of the Billy Bot project is on implementing Step 1 as outlined in the README: **Categorizing a prompt into a Simpson's quote**. This is the foundational capability that will enable the bot to match user inputs with appropriate Simpsons references, with the primary goal of integrating with GitHub for issues, tasks, and pull requests.

Key aspects of this work include:
- Developing the core logic for analyzing prompts
- Implementing the matching algorithm to find relevant Simpsons quotes
- Integrating with the Frinkiac website to retrieve corresponding screen captures through web scraping
- Setting up the command-line interface for the Frinkiac functionality
- Preparing for GitHub integration as the first platform target

## Recent Changes

As the project is in its initial development phase, the recent changes include:
- Setting up the basic project structure with Go modules
- Implementing the CLI framework using Kong
- Creating the command structure for Frinkiac functionality
- Setting up a prototype Smee client for early webhook testing
- Setting up logging with zerolog
- Implementing context-based cancellation for graceful shutdown
- Adding environment variable loading with godotenv
- Adding GitHub Action for code quality checks (gofmt, golint, and tests)
- Fixing all golint issues in the codebase:
  - Adding proper documentation comments to all exported types and functions in the OpenRouter package
  - Fixing a stuttering function name (renamed OpenRouterCall to Call)
  - Updating references to the renamed function

## Next Steps

The immediate next steps for the project are:

1. **Complete the Frinkiac Integration**:
   - Implement web scraping functionality for the Frinkiac website
   - Add quote search functionality using URL queries (e.g., `https://frinkiac.com/?q=quote`)
   - Implement HTML parsing to extract screen captures from responses

2. **Develop GitHub Integration**:
   - Create GitHub API client for interacting with issues, tasks, and pull requests
   - Implement event handling for GitHub webhooks
   - Design response format for GitHub comments

3. **Enhance Prompt Categorization**:
   - Develop more sophisticated matching algorithms
   - Implement confidence scoring for quote matches
   - Add support for fuzzy matching and synonyms

4. **Add Testing**:
   - Write unit tests for core functionality
   - Implement integration tests for API interactions
   - Add end-to-end tests for command execution

5. **Prepare for Future Platform Integrations**:
   - Design modular architecture to support Discord and Slack in the future
   - Document integration points for additional platforms
   - Create abstraction layers for platform-specific functionality

## Active Decisions and Considerations

### Frinkiac Integration Strategy
We've updated our approach for interacting with the Frinkiac website:
- Discovered and implemented the Frinkiac JSON API endpoint (`/api/search`) instead of HTML parsing
- Created structured data types to represent the API response
- Implemented proper JSON parsing and conversion to our internal data structures
- Added unit tests with saved API responses to ensure correct parsing
- Maintained compatibility with existing code that uses the parsed results

**Current Status**: We have implemented a client that successfully sends quote requests to the Frinkiac API (e.g., `https://frinkiac.com/api/search?q=Everything%27s+coming+up+Milhouse%21`) and correctly parses the JSON response. The client now extracts season, episode, and ID information from the API response and constructs the appropriate image paths. This is more robust than the previous HTML parsing approach since it uses a structured API response rather than trying to parse dynamically generated HTML.

### Command Structure
Deciding on the optimal command structure and options for the CLI:
- What parameters should be available for the Frinkiac command
- How to handle different output formats (text, JSON, etc.)
- Whether to add interactive mode capabilities

### Error Handling Strategy
Developing a consistent approach to error handling:
- How to present errors to users in a helpful way
- When to retry failed API requests
- How to log errors for debugging

### Performance Optimization
Considering performance aspects:
- Response time targets for quote matching
- Concurrency model for handling multiple requests
- Resource usage constraints

## Important Patterns and Preferences

### Code Organization
- Preference for clean separation of concerns
- Package structure that reflects the domain model
- Clear interfaces between components

### Code Quality Standards
- All code must be formatted with gofmt
- All code must pass golint checks
- All tests must pass before merging
- GitHub Actions automatically verify these requirements on push and pull requests
- Tests should use the testify package (require and/or assert) for clearer assertions and better error messages

### Error Handling
- Explicit error checking (Go style)
- Contextual error messages
- Appropriate logging at different levels

### Configuration Management
- Environment variables for configuration
- Sensible defaults with override capabilities
- Clear documentation of configuration options

### Command Design
- Consistent command structure
- Self-documenting help text
- Intuitive option naming

### Project Tone and Communication Style
- Maintain a humorous tone that playfully mocks Billy throughout the codebase and user-facing elements
- Include subtle jabs at Billy's expense in error messages, logs, and output
- Ensure the bot occasionally references its superiority over Billy in reliability and performance
- Balance technical functionality with the comedic premise of replacing Billy
- Remember that the project is fundamentally a joke with Billy as the butt of it

## Learnings and Project Insights

### Initial Observations
- The Frinkiac website provides a rich source of Simpsons content but requires careful HTML parsing
- Command-line tools benefit greatly from intuitive, well-documented interfaces
- Go's concurrency model is well-suited for handling webhook events

### Technical Challenges
- Matching natural language prompts to specific quotes requires balancing precision and recall
- Handling the variability in webhook payloads requires robust parsing and validation
- Ensuring consistent performance across different operating environments

### Successful Approaches
- Using structured logging from the start has improved debugging capabilities
- The command pattern provides a clean way to organize CLI functionality
- Context-based cancellation ensures proper resource cleanup

### Areas for Improvement
- More comprehensive documentation would benefit future development
- Additional test coverage would increase confidence in the codebase
- Better error messages would improve the user experience (while maintaining the humorous tone at Billy's expense)
- More consistent integration of the humorous premise throughout the codebase and user experience
- Development of a "Billy would have failed here" metric to highlight the bot's superiority
- Refactoring existing tests to use the testify package for improved readability and better failure messages
