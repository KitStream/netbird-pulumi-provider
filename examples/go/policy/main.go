package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("Example Go policy Group"),
		})
		if err != nil {
			return err
		}

		res, err := index.NewPolicy(ctx, "test-policy", &index.PolicyArgs{
			Name: pulumi.String("Pulumi Go Policy"),
			Enabled: pulumi.Bool(true),
			Rule: &index.PolicyRuleArgs{Action: pulumi.String("accept"), Enabled: pulumi.Bool(true), Name: pulumi.String("rule1"), Sources: pulumi.StringArray{group.ID()}, Destinations: pulumi.StringArray{group.ID()}},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
