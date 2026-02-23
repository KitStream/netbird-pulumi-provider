import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python setup_key Group")

res = netbird.SetupKey("test-setup-key",
    name="Pulumi Python Setup Key",
    type="reusable",
    expiry_seconds=86400,
    auto_groups=[group.id])

pulumi.export("resourceName", res.name)
