#!/usr/bin/env bash
# FORGE CLI Installer — one command to rule them all.
# Usage: curl -fsSL https://raw.githubusercontent.com/victorpothin/forge/main/install.sh | bash

set -e

# --- Colors ---
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

echo ""
echo "╔══════════════════════════════════════╗"
echo "║       FORGE CLI Installer            ║"
echo "║  Focused, Ordered, Restricted,       ║"
echo "║  Guided Execution                    ║"
echo "╚══════════════════════════════════════╝"
echo ""

# --- Check Go ---
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed.${NC}"
    echo ""
    echo "  Install Go first: https://go.dev/doc/install"
    echo ""
    exit 1
fi
echo -e "  ${GREEN}✓${NC} Go found: $(go version | awk '{print $3}')"

# --- Install binary ---
echo -e "\n${CYAN}⟳ Installing forge...${NC}"
go install github.com/victorpothin/forge@latest
echo -e "  ${GREEN}✓${NC} forge installed to $(go env GOPATH)/bin/forge"

# --- PATH setup ---
BIN_DIR="$(go env GOPATH)/bin"
FORGE_BIN="$BIN_DIR/forge"

if ! command -v forge &> /dev/null; then
    echo -e "\n${CYAN}⟳ Adding $BIN_DIR to PATH...${NC}"

    # Detect shell config
    SHELL_RC=""
    if [[ -n "$ZSH_VERSION" ]]; then
        SHELL_RC="$HOME/.zshrc"
    elif [[ -n "$BASH_VERSION" ]]; then
        SHELL_RC="$HOME/.bashrc"
    fi

    if [[ -z "$SHELL_RC" ]]; then
        echo -e "  ${RED}✗ Could not detect shell. Add this to your shell config:${NC}"
        echo ""
        echo "    export PATH=\$PATH:$BIN_DIR"
        echo ""
        exit 1
    fi

    # Add to shell rc if not already there
    if ! grep -q "$BIN_DIR" "$SHELL_RC" 2>/dev/null; then
        echo "" >> "$SHELL_RC"
        echo "# FORGE CLI" >> "$SHELL_RC"
        echo "export PATH=\$PATH:$BIN_DIR" >> "$SHELL_RC"
        echo -e "  ${GREEN}✓${NC} Added to $SHELL_RC"
    else
        echo -e "  ${GREEN}✓${NC} Already in $SHELL_RC"
    fi

    # Export for current session
    export PATH="$PATH:$BIN_DIR"
fi

# --- Verify ---
echo ""
if command -v forge &> /dev/null; then
    echo -e "  ${GREEN}✓${NC} forge is ready: $(forge --version)"
else
    echo -e "  ${YELLOW}⚠${NC} Restart your terminal or run:"
    echo ""
    echo "    export PATH=\$PATH:$BIN_DIR"
    echo ""
    exit 0
fi

echo ""
echo -e "${BOLD}Next steps:${NC}"
echo ""
echo "  forge init --ai qwen -y     # init a project"
echo "  forge doctor                # check health"
echo "  forge --help                # all commands"
echo ""
echo -e "${GREEN}Forge is ready. The fire is lit. 🔥${NC}"
echo ""
