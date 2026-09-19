package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go nameserver_group Group"),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewNameserverGroup(ctx, "test-nameserver-group", &netbird.NameserverGroupArgs{
			Name: pulumi.String("Pulumi Go NS Group"),
			Description: pulumi.String("Pulumi Go NS Group"),
			Enabled: pulumi.Bool(true),
			Nameservers: netbird.NameserverGroupNameserverArray{netbird.NameserverGroupNameserverArgs{Ip: pulumi.String("1.1.1.1"), Port: pulumi.Int(53)}},
			Groups: pulumi.StringArray{group.ID()},
			Primary: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
