import * as netbird from "@kitstream/netbird-pulumi";

const res = new netbird.Group("test-group", {
    name: "Pulumi NodeJS Group",
});

export const resourceName = res.name;
