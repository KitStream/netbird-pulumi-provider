import pulumi
import pulumi_netbird as netbird

res = netbird.Group("test-group",
    name="Pulumi Python Group")

pulumi.export("resourceName", res.name)
