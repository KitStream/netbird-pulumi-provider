import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python route Group")

res = netbird.Route("test-route",
    description="Pulumi Python Route",
    enabled=True,
    network="10.0.0.0/24",
    network_id="test-route",
    peer_groups=[group.id],
    groups=[group.id])

pulumi.export("resourceName", res.description)
