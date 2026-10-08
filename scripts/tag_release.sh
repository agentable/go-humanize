#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <VERSION>" >&2
  exit 1
fi

version="${1#v}"

if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "invalid version: ${1}" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

tag="v${version}"

if [[ "$(git symbolic-ref --quiet --short HEAD || true)" != "main" ]]; then
  echo "release must run from main" >&2
  exit 1
fi

release_commit="$(git rev-parse HEAD)"

if git rev-parse -q --verify "refs/tags/${tag}" >/dev/null; then
  echo "tag already exists: ${tag}" >&2
  exit 1
fi

git tag -a "${tag}" "$release_commit" -m "${tag}"
if ! git push --atomic --no-follow-tags origin "$release_commit:refs/heads/main" "refs/tags/$tag:refs/tags/$tag"; then
  echo "release not published; local tag $tag remains at $release_commit" >&2
  echo "resolve the remote rejection, then retry the atomic push of main and $tag" >&2
  exit 1
fi
