# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Automated versioning pipeline with git-cliff and semantic-release
- GitHub Actions workflow for automatic releases
- Comprehensive release documentation and scripts

### Changed
- Enhanced development guidelines with release process

## [0.1.0] - 2026-02-09

### Added
- Initial release of Serenity-Go framework with Screenplay Pattern
- Core Screenplay Pattern interfaces (Actor, Activity, Question)
- HTTP API testing capabilities with fluent request building
- Expectation system with ensure-style assertions
- TestContext API with automatic error handling
- gomock integration for interface mocking
- Comprehensive CI/CD pipeline with multi-version testing
- Makefile with development commands (test, lint, build)
- GitHub Actions workflow with code coverage integration
- Apache 2.0 license
- Comprehensive documentation for AI agents and developers

### Changed
- Split serenity.go into focused files by responsibility
- Reorganized package structure for better maintainability
- Migrated golangci-lint configuration from v1 to v2

### Documentation
- Added AGENTS.md with comprehensive development guidelines
- Created examples demonstrating framework capabilities
- Added README with badges and getting started guide