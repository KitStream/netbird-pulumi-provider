package main

import (
	netbird "github.com/KitStream/netbird-pulumi-provider/provider"
	"github.com/pulumi/pulumi-terraform-bridge/pf/tfgen"
)

var version string

func main() {
	tfgen.Main("netbird", netbird.Provider(version))
}
