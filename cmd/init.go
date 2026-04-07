package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/victorpothin/forge/internal/templates"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

var (
	green = color.New(color.FgGreen)
	dim   = color.New(color.Faint)
	red   = color.New(color.FgRed, color.Bold)
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
	ui.Banner(version)

	target, err := filepath.Abs(dir)
	if err != nil {
		ui.PrintError("Invalid directory: %v", err)
		return
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		ui.PrintError("Directory does not exist: %s", target)
		return
	}

	// --- Gather config (smart defaults when --ai is set) ---
	hasAI := ai != ""
	ai = gatherAI(ai)
	gateMode = gatherGateMode(gateMode, hasAI)
	model = gatherModel(model, ai, hasAI)

	// --- Show plan ---
	ui.PrintHeader("Configuration")
	ui.PrintConfig("AI model", fmt.Sprintf("%s (%s)", ai, model))
	ui.PrintConfig("Gate mode", gateMode)
	ui.PrintConfig("Target", target)
	fmt.Println()

	skillPath := templates.SkillPathFor(ai)
	ui.PrintHeader("Files to create/overwrite")

	filesToCreate := []string{}
	if !templates.HasForgeMD(target) || force {
		label := "FORGE.md"
		if templates.HasForgeMD(target) {
			label = "FORGE.md (overwrite)"
		}
		filesToCreate = append(filesToCreate, label)
	}
	filesToCreate = append(filesToCreate, ".forgerc.json")
	if !templates.HasSkills(target, ai) || force {
		label := skillPath
		if templates.HasSkills(target, ai) {
			label = skillPath + " (overwrite)"
		}
		filesToCreate = append(filesToCreate, label)
	}
	ui.PrintFileListPending(filesToCreate)
	fmt.Println()

	if !yes && !force {
		ok := false
		prompt := &survey.Confirm{
			Message: "Proceed?",
			Default: true,
		}
		survey.AskOne(prompt, &ok)
		if !ok {
			fmt.Println()
			dim.Println("  Aborted.")
			fmt.Println()
			return
		}
	}

	// --- Execute with animations ---
	fmt.Println()

	// 1. Copy FORGE.md
	spinner := ui.NewSpinner("Copying FORGE.md...")
	spinner.Start()
	time.Sleep(300 * time.Millisecond)
	copied, err := templates.CopyForgeMD(target, force)
	if err != nil {
		spinner.StopWith("%s FORGE.md: %v", red.Sprintf("✗"), err)
	} else if copied {
		spinner.StopWith("%s FORGE.md", green.Sprintf("✓"))
	} else {
		spinner.StopWith("%s FORGE.md (exists)", dim.Sprint("–"))
	}

	// 2. Generate .forgerc.json
	spinner2 := ui.NewSpinner("Generating .forgerc.json...")
	spinner2.Start()
	time.Sleep(200 * time.Millisecond)
	configStr, err := templates.GenerateForgeConfig(ai, model, gateMode)
	if err != nil {
		spinner2.StopWith("%s .forgerc.json: %v", red.Sprintf("✗"), err)
	} else {
		if err := os.WriteFile(filepath.Join(target, ".forgerc.json"), []byte(configStr), 0644); err != nil {
			spinner2.StopWith("%s .forgerc.json: %v", red.Sprintf("✗"), err)
		} else {
			spinner2.StopWith("%s .forgerc.json", green.Sprintf("✓"))
		}
	}

	// 3. Copy skills with progress
	spinner3 := ui.NewSpinner("Copying skill files...")
	spinner3.Start()
	time.Sleep(200 * time.Millisecond)
	skills, err := templates.CopySkills(target, ai, force)
	if err != nil {
		spinner3.StopWith("%s skills: %v", red.Sprintf("✗"), err)
	} else {
		totalLayers := len(templates.SupportedAI())
		_ = totalLayers
		spinner3.StopWith("%s %d skill layers copied", green.Sprintf("✓"), len(skills))
		for _, s := range skills {
			dim.Printf("    %s\n", s)
		}
	}

	// --- Summary ---
	fmt.Println()
	ui.PrintBox(
		red.Sprint("  FORGE initialized"),
		fmt.Sprintf("  in %s", target),
	)
	fmt.Println()
	dim.Println("Next steps:")
	green.Println("  1. Open FORGE.md in your project")
	fmt.Printf("  2. Start a session with your AI (%s)\n", ai)
	dim.Println("  3. Reference FORGE.md and begin Layer 1")
	fmt.Println()
	dim.Println("  Run 'forge doctor' to verify setup")
	fmt.Println()
}

func gatherAI(ai string) string {
	if ai != "" {
		if !templates.IsValidAI(ai) {
			ui.PrintError("Invalid AI model: %s\nSupported: %s", ai, joinAIList())
			os.Exit(1)
		}
		return ai
	}

	var answer string
	prompt := &survey.Select{
		Message: "Which AI model are you using?",
		Options: buildAIOptions(),
	}
	survey.AskOne(prompt, &answer)
	for _, valid := range templates.SupportedAI() {
		if answer == valid {
			return valid
		}
	}
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
			ui.PrintError("Invalid gate mode: %s\nSupported: strict, auto", gateMode)
			os.Exit(1)
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
