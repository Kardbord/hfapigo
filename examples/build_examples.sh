#!/bin/bash

set -e

if ! pushd "$(dirname "${BASH_SOURCE[0]}")"; then
  echo "Failed to enter script directory."
fi

for d in */; do
  pushd "${d}" > /dev/null && echo "Building ${d}" && go build && popd > /dev/null
done