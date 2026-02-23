import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python network Group")

res = netbird.Network("test-network",
    name="Pulumi Python Network")

pulumi.export("resourceName", res.name)
