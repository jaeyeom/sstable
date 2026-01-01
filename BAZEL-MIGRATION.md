# Bazel Migration Guide

This guide helps you complete your Bazel migration after running the migration script.

For small projects, running `bazel run //:gazelle` will generate BUILD.bazel files for your entire codebase, and you're done. However, for larger projects with complex dependencies, Gazelle may generate incorrect BUILD files or encounter issues. In such cases, a gradual migration approach works better—migrate packages one at a time by progressively removing entries from `.bazelignore`, starting from foundational packages with no internal dependencies.

## Migration Order

**Critical**: Migrate packages in dependency order, starting from the foundation.

```
Foundational packages    <- Migrate FIRST (no internal dependencies)
        ^
   Core packages         <- Migrate SECOND (depend on foundational)
        ^
Application packages     <- Migrate LAST (depend on core)
```

### Finding Your Dependency Order

For Go projects:

```bash
# List all packages
go list ./...

# See what a package imports
go list -f '{{.Imports}}' ./pkg/mypackage

# Find packages with no internal dependencies (good starting points)
go list -f '{{if not .Imports}}{{.ImportPath}}{{end}}' ./...
```

For Python projects:

```bash
# Analyze imports (requires pipdeptree or similar)
pipdeptree --local-only

# Or manually check imports in each module
grep -r "^from \. import\|^from \.\." src/
```

## Common Issues and Solutions

### Issue: Gazelle generates incorrect BUILD files

**Solution**: Use Gazelle directives in your root BUILD.bazel:

```starlark
# gazelle:prefix <your-module-path>
# gazelle:exclude vendor
# gazelle:exclude testdata
```

### Issue: Missing dependencies

**Solution**: Ensure all external dependencies are declared:

For Go, run:
```bash
go mod tidy
bazel run //:gazelle -- update-repos -from_file=go.mod
```

For Python, ensure requirements.in is complete:
```bash
pip-compile requirements.in
bazel run //:gazelle
```

### Issue: Build files conflict with existing structure

**Solution**: Use `.bazelignore` to exclude problematic directories temporarily:

```
# .bazelignore
problematic_dir/
legacy_code/
```

### Issue: Tests fail under Bazel

Common causes:
1. **File path assumptions**: Tests assume they run from project root
   - Fix: Use `runfiles` to locate test data
2. **Missing test data**: Data files not included in test target
   - Fix: Add `data` attribute to test rule
3. **Environment differences**: Different env vars in Bazel sandbox
   - Fix: Set `env` attribute or use `--test_env`

### Go-Specific Issues

**Issue**: `go:embed` directives not working

```starlark
go_library(
    name = "mylib",
    srcs = ["lib.go"],
    embedsrcs = ["templates/*"],  # Add embedded files
)
```

**Issue**: Build tags not respected

```starlark
# In BUILD.bazel or via gazelle directive
# gazelle:build_tags integration
```

### Python-Specific Issues

**Issue**: Imports not resolving

Ensure your `BUILD.bazel` has correct `imports`:

```starlark
py_library(
    name = "mylib",
    srcs = ["mylib.py"],
    imports = [".."],  # Adjust import path
)
```

## Verifying Migration

After migrating each package:

```bash
# 1. Build the package
bazel build //path/to/package/...

# 2. Run tests
bazel test //path/to/package/...

# 3. Check that the old build still works (during transition)
go build ./path/to/package/...
go test ./path/to/package/...
# or for Python:
python -m pytest path/to/package/
```

## Gradual Migration Checklist

For each package you migrate:

- [ ] Remove package from `.bazelignore`
- [ ] Run `bazel run //:gazelle`
- [ ] Review generated BUILD.bazel file
- [ ] Run `bazel build //package/...`
- [ ] Run `bazel test //package/...`
- [ ] Commit changes
- [ ] Update dependent packages if needed

## Rollback

If you need to undo the migration:

```bash
# See what changed
git diff

# Revert all Bazel files
git checkout -- MODULE.bazel .bazelrc BUILD.bazel .bazelignore

# Or reset everything
git checkout -- .
```

## Next Steps

Once fully migrated:

1. Update CI/CD to use Bazel
2. Remove old build configurations (Makefile, setup.py, etc.) if no longer needed
3. Consider enabling remote caching for faster builds
4. Explore Bazel's advanced features (aspects, transitions, etc.)

## Resources

- [Bazel Documentation](https://bazel.build/docs)
- [Gazelle Documentation](https://github.com/bazelbuild/bazel-gazelle)
- [rules_go Documentation](https://github.com/bazelbuild/rules_go)
- [rules_python Documentation](https://github.com/bazelbuild/rules_python)
