#!/bin/bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Java: comment out include("lib") — the lib subproject doesn't exist
sed -i.bak 's/^include("lib")/\/\/ include("lib")/' "$ROOT_DIR/sdk/java/settings.gradle"
rm -f "$ROOT_DIR/sdk/java/settings.gradle.bak"

# Node.js: replace ${VERSION} placeholder with 0.0.0-dev so npm/yarn can resolve it,
# add main/types fields pointing to bin/ (tsc output), then regenerate package-lock.json
NODEJS_PKG="$ROOT_DIR/sdk/nodejs/package.json"
sed -i.bak 's/"\${VERSION}"/"0.0.0-dev"/' "$NODEJS_PKG"
rm -f "$NODEJS_PKG.bak"

# Add main and types fields if missing (generated package.json omits them)
if ! grep -q '"main"' "$NODEJS_PKG"; then
    sed -i.bak '/"license":/a\
    "main": "bin/index.js",\
    "types": "bin/index.d.ts",' "$NODEJS_PKG"
    rm -f "$NODEJS_PKG.bak"
fi

(cd "$ROOT_DIR/sdk/nodejs" && npm install --package-lock-only)
