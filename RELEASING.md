# Releasing the NetBird Pulumi Provider

This project is configured for automated multi-language releases via GitHub Actions. A full release is triggered by pushing a version tag (e.g., `v0.0.1`).

## 1. Account Creation & API Tokens

Follow these steps to create the necessary accounts and generate API tokens for each package manager.

### NodeJS (npm)
npm supports **Trusted Publishing** via GitHub Actions, which is more secure than using static tokens.

1.  **Sign up**: [https://www.npmjs.com/signup](https://www.npmjs.com/signup)
2.  **Configure Trusted Publishing**:
    *   **For Personal Accounts**: Log in -> **Settings** -> **Trusted Publishers**.
    *   **For Organizations**: Log in -> Select your **Organization** -> **Settings** -> **Publishing** -> **"Add a new Trusted Publisher"**.
    *   **If the package is already published**: Navigate to the package -> **Settings** -> **Publishing** -> **"Add a new Trusted Publisher"**.
3.  **Fill in details**:
    *   **GitHub Organization/User**: `KitStream`.
    *   **GitHub Repository**: `netbird-pulumi-provider`.
    *   **Workflow Name**: `release.yml`.
    *   **Environment**: (Optional, leave blank if not using GitHub Environments).
4.  Once configured, npm will trust the GitHub Actions workflow to publish without a static `NPM_TOKEN`.

### Python (PyPI)
*   **Sign up**: [https://pypi.org/account/register/](https://pypi.org/account/register/)
*   **Enable 2FA**: Go to [https://pypi.org/manage/account/#two-factor-authentication](https://pypi.org/manage/account/#two-factor-authentication). PyPI requires 2FA to generate API tokens.
*   **Generate Token**: [https://pypi.org/manage/account/token/](https://pypi.org/manage/account/token/)
    *   **Name**: `netbird-pulumi-release`.
    *   **Scope**: Initially set to **"Entire account"**. (Once the project is published once, you can change it to just this project).
    *   **Copy the token** (starts with `pypi-`) for the `PYPI_TOKEN` secret.

### .NET (NuGet)
*   **Login**: [https://www.nuget.org/users/account/LogOn](https://www.nuget.org/users/account/LogOn) (via Microsoft account).
*   **Generate API Key**: [https://www.nuget.org/account/apikeys/create](https://www.nuget.org/account/apikeys/create)
    *   **Key Name**: `netbird-pulumi-release`.
    *   **Scopes**: Ensure **"Push"** is selected.
    *   **Glob Pattern**: Use `*` (or `Pulumi.Netbird`).
    *   **Copy the key** for the `NUGET_PUBLISH_KEY` secret.

### Java (Maven Central / OSSRH)
Maven Central publishing now uses the new Central Portal.
*   **Sign up**: [https://central.sonatype.com/signup](https://central.sonatype.com/signup)
*   **Verify Namespace**: [https://central.sonatype.com/publishing/namespaces](https://central.sonatype.com/publishing/namespaces)
    *   Add namespace: `io.github.kitstream`.
    *   Follow instructions to verify ownership via a temporary repository.
*   **Generate User Token**: [https://central.sonatype.com/account](https://central.sonatype.com/account)
    *   Click **"Generate User Token"**.
    *   The generated **Username** and **Password** are used for `OSSRH_USERNAME` and `OSSRH_TOKEN` respectively.

#### GPG Key for Java Signing
Maven Central requires artifacts to be signed.
1.  **Install GPG**: `brew install gpg` (on Mac).
2.  **Generate Key**: `gpg --full-generate-key` (Use RSA, 4096 bits, no expiration).
3.  **Find Key ID**: `gpg --list-secret-keys --keyid-format=long`. (e.g., `AB12CD34EF56GH78`).
4.  **Publish Public Key**: `gpg --keyserver keyserver.ubuntu.com --send-keys YOUR_KEY_ID`.
5.  **Export Private Key**: `gpg --armor --export-secret-keys YOUR_KEY_ID` (Copy the entire block starting with `-----BEGIN PGP PRIVATE KEY BLOCK-----` for `JAVA_SIGNING_KEY`).
6.  **Passphrase**: The passphrase you used during generation is `JAVA_SIGNING_PASSWORD`.

---

## 2. GitHub Secrets Configuration

Add all the collected tokens and keys as secrets in your GitHub repository settings.

1.  **Navigate to Secrets**: [https://github.com/KitStream/netbird-pulumi-provider/settings/secrets/actions](https://github.com/KitStream/netbird-pulumi-provider/settings/secrets/actions)
2.  **Add New Secrets**: Click **"New repository secret"** for each of the following:

| Name | Description |
| :--- | :--- |
| `NPM_TOKEN` | (Optional) npm Automation Token (Only if not using Trusted Publishing) |
| `PYPI_TOKEN` | PyPI API Token (starts with `pypi-`) |
| `NUGET_PUBLISH_KEY` | NuGet API Key |
| `OSSRH_USERNAME` | Sonatype User Token Username |
| `OSSRH_TOKEN` | Sonatype User Token Password |
| `JAVA_SIGNING_KEY` | Entire GPG Private Key block |
| `JAVA_SIGNING_PASSWORD` | GPG Private Key Passphrase |

---

## 3. Triggering a Release

Once all secrets are added:
1.  Ensure all changes are pushed to `main`.
2.  Tag the release:
    ```bash
    git tag v0.0.1
    git push origin v0.0.1
    ```
3.  Monitor the progress in the **Actions** tab: [https://github.com/KitStream/netbird-pulumi-provider/actions](https://github.com/KitStream/netbird-pulumi-provider/actions)
