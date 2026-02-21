package main

import (
	netbird "github.com/KitStream/pulumi-netbird/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

var version string

func main() {
	tfbridge.Main("netbird", version, netbird.Provider(version), nil)
}
