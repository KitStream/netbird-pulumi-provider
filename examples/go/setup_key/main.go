package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go setup_key Group"),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewSetupKey(ctx, "test-setup-key", &netbird.SetupKeyArgs{
			Name: pulumi.String("Pulumi Go Setup Key"),
			Type: pulumi.String("reusable"),
			ExpirySeconds: pulumi.Int(86400),
			AutoGroups: pulumi.StringArray{group.ID()},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
