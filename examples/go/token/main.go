package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		group, err := netbird.NewGroup(ctx, "example-group", &netbird.GroupArgs{
			Name: pulumi.String("Example Go token Group"),
		})
		if err != nil {
			return err
		}
		_ = group

		user, err := netbird.NewUser(ctx, "example-user", &netbird.UserArgs{
			Email: pulumi.String("pulumi-go-token-test@example.com"),
			Name:  pulumi.String("Pulumi Token Test User"),
			IsServiceUser: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		res, err := netbird.NewPersonalAccessToken(ctx, "test-token", &netbird.PersonalAccessTokenArgs{
			Name: pulumi.String("Pulumi Go Token"),
			ExpirationDays: pulumi.Int(30),
			UserId: user.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.Name)
		return nil
	})
}
