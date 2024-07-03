#!/bin/bash

cd "$(dirname "$0")"

# Get the latest tag
git remote update
latest_tag=$(git describe --abbrev=0 --tags)
echo "Current version: $latest_tag"

# Extract major, minor, and patch versions
IFS='.' read -ra version_parts <<< "$latest_tag"
major="${version_parts[0]}"
minor="${version_parts[1]}"
patch="${version_parts[2]}"

# Increment the patch version
((patch++))

#Construct the new version
new_version="$major.$minor.$patch"

echo "Next version: $new_version"
