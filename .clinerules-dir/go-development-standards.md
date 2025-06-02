# Go Development Standards

## Brief overview
This set of guidelines outlines the development standards for Go projects in this repository. It focuses on code quality, testing, and documentation practices to ensure maintainable and reliable code.

## Code quality requirements
- All code must be formatted with `go fmt` before committing
- All code must pass `golint` checks with no warnings or errors
- All tests must pass with `go test` before completing any task
- These checks must be run before considering any task complete
- Exported types and functions must have proper documentation comments
- Avoid stuttering in function and type names (e.g., use `pkg.Function` not `pkg.PkgFunction`)

## Documentation standards
- All exported types and functions must have documentation comments
- Comments should explain the purpose, not just restate the name
- For complex functions, include examples of usage when appropriate
- Document any non-obvious behavior or edge cases
- Keep comments up-to-date when modifying code

## Testing practices
- Write unit tests for all new functionality
- Use the testify package (require and/or assert) for clearer assertions and better error messages
- Test both success and error cases
- Mock external dependencies when appropriate
- Maintain test fixtures for API responses

## Error handling
- Use explicit error checking (Go style)
- Provide contextual error messages
- Log errors at appropriate levels
- Return errors rather than handling them internally when appropriate

## Project structure
- Maintain clean separation of concerns
- Package structure should reflect the domain model
- Use clear interfaces between components
- Follow Go conventions for directory and file naming

## Development workflow
- Before completing any task, run:
  1. `go fmt ./...` to format all code
  2. `golint ./...` to check for style issues
  3. `go test ./...` to ensure all tests pass
- Update documentation when making significant changes
- Keep commits focused on single concerns
