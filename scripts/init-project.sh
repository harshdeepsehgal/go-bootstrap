#!/usr/bin/env bash
set -euo pipefail

source_project_name="go-bootstrap"
project_name="${1:-}"
module_path="${2:-$project_name}"
script_path="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/$(basename "${BASH_SOURCE[0]}")"
project_root="$(cd "$(dirname "$script_path")/.." && pwd)"

usage() {
  echo "Usage: ./scripts/init-project.sh <project-name> [module-path]" >&2
}

fail() {
  echo "init-project: $*" >&2
  exit 1
}

replace_in_file() {
  local file="$1"
  local old_value="$2"
  local new_value="$3"
  local content

  content="$(<"$file")"
  content="${content//$old_value/$new_value}"
  printf '%s\n' "$content" >"$file"
}

remove_marked_block() {
  local file="$1"
  local start_marker="$2"
  local end_marker="$3"
  local line
  local output=""
  local skipping=0

  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ "$line" == "$start_marker" ]]; then
      skipping=1
      continue
    fi
    if [[ "$line" == "$end_marker" ]]; then
      skipping=0
      continue
    fi
    if ((skipping == 0)); then
      output+="$line"$'\n'
    fi
  done <"$file"

  printf '%s' "$output" >"$file"
}

if (($# < 1 || $# > 2)) || [[ -z "$project_name" ]]; then
  usage
  exit 2
fi
if [[ ! "$project_name" =~ ^[a-z][a-z0-9-]*$ ]] || [[ "$project_name" == *- ]]; then
  fail "project-name must use lowercase letters, numbers, and hyphens"
fi
if [[ "$project_name" == "$source_project_name" ]]; then
  fail "project-name must differ from $source_project_name"
fi
if [[ ! "$module_path" =~ ^[A-Za-z0-9][A-Za-z0-9._/-]*$ ]] ||
  [[ "$module_path" == */ ]] || [[ "$module_path" == *//* ]]; then
  fail "module-path contains unsupported characters"
fi

cd "$project_root"

[[ -f go.mod ]] || fail "go.mod was not found under $project_root"
current_module="$(awk '$1 == "module" { print $2; exit }' go.mod)"
[[ "$current_module" == "$source_project_name" ]] ||
  fail "this repository has already been initialized (module is $current_module)"
[[ -d "internal/$source_project_name" ]] ||
  fail "internal/$source_project_name was not found"
[[ ! -e "internal/$project_name" ]] ||
  fail "internal/$project_name already exists"

while IFS= read -r file; do
  replace_in_file \
    "$file" \
    "$source_project_name/internal/$source_project_name" \
    "$module_path/internal/$project_name"
  replace_in_file "$file" "$source_project_name/" "$module_path/"
  replace_in_file "$file" "\"$source_project_name:\"" "\"$project_name:\""
done < <(find cmd config internal -type f -name '*.go' -print)

replace_in_file go.mod "module $source_project_name" "module $module_path"
replace_in_file Makefile "$source_project_name" "$project_name"
replace_in_file scripts/package.sh "$source_project_name" "$project_name"
replace_in_file infra/local/compose.dependencies.yaml "$source_project_name" "$project_name"
replace_in_file README.md "Go Bootstrap" "$project_name"
replace_in_file README.md "$source_project_name" "$project_name"
remove_marked_block README.md "<!-- init-project:start -->" "<!-- init-project:end -->"

mv "internal/$source_project_name" "internal/$project_name"
find cmd config internal -type f -name '*.go' -exec gofmt -w {} +
rm -- "$script_path"

echo "Initialized $project_name"
echo "Go module: $module_path"
echo "Next: make check"
