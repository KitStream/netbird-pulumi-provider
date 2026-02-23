import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS network Group",
});

const res = new netbird.Network("test-network", {
    name: "Pulumi NodeJS Network",
});

export const resourceName = res.name;
