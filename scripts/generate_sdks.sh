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

# Update npm package name and version
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/"name": "@pulumi\/netbird"/"name": "@kitstream\/netbird-pulumi"/' sdk/nodejs/package.json
  sed -i '' "s/\"\${VERSION}\"/\"$VERSION\"/" sdk/nodejs/package.json
  # Ensure main and types are present (sometimes tfgen misses them or we want to be sure)
  if ! grep -q "\"main\":" sdk/nodejs/package.json; then
    sed -i '' '/"@kitstream\/netbird-pulumi"/a \
    "main": "bin/index.js",\
    "types": "bin/index.d.ts",' sdk/nodejs/package.json
  fi
else
  sed -i 's/"name": "@pulumi\/netbird"/"name": "@kitstream\/netbird-pulumi"/' sdk/nodejs/package.json
  sed -i "s/\"\${VERSION}\"/\"$VERSION\"/" sdk/nodejs/package.json
  if ! grep -q "\"main\":" sdk/nodejs/package.json; then
    sed -i '/"@kitstream\/netbird-pulumi"/a \    "main": "bin/index.js",\n    "types": "bin/index.d.ts",' sdk/nodejs/package.json
  fi
fi

# Ensure version.txt exists for .NET SDK
echo "$VERSION" > sdk/dotnet/version.txt

# Post-process Java SDK to fix Gradle issues
echo "Post-processing Java SDK..."
# 1. Fix Java version to 17 (tfgen defaults to 11)
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/languageVersion = JavaLanguageVersion.of(11)/languageVersion = JavaLanguageVersion.of(17)/' sdk/java/build.gradle
  # 2. Fix non-existent lib project in settings.gradle
  sed -i '' 's/include("lib")/\/\/ include("lib")/' sdk/java/settings.gradle
  # 3. Fix Gradle group and publication
  sed -i '' 's/group = "com.netbird"/group = "io.github.kitstream"/' sdk/java/build.gradle
  sed -i '' 's/groupId = "com.netbird"/groupId = project.group/' sdk/java/build.gradle
  sed -i '' 's/publishRepoUsername = System.getenv("PUBLISH_REPO_USERNAME")/publishRepoUsername = System.getenv("MAVEN_USERNAME")/' sdk/java/build.gradle
  sed -i '' 's/publishRepoPassword = System.getenv("PUBLISH_REPO_PASSWORD")/publishRepoPassword = System.getenv("MAVEN_PASSWORD")/' sdk/java/build.gradle
  # 4. Fix package names in Java source files
  find sdk/java/src/main/java -name "*.java" -exec sed -i '' 's/com\.netbird/io\.github\.kitstream/g' {} +
  # 5. Move files to correct package directory
  if [ -d "sdk/java/src/main/java/com/netbird/netbird" ]; then
    mkdir -p sdk/java/src/main/java/io/github/kitstream/netbird
    mv sdk/java/src/main/java/com/netbird/netbird/* sdk/java/src/main/java/io/github/kitstream/netbird/
    rm -rf sdk/java/src/main/java/com/netbird
  fi
else
  sed -i 's/languageVersion = JavaLanguageVersion.of(11)/languageVersion = JavaLanguageVersion.of(17)/' sdk/java/build.gradle
  # 2. Fix non-existent lib project in settings.gradle
  sed -i 's/include("lib")/\/\/ include("lib")/' sdk/java/settings.gradle
  # 3. Fix Gradle group and publication
  sed -i 's/group = "com.netbird"/group = "io.github.kitstream"/' sdk/java/build.gradle
  sed -i 's/groupId = "com.netbird"/groupId = project.group/' sdk/java/build.gradle
  sed -i 's/publishRepoUsername = System.getenv("PUBLISH_REPO_USERNAME")/publishRepoUsername = System.getenv("MAVEN_USERNAME")/' sdk/java/build.gradle
  sed -i 's/publishRepoPassword = System.getenv("PUBLISH_REPO_PASSWORD")/publishRepoPassword = System.getenv("MAVEN_PASSWORD")/' sdk/java/build.gradle
  # 4. Fix package names in Java source files
  find sdk/java/src/main/java -name "*.java" -exec sed -i 's/com\.netbird/io\.github\.kitstream/g' {} +
  # 5. Move files to correct package directory
  if [ -d "sdk/java/src/main/java/com/netbird/netbird" ]; then
    mkdir -p sdk/java/src/main/java/io/github/kitstream/netbird
    mv sdk/java/src/main/java/com/netbird/netbird/* sdk/java/src/main/java/io/github/kitstream/netbird/
    rm -rf sdk/java/src/main/java/com/netbird
  fi
fi

echo "SDK generation complete."
