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
