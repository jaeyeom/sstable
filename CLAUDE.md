# Claude Code Instructions for sstable

## Quick Commands

- Run `make` to format, test, and lint your code
- Run `make test` for tests only
- Run `make format` for formatting only
- Run `make lint` for linting only


## Go Development

- Run `bazel run //:gazelle` to generate BUILD.bazel files
- Run `golangci-lint run ./...` for Go linting



## Bazel Commands

- `bazel build //...` - Build all targets
- `bazel test //...` - Run all tests
