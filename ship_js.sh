#!/usr/bin/env bash
set -euo pipefail

# ship_js.sh: copy root README into dist_js, bump dist_js package version, and publish to npm

# 1. Copy root README.md to dist_js
cp README.md dist_js/README.md

echo "Copied README.md to dist_js/README.md"

# 2. Bump patch version in dist_js and publish
pushd dist_js > /dev/null
echo "Bumping version in dist_js/package.json..."
npm version patch --no-git-tag-version
echo "Publishing dist_js to npm..."
npm publish
popd > /dev/null
echo "Published dist_js package to npm"
