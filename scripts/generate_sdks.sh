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

# Ensure version.txt exists for .NET SDK
echo "$VERSION" > sdk/dotnet/version.txt

# Post-process package.json for Node.js if needed (though tfgen usually handles some parts)
# sed -i "s/\${VERSION}/$VERSION/g" sdk/nodejs/package.json 2>/dev/null || true

echo "SDK generation complete."
