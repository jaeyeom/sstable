# sstable


[![Copier](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/copier-org/copier/master/img/badge/badge-grayscale-inverted-border-orange.json)](https://github.com/copier-org/copier)


SSTable implementation compatible with https://github.com/mariusaeriksen/sstable

## Prerequisites

- [Bazelisk](https://github.com/bazelbuild/bazelisk) (Bazel version manager)
- Go 1.25.5+
- [golangci-lint](https://golangci-lint.run/)



## Development

```bash
# Format, test, and lint
make

# Individual commands
make format   # Format code
make test     # Run tests
make lint     # Run linters
make fix      # Auto-fix lint issues
```

## Build System

This project uses Bazel for builds:

```bash
bazel build //...    # Build all targets
bazel test //...     # Run all tests
bazel run //:gazelle # Generate BUILD.bazel files

```

## Project Structure

```
sstable/
├── MODULE.bazel     # Bazel dependencies
├── BUILD.bazel      # Root build file
├── Makefile         # Development commands
├── CLAUDE.md        # AI assistant instructions
├── VERSION-BUMPS.md # Guide for upgrading versions
├── go.mod           # Go module
├── .golangci.yml    # Go linting config


└── docs/claude/     # Reference for Claude Code
```
