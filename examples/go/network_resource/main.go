package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go network_resource Group"),
		})
		if err != nil {
			return err
		}

		network, err := netbird.NewNetwork(ctx, "example-network", &netbird.NetworkArgs{
			Name: pulumi.String("Example Go network_resource Network"),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewNetworkResource(ctx, "test-network-resource", &netbird.NetworkResourceArgs{
			Name: pulumi.String("Pulumi Go Net Res"),
			Address: pulumi.String("10.20.0.0/24"),
			NetworkId: network.ID(),
			Groups: pulumi.StringArray{group.ID()},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
