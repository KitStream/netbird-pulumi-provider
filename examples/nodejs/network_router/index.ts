import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS network_router Group",
});

const network = new netbird.Network("example-network", {
    name: "Example NodeJS network_router Network",
});

const res = new netbird.NetworkRouter("test-network-router", {
    networkId: network.id,
    peerGroups: [group.id],
});

export const resourceName = res.networkId;
