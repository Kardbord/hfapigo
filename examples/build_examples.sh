#!/bin/bash

EXPECTED_DIR="hfapigo/examples"

cwd=$(basename "$(dirname "${PWD}")")/$(basename "${PWD}")

if [ "${cwd}" != "${EXPECTED_DIR}" ]; then
  >&2 echo "This script must be run from ${EXPECTED_DIR}"
  exit 1
fi

for d in */; do
  pushd "${d}" > /dev/null && go build
  echo "Building ${d}"
  popd > /dev/null || exit
done