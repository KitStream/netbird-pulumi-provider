---
title: NetBird Installation & Configuration
meta_desc: Information on how to install the NetBird provider for Pulumi.
layout: installation
---

## Installation

The NetBird provider is available as a package in all Pulumi languages:

* JavaScript/TypeScript: [`@kitstream/netbird-pulumi`](https://www.npmjs.com/package/@kitstream/netbird-pulumi)
* Python: [`pulumi_netbird`](https://pypi.org/project/pulumi-netbird/)
* Go: [`github.com/KitStream/netbird-pulumi-provider/sdk/go/index`](https://github.com/KitStream/netbird-pulumi-provider)
* .NET: [`KitStream.Pulumi.Netbird`](https://www.nuget.org/packages/KitStream.Pulumi.Netbird)
* Java: [`io.github.kitstream:netbird`](https://central.sonatype.com/artifact/io.github.kitstream/netbird)

## Configuration

The following configuration options are available for the `netbird` provider:

| Option | Required | Description | Environment Variable |
|--------|----------|-------------|---------------------|
| `netbird:token` | Yes | Admin PAT for NetBird Management Server | `NB_PAT` |
| `netbird:managementUrl` | No | NetBird Management API URL | `NB_MANAGEMENT_URL` |
| `netbird:tenantAccount` | No | Account ID to impersonate | `NB_ACCOUNT` |

### Setting Configuration via Pulumi Config

```bash
pulumi config set --secret netbird:token <your-pat-token>
pulumi config set netbird:managementUrl https://api.netbird.io
```

### Setting Configuration via Environment Variables

```bash
export NB_PAT=<your-pat-token>
export NB_MANAGEMENT_URL=https://api.netbird.io
```

Configuration values defined in Pulumi config take precedence over environment variables.
