package netbird

import (
	"fmt"
	"path/filepath"

	"github.com/netbirdio/terraform-provider-netbird/shim/provider"
	pftfbridge "github.com/pulumi/pulumi-terraform-bridge/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

// all of the token components used below.
const (
	// main packages
	mainPkg = "netbird"
	// modules
	mainMod = "index" // the netbird module
)

// Provider returns additional overlaid control over the Terraform provider
func Provider() tfbridge.ProviderInfo {
	// Instantiate the Terraform provider
	p := pftfbridge.ShimProvider(provider.New("0.0.7")())

	// Create a Pulumi provider control struct
	prov := tfbridge.ProviderInfo{
		P:           p,
		Name:        "netbird",
		GitHubOrg:   "netbirdio",
		Description: "A Pulumi package for creating and managing netbird cloud resources.",
		Keywords:    []string{"pulumi", "netbird"},
		License:     "Apache-2.0",
		Homepage:    "https://netbird.io",
		Repository:  "https://github.com/KitStream/pulumi-netbird",
		Config:      map[string]*tfbridge.SchemaInfo{
			// Add any custom config mapping here
		},
		Resources: map[string]*tfbridge.ResourceInfo{
			"netbird_account_settings": {Tok: tfbridge.MakeResource(mainPkg, mainMod, "AccountSettings")},
			"netbird_dns_settings":     {Tok: tfbridge.MakeResource(mainPkg, mainMod, "DnsSettings")},
			"netbird_group":            {Tok: tfbridge.MakeResource(mainPkg, mainMod, "Group")},
			"netbird_nameserver_group": {Tok: tfbridge.MakeResource(mainPkg, mainMod, "NameserverGroup")},
			"netbird_network":          {Tok: tfbridge.MakeResource(mainPkg, mainMod, "Network")},
			"netbird_network_resource": {Tok: tfbridge.MakeResource(mainPkg, mainMod, "NetworkResource")},
			"netbird_network_router":   {Tok: tfbridge.MakeResource(mainPkg, mainMod, "NetworkRouter")},
			"netbird_peer":             {Tok: tfbridge.MakeResource(mainPkg, mainMod, "Peer")},
			"netbird_policy":           {Tok: tfbridge.MakeResource(mainPkg, mainMod, "Policy")},
			"netbird_posture_check":    {Tok: tfbridge.MakeResource(mainPkg, mainMod, "PostureCheck")},
			"netbird_route":            {Tok: tfbridge.MakeResource(mainPkg, mainMod, "Route")},
			"netbird_setup_key":        {Tok: tfbridge.MakeResource(mainPkg, mainMod, "SetupKey")},
			"netbird_token":            {Tok: tfbridge.MakeResource(mainPkg, mainMod, "Token")},
			"netbird_user":             {Tok: tfbridge.MakeResource(mainPkg, mainMod, "User")},
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			"netbird_account_settings": {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getAccountSettings")},
			"netbird_dns_settings":     {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getDnsSettings")},
			"netbird_group":            {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getGroup")},
			"netbird_nameserver_group": {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getNameserverGroup")},
			"netbird_network":          {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getNetwork")},
			"netbird_network_resource": {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getNetworkResource")},
			"netbird_network_router":   {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getNetworkRouter")},
			"netbird_peer":             {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getPeer")},
			"netbird_peers":            {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getPeers")},
			"netbird_policy":           {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getPolicy")},
			"netbird_posture_check":    {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getPostureCheck")},
			"netbird_route":            {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getRoute")},
			"netbird_setup_key":        {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getSetupKey")},
			"netbird_token":            {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getToken")},
			"netbird_user":             {Tok: tfbridge.MakeDataSource(mainPkg, mainMod, "getUser")},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			Dependencies: map[string]string{
				"@pulumi/pulumi": "^3.0.0",
			},
			DevDependencies: map[string]string{
				"@types/node": "^10.0.0", // so we can access native node types.
				"@types/mime": "^2.0.0",
			},
		},
		Python: &tfbridge.PythonInfo{
			Requires: map[string]string{
				"pulumi": ">=3.0.0,<4.0.0",
			},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/KitStream/pulumi-%[1]s/sdk/", mainPkg),
				"go",
				mainMod,
			),
			GenerateResourceContainerTypes: true,
		},
		CSharp: &tfbridge.CSharpInfo{
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
		Java: &tfbridge.JavaInfo{
			BasePackage: "com.netbird",
		},
	}

	return prov
}
