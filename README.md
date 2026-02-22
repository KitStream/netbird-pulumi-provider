# NetBird Pulumi Provider

A Pulumi provider for managing resources in [NetBird](https://netbird.io).

This provider is built on top of the [NetBird Terraform Provider](https://github.com/netbirdio/terraform-provider-netbird) using the [Pulumi Terraform Bridge](https://github.com/pulumi/pulumi-terraform-bridge).

## Supported Languages

The following SDKs are available and automatically updated:

- **NodeJS** (`@kitstream/netbird-pulumi` on npm)
- **Python** (`pulumi_netbird` on PyPI)
- **Go** (`github.com/KitStream/netbird-pulumi-provider/sdk/go/index`)
- **.NET** (`KitStream.Pulumi.Netbird` on NuGet)
- **Java** (`io.github.kitstream:netbird` on Maven Central)

## Installation

To use this provider, install the appropriate SDK for your language. For example:

### NodeJS
```bash
npm install @kitstream/netbird-pulumi
```

### Python
```bash
pip install pulumi_netbird
```

### Go
```bash
go get github.com/KitStream/netbird-pulumi-provider/sdk/go/index
```

## Configuration

The provider requires the following configuration settings:

| Variable | Description | Environment Variable |
| :--- | :--- | :--- |
| `token` | Admin PAT for NetBird Management Server | `NB_PAT` |
| `managementUrl` | NetBird Management API URL | `NB_MANAGEMENT_URL` |
| `tenantAccount` | Account ID to impersonate | `NB_ACCOUNT` |

You can set these using `pulumi config set`:

```bash
pulumi config set netbird:token YOUR_PAT --secret
pulumi config set netbird:managementUrl https://api.netbird.io
```

## Local Development

If you are developing or testing the provider locally, you can use the provided Python virtual environment to avoid conflicts with your system Python (e.g., Homebrew).

### Python Setup
A virtual environment is configured in the root directory. To activate it:

```bash
source .venv/bin/activate
```

This environment has the local Python SDK installed in editable mode, along with all necessary dependencies for running the examples.

## Documentation

For more information on the available resources and data sources, see the [Pulumi Registry](https://www.pulumi.com/registry/packages/netbird/).

## Releasing

For instructions on how to release new versions of this provider, see [RELEASING.md](./RELEASING.md).
