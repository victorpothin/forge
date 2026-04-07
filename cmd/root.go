package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// version is set via -ldflags at build time. Falls back to auto-detection.
var version = ""

func detectVersion() string {
	// If set via ldflags, use it
	if version != "" {
		return version
	}

	// Try git describe (works when building from repo)
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	if err == nil {
		tag := strings.TrimSpace(string(out))
		if strings.HasPrefix(tag, "v") {
			return strings.TrimPrefix(tag, "v")
		}
		return tag
	}

	// Fallback
	return "dev"
}

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "FORGE: Focused, Ordered, Restricted, Guided Execution",
	Long: `FORGE is a structured method for working with AI on software projects.
It ensures the AI understands your intent, respects your restrictions,
and executes only what was asked — nothing more.

AI-agnostic. Works with Qwen, Claude, Gemini, GPT, or any coding assistant.`,
	Version: detectVersion(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("forge-cli {{printf \"v%s\\n\" .Version}}")
}
