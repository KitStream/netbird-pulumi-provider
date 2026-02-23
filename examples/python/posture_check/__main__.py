import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python posture_check Group")

res = netbird.PostureCheck("test-posture-check",
    name="Pulumi Python Posture Check",
    os_version_check=netbird.PostureCheckOsVersionCheckArgs(darwin_min_version="1.0.0"))

pulumi.export("resourceName", res.name)
