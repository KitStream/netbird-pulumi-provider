package main

import (
	"context"
	_ "embed"

	netbird "github.com/KitStream/netbird-pulumi-provider/provider"
	"github.com/pulumi/pulumi-terraform-bridge/pf/tfbridge"
)

var version string

//go:embed schema.json
var schema []byte

//go:embed bridge-metadata.json
var bridgeMetadata []byte

func main() {
	tfbridge.Main(context.Background(), "netbird", netbird.Provider(version), tfbridge.ProviderMetadata{
		PackageSchema:  schema,
		BridgeMetadata: bridgeMetadata,
	})
}
