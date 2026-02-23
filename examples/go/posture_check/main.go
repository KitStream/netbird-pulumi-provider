package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("Example Go posture_check Group"),
		})
		if err != nil {
			return err
		}
		_ = group

		res, err := index.NewPostureCheck(ctx, "test-posture-check", &index.PostureCheckArgs{
			Name: pulumi.String("Pulumi Go Posture Check"),
			OsVersionCheck: &index.PostureCheckOsVersionCheckArgs{DarwinMinVersion: pulumi.String("1.0.0")},
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
