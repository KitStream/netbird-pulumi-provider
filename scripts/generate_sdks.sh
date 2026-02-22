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
  # 3. Fix Gradle group and publication
  sed -i '' 's/group = "com.netbird"/group = "io.github.kitstream"/' sdk/java/build.gradle
  sed -i '' 's/groupId = "com.netbird"/groupId = project.group/' sdk/java/build.gradle
  sed -i '' 's/publishRepoUsername = System.getenv("PUBLISH_REPO_USERNAME")/publishRepoUsername = System.getenv("MAVEN_USERNAME")/' sdk/java/build.gradle
  sed -i '' 's/publishRepoPassword = System.getenv("PUBLISH_REPO_PASSWORD")/publishRepoPassword = System.getenv("MAVEN_PASSWORD")/' sdk/java/build.gradle
else
  sed -i 's/languageVersion = JavaLanguageVersion.of(11)/languageVersion = JavaLanguageVersion.of(17)/' sdk/java/build.gradle
  # 2. Fix non-existent lib project in settings.gradle
  sed -i 's/include("lib")/\/\/ include("lib")/' sdk/java/settings.gradle
  # 3. Fix Gradle group and publication
  sed -i 's/group = "com.netbird"/group = "io.github.kitstream"/' sdk/java/build.gradle
  sed -i 's/groupId = "com.netbird"/groupId = project.group/' sdk/java/build.gradle
  sed -i 's/publishRepoUsername = System.getenv("PUBLISH_REPO_USERNAME")/publishRepoUsername = System.getenv("MAVEN_USERNAME")/' sdk/java/build.gradle
  sed -i 's/publishRepoPassword = System.getenv("PUBLISH_REPO_PASSWORD")/publishRepoPassword = System.getenv("MAVEN_PASSWORD")/' sdk/java/build.gradle
fi

echo "SDK generation complete."
