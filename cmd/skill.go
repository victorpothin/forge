package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/victorpothin/forge/internal/skillpkg"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

var (
	skillCyan  = color.New(color.FgCyan, color.Bold)
	skillWhite = color.New(color.FgWhite)
)

func init() {
	skillCmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage FORGE skills",
		Long:  "Import, list, and remove skills for specific FORGE layers.",
	}

	// forge skill add
	var addFrom, addLayer, addDir string
	var addForce bool
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Import a skill into a FORGE layer",
		Long: `Import a skill from a local directory, URL, or Git repository.

Examples:
  # Interactive wizard
  forge skill add

  # From a local directory
  forge skill add --from ./my-skill --layer execution

  # From a GitHub repo
  forge skill add --from github.com/user/my-skill --layer planning

  # From a URL (tarball)
  forge skill add --from https://example.com/skill.tar.gz --layer testing

  # Force overwrite existing skill
  forge skill add --from ./my-skill --layer execution --force`,
		Run: func(cmd *cobra.Command, args []string) {
			runSkillAdd(addFrom, addLayer, addDir, addForce)
		},
	}

	addCmd.Flags().StringVar(&addFrom, "from", "", "Source: local path, URL, or GitHub repo")
	addCmd.Flags().StringVarP(&addLayer, "layer", "l", "", "Target layer: context, problem, locked-path, planning, execution, testing, docs")
	addCmd.Flags().StringVarP(&addDir, "dir", "d", ".", "Target project directory")
	addCmd.Flags().BoolVar(&addForce, "force", false, "Overwrite existing skill")
	skillCmd.AddCommand(addCmd)

	// forge skill list
	var listDir string
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List skills by layer",
		Long: `Show all imported skills organized by FORGE layer.

Examples:
  forge skill list
  forge skill list --dir /path/to/project`,
		Run: func(cmd *cobra.Command, args []string) {
			runSkillList(listDir)
		},
	}

	listCmd.Flags().StringVarP(&listDir, "dir", "d", ".", "Target project directory")
	skillCmd.AddCommand(listCmd)

	rootCmd.AddCommand(skillCmd)
}

func runSkillAdd(from, layer, dir string, force bool) {
	target, err := filepath.Abs(dir)
	if err != nil {
		ui.PrintError("Invalid directory: %v", err)
		return
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		ui.PrintError("Directory does not exist: %s", target)
		return
	}

	// Determine skill base path
	skillBase := detectSkillBase(target)

	// Gather source
	if from == "" {
		from = gatherSource()
	}

	// Gather layer
	if layer == "" {
		layer = gatherLayer()
	}

	// Confirm
	fmt.Println()
	ui.PrintHeader("Import Skill")
	ui.PrintConfig("Source", from)
	ui.PrintConfig("Layer", layer)
	ui.PrintConfig("Target", skillBase+"/"+layer+"/")
	fmt.Println()

	if !force {
		ok := false
		survey.AskOne(&survey.Confirm{Message: "Proceed?", Default: true}, &ok)
		if !ok {
			ui.Dim.Println("  Aborted.")
			return
		}
	}

	// Execute
	fmt.Println()
	spinner := ui.NewSpinner("Fetching skill...")
	spinner.Start()

	skill, err := skillpkg.Add(from, layer, skillBase, force)
	if err != nil {
		spinner.StopWith("✗ %s", color.New(color.FgRed).Sprintf("Failed: %v", err))
		return
	}

	spinner.StopWith("✓ Skill %s imported to %s",
		color.New(color.FgGreen, color.Bold).Sprint(skill.Name),
		ui.Dim.Sprintf("%s/%s/", layer, skill.Name))

	// Show what's inside
	entries, _ := os.ReadDir(skill.Path)
	if len(entries) > 0 {
		ui.Dim.Println("\n  Contents:")
		for _, e := range entries {
			if e.IsDir() {
				ui.Dim.Printf("    📁 %s/\n", e.Name())
			} else {
				ui.Dim.Printf("    📄 %s\n", e.Name())
			}
		}
	}

	fmt.Println()
	ui.Dim.Println("The skill is now available for use in the FORGE planning and execution layers.")
	fmt.Println()
}

func runSkillList(dir string) {
	target, err := filepath.Abs(dir)
	if err != nil {
		ui.PrintError("Invalid directory: %v", err)
		return
	}

	skillBase := detectSkillBase(target)

	fmt.Println()
	skillCyan.Printf("  FORGE Skills — %s\n\n", skillBase)

	skills, err := skillpkg.List(skillBase)
	if err != nil {
		ui.PrintError("Failed to list skills: %v", err)
		return
	}

	totalCount := 0
	for _, layer := range skillpkg.ValidLayers {
		layerSkills := skills[layer]
		count := len(layerSkills)
		totalCount += count

		if count == 0 {
			ui.Dim.Printf("  %-14s  (empty)\n", layer+":")
			continue
		}

		ui.GreenB.Printf("  %-14s  ", layer+":")
		for i, s := range layerSkills {
			if i > 0 {
				ui.Dim.Print(", ")
			}
			skillWhite.Print(s.Name)
		}
		fmt.Println()
	}

	fmt.Println()
	ui.Dim.Printf("  Total: %d skill(s)\n", totalCount)
	fmt.Println()
}

// -- Helpers -----------------------------------------------------------------

func detectSkillBase(target string) string {
	// Default based on .forgerc.json or common patterns
	paths := []string{
		".qwen/skills/forge",
		".claude/skills/forge",
		".gemini/skills/forge",
		".gpt/skills/forge",
		"skills/forge",
	}

	for _, p := range paths {
		if _, err := os.Stat(filepath.Join(target, p)); err == nil {
			return filepath.Join(target, p)
		}
	}

	// Fallback: first existing dir pattern
	return filepath.Join(target, "skills/forge")
}

func gatherSource() string {
	var answer string
	survey.AskOne(&survey.Input{
		Message: "Skill source (local path, URL, or github.com/user/repo):",
	}, &answer)
	return answer
}

func gatherLayer() string {
	var answer string
	opts := make([]string, len(skillpkg.ValidLayers))
	for i, l := range skillpkg.ValidLayers {
		opts[i] = l
	}
	survey.AskOne(&survey.Select{
		Message: "Target layer:",
		Options: opts,
	}, &answer)
	return answer
}
