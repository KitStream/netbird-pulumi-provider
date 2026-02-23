import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python nameserver_group Group")

res = netbird.NameserverGroup("test-nameserver-group",
    name="Pulumi Python NS Group",
    description="Pulumi Python NS Group",
    enabled=True,
    nameservers=[netbird.NameserverGroupNameserverArgs(ip="1.1.1.1", port=53)],
    groups=[group.id],
    primary=True)

pulumi.export("resourceName", res.name)
