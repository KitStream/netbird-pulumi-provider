import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS route Group",
});

const res = new netbird.Route("test-route", {
    description: "Pulumi NodeJS Route",
    enabled: true,
    network: "10.0.0.0/24",
    networkId: "test-route",
    peerGroups: [group.id],
    groups: [group.id],
});

export const resourceName = res.description;
