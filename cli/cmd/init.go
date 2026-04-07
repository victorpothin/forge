package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/forge-cli/forge/internal/templates"
	"github.com/spf13/cobra"
)

func init() {
	var initDir, ai, model, gateMode string
	var yes, force bool

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize FORGE method in a project",
		Long: `Initialize the FORGE method in the current (or specified) project.

This command creates:
  - FORGE.md          The master template with 7 execution layers
  - .forgerc.json     Configuration file for AI model and gate behavior
  - skills/           All skill files for the selected AI model

Examples:
  forge init                    # Interactive wizard
  forge init --ai qwen          # Non-interactive with specific AI
  forge init --ai claude -y     # Skip confirmation
  forge init --ai gpt --force   # Overwrite existing files`,
		Run: func(cmd *cobra.Command, args []string) {
			runInit(initDir, ai, model, gateMode, yes, force)
		},
	}

	initCmd.Flags().StringVarP(&initDir, "dir", "d", ".", "Target directory")
	initCmd.Flags().StringVarP(&ai, "ai", "a", "", "AI model: qwen, claude, gpt, gemini, custom")
	initCmd.Flags().StringVarP(&model, "model", "m", "", "Model name for documentation (e.g. 'qwen-code')")
	initCmd.Flags().StringVarP(&gateMode, "gate-mode", "g", "", "Gate mode: strict (default) or auto")
	initCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts")
	initCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing FORGE files")

	rootCmd.AddCommand(initCmd)
}

func runInit(dir, ai, model, gateMode string, yes, force bool) {
	banner()

	target, err := filepath.Abs(dir)
	if err != nil {
		fail("Invalid directory: %v", err)
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		fail("Directory does not exist: %s", target)
	}

	// --- Gather config (smart defaults when --ai is set) ---
	hasAI := ai != ""
	ai = gatherAI(ai)
	gateMode = gatherGateMode(gateMode, hasAI)
	model = gatherModel(model, ai, hasAI)

	// --- Show plan ---
	fmt.Println("📋 Configuration:")
	fmt.Printf("  AI model:   %s (%s)\n", ai, model)
	fmt.Printf("  Gate mode:  %s\n", gateMode)
	fmt.Printf("  Target:     %s\n", target)
	fmt.Println()

	skillPath := templates.SkillPathFor(ai)
	fmt.Println("📁 Files to create/overwrite:")

	if !templates.HasForgeMD(target) || force {
		label := "FORGE.md"
		if templates.HasForgeMD(target) {
			label = "FORGE.md (overwrite)"
		}
		fmt.Printf("  ✅ %s\n", label)
	}

	fmt.Println("  ✅ .forgerc.json")

	if !templates.HasSkills(target, ai) || force {
		label := skillPath
		if templates.HasSkills(target, ai) {
			label = skillPath + " (overwrite)"
		}
		fmt.Printf("  ✅ %s\n", label)
	}
	fmt.Println()

	if !yes && !force {
		ok := false
		prompt := &survey.Confirm{
			Message: "Proceed?",
			Default: true,
		}
		survey.AskOne(prompt, &ok)
		if !ok {
			fmt.Println("\n⚠️  Aborted.")
			return
		}
	}

	// --- Execute ---
	fmt.Println()

	// 1. Copy FORGE.md
	copied, err := templates.CopyForgeMD(target, force)
	if err != nil {
		fmt.Printf("  ⚠️  %v\n", err)
	} else if copied {
		fmt.Println("  ✓ FORGE.md")
	}

	// 2. Generate .forgerc.json
	configStr, err := templates.GenerateForgeConfig(ai, model, gateMode)
	if err != nil {
		fail("Failed to generate config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, ".forgerc.json"), []byte(configStr), 0644); err != nil {
		fail("Failed to write .forgerc.json: %v", err)
	}
	fmt.Println("  ✓ .forgerc.json")

	// 3. Copy skills
	skills, err := templates.CopySkills(target, ai, force)
	if err != nil {
		fmt.Printf("  ⚠️  %v\n", err)
	} else {
		for _, s := range skills {
			fmt.Printf("  ✓ %s\n", s)
		}
	}

	// --- Summary ---
	fmt.Println()
	fmt.Println("┌──────────────────────────────────────┐")
	fmt.Println("│   ✅ FORGE initialized               │")
	fmt.Printf("│   in %-36s│\n", truncatePath(target, 36))
	fmt.Println("└──────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Open FORGE.md in your project")
	fmt.Printf("  2. Start a session with your AI (%s)\n", ai)
	fmt.Println("  3. Reference FORGE.md and begin Layer 1")
	fmt.Println()
	fmt.Println("  Run 'forge doctor' to verify setup")
	fmt.Println()
}

func gatherAI(ai string) string {
	if ai != "" {
		if !templates.IsValidAI(ai) {
			fail("Invalid AI model: %s\nSupported: %s", ai, joinAIList())
		}
		return ai
	}

	var answer string
	prompt := &survey.Select{
		Message: "Which AI model are you using?",
		Options: buildAIOptions(),
	}
	survey.AskOne(prompt, &answer)
	// Extract the AI name (before " — ")
	for _, valid := range templates.SupportedAI() {
		if answer == valid {
			return valid
		}
	}
	// Fallback: split on " — "
	for _, valid := range templates.SupportedAI() {
		if len(answer) > len(valid) && answer[:len(valid)] == valid {
			return valid
		}
	}
	return answer
}

func gatherGateMode(gateMode string, hasAIFlag bool) string {
	if gateMode != "" {
		if !templates.IsValidGateMode(gateMode) {
			fail("Invalid gate mode: %s\nSupported: strict, auto", gateMode)
		}
		return gateMode
	}
	if hasAIFlag {
		return "strict"
	}

	var answer string
	prompt := &survey.Select{
		Message: "Gate mode?",
		Options: []string{
			"strict",
			"auto",
		},
		Default: "strict",
		Help:    "strict = AI stops at each gate; auto = AI proceeds automatically",
	}
	survey.AskOne(prompt, &answer)
	return answer
}

func gatherModel(model, ai string, hasAIFlag bool) string {
	if model != "" {
		return model
	}
	if hasAIFlag {
		defaultModel := templates.AIDefaults[ai]
		if defaultModel != "" {
			return defaultModel
		}
		return fmt.Sprintf("%s-default", ai)
	}

	defaultModel := templates.AIDefaults[ai]
	if defaultModel == "" {
		defaultModel = fmt.Sprintf("%s-default", ai)
	}

	var answer string
	prompt := &survey.Input{
		Message: "Model name:",
		Default: defaultModel,
	}
	survey.AskOne(prompt, &answer)
	if answer == "" {
		return defaultModel
	}
	return answer
}

func buildAIOptions() []string {
	var opts []string
	for _, ai := range templates.SupportedAI() {
		opts = append(opts, fmt.Sprintf("%s — %s", ai, templates.AIDescriptions[ai]))
	}
	return opts
}

func joinAIList() string {
	result := ""
	for i, ai := range templates.SupportedAI() {
		if i > 0 {
			result += ", "
		}
		result += ai
	}
	return result
}

func truncatePath(p string, maxLen int) string {
	if len(p) <= maxLen {
		return p
	}
	// Show last maxLen-3 chars with ...
	return "..." + p[len(p)-(maxLen-3):]
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\n❌ Error: "+format+"\n", args...)
	os.Exit(1)
}
