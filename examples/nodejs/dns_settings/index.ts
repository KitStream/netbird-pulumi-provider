import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS dns_settings Group",
});

const res = new netbird.DnsSettings("test-dns-settings", {
    disabledManagementGroups: [group.id],
});

export const resourceName = res.disabledManagementGroups;
