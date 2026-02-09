# Releasing Serenity-Go

This document explains the automated release process for Serenity-Go project.

## 🚀 Overview

Serenity-Go uses semantic versioning with automated releases based on conventional commits. The release process is fully automated using GitHub Actions, with local tools for testing and validation.

## 📋 Versioning Rules

The project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html) with automatic version bumps based on commit types:

| Commit Type | Description | Version Increment |
|-------------|-------------|-------------------|
| `feat` | New feature | **MINOR** (0.1.0 → 0.2.0) |
| `fix` | Bug fix | **PATCH** (0.1.0 → 0.1.1) |
| `BREAKING CHANGE` | Breaking changes | **MAJOR** (0.1.0 → 1.0.0) |
| `docs`, `style`, `refactor`, `test`, `chore` | Maintenance | **No change** |

## 🔄 Automated Process

### GitHub Actions Workflow

On every push to `main`/`master` branch:
1. **Runs tests and linting** - ensures code quality
2. **Analyzes commits** - determines next version
3. **Generates changelog** - updates CHANGELOG.md
4. **Creates GitHub Release** - with automatic notes
5. **Creates git tag** - marks the release version

### Trigger Conditions
- Push to `main` or `master` branch
- Commits contain conventional commit types
- All tests and linting pass
- Commit message doesn't contain `[skip ci]`

## 🛠️ Local Development Tools

### Release Script

Use the provided script for local testing and validation:

```bash
# Preview next release without changes
./scripts/release.sh preview

# Prepare release (run tests, generate changelog)
./scripts/release.sh prepare

# Create and push release
./scripts/release.sh release

# Show help
./scripts/release.sh help
```

### Makefile Commands

Convenient commands for release management:

```bash
# Preview changelog (dry run)
make release-dry

# Prepare release (tests, lint, changelog)
make release-prepare

# Create release (commit changelog, trigger workflow)
make release
```

## 📝 Example Workflow

### 1. Making Changes

```bash
# Create a feature
git checkout -b feature/new-api
# ... make changes ...
git add .
git commit -m "feat: add new API endpoint for data processing"
git push origin feature/new-api
# Create PR and merge to main
```

### 2. Creating Release

After merging to main:

```bash
# Check what will be released
./scripts/release.sh preview

# Output:
# 🎯 Next version will be: v0.2.0
# 📋 Changes:
# ### 🚀 Features
# - Add new API endpoint for data processing
```

### 3. Automated Release

Push to main branch and GitHub Actions will:
1. Detect the `feat` commit
2. Bump version to v0.2.0
3. Update CHANGELOG.md
4. Create GitHub Release with notes
5. Create and push git tag

## 🏷️ Release Examples

### Initial Release (v0.1.0)
```bash
# First release with all current changes
git tag v0.1.0
git push origin v0.1.0
```

### Feature Release (v0.2.0)
```bash
# After feat commits
git log --oneline v0.1.0..HEAD
# feat: add new API endpoint
# feat: improve error handling

# Result: v0.2.0
```

### Patch Release (v0.2.1)
```bash
# After fix commits
git log --oneline v0.2.0..HEAD
# fix: resolve nil pointer in API client

# Result: v0.2.1
```

## 🔧 Configuration Files

### `cliff.toml`
Configuration for [git-cliff](https://github.com/orhun/git-cliff) changelog generator:
- Groups commits by type (Features, Bug Fixes, Documentation, etc.)
- Uses semantic versioning logic
- Generates beautiful changelog format

### `.semrelrc`
Configuration for semantic-release:
- Defines maintainer information
- Configures release plugins
- Sets up repository integration

### `.github/workflows/release.yml`
GitHub Actions workflow:
- Runs on push to main/master
- Installs required tools
- Performs automated version bumping
- Creates GitHub releases

## 🚨 Troubleshooting

### Release Not Triggered
- **Check**: Commit messages follow conventional format
- **Check**: No `[skip ci]` in commit message
- **Check**: All tests and linting pass
- **Check**: Workflow has proper permissions

### Version Not Bumped
- **Check**: Commits contain versionable changes (`feat`, `fix`)
- **Check**: No existing tag for same version
- **Check**: GitHub Actions workflow runs successfully

### Changelog Not Updated
- **Check**: `cliff.toml` configuration
- **Check**: Git history has conventional commits
- **Check**: Workflow has permissions to update repo

### Manual Override
If automation fails, create release manually:

```bash
# Determine version
./scripts/release.sh preview

# Create tag
git tag v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# Create GitHub release manually
```

## 📊 Current Status

- ✅ Automated version bumping
- ✅ Conventional commits support
- ✅ Automated changelog generation
- ✅ GitHub Actions integration
- ✅ Semantic versioning
- ✅ Local testing tools
- ✅ First release (v0.1.0) created
- ✅ Automated workflow tested

## 🎯 Best Practices

1. **Write good commit messages**
   ```bash
   feat: add new API endpoint
   fix: resolve authentication issue
   docs: update installation guide
   ```

2. **Test locally before pushing**
   ```bash
   ./scripts/release.sh preview
   make release-dry
   ```

3. **Monitor GitHub Actions**
   - Check workflow runs after pushing
   - Verify releases are created correctly
   - Watch for any failures

4. **Keep changelog clean**
   - Use conventional commit types
   - Write clear, descriptive messages
   - Avoid merging unrelated changes

## 🔄 Migration from Manual

Previously: Manual version bumping, manual changelog updates
Now: Fully automated based on conventional commits

Benefits:
- ✅ No manual version management
- ✅ Consistent changelog format
- ✅ Automatic GitHub releases
- ✅ Semantic versioning
- ✅ Local testing capabilities