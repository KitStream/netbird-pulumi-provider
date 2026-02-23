import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python network_resource Group")

network = netbird.Network("example-network",
    name="Example Python network_resource Network")

res = netbird.NetworkResource("test-network-resource",
    name="Pulumi Python Net Res",
    address="10.20.0.0/24",
    network_id=network.id,
    groups=[group.id])

pulumi.export("resourceName", res.name)
