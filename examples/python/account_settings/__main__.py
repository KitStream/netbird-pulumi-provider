import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python account_settings Group")

res = netbird.AccountSettings("test-account-settings",
    peer_approval_enabled=True)

pulumi.export("resourceName", res.peerapprovalenabled)
