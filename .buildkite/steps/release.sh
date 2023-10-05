#!/usr/bin/env bash

set -euo pipefail

if [ -z "$BUILDKITE_TAG" ]; then
  echo "Releases only run when a new tag is pushed to github.com"
  exit 0
else
  echo "Preparing to release: ${BUILDKITE_TAG}"
fi

echo "--- Checking GitHub token"

if [ -z "$GITHUB_TOKEN" ]; then
  echo "\$GITHUB_TOKEN env variable must contain a Github API token with permission to create releases in buildkite/ecr-scan-results-buildkite-plugin"
  exit 1
fi

echo "--- Running goreleaser"

goreleaser release --clean