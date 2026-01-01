# Version Bump Guide

This guide explains how to upgrade versions for Bazel, Go, Python, and other dependencies in your project.

## Version Categories

### User-Configurable Versions

These versions are stored in `.copier-answers.yml` and can be customized per-project:

| Component | File(s)                  | Copier Variable  |
|-----------|--------------------------|------------------|
| Bazel     | `.bazeliskrc`            | `bazel_version`  |
| Go SDK    | `MODULE.bazel`, `go.mod` | `go_version`     |
| Python    | `MODULE.bazel`           | `python_version` |

### Template-Managed Versions

These versions are maintained in the template and updated automatically when you update from the template:

| Component           | File           | Updated Via     |
|---------------------|----------------|-----------------|
| `rules_go`          | `MODULE.bazel` | Template update |
| `gazelle`           | `MODULE.bazel` | Template update |
| `rules_python`      | `MODULE.bazel` | Template update |
| `aspect_rules_lint` | `MODULE.bazel` | Template update |
| `buildifier`        | `MODULE.bazel` | Template update |

**To get the latest template-managed versions**, simply update from the template periodically:

```bash
# From the template repository
./update-project.sh /path/to/your/project
```

## Upgrading User-Configurable Versions

### Recommended: Edit `.copier-answers.yml`

Edit `.copier-answers.yml` to change your preferred versions:

```yaml
# .copier-answers.yml
bazel_version: "8.6.0"
go_version: "1.25.0"
python_version: "3.13"
```

Then commit and run the update:

```bash
# Copier requires a clean working directory
git add .copier-answers.yml
git commit -m "chore: update version settings"

# From the template repository, update project files
./update-project.sh /path/to/your/project
```

This approach ensures your version preferences persist through future template updates.

### Alternative: Direct Edit

For quick testing or one-off changes, edit the files directly:

**Bazel version** (`.bazeliskrc`):
```
USE_BAZEL_VERSION=8.6.0
```

**Go version** (`MODULE.bazel` and `go.mod`):
```python
# MODULE.bazel
go_sdk.download(version = "1.25.0")
```
```
# go.mod
go 1.25
```

**Python version** (`MODULE.bazel`):
```python
pip_ext.parse(
    hub_name = "py_deps",
    python_version = "3.13",
    ...
)
```

**Note:** Direct edits may be overwritten by future template updates unless you also update `.copier-answers.yml`.

After any version change, verify with:
```bash
bazel build //...
bazel test //...
```

## Checking for Updates

**Bazel:**
```bash
curl -s https://api.github.com/repos/bazelbuild/bazel/releases/latest | grep tag_name
```

**Go:**
```bash
curl -s https://go.dev/VERSION?m=text
```

**Bazel modules:**
```bash
# Or visit https://registry.bazel.build/ to see latest versions
```

## Common Issues

### Go Version Mismatch

If `go.mod` and `MODULE.bazel` have different Go versions:
```bash
# Ensure both files use the same version
# go.mod
go 1.24

# MODULE.bazel
go_sdk.download(version = "1.24.2")  # Patch version should match or be compatible
```

### Bazel Cache Issues After Upgrade

After upgrading Bazel or toolchains:
```bash
# Clean Bazel cache if you encounter issues
bazel clean --expunge
bazel build //...
```
