# CodeSync

[![CI](https://github.com/denesbeck/code-sync/actions/workflows/main.yml/badge.svg)](https://github.com/denesbeck/code-sync/actions/workflows/main.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24.4-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight version control system inspired by Git, built from scratch in Go. CodeSync implements core version control concepts including staging, commits, branching, and history tracking.

## Overview

CodeSync is an educational project that demonstrates the fundamental principles behind modern version control systems. It provides a simplified implementation of Git-like functionality, making it easier to understand how version control works under the hood.

**Key Features:**
- Stage and commit file changes
- Branch management (create, switch, delete)
- Commit history tracking
- File status monitoring
- User configuration management
- Isolated testing environment

## Prerequisites

- Go 1.24.4 or higher
- Unix-like environment (Linux, macOS) or Windows with Git Bash

## Installation

1. Clone the repository:

```bash
git clone https://github.com/denesbeck/code-sync.git
cd code-sync
```

2. Install dependencies:

```bash
go mod download
```

3. Build the binary:

```bash
go build -o csync ./cmd/csync
```

4. (Optional) Add to PATH:

```bash
# Add to your shell profile (.bashrc, .zshrc, etc.)
export PATH="$PATH:/path/to/code-sync"
```

## Usage

### Initialize a Repository

```bash
./csync init
```

### Configure User Settings

```bash
./csync config set username "Your Name"
./csync config set email "your.email@example.com"
./csync config set default-branch "main"
```

### Basic Workflow

```bash
# Check file status
./csync status

# Add files to staging area
./csync add file1.txt file2.txt

# Commit changes
./csync commit -m "Initial commit"

# View commit history
./csync history
```

### Branch Management

```bash
# Create and switch to new branch
./csync branch new feature-branch

# List branches
./csync branch current

# Switch branches
./csync branch switch main

# Delete a branch
./csync branch drop feature-branch
```

## Available Commands

| Command    | Description                                                       |
|------------|-------------------------------------------------------------------|
| `init`     | Initialize the CodeSync version control system                    |
| `add`      | Add files to the staging area                                     |
| `remove`   | Remove files from the staging area                                |
| `commit`   | Commit staged changes with a message                              |
| `status`   | Display staged, tracked, and untracked files                      |
| `history`  | List all commits for the current branch                           |
| `branch`   | Manage branches (new, drop, switch, default, current)             |
| `workdir`  | List files in the current working directory state                 |
| `config`   | Get or set configuration values (username, email, default-branch) |
| `purge`    | Remove CodeSync and all its data (irreversible)                   |

For detailed command usage, run:

```bash
./csync [command] --help
```

## Development

### Development Setup

Install Git hooks to ensure code quality:

```bash
make install-hooks
```

This installs a pre-commit hook that automatically runs:
- `go vet ./...` - Lints code for common issues
- `gofmt -l .` - Checks code formatting

If formatting issues are detected, fix them with:

```bash
gofmt -w .
```

To remove the hooks:

```bash
make uninstall-hooks
```

### Running Tests

Run the test suite using the provided script:

```bash
bash ./scripts/run-tests.sh
```

**Important:** Always use `run-tests.sh` instead of `go test` directly. The script sets the `CSYNC_ENV=test` environment variable, which ensures tests run in an isolated namespace to prevent conflicts with your actual `.csync` directory.

### Project Structure

```
code-sync/
├── cmd/csync/          # CLI application and commands
├── scripts/            # Build and test scripts
├── .github/workflows/  # CI/CD configuration
└── go.mod              # Go module dependencies
```

## Built With

- [Go](https://go.dev/) - Programming language
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [pterm](https://github.com/pterm/pterm) - Terminal output styling
- [fatih/color](https://github.com/fatih/color) - Color output

## How It Works

CodeSync stores version control data in a `.csync` directory at the root of your project:

- **Staging area**: Tracks files prepared for commit
- **Commits**: Stores snapshots of file states with metadata
- **Branches**: Maintains separate lines of development
- **Configuration**: Stores user settings and repository configuration

Unlike Git, CodeSync uses a simpler file-based storage system and YAML for metadata, making the internals easier to understand and inspect.

## Limitations

CodeSync is designed for educational purposes and lacks several features found in production version control systems:

- No remote repository support
- No merge conflict resolution
- No diff visualization
- No file compression or delta storage
- Limited to local repositories

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Git's architecture and design principles
- Built as a learning project to understand version control internals
