#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    VERSION="0.0.1"
fi

echo "Generating SDKs with version: $VERSION"

./bin/pulumi-tfgen-netbird schema --out provider/cmd/pulumi-resource-netbird
./bin/pulumi-tfgen-netbird go --out sdk/go
./bin/pulumi-tfgen-netbird nodejs --out sdk/nodejs
./bin/pulumi-tfgen-netbird python --out sdk/python
./bin/pulumi-tfgen-netbird dotnet --out sdk/dotnet
./bin/pulumi-tfgen-netbird java --out sdk/java

echo "SDK generation complete."

# Apply post-generation fixes
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [ -x "$SCRIPT_DIR/fix_sdks.sh" ]; then
    echo "Applying post-generation SDK fixes..."
    "$SCRIPT_DIR/fix_sdks.sh"
fi
