#!/usr/bin/env bash

# NOTE: Configuration options may differ between development and CI.
# This is a dev utility script, and should not be run in CI.

set -e

pushd "$(dirname "${BASH_SOURCE[0]}")/.." >/dev/null

export GOCACHE=/tmp/hfgo-cache

echo "Formatting code..."
gofmt -s -w .

echo "Tidying module files..."
go mod tidy

echo "Vetting..."
go vet ./...

echo "Linting..."
golangci-lint config verify
golangci-lint run --fix --disable godox ./...

echo "Building..."
go build ./...

echo "Running tests..."
go test ./...

echo "Checking for race conditions..."
go test -race ./...

echo "Reporting code coverage..."
go test -cover ./...

