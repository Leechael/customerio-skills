#!/usr/bin/env bash
set -euo pipefail

force="false"
if [[ "${1:-}" == "--force" ]]; then
  force="true"
  shift
fi

target="${1:-}"
if [[ -z "$target" ]]; then
  echo "Usage: scripts/init-prek.sh [--force] <target-repo-dir>" >&2
  exit 1
fi

if [[ ! -d "$target" ]]; then
  echo "target repo dir does not exist: $target" >&2
  exit 1
fi

src="$(dirname "$0")/../prek.toml"
dest="$target/prek.toml"
if [[ -f "$dest" && "$force" != "true" ]]; then
  echo "prek.toml already exists: $dest" >&2
  echo "Use --force to overwrite." >&2
  exit 1
fi

cp "$src" "$dest"
echo "Copied $dest"
