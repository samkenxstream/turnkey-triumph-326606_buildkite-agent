#!/bin/bash

set -euo pipefail

echo '--- Downloading built agent'

buildkite-agent artifact download pkg/buildkite-agent-linux-386 pkg

# Grab the version of the binary while we're here (we need it if we deploy this
# commit to GitHub)
echo '+++ Extracting agent version from binary'

FULL_AGENT_VERSION=`pkg/buildkite-agent-linux-386 --version`
AGENT_VERSION=$(echo $FULL_AGENT_VERSION | sed 's/buildkite-agent version //' | sed -E 's/\, build .+//')
BUILD_VERSION=$(echo $FULL_AGENT_VERSION | sed 's/buildkite-agent version .*, build //')

echo "Full agent version: $FULL_AGENT_VERSION"
echo "Agent version: $AGENT_VERSION"
echo "Build version: $BUILD_VERSION"

buildkite-agent meta-data set "agent-version" "$AGENT_VERSION"
buildkite-agent meta-data set "agent-version-full" "$FULL_AGENT_VERSION"
buildkite-agent meta-data set "agent-version-build" "$BUILD_VERSION"