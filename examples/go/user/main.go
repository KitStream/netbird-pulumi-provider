package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go user Group"),
		})
		if err != nil {
			return err
		}
		_ = group

		res, err := netbird.NewUser(ctx, "test-user", &netbird.UserArgs{
			Email: pulumi.String("pulumi-Go@example.com"),
			Name: pulumi.String("Pulumi Go User"),
			IsServiceUser: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Email)
		return nil
	})
}
