package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Arg struct {
	Name  string
	Value string
	Type  string // string, bool, int, string_array
}

type Resource struct {
	Name       string
	GoName     string
	PythonName string
	NodeJSName string
	DotNetName string
	Args       []Arg
	APIPath    string
	CheckField string
	CheckValue string
}

func main() {
	resources := []Resource{
		{
			Name:       "group",
			GoName:     "Group",
			PythonName: "Group",
			NodeJSName: "Group",
			DotNetName: "Group",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Group\"", Type: "string"},
			},
			APIPath:    "/api/groups",
			CheckField: "Name",
			CheckValue: "Pulumi %s Group",
		},
		{
			Name:       "setup_key",
			GoName:     "SetupKey",
			PythonName: "SetupKey",
			NodeJSName: "SetupKey",
			DotNetName: "SetupKey",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Setup Key\"", Type: "string"},
				{Name: "Type", Value: "\"reusable\"", Type: "string"},
				{Name: "ExpirySeconds", Value: "86400", Type: "int"},
				{Name: "AutoGroups", Value: "", Type: "string_array"},
			},
			APIPath:    "/api/setup-keys",
			CheckField: "Name",
			CheckValue: "Pulumi %s Setup Key",
		},
		{
			Name:       "user",
			GoName:     "User",
			PythonName: "User",
			NodeJSName: "User",
			DotNetName: "User",
			Args: []Arg{
				{Name: "Email", Value: "\"pulumi-%s@example.com\"", Type: "string"},
				{Name: "Name", Value: "\"Pulumi %s User\"", Type: "string"},
				{Name: "IsServiceUser", Value: "true", Type: "bool"},
			},
			APIPath:    "/api/users",
			CheckField: "Email",
			CheckValue: "pulumi-%s@example.com",
		},
		{
			Name:       "dns_settings",
			GoName:     "DnsSettings",
			PythonName: "DnsSettings",
			NodeJSName: "DnsSettings",
			DotNetName: "DnsSettings",
			Args: []Arg{
				{Name: "DisabledManagementGroups", Value: "", Type: "string_array"},
			},
			APIPath:    "/api/dns/settings",
			CheckField: "DisabledManagementGroups",
			CheckValue: "all", // Still use all as check value? No, I'll use Example Group's ID
		},
		{
			Name:       "account_settings",
			GoName:     "AccountSettings",
			PythonName: "AccountSettings",
			NodeJSName: "AccountSettings",
			DotNetName: "AccountSettings",
			Args: []Arg{
				{Name: "PeerApprovalEnabled", Value: "true", Type: "bool"},
			},
			APIPath:    "/api/accounts",
			CheckField: "PeerApprovalEnabled",
			CheckValue: "true",
		},
		{
			Name:       "network",
			GoName:     "Network",
			PythonName: "Network",
			NodeJSName: "Network",
			DotNetName: "Network",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Network\"", Type: "string"},
			},
			APIPath:    "/api/networks",
			CheckField: "Name",
			CheckValue: "Pulumi %s Network",
		},
		{
			Name:       "route",
			GoName:     "Route",
			PythonName: "Route",
			NodeJSName: "Route",
			DotNetName: "Route",
			Args: []Arg{
				{Name: "Description", Value: "\"Pulumi %s Route\"", Type: "string"},
				{Name: "Enabled", Value: "true", Type: "bool"},
				{Name: "Network", Value: "\"10.0.0.0/24\"", Type: "string"},
				{Name: "NetworkId", Value: "\"test-route\"", Type: "string"},
				{Name: "PeerGroups", Value: "", Type: "string_array"},
				{Name: "Groups", Value: "", Type: "string_array"},
			},
			APIPath:    "/api/routes",
			CheckField: "Description",
			CheckValue: "Pulumi %s Route",
		},
		{
			Name:       "posture_check",
			GoName:     "PostureCheck",
			PythonName: "PostureCheck",
			NodeJSName: "PostureCheck",
			DotNetName: "PostureCheck",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Posture Check\"", Type: "string"},
				{Name: "OsVersionCheck", Value: "", Type: "os_version_check"},
			},
			APIPath:    "/api/posture-checks",
			CheckField: "Name",
			CheckValue: "Pulumi %s Posture Check",
		},
		{
			Name:       "nameserver_group",
			GoName:     "NameserverGroup",
			PythonName: "NameserverGroup",
			NodeJSName: "NameserverGroup",
			DotNetName: "NameserverGroup",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s NS Group\"", Type: "string"},
				{Name: "Description", Value: "\"Pulumi %s NS Group\"", Type: "string"},
				{Name: "Enabled", Value: "true", Type: "bool"},
				{Name: "Nameservers", Value: "", Type: "nameserver_array"},
				{Name: "Groups", Value: "", Type: "string_array"},
				{Name: "Primary", Value: "true", Type: "bool"},
			},
			APIPath:    "/api/dns/nameservers",
			CheckField: "Name",
			CheckValue: "Pulumi %s NS Group",
		},
		{
			Name:       "network_resource",
			GoName:     "NetworkResource",
			PythonName: "NetworkResource",
			NodeJSName: "NetworkResource",
			DotNetName: "NetworkResource",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Net Res\"", Type: "string"},
				{Name: "Address", Value: "\"10.20.0.0/24\"", Type: "string"},
				{Name: "NetworkId", Value: "network.id", Type: "dependency"},
				{Name: "Groups", Value: "", Type: "string_array"},
			},
			APIPath:    "/api/networks/%s/resources", // Needs special handling in test
			CheckField: "Name",
			CheckValue: "Pulumi %s Net Res",
		},
		{
			Name:       "network_router",
			GoName:     "NetworkRouter",
			PythonName: "NetworkRouter",
			NodeJSName: "NetworkRouter",
			DotNetName: "NetworkRouter",
			Args: []Arg{
				{Name: "NetworkId", Value: "network.id", Type: "dependency"},
				{Name: "PeerGroups", Value: "", Type: "string_array"},
			},
			APIPath:    "/api/networks/%s/routers", // Needs special handling in test
			CheckField: "NetworkId",
			CheckValue: "test-network",
		},
		{
			Name:       "token",
			GoName:     "PersonalAccessToken",
			PythonName: "PersonalAccessToken",
			NodeJSName: "PersonalAccessToken",
			DotNetName: "PersonalAccessToken",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Token\"", Type: "string"},
				{Name: "ExpirationDays", Value: "30", Type: "int"},
				{Name: "UserId", Value: "user.id", Type: "dependency"},
			},
			APIPath:    "/api/users",
			CheckField: "Name",
			CheckValue: "Pulumi %s Token",
		},
		{
			Name:       "policy",
			GoName:     "Policy",
			PythonName: "Policy",
			NodeJSName: "Policy",
			DotNetName: "Policy",
			Args: []Arg{
				{Name: "Name", Value: "\"Pulumi %s Policy\"", Type: "string"},
				{Name: "Enabled", Value: "true", Type: "bool"},
				{Name: "Rule", Value: "", Type: "policy_rule"},
			},
			APIPath:    "/api/policies",
			CheckField: "Name",
			CheckValue: "Pulumi %s Policy",
		},
	}

	languages := []string{"go", "nodejs", "python", "dotnet", "java"}

	for _, lang := range languages {
		for _, res := range resources {
			dir := filepath.Join("examples", lang, res.Name)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Error creating directory %s: %v\n", dir, err)
				continue
			}

			generatePulumiYaml(dir, lang, res.Name)

			switch lang {
			case "go":
				generateGo(dir, res, lang)
			case "nodejs":
				generateNodeJS(dir, res, lang)
			case "python":
				generatePython(dir, res, lang)
			case "dotnet":
				generateDotNet(dir, res, lang)
			case "java":
				generateJava(dir, res, lang)
			}
		}
	}
}

