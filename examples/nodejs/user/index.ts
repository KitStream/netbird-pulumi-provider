import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS user Group",
});

const res = new netbird.User("test-user", {
    email: "pulumi-NodeJS@example.com",
    name: "Pulumi NodeJS User",
    isServiceUser: true,
});

export const resourceName = res.email;
