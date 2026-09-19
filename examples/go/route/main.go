package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go route Group"),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewRoute(ctx, "test-route", &netbird.RouteArgs{
			Description: pulumi.String("Pulumi Go Route"),
			Enabled: pulumi.Bool(true),
			Network: pulumi.String("10.0.0.0/24"),
			NetworkId: pulumi.String("test-route"),
			PeerGroups: pulumi.StringArray{group.ID()},
			Groups: pulumi.StringArray{group.ID()},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Description)
		return nil
	})
}
