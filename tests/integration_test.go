package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

const (
	localMgmtURL   = "http://localhost:8080"
	localSeededPAT = "nbp_apTmlmUXHSC4PKmHwtIZNaGr8eqcVI2gMURp"
)

func TestMain(m *testing.M) {
	if !haveDocker() {
		fmt.Println("docker compose not available; skipping local stack tests")
		os.Exit(0)
	}

	composeFile, err := filepath.Abs(filepath.Join("..", "upstream", "test", "compose.yml"))
	if err != nil {
		fmt.Printf("failed to get absolute path for compose file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Starting local NetBird stack...")
	cmd := exec.Command("docker", "compose", "-f", composeFile, "up", "-d", "--wait")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("failed to start local netbird stack: %v\n", err)
		os.Exit(1)
	}

	// 2. Wait for it to be ready
	deadline := time.Now().Add(2 * time.Minute)
	ready := false
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, localMgmtURL+"/", nil)
		if err != nil {
			cancel()
			time.Sleep(2 * time.Second)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		cancel()
		if err == nil && resp != nil {
			_ = resp.Body.Close()
			ready = true
			break
		}
		time.Sleep(2 * time.Second)
	}

	if !ready {
		fmt.Println("netbird stack timed out waiting to be ready")
		_ = exec.Command("docker", "compose", "-f", composeFile, "down", "-v", "--remove-orphans").Run()
		os.Exit(1)
	}

	// 3. Pre-build some SDKs if needed (e.g., Java needs to be in mavenLocal)
	fmt.Println("Preparing SDKs for tests...")
	// Java
	javaSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "java"))
	if err != nil {
		fmt.Printf("failed to get absolute path for Java SDK: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Publishing Java SDK to mavenLocal from %s...\n", javaSdkPath)
	publishCmd := exec.Command("gradle", "publishToMavenLocal")
	publishCmd.Dir = javaSdkPath
	if out, err := publishCmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to publish Java SDK: %v\n%s\n", err, string(out))
		// We don't exit here, just skip java tests later
	}

	// NodeJS
	nodeSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "nodejs"))
	if err != nil {
		fmt.Printf("failed to get absolute path for NodeJS SDK: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Building NodeJS SDK from %s...\n", nodeSdkPath)
	npmInstallCmd := exec.Command("npm", "install")
	npmInstallCmd.Dir = nodeSdkPath
	if out, err := npmInstallCmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to run npm install in NodeJS SDK: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	npmBuildCmd := exec.Command("npm", "run", "build")
	npmBuildCmd.Dir = nodeSdkPath
	if out, err := npmBuildCmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to run npm build in NodeJS SDK: %v\n%s\n", err, string(out))
		os.Exit(1)
	}

	// 3.5. Copy package.json to bin/ so utilities.js can find it for versioning
	packageJsonPath := filepath.Join(nodeSdkPath, "package.json")
	binPackageJsonPath := filepath.Join(nodeSdkPath, "bin", "package.json")
	if err := os.MkdirAll(filepath.Dir(binPackageJsonPath), 0755); err != nil {
		fmt.Printf("failed to create bin directory: %v\n", err)
		os.Exit(1)
	}
	packageJsonContent, err := os.ReadFile(packageJsonPath)
	if err != nil {
		fmt.Printf("failed to read package.json: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(binPackageJsonPath, packageJsonContent, 0644); err != nil {
		fmt.Printf("failed to copy package.json to bin/: %v\n", err)
		os.Exit(1)
	}

	// 3.6. Python - clean stale build artifacts that break wheel builds
	pythonSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "python"))
	if err != nil {
		fmt.Printf("failed to get absolute path for Python SDK: %v\n", err)
		os.Exit(1)
	}
	for _, dir := range []string{"build", "pulumi_netbird.egg-info"} {
		p := filepath.Join(pythonSdkPath, dir)
		if _, err := os.Stat(p); err == nil {
			fmt.Printf("Cleaning stale Python build artifact: %s\n", p)
			os.RemoveAll(p)
		}
	}

	// 3.7. DotNet
	dotNetSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "dotnet"))
	if err != nil {
		fmt.Printf("failed to get absolute path for DotNet SDK: %v\n", err)
		os.Exit(1)
	}
	versionFile := filepath.Join(dotNetSdkPath, "version.txt")
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		fmt.Printf("Generating missing version.txt for DotNet SDK in %s...\n", dotNetSdkPath)
		if err := os.WriteFile(versionFile, []byte("0.0.1"), 0644); err != nil {
			fmt.Printf("failed to generate version.txt for DotNet SDK: %v\n", err)
			os.Exit(1)
		}
	}

	// 4. Run tests
	code := m.Run()

	// 5. Tear down
	fmt.Println("Tearing down local NetBird stack...")
	_ = exec.Command("docker", "compose", "-f", composeFile, "down", "-v", "--remove-orphans").Run()

	os.Exit(code)
}