func generateJava(dir string, res Resource, lang string) {
	args := ""
	for _, arg := range res.Args {
		val := arg.Value
		if arg.Type == "string" && strings.Contains(arg.Value, "%s") {
			val = fmt.Sprintf(arg.Value, "Java")
		}

		if strings.HasSuffix(arg.Name, "Groups") && val == "" {
			args += "                ." + strings.ToLower(arg.Name[:1]) + arg.Name[1:] + "(group.id().applyValue(List::of))\n"
			continue
		}

		k := strings.ToLower(arg.Name[:1]) + arg.Name[1:]
		switch arg.Type {
		case "string":
			args += "                ." + k + "(" + val + ")\n"
		case "bool":
			args += "                ." + k + "(" + val + ")\n"
		case "int":
			args += "                ." + k + "(" + val + ")\n"
		case "string_array":
			args += "                ." + k + "(List.of(" + val + "))\n"
		case "nameserver_array":
			args += "                ." + k + "(List.of(NameserverGroupNameserverArgs.builder().ip(\"1.1.1.1\").port(53).build()))\n"
		case "policy_rule":
			args += "                ." + k + "(PolicyRuleArgs.builder().action(\"accept\").enabled(true).name(\"rule1\").sources(List.of(group.id())).destinations(List.of(group.id())).build())\n"
		case "os_version_check":
			args += "                ." + k + "(PostureCheckOsVersionCheckArgs.builder().darwinMinVersion(\"1.0.0\").build())\n"
		case "dependency":
			javaVal := strings.Replace(val, ".id", ".id()", 1)
			args += "                ." + k + "(" + javaVal + ")\n"
		}
	}

	groupName := fmt.Sprintf("Example Java %s Group", res.Name)
	networkName := fmt.Sprintf("Example Java %s Network", res.Name)
	userName := fmt.Sprintf("pulumi-java-%s-test@example.com", res.Name)

	prefix := ""
	if res.Name != "group" {
		prefix = fmt.Sprintf(`            var group = new Group("example-group", GroupArgs.builder()
                .name("%s")
                .build());

`, groupName)
		if strings.Contains(args, "network.id") {
			prefix += fmt.Sprintf(`            var network = new Network("example-network", NetworkArgs.builder()
                .name("%s")
                .build());

`, networkName)
		}

		if strings.Contains(args, "user.id") {
			prefix += fmt.Sprintf(`            var user = new User("example-user", UserArgs.builder()
                .email("%s")
                .name("Pulumi Token Test User")
                .isServiceUser(true)
                .build());

`, userName)
		}
	}

	content := fmt.Sprintf(`package myproject;

import com.pulumi.Pulumi;
import io.github.kitstream.netbird.*;
import java.util.List;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
%s            var res = new %s("test-%s", %sArgs.builder()
%s                .build());

            ctx.export("resourceName", res.%s());
        });
    }
}
`, prefix, res.GoName, strings.ReplaceAll(res.Name, "_", "-"), res.GoName, args, strings.ToLower(res.CheckField[:1])+res.CheckField[1:])

	// Write App.java to proper Java source directory
	javaSrcDir := filepath.Join(dir, "src", "main", "java", "myproject")
	os.MkdirAll(javaSrcDir, 0755)
	os.WriteFile(filepath.Join(javaSrcDir, "App.java"), []byte(content), 0644)

	// Generate build.gradle
	buildGradle := fmt.Sprintf(`plugins {
    id("java")
    id("application")
}

group = "myproject"
version = "0.1.0"

repositories {
    mavenCentral()
    mavenLocal()
}

dependencies {
    implementation("com.pulumi:pulumi:1.0.0")
    implementation("io.github.kitstream:netbird:0.0.1")
}

application {
    mainClass = "myproject.App"
}
`)
	os.WriteFile(filepath.Join(dir, "build.gradle"), []byte(buildGradle), 0644)

	// Generate settings.gradle
	settingsGradle := fmt.Sprintf("rootProject.name = 'netbird-java-%s'\n", strings.ReplaceAll(res.Name, "_", "-"))
	os.WriteFile(filepath.Join(dir, "settings.gradle"), []byte(settingsGradle), 0644)
}

