#!/usr/bin/env bash
set -e

pushd "$(dirname "${BASH_SOURCE[0]}")/.." >/dev/null

echo "Formatting code..."
gofmt -s -w .

echo "Tidying module files..."
go mod tidy

echo "Vetting..."
go vet ./...

echo "Building ..."
go build ./...

echo "Running tests..."
go test ./...

echo "Checking for race conditions..."
go test -race ./...

echo "Reporting code coverage..."
go test -cover ./...

