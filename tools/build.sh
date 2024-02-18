#!/usr/bin/env bash
set -e

pushd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null

./build-examples.sh

pushd .. >/dev/null
echo "Formatting code..."
go fmt ./...
echo "Building $(basename "$(pwd)")"
go build ./...
echo "Vetting..."
go vet ./...
echo "Running tests..."
go test ./...
echo "Done."