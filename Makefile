.PHONY: all check check-format format test lint fix clean test-bazel tidy

all: tidy format test fix

tidy: MODULE.bazel.lock go.sum

MODULE.bazel.lock: MODULE.bazel go.sum
	bazel mod tidy

check: check-format test lint

check-format:
	bazel test //tools/format/...:all

format:
	bazel run //tools/format

test: test-bazel

lint: lint-go

# Target fix is best-effort autofix for lint issues. If autofix is not
# available, it still runs lint checks.
fix: fix-go

test-bazel:
	bazel test //...

# Go lint and fix isn't integrated with bazel yet. Nogo is a good option.
.PHONY: lint-go fix-go

lint-go: go.sum
	golangci-lint run ./...

fix-go: go.sum
	golangci-lint run --fix ./...

go.sum: go.mod
	bazel run @rules_go//go -- mod tidy


clean:
	bazel clean --async
