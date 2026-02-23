import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS account_settings Group",
});

const res = new netbird.AccountSettings("test-account-settings", {
    peerApprovalEnabled: true,
});

export const resourceName = res.peerApprovalEnabled;