func generatePulumiYaml(dir, lang, resName string) {
	content := fmt.Sprintf("name: netbird-%s-%s\ndescription: A Pulumi %s program to create a netbird %s\nruntime: %s\n", lang, resName, lang, resName, lang)
	os.WriteFile(filepath.Join(dir, "Pulumi.yaml"), []byte(content), 0644)
}

func generateGo(dir string, res Resource, lang string) {
	args := ""
	for _, arg := range res.Args {
		val := arg.Value
		if arg.Type == "string" && strings.Contains(arg.Value, "%s") {
			val = fmt.Sprintf(arg.Value, "Go")
		}

		if strings.HasSuffix(arg.Name, "Groups") && val == "" {
			args += "\t\t\t" + arg.Name + ": pulumi.StringArray{group.ID()},\n"
			continue
		}

		switch arg.Type {
		case "string":
			args += "\t\t\t" + arg.Name + ": pulumi.String(" + val + "),\n"
		case "bool":
			args += "\t\t\t" + arg.Name + ": pulumi.Bool(" + val + "),\n"
		case "int":
			args += "\t\t\t" + arg.Name + ": pulumi.Int(" + val + "),\n"
		case "string_array":
			args += "\t\t\t" + arg.Name + ": pulumi.StringArray{pulumi.String(" + val + ")},\n"
		case "nameserver_array":
			args += "\t\t\t" + arg.Name + ": index.NameserverGroupNameserverArray{index.NameserverGroupNameserverArgs{Ip: pulumi.String(\"1.1.1.1\"), Port: pulumi.Int(53)}},\n"
		case "policy_rule":
			args += "\t\t\t" + arg.Name + ": &index.PolicyRuleArgs{Action: pulumi.String(\"accept\"), Enabled: pulumi.Bool(true), Name: pulumi.String(\"rule1\"), Sources: pulumi.StringArray{group.ID()}, Destinations: pulumi.StringArray{group.ID()}},\n"
		case "os_version_check":
			args += "\t\t\t" + arg.Name + ": &index.PostureCheckOsVersionCheckArgs{DarwinMinVersion: pulumi.String(\"1.0.0\")},\n"
		case "dependency":
			goVal := strings.Replace(val, ".id", ".ID()", 1)
			args += "\t\t\t" + arg.Name + ": " + goVal + ",\n"
		}
	}

	groupName := fmt.Sprintf("Example Go %s Group", res.Name)
	networkName := fmt.Sprintf("Example Go %s Network", res.Name)
	userName := fmt.Sprintf("pulumi-go-%s-test@example.com", res.Name)

	prefix := ""
	if res.Name != "group" {
		groupUsed := strings.Contains(args, "group.")
		groupSuffix := "\n"
		if !groupUsed {
			groupSuffix = "		_ = group\n\n"
		}
		prefix = fmt.Sprintf(`		group, err := index.NewGroup(ctx, "example-group", &index.GroupArgs{
			Name: pulumi.String("%s"),
		})
		if err != nil {
			return err
		}
`, groupName) + groupSuffix

		if strings.Contains(args, "network.") {
			prefix += fmt.Sprintf(`		network, err := index.NewNetwork(ctx, "example-network", &index.NetworkArgs{
			Name: pulumi.String("%s"),
		})
		if err != nil {
			return err
		}

`, networkName)
		}

		if strings.Contains(args, "user.") {
			prefix += fmt.Sprintf(`		user, err := index.NewUser(ctx, "example-user", &index.UserArgs{
			Email: pulumi.String("%s"),
			Name:  pulumi.String("Pulumi Token Test User"),
			IsServiceUser: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

`, userName)
		}
	}

	content := fmt.Sprintf(`package main

import (
	"github.com/KitStream/netbird-pulumi-provider/sdk/go/index"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
%s		res, err := index.New%s(ctx, "test-%s", &index.%sArgs{
%s		})
		if err != nil {
			return err
		}

		ctx.Export("resourceName", res.%s)
		return nil
	})
}
`, prefix, res.GoName, strings.ReplaceAll(res.Name, "_", "-"), res.GoName, args, strings.Title(res.CheckField))

	os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0644)

	// Copy go.mod from minimal and update it
	minimalGoMod, _ := os.ReadFile("examples/go/minimal/go.mod")
	newGoMod := strings.Replace(string(minimalGoMod), "minimal", res.Name, 1)
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(newGoMod), 0644)
}

