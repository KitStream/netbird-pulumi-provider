import pulumi
import pulumi_netbird as netbird

group = netbird.Group("example-group",
    name="Example Python dns_settings Group")

res = netbird.DnsSettings("test-dns-settings",
    disabled_management_groups=[group.id])

pulumi.export("resourceName", res.disabledmanagementgroups)
