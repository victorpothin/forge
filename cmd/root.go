package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.5.1"

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "FORGE: Focused, Ordered, Restricted, Guided Execution",
	Long: `FORGE is a structured method for working with AI on software projects.
It ensures the AI understands your intent, respects your restrictions,
and executes only what was asked — nothing more.

AI-agnostic. Works with Qwen, Claude, Gemini, GPT, or any coding assistant.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`forge-cli {{printf "v%s\n" .Version}}`)
}

// fail prints erro e sai com exit code 1.
func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\n❌ Error: "+format+"\n", args...)
	os.Exit(1)
}
