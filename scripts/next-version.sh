#!/usr/bin/env bash
# Prints the next semantic version (without the leading "v") based on the
# Conventional Commit messages since the latest v* tag:
#   BREAKING CHANGE / "type!:"  -> major
#   feat:                       -> minor
#   anything else               -> patch
# With no previous tag the first release is 1.0.0.
set -euo pipefail

latest=$(git describe --tags --abbrev=0 --match 'v[0-9]*' 2>/dev/null || true)
if [ -z "$latest" ]; then
  echo "1.0.0"
  exit 0
fi

IFS=. read -r major minor patch <<<"${latest#v}"
log=$(git log "$latest"..HEAD --format=%B)

if grep -qE '^[a-z]+(\([^)]*\))?!:|BREAKING[ -]CHANGE' <<<"$log"; then
  echo "$((major + 1)).0.0"
elif grep -qE '^feat(\([^)]*\))?:' <<<"$log"; then
  echo "${major}.$((minor + 1)).0"
else
  echo "${major}.${minor}.$((patch + 1))"
fi
