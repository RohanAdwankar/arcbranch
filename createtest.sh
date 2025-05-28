#!/bin/bash
set -e

# Go to the parent directory of arcbranch
cd "$(dirname "$PWD")"

# Create a temp directory for the test repo in the parent folder
TESTDIR="arcbranch-test-$(date +%s)"
mkdir "$TESTDIR"
cd "$TESTDIR"

# Initialize a new git repo
git init

# Copy examples/ contents from the arcbranch folder below
cp -r ../arcbranch/examples/* .

# Add and commit everything
git add .
git commit -m "Initial commit with example files"

# Run arcbranch 4
arcbranch 4

echo "Test repo created and arcbranch 4 run in $PWD"