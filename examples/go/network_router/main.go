package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("Example Go network_router Group"),
		})
		if err != nil {
			return err
		}

		network, err := index.NewNetwork(ctx, "example-network", &index.NetworkArgs{
			Name: pulumi.String("Example Go network_router Network"),
		})
		if err != nil {
			return err
		}

		res, err := index.NewNetworkRouter(ctx, "test-network-router", &index.NetworkRouterArgs{
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
