package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/victorpothin/forge/internal/templates"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

var allLayers = []string{
	"context", "problem", "locked-path", "planning",
	"execution", "testing", "documentation",
}

var layerDescriptions = map[string]string{
	"context":       "Extract and validate project understanding",
	"problem":       "Identify and prioritize what needs solving",
	"locked-path":   "Register restrictions and off-limits approaches",
	"planning":      "Decompose problems into ordered tasks",
	"execution":     "Task-by-task implementation, strictly scoped",
	"testing":       "Validation built into each task delivery",
	"documentation": "Generate docs from what was actually built",
}

func init() {
	var initDir, ai, model, gateMode, layersStr string
	var yes, force bool

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize FORGE method in a project",
		Long: `Initialize the FORGE method in the current (or specified) project.

Examples:
  forge init                           # Interactive wizard
  forge init --ai qwen -y              # Non-interactive, all layers
  forge init --ai claude --layers context,problem,execution -y`,
		Run: func(cmd *cobra.Command, args []string) {
			var layers []string
			if layersStr != "" {
				layers = strings.Split(layersStr, ",")
				for i := range layers {
					layers[i] = strings.TrimSpace(layers[i])
				}
			}
			runInit(initDir, ai, model, gateMode, layers, yes, force)
		},
	}

	initCmd.Flags().StringVarP(&initDir, "dir", "d", ".", "Target directory")
	initCmd.Flags().StringVarP(&ai, "ai", "a", "", "AI model: qwen, claude, gpt, gemini, custom")
	initCmd.Flags().StringVarP(&model, "model", "m", "", "Model name for documentation")
	initCmd.Flags().StringVarP(&gateMode, "gate-mode", "g", "", "Gate mode: strict (default) or auto")
	initCmd.Flags().StringVarP(&layersStr, "layers", "L", "", "Comma-separated layers (default: all 7)")
	initCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts")
	initCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing FORGE files")

	rootCmd.AddCommand(initCmd)
}

func runInit(dir, ai, model, gateMode string, selLayers []string, yes, force bool) {
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

	// --- Gather config ---
	hasAI := ai != ""
	ai = gatherAI(ai)
	gateMode = gatherGateMode(gateMode, hasAI)
	model = gatherModel(model, ai, hasAI)

	// Layer selection
	if len(selLayers) == 0 {
		selLayers = gatherLayers()
	}

	// --- Show plan ---
	ui.PrintHeader("Configuration")
	ui.PrintConfig("AI model", fmt.Sprintf("%s (%s)", ai, model))
	ui.PrintConfig("Gate mode", gateMode)
	ui.PrintConfig("Target", target)
	fmt.Println()

	ui.PrintHeader("Layers")
	for _, l := range allLayers {
		if contains(selLayers, l) {
			fmt.Printf("  ✓ %s\n", l)
		}
	}
	fmt.Println()

	skillPath := templates.SkillPathFor(ai)
	ui.PrintHeader("Files to create")

	filesToCreate := []string{}
	if !templates.HasForgeMD(target) || force {
		label := "FORGE.md"
		if templates.HasForgeMD(target) {
			label = "FORGE.md (overwrite)"
		}
		filesToCreate = append(filesToCreate, label)
	}
	forgeConfigPath := templates.ForgeConfigPath(target, ai)
	filesToCreate = append(filesToCreate, "forgerc.json (inside "+templates.AIDir(ai)+"/)")
	if !templates.HasSkills(target, ai) || force {
		label := skillPath
		if templates.HasSkills(target, ai) {
			label = skillPath + " (overwrite)"
		}
		filesToCreate = append(filesToCreate, label)
	}
	ui.PrintFileListPending(filesToCreate)
	fmt.Println()
	ui.Dim.Println("  Press Ctrl+C at any prompt to quit.")
	fmt.Println()

	if !yes && !force {
		ok := false
		ui.Ask(&survey.Confirm{Message: "Proceed?", Default: true}, &ok)
		if !ok {
			fmt.Println()
			ui.Dim.Println("  Aborted.")
			fmt.Println()
			return
		}
	}

	// --- Execute ---
	fmt.Println()

	spinner := ui.NewSpinner("Copying FORGE.md")
	spinner.Start()
	time.Sleep(300 * time.Millisecond)
	copied, err := templates.CopyForgeMD(target, force)
	if err != nil {
		spinner.StopWith("✗ FORGE.md: %v", err)
	} else if copied {
		spinner.StopWith("✓ FORGE.md")
	} else {
		spinner.StopWith("– FORGE.md (skipped)")
	}

	spinner2 := ui.NewSpinner("Generating .forgerc.json")
	spinner2.Start()
	time.Sleep(200 * time.Millisecond)

	// Build layers map from selection
	layersMap := make(map[string]bool)
	for _, l := range allLayers {
		layersMap[l] = contains(selLayers, l)
	}

	configStr, err := templates.GenerateForgeConfigWithLayers(ai, model, gateMode, layersMap)
	if err != nil {
		spinner2.StopWith("✗ .forgerc.json: %v", err)
	} else {
		if err := os.WriteFile(forgeConfigPath, []byte(configStr), 0644); err != nil {
			spinner2.StopWith("✗ .forgerc.json: %v", err)
		} else {
			spinner2.StopWith("✓ .forgerc.json")
		}
	}

	spinner3 := ui.NewSpinner("Copying skill files")
	spinner3.Start()
	time.Sleep(200 * time.Millisecond)
	skills, err := templates.CopySkills(target, ai, force)
	if err != nil {
		spinner3.StopWith("✗ skills: %v", err)
	} else {
		spinner3.StopWith("✓ %d skill layers copied", len(skills))
		for _, s := range skills {
			ui.Dim.Printf("    %s\n", s)
		}
	}

	// --- Summary ---
	fmt.Println()
	ui.PrintSuccessBox(target)
	fmt.Println()
	ui.Dim.Println("Next steps:")
	fmt.Println("  1. Open FORGE.md in your project")
	fmt.Printf("  2. Start a session with your AI (%s)\n", ai)
	ui.Dim.Println("  3. Reference FORGE.md and begin Layer 1")
	fmt.Println()
	ui.Dim.Println("  Run 'forge edit' to add/remove layers later")
	fmt.Println()
}

func gatherLayers() []string {
	opts := make([]string, len(allLayers))
	for i, l := range allLayers {
		opts[i] = fmt.Sprintf("%s — %s", l, layerDescriptions[l])
	}

	var selected []string
	ui.Ask(&survey.MultiSelect{
		Message: "Which FORGE layers do you want?",
		Options: opts,
		Default: opts, // all checked
	}, &selected)

	// Extract layer names
	var layers []string
	for _, s := range selected {
		// Take the part before " — "
		idx := strings.Index(s, " — ")
		if idx > 0 {
			layers = append(layers, s[:idx])
		} else {
			layers = append(layers, s)
		}
	}
	return layers
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
	ui.Ask(&survey.Select{
		Message: "Which AI model are you using?",
		Options: buildAIOptions(),
	}, &answer)
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
	ui.Ask(&survey.Select{
		Message: "Gate mode?",
		Options: []string{"strict", "auto"},
		Default: "strict",
	}, &answer)
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
	ui.Ask(&survey.Input{
		Message: "Model name:",
		Default: defaultModel,
	}, &answer)
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