func generateNodeJS(dir string, res Resource, lang string) {
	args := ""
	for _, arg := range res.Args {
		val := arg.Value
		if arg.Type == "string" && strings.Contains(arg.Value, "%s") {
			val = fmt.Sprintf(arg.Value, "NodeJS")
		}

		if strings.HasSuffix(arg.Name, "Groups") && val == "" {
			args += "    " + strings.ToLower(arg.Name[:1]) + arg.Name[1:] + ": [group.id],\n"
			continue
		}

		if arg.Type == "string_array" {
			val = "[" + val + "]"
		}
		k := strings.ToLower(arg.Name[:1]) + arg.Name[1:]
		switch arg.Type {
		case "nameserver_array":
			args += "    " + k + ": [{ ip: \"1.1.1.1\", port: 53 }],\n"
		case "policy_rule":
			args += "    " + k + ": { action: \"accept\", enabled: true, name: \"rule1\", sources: [group.id], destinations: [group.id] },\n"
		case "os_version_check":
			args += "    " + k + ": { darwinMinVersion: \"1.0.0\" },\n"
		default:
			args += "    " + k + ": " + val + ",\n"
		}
	}

	groupName := fmt.Sprintf("Example NodeJS %s Group", res.Name)
	networkName := fmt.Sprintf("Example NodeJS %s Network", res.Name)
	userName := fmt.Sprintf("pulumi-nodejs-%s-test@example.com", res.Name)

	prefix := ""
	if res.Name != "group" {
		prefix = fmt.Sprintf(`const group = new netbird.Group("example-group", {
    name: "%s",
});

`, groupName)
		if strings.Contains(args, "network.id") {
			prefix += fmt.Sprintf(`const network = new netbird.Network("example-network", {
    name: "%s",
});

`, networkName)
		}

		if strings.Contains(args, "user.id") {
			prefix += fmt.Sprintf(`const user = new netbird.User("example-user", {
    email: "%s",
    name: "Pulumi Token Test User",
    isServiceUser: true,
});

`, userName)
		}
	}

	content := fmt.Sprintf(`import * as netbird from "@kitstream/netbird-pulumi";

%sconst res = new netbird.%s("test-%s", {
%s});

export const resourceName = res.%s;
`, prefix, res.NodeJSName, strings.ReplaceAll(res.Name, "_", "-"), args, strings.ToLower(res.CheckField[:1])+res.CheckField[1:])

	os.WriteFile(filepath.Join(dir, "index.ts"), []byte(content), 0644)

	// Copy package.json from minimal and update it
	minimalPkg, _ := os.ReadFile("examples/nodejs/minimal/package.json")
	newPkg := strings.Replace(string(minimalPkg), "minimal-nodejs", "netbird-nodejs-"+res.Name, 1)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(newPkg), 0644)

	// Copy tsconfig.json
	tsconfig, _ := os.ReadFile("examples/nodejs/minimal/tsconfig.json")
	os.WriteFile(filepath.Join(dir, "tsconfig.json"), tsconfig, 0644)
}

