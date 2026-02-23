import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS policy Group",
});

const res = new netbird.Policy("test-policy", {
    name: "Pulumi NodeJS Policy",
    enabled: true,
    rule: { action: "accept", enabled: true, name: "rule1", sources: [group.id], destinations: [group.id] },
});

export const resourceName = res.name;
