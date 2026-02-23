package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("Example Go nameserver_group Group"),
		})
		if err != nil {
			return err
		}

		res, err := index.NewNameserverGroup(ctx, "test-nameserver-group", &index.NameserverGroupArgs{
			Name: pulumi.String("Pulumi Go NS Group"),
			Description: pulumi.String("Pulumi Go NS Group"),
			Enabled: pulumi.Bool(true),
			Nameservers: index.NameserverGroupNameserverArray{index.NameserverGroupNameserverArgs{Ip: pulumi.String("1.1.1.1"), Port: pulumi.Int(53)}},
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