func generatePython(dir string, res Resource, lang string) {
	args := ""
	for _, arg := range res.Args {
		val := arg.Value
		if arg.Type == "string" && strings.Contains(arg.Value, "%s") {
			val = fmt.Sprintf(arg.Value, "Python")
		}
		// Simple camel to snake
		snakeK := ""
		for i, r := range arg.Name {
			if i > 0 && r >= 'A' && r <= 'Z' {
				snakeK += "_"
			}
			snakeK += strings.ToLower(string(r))
		}

		if strings.HasSuffix(arg.Name, "Groups") && val == "" {
			args += "    " + snakeK + "=[group.id],\n"
			continue
		}

		if arg.Type == "string_array" {
			val = "[" + val + "]"
		}
		if arg.Type == "bool" {
			val = strings.Title(val)
		}
		switch arg.Type {
		case "nameserver_array":
			args += "    " + snakeK + "=[netbird.NameserverGroupNameserverArgs(ip=\"1.1.1.1\", port=53)],\n"
		case "policy_rule":
			args += "    " + snakeK + "=netbird.PolicyRuleArgs(action=\"accept\", enabled=True, name=\"rule1\", sources=[group.id], destinations=[group.id]),\n"
		case "os_version_check":
			args += "    " + snakeK + "=netbird.PostureCheckOsVersionCheckArgs(darwin_min_version=\"1.0.0\"),\n"
		default:
			args += "    " + snakeK + "=" + val + ",\n"
		}
	}

	groupName := fmt.Sprintf("Example Python %s Group", res.Name)
	networkName := fmt.Sprintf("Example Python %s Network", res.Name)
	userName := fmt.Sprintf("pulumi-python-%s-test@example.com", res.Name)

	prefix := ""
	if res.Name != "group" {
		prefix = fmt.Sprintf(`group = netbird.Group("example-group",
    name="%s")

`, groupName)
		if strings.Contains(args, "network_id=network.id") {
			prefix += fmt.Sprintf(`network = netbird.Network("example-network",
    name="%s")

`, networkName)
		}

		if strings.Contains(args, "user_id=user.id") {
			prefix += fmt.Sprintf(`user = netbird.User("example-user",
    email="%s",
    name="Pulumi Token Test User",
    is_service_user=True)

`, userName)
		}
	}

	content := fmt.Sprintf(`import pulumi
import pulumi_netbird as netbird

%sres = netbird.%s("test-%s",
%s)

pulumi.export("resourceName", res.%s)
`, prefix, res.PythonName, strings.ReplaceAll(res.Name, "_", "-"), strings.TrimSuffix(args, ",\n"), strings.ToLower(res.CheckField))

	os.WriteFile(filepath.Join(dir, "__main__.py"), []byte(content), 0644)

	// Copy requirements.txt
	reqs, _ := os.ReadFile("examples/python/minimal/requirements.txt")
	os.WriteFile(filepath.Join(dir, "requirements.txt"), reqs, 0644)
}