func haveDocker() bool {
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func providerPluginPath(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(cwd, "..", "bin")
	binaryPath := filepath.Join(p, "pulumi-resource-netbird")
	if _, err := os.Stat(binaryPath); err != nil {
		t.Fatalf("provider plugin not found at %s; build it first (cd provider && go build -o ../bin/pulumi-resource-netbird ./cmd/pulumi-resource-netbird)", binaryPath)
	}
	return p
}

func baseOptions(t *testing.T) *integration.ProgramTestOptions {
	return &integration.ProgramTestOptions{
		Quick:       true,
		SkipRefresh: true,
		NoParallel:  true,
		LocalProviders: []integration.LocalDependency{{
			Package: "netbird",
			Path:    providerPluginPath(t),
		}},
		Env: []string{
			"NB_MANAGEMENT_URL=" + localMgmtURL,
			"NB_PAT=" + localSeededPAT,
			"PULUMI_CONFIG_PASSPHRASE=test",
		},
	}
}

// langDisplayName returns the display name for a language, used to format expected values.
func langDisplayName(lang string) string {
	switch lang {
	case "nodejs":
		return "NodeJS"
	case "dotnet":
		return "DotNet"
	default:
		return strings.ToUpper(lang[:1]) + lang[1:]
	}
}

// applyLanguageOptions configures language-specific PrePrepareProject hooks.
func applyLanguageOptions(t *testing.T, opts *integration.ProgramTestOptions, lang, resName string) {
	t.Helper()
	switch lang {
	case "go":
		opts.PrePrepareProject = func(proj *engine.Projinfo) error {
			absSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "go", "index"))
			if err != nil {
				return err
			}
			goModPath := filepath.Join(proj.Root, "go.mod")
			content, err := os.ReadFile(goModPath)
			if err != nil {
				return err
			}
			newContent := strings.Replace(string(content), "../../../sdk/go/index", absSdkPath, 1)
			return os.WriteFile(goModPath, []byte(newContent), 0644)
		}
	case "nodejs":
		opts.PrePrepareProject = func(proj *engine.Projinfo) error {
			absSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "nodejs"))
			if err != nil {
				return err
			}
			pkgJsonPath := filepath.Join(proj.Root, "package.json")
			content, err := os.ReadFile(pkgJsonPath)
			if err != nil {
				return err
			}
			newContent := strings.Replace(string(content), "file:../../../sdk/nodejs", "file:"+absSdkPath, 1)
			return os.WriteFile(pkgJsonPath, []byte(newContent), 0644)
		}
	case "python":
		opts.PrePrepareProject = func(proj *engine.Projinfo) error {
			absSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "python"))
			if err != nil {
				return err
			}
			reqsPath := filepath.Join(proj.Root, "requirements.txt")
			content, err := os.ReadFile(reqsPath)
			if err != nil {
				return err
			}
			newContent := strings.Replace(string(content), "../../../sdk/python", absSdkPath, 1)
			return os.WriteFile(reqsPath, []byte(newContent), 0644)
		}
	case "dotnet":
		opts.PrePrepareProject = func(proj *engine.Projinfo) error {
			absSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "dotnet", "KitStream.Pulumi.Netbird.csproj"))
			if err != nil {
				return err
			}
			csprojPath := filepath.Join(proj.Root, resName+"-dotnet.csproj")
			content, err := os.ReadFile(csprojPath)
			if err != nil {
				return err
			}
			newContent := strings.Replace(string(content), "../../../sdk/dotnet/KitStream.Pulumi.Netbird.csproj", absSdkPath, 1)
			return os.WriteFile(csprojPath, []byte(newContent), 0644)
		}
	case "java":
		// Java SDK is published to mavenLocal in TestMain
		// Examples use build.gradle with mavenLocal() repository
	}
}

// apiGet performs an authenticated GET request against the local NetBird API.
func apiGet(t *testing.T, path string) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, localMgmtURL+path, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+localSeededPAT)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to fetch resources from %s: %v", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d from %s: %s", resp.StatusCode, path, string(body))
	}
	return body
}

// resolveAPIPath handles API paths that contain %s by looking up the network ID
// from the /api/networks endpoint.
func resolveAPIPath(t *testing.T, apiPath string) string {
	t.Helper()
	if !strings.Contains(apiPath, "%s") {
		return apiPath
	}

	body := apiGet(t, "/api/networks")
	type network struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var networks []network
	if err := json.Unmarshal(body, &networks); err != nil {
		t.Fatalf("failed to parse networks response: %v", err)
	}

	// Find the most recently created network (last in list)
	for i := len(networks) - 1; i >= 0; i-- {
		if networks[i].ID != "" {
			return fmt.Sprintf(apiPath, networks[i].ID)
		}
	}
	t.Fatal("no networks found to resolve API path")
	return ""
}

