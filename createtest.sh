#!/bin/bash
set -e

# Create a temp directory for the test repo
TESTDIR="arcbranch-test-$(date +%s)"
mkdir "$TESTDIR"
cd "$TESTDIR"

# Initialize a new git repo
git init

# Copy examples/ contents from the parent directory
cp -r ../examples/* .

# Add and commit everything
git add .
git commit -m "Initial commit with example files"

echo "Test repo created in $PWD"
echo "You can now run: arcbranch 4"