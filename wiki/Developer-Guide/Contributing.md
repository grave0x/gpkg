# Contributing to gpkg

Thank you for your interest in contributing to gpkg! We welcome contributions of all kinds - bug reports, feature requests, documentation improvements, and code contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Communication](#communication)

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions. We're building a welcoming community.

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Git
- SQLite3 development libraries
- Basic understanding of package management concepts

### Setting Up Development Environment

1. **Fork and Clone**

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/gpkg.git
cd gpkg
```

2. **Add Upstream Remote**

```bash
git remote add upstream https://github.com/grave0x/gpkg.git
git fetch upstream
```

3. **Install Dependencies**

```bash
go mod download
go mod verify
```

4. **Build and Test**

```bash
# Build
go build -o bin/gpkg .

# Run tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
```

See the [Development Setup](Development-Setup) guide for detailed instructions.

## Development Workflow

### 1. Pick or Create an Issue

- Browse [open issues](https://github.com/grave0x/gpkg/issues)
- Look for issues labeled `good-first-issue` for beginner-friendly tasks
- Comment on an issue to claim it or ask for assignment
- If proposing a new feature, create an issue first to discuss

### 2. Create a Branch

Use descriptive branch names:

```bash
# Feature branches
git checkout -b feat/add-package-signing

# Bug fix branches
git checkout -b fix/checksum-verification

# Documentation branches
git checkout -b docs/improve-readme
```

### 3. Make Changes

- **Keep commits small and focused** - One logical change per commit
- **Write clear commit messages** - Follow conventional commits format:
  ```
  feat: add package signing support
  fix: correct checksum verification logic
  docs: update installation instructions
  test: add unit tests for manifest parser
  ```

### 4. Run Tests and Linters

Before committing:

```bash
# Format code
gofmt -s -w .

# Run linter
golangci-lint run

# Run tests
go test ./...

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### 5. Commit and Push

```bash
git add .
git commit -m "feat: add package signing support"
git push origin feat/add-package-signing
```

### 6. Open Pull Request

- Provide clear title and description
- Link to related issue: `Closes #123`
- Add screenshots/examples if relevant
- Request review from maintainers

## Coding Standards

### Go Style Guidelines

- **Follow** [Effective Go](https://golang.org/doc/effective_go.html)
- **Use** `gofmt` for formatting (non-negotiable)
- **Run** `golangci-lint` and fix all issues
- **Write** idiomatic Go code
- **Minimize** dependencies - only add when necessary

### Code Organization

- Keep functions small and focused (< 50 lines when possible)
- Use meaningful variable and function names
- Group related functionality in packages
- Avoid global state where possible

### Error Handling

```go
// Good: Return errors, don't panic
func doSomething() error {
    if err := validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    return nil
}

// Bad: Panicking in library code
func doSomething() {
    if err := validate(); err != nil {
        panic(err)  // Don't do this!
    }
}
```

### CLI Design Principles

- **Consistency**: Follow existing command patterns
- **Script-friendly**: Support `--json` output for automation
- **Interactive**: Provide helpful prompts and confirmations
- **Non-interactive**: Support `--yes` flag for automation
- **Dry-run**: Support `--dry-run` for safe preview
- **Informative**: Provide clear success/error messages

## Testing Requirements

### Unit Tests

- Write unit tests for all new functionality
- Aim for **>70% code coverage** on new code
- Use table-driven tests for multiple test cases
- Mock external dependencies

Example:

```go
func TestManifestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid manifest", validManifest, false},
        {"missing name", missingNameManifest, true},
        {"invalid checksum", invalidChecksumManifest, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateManifest(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateManifest() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

- Add integration tests for end-to-end workflows
- Located in `tests/integration/`
- Test actual command execution and file system operations

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/manifest

# With race detector
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
```

See the [Testing](Testing) guide for more details.

## Documentation

### Code Documentation

- Add godoc comments for exported functions, types, and packages
- Include examples in documentation where helpful

```go
// ParseManifest parses a YAML manifest file and validates its contents.
// It returns an error if the manifest is invalid or cannot be parsed.
//
// Example:
//     manifest, err := ParseManifest("./manifest.yaml")
//     if err != nil {
//         return err
//     }
func ParseManifest(path string) (*Manifest, error) {
    // Implementation
}
```

### README and Wiki

- Update `README.md` when adding user-facing features
- Update wiki pages for significant changes
- Add examples to `examples/` directory
- Keep documentation in sync with code

## Pull Request Process

### Before Opening a PR

- [ ] Tests pass locally
- [ ] Code is formatted (`gofmt`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main

### PR Description Template

```markdown
## Description
Brief description of changes

## Related Issue
Closes #123

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe how you tested the changes

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] Follows coding standards
```

### Review Process

1. **Automated Checks**: CI runs tests, linting, security scans
2. **Code Review**: Maintainers review code and provide feedback
3. **Iterate**: Address feedback and update PR
4. **Approval**: Maintainer approves changes
5. **Merge**: Maintainer merges PR

## Issue Guidelines

### Reporting Bugs

Use this template:

```markdown
**Description**
Clear description of the bug

**Steps to Reproduce**
1. Step one
2. Step two
3. See error

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- gpkg version: 
- OS: 
- Go version:

**Additional Context**
Any other relevant information
```

### Requesting Features

```markdown
**Problem**
Describe the problem this feature would solve

**Proposed Solution**
Your proposed solution

**Alternatives Considered**
Other solutions you've considered

**Additional Context**
Any other relevant information
```

## Communication

### Where to Ask Questions

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and discussions
- **Pull Requests**: Code review and implementation discussion

### Getting Help

- Check existing [documentation](../README.md)
- Search [existing issues](https://github.com/grave0x/gpkg/issues)
- Ask in [Discussions](https://github.com/grave0x/gpkg/discussions)
- Join development discussions in PRs

### Suggesting Improvements

We love suggestions! Open an issue or discussion to:
- Propose new features
- Suggest architectural improvements
- Recommend better approaches
- Share ideas for the roadmap

## Starter Tasks

Looking for a good first contribution? Look for issues labeled:

- `good-first-issue` - Beginner-friendly tasks
- `help-wanted` - Tasks where we need help
- `documentation` - Documentation improvements
- `enhancement` - Feature additions

## Recognition

Contributors are recognized in:
- GitHub Contributors page
- Release notes for significant contributions
- Wiki acknowledgments section

## Questions?

Don't hesitate to ask! We're here to help.

- Open an [issue](https://github.com/grave0x/gpkg/issues)
- Start a [discussion](https://github.com/grave0x/gpkg/discussions)
- Comment on an existing issue

Thank you for contributing to gpkg! 🎉

---

**See Also:**
- [Development Setup](Development-Setup) - Set up your dev environment
- [Architecture](Architecture) - Understand the codebase
- [Testing](Testing) - Testing guidelines
- [Release Process](Release-Process) - How releases work
