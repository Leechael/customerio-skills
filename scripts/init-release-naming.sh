#!/usr/bin/env bash
set -euo pipefail

force="false"
if [[ "${1:-}" == "--force" ]]; then
  force="true"
  shift
fi

target="${1:-}"
if [[ -z "$target" ]]; then
  echo "Usage: scripts/init-release-naming.sh [--force] <target-repo-dir>" >&2
  exit 1
fi

if [[ ! -d "$target" ]]; then
  echo "target repo dir does not exist: $target" >&2
  exit 1
fi

dest="$target/release-naming.env"
if [[ -f "$dest" && "$force" != "true" ]]; then
  echo "release-naming.env already exists: $dest" >&2
  echo "Use --force to overwrite." >&2
  exit 1
fi

cp "$(dirname "$0")/../release-naming.env" "$dest"
echo "Copied $dest"
