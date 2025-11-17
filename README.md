# CodeSync 🚀

CodeSync is a version control system inspired by Git, implemented in Golang.

## Tech Stack 🛠

| **Category**        | **Technology**                                 |
| ------------------- | ---------------------------------------------- |
| **Core**            | Golang v1.24                                   |
| **Version Control** | Custom version control logic (inspired by Git) |
| **Utilities**       | Cobra (CLI framework)                          |

## Commands ⚙️

| **Category** | **Technology**                                                    | **Tested** |
| ------------ | ----------------------------------------------------------------- | ---------- |
| `add`        | Add the selected files to the staging area                        | ✅         |
| `branch`     | Branch management (`new`, `drop`, `switch`, `default`, `current`) | ✅         |
| `commit`     | Commit the staged files                                           | ✅         |
| `config`     | Config management (`get\|set <default-branch\|email\|username`)   | ✅         |
| `history`    | List all commits for the current branch                           | ✅         |
| `init`       | Initialize the CSync version control system                       | ✅         |
| `purge`      | Purge CSync and all its data. THIS COMMAND IS IRREVERSIBLE!       | ✅         |
| `remove`     | Remove the selected files from the staging area                   | ✅         |
| `status`     | List the files that are staged for commit, tracked, untracked     | ✅         |
| `workdir`    | List the files that are committed                                 | ✅         |

## License 📜

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Getting Started 💻

1. Clone the repository:

```bash
git clone https://github.com/denesbeck/code-sync.git
```

2. Change the directory:

```bash
cd code-sync
```

3. Install the dependencies:

```bash
go mod tidy
```

4. Build the project:

```bash
go build -o csync ./cmd/csync
```

5. Run the project:

```bash
./csync
```

## Running Tests 🧪

To run the test suite, use the provided test script:

```bash
cd cmd/csync
./run-tests.sh
```

**Important:** Tests must be run using `run-tests.sh` instead of `go test` directly. The script sets the required `CSYNC_ENV=test` environment variable, which ensures tests run in an isolated namespace to prevent conflicts with your actual `.csync` directory.