func generateDotNet(dir string, res Resource, lang string) {
	args := ""
	for _, arg := range res.Args {
		val := arg.Value
		if arg.Type == "string" && strings.Contains(arg.Value, "%s") {
			val = fmt.Sprintf(arg.Value, "DotNet")
		}

		if strings.HasSuffix(arg.Name, "Groups") && val == "" {
			args += "        " + arg.Name + " = { group.Id },\n"
			continue
		}

		if arg.Type == "string_array" {
			val = "new[] { " + val + " }"
		}
		switch arg.Type {
		case "nameserver_array":
			args += "        " + arg.Name + " = { new NameserverGroupNameserverArgs { Ip = \"1.1.1.1\", Port = 53 } },\n"
		case "policy_rule":
			args += "        " + arg.Name + " = new PolicyRuleArgs { Action = \"accept\", Enabled = true, Name = \"rule1\", Sources = { group.Id }, Destinations = { group.Id } },\n"
		case "os_version_check":
			args += "        " + arg.Name + " = new PostureCheckOsVersionCheckArgs { DarwinMinVersion = \"1.0.0\" },\n"
		case "dependency":
			dotnetVal := strings.Replace(val, ".id", ".Id", 1)
			args += "        " + arg.Name + " = " + dotnetVal + ",\n"
		default:
			args += "        " + arg.Name + " = " + val + ",\n"
		}
	}

	groupName := fmt.Sprintf("Example DotNet %s Group", res.Name)
	networkName := fmt.Sprintf("Example DotNet %s Network", res.Name)
	userName := fmt.Sprintf("pulumi-dotnet-%s-test@example.com", res.Name)

	prefix := ""
	if res.Name != "group" {
		prefix = fmt.Sprintf(`    var group = new Group("example-group", new GroupArgs
    {
        Name = "%s",
    });

`, groupName)
		if strings.Contains(args, "NetworkId = network.Id") {
			prefix += fmt.Sprintf(`    var network = new Network("example-network", new NetworkArgs
    {
        Name = "%s",
    });

`, networkName)
		}

		if strings.Contains(args, "UserId = user.Id") {
			prefix += fmt.Sprintf(`    var user = new User("example-user", new UserArgs
    {
        Email = "%s",
        Name = "Pulumi Token Test User",
        IsServiceUser = true,
    });

`, userName)
		}
	}

	content := fmt.Sprintf(`using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
%s    var res = new %s("test-%s", new %sArgs
    {
%s    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.%s,
    };
});
`, prefix, res.DotNetName, strings.ReplaceAll(res.Name, "_", "-"), res.DotNetName, args, strings.Title(res.CheckField))

	os.WriteFile(filepath.Join(dir, "Program.cs"), []byte(content), 0644)

	// Copy csproj
	csproj, _ := os.ReadFile("examples/dotnet/minimal/minimal-dotnet.csproj")
	newCsproj := strings.Replace(string(csproj), "minimal-dotnet", "netbird-dotnet-"+res.Name, 1)
	os.WriteFile(filepath.Join(dir, res.Name+"-dotnet.csproj"), []byte(newCsproj), 0644)
}
