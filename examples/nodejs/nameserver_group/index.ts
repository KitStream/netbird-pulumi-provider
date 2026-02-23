import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS nameserver_group Group",
});

const res = new netbird.NameserverGroup("test-nameserver-group", {
    name: "Pulumi NodeJS NS Group",
    description: "Pulumi NodeJS NS Group",
    enabled: true,
    nameservers: [{ ip: "1.1.1.1", port: 53 }],
    groups: [group.id],
    primary: true,
});

export const resourceName = res.name;
