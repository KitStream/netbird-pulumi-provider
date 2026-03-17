#!/bin/bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Java: comment out include("lib") — the lib subproject doesn't exist
sed -i.bak 's/^include("lib")/\/\/ include("lib")/' "$ROOT_DIR/sdk/java/settings.gradle"
rm -f "$ROOT_DIR/sdk/java/settings.gradle.bak"

# Node.js: regenerate package-lock.json
(cd "$ROOT_DIR/sdk/nodejs" && npm install --package-lock-only)
