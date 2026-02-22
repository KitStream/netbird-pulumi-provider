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

	// DotNet
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

	// NodeJS Example
	nodeExamplePath, err := filepath.Abs(filepath.Join("..", "examples", "nodejs", "minimal"))
	if err != nil {
		fmt.Printf("failed to get absolute path for NodeJS example: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Installing dependencies for NodeJS example in %s...\n", nodeExamplePath)
	npmInstallExampleCmd := exec.Command("npm", "install")
	npmInstallExampleCmd.Dir = nodeExamplePath
	if out, err := npmInstallExampleCmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to run npm install in NodeJS example: %v\n%s\n", err, string(out))
		// We don't exit here as it's not strictly necessary for the test to pass (Pulumi will re-install in temp dir),
		// but it's good for clearing IDE warnings.
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

func validateGroupName(expectedName string) func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
	return func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
		if v, ok := stack.Outputs["groupName"]; !ok || v == nil {
			t.Fatalf("expected output groupName, got: %v", stack.Outputs)
		}

		// Verify that the group exists in NetBird API
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, localMgmtURL+"/api/groups", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+localSeededPAT)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to fetch groups: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unexpected status code %d and failed to read body: %v", resp.StatusCode, err)
			}
			t.Fatalf("unexpected status code %d: %s", resp.StatusCode, string(body))
		}

		var groups []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
			t.Fatalf("failed to decode groups: %v", err)
		}

		found := false
		for _, g := range groups {
			if g["name"] == expectedName {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("expected group %q not found in NetBird API", expectedName)
		}
	}
}

func TestGo_Minimal_LocalStack(t *testing.T) {
	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", "go", "minimal")
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
	opts.ExtraRuntimeValidation = validateGroupName("Pulumi Go Test Group")

	integration.ProgramTest(t, opts)
}

func TestNodeJS_Minimal_LocalStack(t *testing.T) {
	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", "nodejs", "minimal")
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
	opts.ExtraRuntimeValidation = validateGroupName("Pulumi NodeJS Test Group")

	integration.ProgramTest(t, opts)
}

func TestPython_Minimal_LocalStack(t *testing.T) {
	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", "python", "minimal")
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
	opts.ExtraRuntimeValidation = validateGroupName("Pulumi Python Test Group")

	integration.ProgramTest(t, opts)
}

func TestDotNet_Minimal_LocalStack(t *testing.T) {
	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", "dotnet", "minimal")
	opts.PrePrepareProject = func(proj *engine.Projinfo) error {
		absSdkPath, err := filepath.Abs(filepath.Join("..", "sdk", "dotnet", "Pulumi.Netbird.csproj"))
		if err != nil {
			return err
		}
		csprojPath := filepath.Join(proj.Root, "minimal-dotnet.csproj")
		content, err := os.ReadFile(csprojPath)
		if err != nil {
			return err
		}
		newContent := strings.Replace(string(content), "../../../sdk/dotnet/Pulumi.Netbird.csproj", absSdkPath, 1)
		return os.WriteFile(csprojPath, []byte(newContent), 0644)
	}
	opts.ExtraRuntimeValidation = validateGroupName("Pulumi DotNet Test Group")

	integration.ProgramTest(t, opts)
}

func TestJava_Minimal_LocalStack(t *testing.T) {
	opts := baseOptions(t)
	opts.Dir = filepath.Join("..", "examples", "java", "minimal")
	opts.ExtraRuntimeValidation = validateGroupName("Pulumi Java Test Group")

	integration.ProgramTest(t, opts)
}
