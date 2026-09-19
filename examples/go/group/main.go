package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		res, err := netbird.NewGroup(ctx, "test-group", &netbird.GroupArgs{
			Name: pulumi.String("Pulumi Go Group"),
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
