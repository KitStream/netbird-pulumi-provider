import * as netbird from "@kitstream/netbird-pulumi";

const group = new netbird.Group("example-group", {
    name: "Example NodeJS posture_check Group",
});

const res = new netbird.PostureCheck("test-posture-check", {
    name: "Pulumi NodeJS Posture Check",
    osVersionCheck: { darwinMinVersion: "1.0.0" },
});

export const resourceName = res.name;
