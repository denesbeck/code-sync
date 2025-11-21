# Contributing to Nexio

Thank you for your interest in contributing to Nexio! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful, inclusive, and considerate in all interactions.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/nexio.git
   cd nexio
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/denesbeck/nexio.git
   ```
4. Create a new branch for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.24.4 or higher
- Git
- Make (optional, for using Makefile targets)

### Install Dependencies

```bash
go mod download
```

### Build the Project

```bash
go build -o nexio ./cmd/nexio
```

### Install Git Hooks

We use Git hooks to ensure code quality. Install them with:

```bash
make install-hooks
```

This will automatically run linting and formatting checks before each commit.

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug fixes**: Fix issues reported in the issue tracker
- **New features**: Implement new functionality
- **Documentation**: Improve README, code comments, or add examples
- **Tests**: Add or improve test coverage
- **Code quality**: Refactoring, performance improvements
- **UI/UX**: Enhance the CLI interface and user experience

### Finding Work

- Check the [issue tracker](https://github.com/denesbeck/nexio/issues) for open issues
- Look for issues labeled `good first issue` or `help wanted`
- Comment on an issue to express interest before starting work
- Feel free to propose new features or improvements

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting (enforced by pre-commit hook)
- Run `go vet` to catch common mistakes (enforced by pre-commit hook)
- Write clear, descriptive variable and function names
- Keep functions small and focused on a single responsibility

### Code Organization

- Place all CLI commands in `cmd/nexio/`
- Keep related functionality grouped together
- Use meaningful package names
- Avoid circular dependencies

### Documentation

- Add comments to exported functions, types, and constants
- Use GoDoc conventions for documentation comments
- Include usage examples in comments where helpful
- Update README.md if adding user-facing features

### Error Handling

- Return errors rather than panicking
- Provide clear, actionable error messages
- Use error wrapping for context: `fmt.Errorf("context: %w", err)`
- Handle all error returns explicitly

## Commit Guidelines

### Commit Message Format

Use clear, descriptive commit messages following this format:

```
<type>: <subject>

<body (optional)>

<footer (optional)>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring without functionality changes
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks, dependency updates

**Examples:**

```
feat: add commit message validation

Implements automatic validation of commit messages to ensure
they follow the project's conventions.

Closes #42
```

```
fix: prevent panic when staging non-existent files

Added validation to check file existence before staging.
Returns clear error message to the user.
```

### Commit Best Practices

- Make atomic commits (one logical change per commit)
- Write clear, descriptive commit messages
- Reference issue numbers in commit messages (e.g., "Fixes #123")
- Keep commits focused and avoid mixing unrelated changes

## Pull Request Process

### Before Submitting

1. **Ensure your code builds**: Run `go build -o nexio ./cmd/nexio`
2. **Run tests**: Execute `bash ./scripts/run-tests.sh`
3. **Run linting**: Execute `go vet ./...`
4. **Format code**: Run `gofmt -w .`
5. **Update documentation**: Update README.md if needed
6. **Test manually**: Verify your changes work as expected

### Submitting a Pull Request

1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Open a Pull Request from your branch to `main`

3. Fill out the PR template completely

4. Link related issues using keywords (e.g., "Fixes #123", "Relates to #456")

5. Wait for CI checks to pass

6. Respond to review feedback promptly

### PR Requirements

- [ ] All tests pass
- [ ] Code is formatted with `gofmt`
- [ ] No linting errors from `go vet`
- [ ] Documentation updated if needed
- [ ] Commit messages follow guidelines
- [ ] PR description clearly explains changes

### Review Process

- Maintainers will review your PR within a few days
- Address feedback by pushing new commits to your branch
- Once approved, a maintainer will merge your PR
- Feel free to ask questions or request clarification on feedback

## Testing

### Running Tests

Execute the test suite:

```bash
bash ./scripts/run-tests.sh
```

**Important**: Always use the test script, not `go test` directly. The script sets required environment variables.

### Writing Tests

- Write tests for new functionality
- Use table-driven tests where appropriate
- Test both success and error cases
- Keep tests focused and independent
- Use descriptive test names: `TestFunctionName_Scenario_ExpectedBehavior`

### Test Coverage

- Aim for high test coverage on new code
- Don't sacrifice test quality for coverage numbers
- Focus on testing critical paths and edge cases

## Reporting Bugs

Found a bug? Please help us fix it!

1. **Check existing issues** to avoid duplicates
2. **Use the bug report template** when creating a new issue
3. **Provide detailed information**:
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment details (OS, Go version)
   - Error messages or logs

## Suggesting Features

We welcome feature suggestions!

1. **Check existing issues** for similar requests
2. **Use the feature request template**
3. **Describe the problem** you're trying to solve
4. **Explain your proposed solution**
5. **Consider alternatives** and trade-offs

## Questions?

- Open a [discussion](https://github.com/denesbeck/nexio/discussions) for general questions
- Comment on relevant issues for specific questions
- Check existing documentation and issues first

## License

By contributing, you agree that your contributions will be licensed under the same MIT License that covers the project.

---

Thank you for contributing to Nexio!
