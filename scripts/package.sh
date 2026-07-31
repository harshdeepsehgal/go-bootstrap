#!/usr/bin/env bash
set -euo pipefail

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
archive_name="go-bootstrap.zip"
dist_dir="$project_root/dist"
temporary_archive="$(mktemp -t go-bootstrap.XXXXXX).zip"

cleanup() {
  rm -f "$temporary_archive"
}
trap cleanup EXIT

cd "$project_root"
mkdir -p "$dist_dir"

zip -q -r "$temporary_archive" . \
  -x '.git/*' \
  -x '.vscode/*' \
  -x '.DS_Store' \
  -x '*/.DS_Store' \
  -x '.env' \
  -x '*/.env' \
  -x 'bin/*' \
  -x 'coverage.out' \
  -x 'coverage.html' \
  -x 'dist/*' \
  -x '*.test' \
  -x '*.prof'

mv -f "$temporary_archive" "$dist_dir/$archive_name"
echo "Created $dist_dir/$archive_name"
