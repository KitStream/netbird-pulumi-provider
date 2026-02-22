import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("test-group", {
    name: "Pulumi NodeJS Test Group",
});

export const groupName = group.name;
