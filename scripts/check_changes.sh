#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

latest_tag="$(git tag --list 'v*' --sort=-version:refname | head -n1)"

if [[ -z "$latest_tag" ]]; then
  echo "No release tag found; a new release tag is needed."
  exit 0
fi

if git diff --quiet "$latest_tag"..HEAD -- \
  '*.go' \
  'go.mod' \
  'go.sum' \
  '.golangci.version' \
  '.govulncheck.version' \
  'README.md' \
  'CLAUDE.md' \
  'SPECS/*.md' \
  'Taskfile.yml' \
  'scripts/*.sh'; then
  :
else
  echo "Release-relevant files changed since ${latest_tag}; a new release tag is needed."
  exit 0
fi

echo "No release-relevant changes found since ${latest_tag}."
exit 1
