package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("Example Go setup_key Group"),
		})
		if err != nil {
			return err
		}

		res, err := index.NewSetupKey(ctx, "test-setup-key", &index.SetupKeyArgs{
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
