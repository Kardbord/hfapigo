# Contributing to hfgo

Thank you for your interest in contributing to hfgo!
This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions.
We're committed to providing a welcoming environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.25 or later
- Git
- Hugging Face API key (for running integration tests)
  - Get one at https://huggingface.co/settings/tokens
  - Set as `HUGGING_FACE_TOKEN` environment variable

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/hfgo.git
   cd hfgo
   ```
3. Add the upstream remote to track the original repository:
   ```bash
   git remote add upstream https://github.com/Kardbord/hfgo.git
   ```
4. Verify your remotes are configured correctly:
   ```bash
   git remote -v
   ```
   You should see:
   ```
   origin    https://github.com/your-username/hfgo.git (fetch)
   origin    https://github.com/your-username/hfgo.git (push)
   upstream  https://github.com/Kardbord/hfgo.git (fetch)
   upstream  https://github.com/Kardbord/hfgo.git (push)
   ```
5. Fetch the latest changes from upstream:
   ```bash
   git fetch upstream
   ```
6. Create a feature branch from the latest upstream main:
   ```bash
   git checkout -b your-feature-or-fix upstream/main
   ```

### Building and Testing

Run the build script to format, lint, and test:

```bash
./tools/build.sh
```

**Note**: The build script runs integration tests, which will make API calls to
Hugging Face and may incur costs on your account. If you want to skip integration
tests during development, run the individual commands below instead.

Or run individual commands:

```bash
# Format code
gofmt -s -w .

# Tidy dependencies
go mod tidy

# Vet code
go vet ./...

# Lint
golangci-lint run --fix ./...

# Build
go build ./...

# Unit tests
go test -timeout 600s -v ./...

# Race condition detection
go test -race -timeout 600s -v ./...

# Integration tests (requires HUGGING_FACE_TOKEN)
go test -tags=integration -timeout 600s -v ./...
```

## Pull Requests

### Before Creating a PR

1. Ensure your branch is up to date with main:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Run all checks locally:
   ```bash
   ./tools/build.sh
   ```

3. Write or update tests for your changes

4. Update documentation if needed

### Creating a PR

1. Push your feature branch to your fork
2. Create a PR on GitHub with:
   - **Title**: Follow [Conventional Commits](https://www.conventionalcommits.org/) format (enforced by PR title check workflow)
     - Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`
     - Examples: `feat: add support for tool calling`, `fix: handle nil context in streaming`
     - Breaking changes: Use `feat!:` or `fix!:` to indicate breaking changes
   - **Description**: Explain what changes and why
   - **References**: Link to related issues (e.g., "Fixes #123")

### PR Title Examples

```
feat: add support for tool calling
fix: handle nil context in streaming
feat!: remove deprecated Chat.Complete method
docs: update API documentation
refactor: reorganize internal packages
test: add integration tests for streaming
```

### Code Review

- Be responsive to feedback
- Make requested changes in new commits (don't force push during review)
- Ask for clarification if feedback is unclear
- Be polite to reviewers (they should also be polite to you)

## Release Process

hfgo uses an **automated release workflow** based on
[Conventional Commits](https://www.conventionalcommits.org/)
in the PR title.

### How Releases Work

1. **PR is merged to main** with a conventional commit title
2. **release-please analyzes PR titles** since the last release
3. **Release PR is created automatically** with:
   - Updated version in `internal/sdkversion/version.go`
   - Updated `go.mod` (for major version bumps)
   - Auto-generated changelog
4. **Maintainer reviews and merges Release PR**
5. **GitHub release is created automatically** with:
   - Git tag (e.g., v4.1.0)
   - Release notes with changelog
   - Download links

### Version Bumping

Versions are determined automatically based on PR titles:

- **Major (v4.0.0 → v5.0.0)**: Breaking API changes (`feat!:` or `fix!:`)
- **Minor (v4.0.0 → v4.1.0)**: New features (`feat:`)
- **Patch (v4.0.0 → v4.0.1)**: Bug fixes (`fix:`) or refactoring
- **No bump**: Documentation, tests, or chore commits

### Your Role in Releases

As a contributor, your role is:
- ✅ Write PR titles following Conventional Commits format
- ✅ Mark breaking changes with `!:` if applicable
- ✅ Ensure PR title accurately describes your changes

Maintainers handle:
- ✅ Reviewing Release PRs
- ✅ Merging Release PRs
- ✅ Managing actual releases

## Testing

### Build Tags

Test files must include appropriate Go build tags to ensure they run in the correct context:

- **Unit Tests**: Add `//go:build !integration` at the top of the file
  - These tests run by default with `go test ./...`
  - Should not require external services or API calls

- **Integration Tests**: Add `//go:build integration` at the top of the file
  - These tests only run with `-tags=integration`
  - Require a valid Hugging Face API token
  - Should be marked with the `_integration_test.go` filename suffix

**Example unit test**:
```go
//go:build !integration

package hfgo

import "testing"

func TestSomething(t *testing.T) {
    // Test code
}
```

**Example integration test**:
```go
//go:build integration

package hfgo

import "testing"

func TestIntegration(t *testing.T) {
    // Integration test code
}
```

### Unit Tests

All code changes should include tests:

```bash
go test -timeout 600s -v ./...
```

### Integration Tests

Integration tests require a valid Hugging Face API token:

```bash
HUGGING_FACE_TOKEN=your_token go test -tags=integration -timeout 600s -v ./...
```

### Race Condition Detection

Always run race detection before submitting:

```bash
go test -race -timeout 600s -v ./...
```

## Documentation

### Code Comments

- Document all exported functions, types, and fields
- Use clear, concise language
- Explain the "why" not just the "what"
- Include examples for complex functions

### README and Docs

- Update README.md if adding new features
- Update examples/ if adding new features
- Keep documentation in sync with code

## Code Style

The project uses:
- **gofmt** for formatting
- **gofumpt** for additional formatting rules
- **goimports** for import organization
- **golines** for line length management
- **golangci-lint** for linting

All of these are run by `./tools/build.sh`.

## Concurrency and Safety

The SDK prioritizes **concurrency safety**:

- Clients are immutable value types
- Services are lightweight snapshots of client options
- No shared mutable state between goroutines
- All public APIs are safe for concurrent use

When making changes:
- Don't introduce mutable state in clients or services
- Test with `go test -race` to catch race conditions
- Document concurrency guarantees in comments

## Common Mistakes to Avoid

1. **Mutating requests**: Don't modify requests passed to methods by pointer after passing them
2. **Not closing streams**: Always call `Close()` on `ChatStream` or `RawStream`
3. **Sharing mutable HTTP clients**: If injecting HTTP clients, ensure thread-safety
4. **Breaking API contracts**: Changing function signatures is a breaking change
5. **Ignoring context**: Always respect context cancellation and timeouts

## Need Help?

- Check existing issues and discussions
- Read the README.md and AGENTS.md
- Review examples in the examples/ directory
- Don't be afraid to open a new issue

## Recognition

Contributors are recognized in:
- GitHub contributor graph
- Release notes (via commit attribution)

Thank you for contributing to hfgo! 🙏
