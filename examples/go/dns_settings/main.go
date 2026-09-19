package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go dns_settings Group"),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewDnsSettings(ctx, "test-dns-settings", &netbird.DnsSettingsArgs{
			DisabledManagementGroups: pulumi.StringArray{group.ID()},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.DisabledManagementGroups)
		return nil
	})
}
