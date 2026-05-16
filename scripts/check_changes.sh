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

filter_go_mod() {
  sed -E '/^(module|go|toolchain) /d'
}

current_go_mod="$(filter_go_mod < go.mod)"
tagged_go_mod="$(git show "${latest_tag}:go.mod" 2>/dev/null | filter_go_mod || true)"

if [[ "$current_go_mod" != "$tagged_go_mod" ]]; then
  echo "go.mod dependencies changed since ${latest_tag}; a new release tag is needed."
  exit 0
fi

current_go_sum="$(cat go.sum 2>/dev/null || true)"
tagged_go_sum="$(git show "${latest_tag}:go.sum" 2>/dev/null || true)"

if [[ "$current_go_sum" != "$tagged_go_sum" ]]; then
  echo "go.sum changed since ${latest_tag}; a new release tag is needed."
  exit 0
fi

echo "No release-relevant changes found since ${latest_tag}."
exit 1
