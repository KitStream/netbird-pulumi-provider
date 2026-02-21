package main

import (
	netbird "github.com/KitStream/pulumi-netbird/provider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfgen"
)

var version string

func main() {
	tfgen.Main("netbird", version, netbird.Provider(version))
}
