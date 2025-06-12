# speakeasy-helpers

A collection of command-line helper utilities for working with Speakeasy generated API clients.

[![CI](https://github.com/rwjblue-glean/speakeasy-helpers/actions/workflows/ci.yml/badge.svg)](https://github.com/rwjblue-glean/speakeasy-helpers/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rwjblue-glean/speakeasy-helpers)](https://goreportcard.com/report/github.com/rwjblue-glean/speakeasy-helpers)

## Overview

`speakeasy-helpers` is a CLI tool that provides a set of utilities to help you work more effectively with Speakeasy generated API clients.

## Installation

### mise

Installation via [mise](https://mise.jdx.dev/) is very straightforward. If you have mise installed, you can run:

```bash
mise use ubi:rwjblue-glean/speakeasy-helpers
```

### Download Pre-built Binaries

Download the latest release for your platform from the [GitHub Releases page](https://github.com/rwjblue-glean/speakeasy-helpers/releases).

### Build from Source

```bash
git clone https://github.com/rwjblue-glean/speakeasy-helpers.git
cd speakeasy-helpers
go build -o dist/speakeasy-helpers
```

## Usage

```bash
# Show help
speakeasy-helpers --help

# Show version
speakeasy-helpers --version
```

## Development

### Prerequisites

- Go 1.24 or later
- [mise](https://mise.jdx.dev/) (optional, for tool management)

### Building

```bash
# Using mise
mise run build

# Or directly with go
go build -o bin/speakeasy-helpers
```

### Testing

```bash
# Using mise
mise run test

# Or directly with go
go test ./...
```

### Releasing

This project uses [GoReleaser](https://goreleaser.com/) for automated releases.

To create a new release:

1. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. The GitHub Actions workflow will automatically build and publish the release.

To test the release configuration locally:

```bash
# Test the build
./scripts/test-release.sh
```

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
