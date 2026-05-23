#!/usr/bin/env bash

# NOTE: Configuration options may differ between development and CI.
# This is a dev utility script, and should not be run in CI.

set -e

pushd "$(dirname "${BASH_SOURCE[0]}")/.." >/dev/null

export GOCACHE=/tmp/hfgo-cache

TEST_CMD="go test"
if command -v gotestsum &>/dev/null; then
  TEST_CMD="gotestsum --format-icons hivis --format pkgname-and-test-fails --"
fi

RUN_INTEGRATION_TESTS=
if [[ "${1}" = "-i" ]]; then
  RUN_INTEGRATION_TESTS="true"
fi

echo "Formatting code..."
gofmt -s -w .
echo

echo "Tidying module files..."
go mod tidy
echo

echo "Vetting..."
go vet ./...
echo

echo "Linting..."
golangci-lint config verify
golangci-lint run --fix --disable godox ./...
echo

echo "Checking for vulnerabilities..."
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
echo

echo "Building..."
go build ./...
echo

echo "Running unit tests..."
${TEST_CMD} -shuffle=on ./...
echo

echo "Checking for race conditions..."
${TEST_CMD} -shuffle=on -race ./...
echo

if [[ "${RUN_INTEGRATION_TESTS}" = "true" ]]; then
  echo "Running integration tests..."
  ${TEST_CMD} -shuffle=on -tags=integration ./...
  echo
fi

echo "Reporting code coverage..."
${TEST_CMD} -shuffle=on -cover ./...
echo

popd >/dev/null