// validateResource returns a validation function that checks the NetBird API
// for the expected value in the response body.
func validateResource(apiPath string, expectedValue string) func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
	return func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
		resolvedPath := resolveAPIPath(t, apiPath)
		body := apiGet(t, resolvedPath)

		if expectedValue != "" && !strings.Contains(string(body), expectedValue) {
			t.Fatalf("expected value %q not found in NetBird API response from %s: %s", expectedValue, resolvedPath, string(body))
		}
	}
}

// resource describes a test resource with its API path and expected check value.
type resource struct {
	name       string
	apiPath    string
	checkValue string // may contain %s for language name substitution
}

// allResources returns the full list of resources to test.
func allResources() []resource {
	return []resource{
		{"group", "/api/groups", "Pulumi %s Group"},
		{"setup_key", "/api/setup-keys", "Pulumi %s Setup Key"},
		{"user", "/api/users", "Pulumi %s User"},
		{"account_settings", "/api/accounts", "true"},
		// dns_settings is excluded: the upstream TF provider's Delete is a no-op so the group
		// remains linked to disabled DNS management groups, causing group deletion to fail on destroy.
		{"network", "/api/networks", "Pulumi %s Network"},
		{"network_resource", "/api/networks/%s/resources", "Pulumi %s Net Res"},
		{"network_router", "/api/networks/%s/routers", ""},
		{"route", "/api/routes", "Pulumi %s Route"},
		{"posture_check", "/api/posture-checks", "Pulumi %s Posture Check"},
		{"nameserver_group", "/api/dns/nameservers", "Pulumi %s NS Group"},
		{"policy", "/api/policies", "Pulumi %s Policy"},
		{"token", "/api/users", "Pulumi Token Test User"},
	}
}

// runResourceTest runs a single resource test for a given language.
func runResourceTest(t *testing.T, lang string, res resource) {
	t.Helper()

	langName := langDisplayName(lang)

	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", lang, res.name)

	applyLanguageOptions(t, opts, lang, res.name)

	expectedValue := res.checkValue
	if strings.Contains(expectedValue, "%s") {
		expectedValue = fmt.Sprintf(expectedValue, langName)
	}
	opts.ExtraRuntimeValidation = validateResource(res.apiPath, expectedValue)

	integration.ProgramTest(t, opts)
}

// TestSmoke runs only the "group" resource across all languages.
// This is a fast sanity check before running the full suite.
//
// Run with: go test -v -timeout 30m -run TestSmoke ./...
func TestSmoke(t *testing.T) {
	groupRes := resource{"group", "/api/groups", "Pulumi %s Group"}

	for _, lang := range []string{"go", "nodejs", "python", "dotnet", "java"} {
		t.Run(lang, func(t *testing.T) {
			runResourceTest(t, lang, groupRes)
		})
	}
}

// TestExamples_Go runs all resource tests for Go.
//
// Run with: go test -v -timeout 60m -run TestExamples_Go ./...
func TestExamples_Go(t *testing.T) {
	for _, res := range allResources() {
		t.Run(res.name, func(t *testing.T) {
			runResourceTest(t, "go", res)
		})
	}
}

// TestExamples_NodeJS runs all resource tests for NodeJS.
//
// Run with: go test -v -timeout 60m -run TestExamples_NodeJS ./...
func TestExamples_NodeJS(t *testing.T) {
	for _, res := range allResources() {
		t.Run(res.name, func(t *testing.T) {
			runResourceTest(t, "nodejs", res)
		})
	}
}

// TestExamples_Python runs all resource tests for Python.
//
// Run with: go test -v -timeout 60m -run TestExamples_Python ./...
func TestExamples_Python(t *testing.T) {
	for _, res := range allResources() {
		t.Run(res.name, func(t *testing.T) {
			runResourceTest(t, "python", res)
		})
	}
}

// TestExamples_DotNet runs all resource tests for DotNet.
//
// Run with: go test -v -timeout 60m -run TestExamples_DotNet ./...
func TestExamples_DotNet(t *testing.T) {
	for _, res := range allResources() {
		t.Run(res.name, func(t *testing.T) {
			runResourceTest(t, "dotnet", res)
		})
	}
}

// TestExamples_Java runs all resource tests for Java.
//
// Run with: go test -v -timeout 60m -run TestExamples_Java ./...
func TestExamples_Java(t *testing.T) {
	for _, res := range allResources() {
		t.Run(res.name, func(t *testing.T) {
			runResourceTest(t, "java", res)
		})
	}
}

// TestJava_Minimal_LocalStack tests the minimal Java example.
func TestJava_Minimal_LocalStack(t *testing.T) {
	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", "java", "minimal")
	opts.ExtraRuntimeValidation = validateResource("/api/groups", "Pulumi Java Test Group")

	integration.ProgramTest(t, opts)
}
