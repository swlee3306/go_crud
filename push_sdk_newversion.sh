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

# Construct the new version

echo "Next version: $new_version"


# Commit and push
ARGMSG="$(date) $(date +%T)"
NEW_VERSION="$major.$minor.$patch"

# Check if at least one argument is provided
if [ "$#" -ge 1 ]; then
    ARGMSG="$1"
fi

# test
echo "initiating..."
echo "ARGMSG: $ARGMSG"
echo "New version: $NEW_VERSION"

echo "pushing..."

git remote update
git add .
git commit -a -m "$ARGMSG"
git tag "$NEW_VERSION"
git push origin
git push --tags

echo "done."
sleep 3