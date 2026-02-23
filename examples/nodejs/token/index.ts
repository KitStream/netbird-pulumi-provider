import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS token Group",
});

const user = new netbird.User("example-user", {
    email: "pulumi-nodejs-token-test@example.com",
    name: "Pulumi Token Test User",
    isServiceUser: true,
});

const res = new netbird.PersonalAccessToken("test-token", {
    name: "Pulumi NodeJS Token",
    expirationDays: 30,
    userId: user.id,
});

export const resourceName = res.name;
