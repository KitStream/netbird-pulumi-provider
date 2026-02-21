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

# Post-process Java SDK to fix Gradle issues
echo "Post-processing Java SDK..."
# 1. Fix Java version to 17 (tfgen defaults to 11)
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/languageVersion = JavaLanguageVersion.of(11)/languageVersion = JavaLanguageVersion.of(17)/' sdk/java/build.gradle
  # 2. Fix non-existent lib project in settings.gradle
  sed -i '' 's/include("lib")/\/\/ include("lib")/' sdk/java/settings.gradle
else
  sed -i 's/languageVersion = JavaLanguageVersion.of(11)/languageVersion = JavaLanguageVersion.of(17)/' sdk/java/build.gradle
  # 2. Fix non-existent lib project in settings.gradle
  sed -i 's/include("lib")/\/\/ include("lib")/' sdk/java/settings.gradle
fi

echo "SDK generation complete."
