import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python policy Group")

res = netbird.Policy("test-policy",
    name="Pulumi Python Policy",
    enabled=True,
    rule=netbird.PolicyRuleArgs(action="accept", enabled=True, name="rule1", sources=[group.id], destinations=[group.id]))

pulumi.export("resourceName", res.name)
