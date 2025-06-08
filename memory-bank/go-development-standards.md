# Go Development Standards

This document outlines the coding standards and patterns preferred for Go development in this project.

## Code Organization Patterns

### Functional vs Object-Oriented Style

This project prefers **functional style** over object-oriented patterns for external API clients and business logic.

#### Preferred: Functional Style

Functions should take dependencies as parameters rather than embedding them in structs.

```go
// ✅ Preferred: Functional approach
func DoSomeBusinessCode(ctx context.Context, client *http.Client, config Config, param string) (BusinessResults, error) {
    // Implementation uses the client parameter
    resp, err := client.Do(req)
    // ...
}
```

**Benefits:**
- Explicit dependencies make testing easier
- No hidden state or configuration
- Clear function signatures show exactly what's needed
- Easy to mock dependencies in tests
- Follows Go's preference for composition over inheritance

#### Avoid: Object-Oriented Style

Avoid embedding dependencies in struct fields and using methods.

```go
// ❌ Avoid: Object-oriented approach
type BusinessClient struct {
    client *http.Client
    config Config
}

func (c *BusinessClient) DoSomeBusinessCode(ctx context.Context, param string) (BusinessResults, error) {
    // Implementation uses c.client
    resp, err := c.client.Do(req)
    // ...
}
```

**Why to avoid:**
- Hidden dependencies make testing harder
- Stateful objects can lead to unexpected behavior
- Configuration is less explicit
- Harder to reason about what the function needs

### Configuration Patterns

When functions need configuration, use explicit config parameters:

```go
// Configuration struct
type Config struct {
    BaseURL string
    Timeout time.Duration
}

// Default configuration function
func DefaultConfig() Config {
    return Config{
        BaseURL: "https://api.example.com",
        Timeout: 30 * time.Second,
    }
}

// Helper function for common client setup
func NewHTTPClient() *http.Client {
    return &http.Client{
        Timeout: 10 * time.Second,
    }
}
```

### Refactoring from OO to Functional

When refactoring existing object-oriented code to functional style:

1. **Identify the embedded dependencies** (e.g., `http.Client`, database connections)
2. **Extract configuration** into a separate `Config` struct
3. **Convert methods to functions** that take dependencies as parameters
4. **Create helper functions** for common setup (e.g., `NewHTTPClient()`)
5. **Update all calling code** to pass dependencies explicitly
6. **Update tests** to use the new function signatures

#### Example Refactoring

**Before (Object-Oriented):**
```go
type APIClient struct {
    client *http.Client
    baseURL string
}

func New(opts ...Option) *APIClient {
    return &APIClient{
        client: &http.Client{Timeout: 10*time.Second},
        baseURL: "https://api.example.com",
    }
}

func (c *APIClient) GetData(ctx context.Context, id string) (*Data, error) {
    // Uses c.client and c.baseURL
}
```

**After (Functional):**
```go
type Config struct {
    BaseURL string
}

func DefaultConfig() Config {
    return Config{BaseURL: "https://api.example.com"}
}

func NewHTTPClient() *http.Client {
    return &http.Client{Timeout: 10*time.Second}
}

func GetData(ctx context.Context, client *http.Client, config Config, id string) (*Data, error) {
    // Uses client and config parameters
}
```

## Testing

Functional style makes testing much easier:

```go
func TestGetData(t *testing.T) {
    // Easy to mock the HTTP client
    mockClient := &http.Client{
        Transport: &mockTransport{
            response: &http.Response{...},
        },
    }
    
    config := Config{BaseURL: "https://test.example.com"}
    
    result, err := GetData(ctx, mockClient, config, "test-id")
    // Assert results
}
```

## Error Handling

Follow Go's explicit error handling patterns:

- Return errors as the last return value
- Use `fmt.Errorf` to wrap errors with context
- Don't panic in business logic; return errors instead

## Documentation

- Use Go doc comments for all exported functions
- Explain what the function does, not how it does it
- Document parameters and return values when not obvious
- Include usage examples for complex functions

## Package Organization

- Keep related functionality together in packages
- Use clear, descriptive package names
- Avoid circular dependencies
- Export only what needs to be public

## Git Workflow and File Management

### Prefer Git Moves for Better History

When refactoring code that involves moving functionality between files or packages, prefer using `git mv` to preserve git history and make the changes easier to track.

#### When Moving Files
```bash
# ✅ Preferred: Use git mv to preserve history
git mv old/path/file.go new/path/file.go
# Then edit the file to update package declaration and imports
```

#### When Moving Functions Between Files
If moving functions between existing files:

1. **First**: Use `git mv` to create the new file structure if needed
2. **Then**: Move the functions using normal edit operations
3. **Document**: In commit messages, clearly describe what was moved where

#### Benefits of Git Moves
- **Preserves file history**: `git log --follow` can track the file across moves
- **Better diffs**: Git can better understand what changed vs. what moved
- **Easier reviews**: Reviewers can see the relationship between old and new locations
- **Cleaner history**: Shows intent to move rather than delete/create

#### Example Workflow
```bash
# Moving a package from old location to new location
git mv pkg/old/client.go pkg/new/client.go
# Edit the file to update package declaration
# Update imports in files that reference it
git add -A
git commit -m "refactor: move client from pkg/old to pkg/new"
```

## Code Style

- Use `gofmt` for formatting
- Use `golint` for style checking
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Keep functions focused and single-purpose
- Use meaningful variable names