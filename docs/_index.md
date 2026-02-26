---
title: NetBird
meta_desc: Provides an overview of the NetBird Provider for Pulumi.
layout: overview
---

The NetBird provider for Pulumi can be used to provision and manage
[NetBird](https://netbird.io) cloud resources. NetBird is a zero-trust
networking platform that enables secure peer-to-peer connectivity.

Use this provider to manage groups, networks, routes, DNS settings, policies,
setup keys, and more through Pulumi's infrastructure-as-code framework.

## Example

{{< chooser language "typescript,python,go,csharp,java" >}}

{{% choosable language typescript %}}

```typescript
import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("test-group", {
    name: "Pulumi NodeJS Test Group",
});

export const groupName = group.name;
```

{{% /choosable %}}

{{% choosable language python %}}

```python
import pulumi
import pulumi_netbird as netbird

group = netbird.Group("test-group",
    name="Pulumi Python Test Group")

pulumi.export("groupName", group.name)
```

{{% /choosable %}}

{{% choosable language go %}}

```go
package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "test-group", &index.GroupArgs{
			Name: pulumi.String("Pulumi Go Test Group"),
		})
		if err != nil {
			return err
		}

		ctx.Export("groupName", group.Name)
		return nil
	})
}
```

{{% /choosable %}}

{{% choosable language csharp %}}

```csharp
using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() =>
{
    var group = new Group("test-group", new GroupArgs
    {
        Name = "Pulumi DotNet Test Group",
    });

    return new Dictionary<string, object?>
    {
        ["groupName"] = group.Name,
    };
});
```

{{% /choosable %}}

{{% choosable language java %}}

```java
package com.minimal;

import com.pulumi.Pulumi;
import io.github.kitstream.netbird.Group;
import io.github.kitstream.netbird.GroupArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var group = new Group("test-group", GroupArgs.builder()
                .name("Pulumi Java Test Group")
                .build());

            ctx.export("groupName", group.name());
        });
    }
}
```

{{% /choosable %}}

{{< /chooser >}}
