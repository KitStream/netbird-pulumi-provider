import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python user Group")

res = netbird.User("test-user",
    email="pulumi-Python@example.com",
    name="Pulumi Python User",
    is_service_user=True)

pulumi.export("resourceName", res.email)
