# Releasing the NetBird Pulumi Provider

This project is configured for automated multi-language releases via GitHub Actions. A full release is triggered by pushing a version tag (e.g., `v0.0.1`).

## 1. Account Creation & API Tokens

Follow these steps to create the necessary accounts and generate API tokens for each package manager.

### NodeJS (npm)
npm supports **Trusted Publishing** via GitHub Actions, which is more secure than using static tokens.

1.  **Sign up**: [https://www.npmjs.com/signup](https://www.npmjs.com/signup)
2.  **Ensure Org Ownership**: You must be an **Owner/Admin** of the npm organization.
3.  **Configure Trusted Publishing**:
    *   **For Personal Accounts**: Log in -> **Settings** (from top right avatar) -> **Trusted Publishers**.
    *   **For Organizations**:
        1. Click your **profile picture** (top right) -> **"Organizations"**.
        2. Select the organization where you want to publish (e.g., `KitStream`).
        3. Look for the **"Settings"** tab (usually found after "Packages", "Members", "Teams").
        4. **Troubleshooting: If the "Settings" tab is missing**:
            *   Confirm you are an **Owner** of the organization. (npm sometimes only shows "Settings" to Owners, while Admins only see "Members", "Teams", "Billing").
            *   **Direct URL**: Try navigating directly to `https://www.npmjs.com/org/YOUR_ORG_NAME/settings/publishing`. (e.g., replace `YOUR_ORG_NAME` with `kitstream`).
            *   **Fallback Method**: If you cannot see organization-level settings, npm often hides them for organizations with zero packages. Follow the **Initial Publish** method below.
        5. In the left sidebar or sub-menu, click **"Publishing"** or **"Trusted Publishers"**.
        6. Click **"Add a new Trusted Publisher"**.
    *   **Initial Publish (Fallback Method)**:
        If you cannot find the organization-level settings, perform the first publish using a classic token:
        1. Create an **Automation Token** at [npmjs.com/settings/tokens/new](https://www.npmjs.com/settings/tokens/new).
        2. Add it to GitHub Secrets as `NPM_TOKEN`.
        3. Trigger the first release (see [Triggering a Release](#3-triggering-a-release)).
        4. Once the package exists (e.g., `@kitstream/netbird`), go to the package page on npm -> **Settings** -> **Publishing** to configure Trusted Publishing.
4.  **Verify Package Scope**: Ensure the `name` in `sdk/nodejs/package.json` matches your organization (e.g., `@kitstream/netbird` if your org is `kitstream`).
5.  **Fill in details**:
    *   **GitHub Organization/User**: `KitStream`.
    *   **GitHub Repository**: `netbird-pulumi-provider`.
    *   **Workflow Name**: `release.yml`.
    *   **Environment**: (Optional, leave blank if not using GitHub Environments).
6.  **Done**: Once configured, npm will trust the GitHub Actions workflow to publish without a static `NPM_TOKEN`.

### Python (PyPI)
1.  **Sign up**: [https://pypi.org/account/register/](https://pypi.org/account/register/)
2.  **Organization Setup**:
    *   Ensure you have an organization named `KitStream`: [https://pypi.org/org/kitstream/](https://pypi.org/org/kitstream/)
    *   You must be an **Owner/Admin** to manage the organization's projects.
3.  **Enable 2FA**: Go to [https://pypi.org/manage/account/#two-factor-authentication](https://pypi.org/manage/account/#two-factor-authentication). PyPI requires 2FA to generate API tokens.
4.  **Generate Token**: [https://pypi.org/manage/account/token/](https://pypi.org/manage/account/token/)
    *   **Name**: `netbird-pulumi-release`.
    *   **Scope**: Initially set to **"Entire account"**. (Once the project is published once, you can change it to just this project).
    *   **Copy the token** (starts with `pypi-`) for the `PYPI_TOKEN` secret.
5.  **Assign Project to Organization**: After the initial publish, go to the project's **Settings** -> **Ownership** and transfer it to the `KitStream` organization.

### .NET (NuGet)
1.  **Login**: [https://www.nuget.org/users/account/LogOn](https://www.nuget.org/users/account/LogOn) (via Microsoft account).
2.  **Organization Setup**:
    *   Ensure the `KitStream` organization exists on NuGet: [https://www.nuget.org/organizations/KitStream](https://www.nuget.org/organizations/KitStream)
    *   Ensure you are a member of the organization with **Owner** or **Admin** privileges.
3.  **Generate API Key**: [https://www.nuget.org/account/apikeys/create](https://www.nuget.org/account/apikeys/create)
    *   **Key Name**: `netbird-pulumi-release`.
    *   **Key Owner**: Select the **KitStream** organization (not your personal account).
    *   **Scopes**: Ensure **"Push"** is selected.
    *   **Glob Pattern**: Use `*` (or `Pulumi.Netbird`).
    *   **Copy the key** for the `NUGET_PUBLISH_KEY` secret.

### Java (Maven Central / OSSRH)
Maven Central publishing now uses the new Central Portal.
1.  **Sign up**: [https://central.sonatype.com/signup](https://central.sonatype.com/signup)
2.  **Verify Namespace (Organization Identifier)**: [https://central.sonatype.com/publishing/namespaces](https://central.sonatype.com/publishing/namespaces)
    *   Add namespace: `io.github.kitstream`. This represents the **KitStream** organization on Maven Central.
    *   Follow instructions to verify ownership via a temporary repository.
3.  **Generate User Token**: [https://central.sonatype.com/account](https://central.sonatype.com/account)
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
