package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go network_router Group"),
		})
		if err != nil {
			return err
		}

		network, err := netbird.NewNetwork(ctx, "example-network", &netbird.NetworkArgs{
			Name: pulumi.String("Example Go network_router Network"),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewNetworkRouter(ctx, "test-network-router", &netbird.NetworkRouterArgs{
			NetworkId: network.ID(),
			PeerGroups: pulumi.StringArray{group.ID()},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.NetworkId)
		return nil
	})
}
