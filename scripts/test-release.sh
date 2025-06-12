#!/bin/bash

# Script to test GoReleaser configuration locally
# This will build binaries without creating a release

set -e

echo "🚀 Testing GoReleaser configuration..."

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "❌ GoReleaser is not installed. Please install it first:"
    echo "   brew install goreleaser/tap/goreleaser"
    echo "   or visit: https://goreleaser.com/install/"
    exit 1
fi

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf dist/

# Run goreleaser in snapshot mode (no git tag required)
echo "📦 Building snapshot release..."
goreleaser build --snapshot --clean

echo "✅ Build completed successfully!"
echo ""
echo "📁 Built binaries are in the dist/ directory:"
ls -la dist/

echo ""
echo "🎉 You can now test the binaries:"
echo "   ./dist/speakeasy-helpers_linux_amd64_v1/speakeasy-helpers --help"
echo "   ./dist/speakeasy-helpers_darwin_amd64_v1/speakeasy-helpers --help"
echo "   ./dist/speakeasy-helpers_darwin_arm64/speakeasy-helpers --help"
