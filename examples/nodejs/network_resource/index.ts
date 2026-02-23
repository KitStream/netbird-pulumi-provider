import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS network_resource Group",
});

const network = new netbird.Network("example-network", {
    name: "Example NodeJS network_resource Network",
});

const res = new netbird.NetworkResource("test-network-resource", {
    name: "Pulumi NodeJS Net Res",
    address: "10.20.0.0/24",
    networkId: network.id,
    groups: [group.id],
});

export const resourceName = res.name;
