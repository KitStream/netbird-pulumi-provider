import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS setup_key Group",
});

const res = new netbird.SetupKey("test-setup-key", {
    name: "Pulumi NodeJS Setup Key",
    type: "reusable",
    expirySeconds: 86400,
    autoGroups: [group.id],
});

export const resourceName = res.name;
