import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python network_router Group")

network = netbird.Network("example-network",
    name="Example Python network_router Network")

res = netbird.NetworkRouter("test-network-router",
    network_id=network.id,
    peer_groups=[group.id])

pulumi.export("resourceName", res.networkid)
