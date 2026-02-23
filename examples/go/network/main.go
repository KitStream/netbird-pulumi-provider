package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("Example Go network Group"),
		})
		if err != nil {
			return err
		}
		_ = group

		res, err := index.NewNetwork(ctx, "test-network", &index.NetworkArgs{
			Name: pulumi.String("Pulumi Go Network"),
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
