#!/usr/bin/env bash
# setup_upstream.sh - Prepare the upstream submodule for building.
#
# The upstream directory is a git submodule pointing at the original
# netbirdio/terraform-provider-netbird repository.  Two local
# modifications are needed that cannot be pushed there:
#
#   1. A thin "shim" package that re-exports internal/provider.New so
#      the Pulumi-Terraform bridge can import it (Go's "internal"
#      package rule prevents direct imports from outside the module).
#      The shim MUST live inside upstream/ for this reason.
#
#   2. Small test-data patches (encryption key format, service-user
#      flag, unique peer keys) required by our integration tests.
#      Without these the management server fails to start.
#
# Because .gitmodules sets `ignore = dirty`, these changes never make
# the parent repo appear dirty.

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
UPSTREAM_DIR="$SCRIPT_DIR/upstream"

# -- 1. Ensure submodule is initialised ------------------------------------
if [ ! -f "$UPSTREAM_DIR/go.mod" ]; then
    echo "Initialising upstream submodule..."
    git -C "$SCRIPT_DIR" submodule update --init upstream
fi

# -- 2. Generate the thin shim ---------------------------------------------
# A single Go file that re-exports internal/provider.New.  This is the
# minimal surface needed by the Pulumi bridge and avoids copying ~50 files.
SHIM_DIR="$UPSTREAM_DIR/shim/provider"
SHIM_FILE="$SHIM_DIR/shim.go"

if [ ! -f "$SHIM_FILE" ]; then
    echo "Generating upstream shim (shim/provider/shim.go)..."
    mkdir -p "$SHIM_DIR"
    cat > "$SHIM_FILE" <<'EOF'
// Package provider re-exports the internal provider constructor so it can be
// consumed by the Pulumi-Terraform bridge (which lives in a separate Go module
// and therefore cannot import an "internal" package directly).
package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"
	internal "github.com/netbirdio/terraform-provider-netbird/internal/provider"
)

// New returns a factory function that creates a new instance of the NetBird
// Terraform provider for the given version string.
func New(version string) func() provider.Provider {
	return internal.New(version)
}
EOF
else
    echo "Upstream shim already present."
fi

# -- 3. Apply test-data patches (idempotent) --------------------------------
echo "Applying upstream test patches..."
cd "$UPSTREAM_DIR"
if git diff --quiet -- test/management.json test/seed_database.sql 2>/dev/null; then
    # Files are unmodified - apply the patch
    git apply "$SCRIPT_DIR/upstream.patch" 2>/dev/null || true
else
    echo "  (patches already applied)"
fi

echo "Upstream ready."

