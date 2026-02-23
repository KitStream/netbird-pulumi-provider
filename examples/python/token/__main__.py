import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python token Group")

user = netbird.User("example-user",
    email="pulumi-python-token-test@example.com",
    name="Pulumi Token Test User",
    is_service_user=True)

res = netbird.PersonalAccessToken("test-token",
    name="Pulumi Python Token",
    expiration_days=30,
    user_id=user.id)

pulumi.export("resourceName", res.name)
