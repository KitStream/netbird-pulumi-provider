package main

import (
	netbird "github.com/KitStream/pulumi-netbird/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

func main() {
	tfbridge.Main("netbird", "0.0.1", netbird.Provider(), nil)
}
