package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/victorpothin/forge/internal/skillpkg"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	var addFrom, addLayer, addDir string
	var addForce bool

	skillCmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage FORGE skills",
		Long:  "Import, list, and remove skills for specific FORGE layers.",
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Import a skill into a FORGE layer",
		Long: `Import a skill from a local directory, URL, or Git repository.

Examples:
  forge skill add
  forge skill add --from ./my-skill --layer execution
  forge skill add --from github.com/user/my-skill --layer planning`,
		Run: func(cmd *cobra.Command, args []string) {
			runSkillAdd(addFrom, addLayer, addDir, addForce)
		},
	}

	addCmd.Flags().StringVar(&addFrom, "from", "", "Source: local path, URL, or GitHub repo")
	addCmd.Flags().StringVarP(&addLayer, "layer", "l", "", "Target layer")
	addCmd.Flags().StringVarP(&addDir, "dir", "d", ".", "Target project directory")
	addCmd.Flags().BoolVar(&addForce, "force", false, "Overwrite existing skill")
	skillCmd.AddCommand(addCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List skills by layer",
		Long: `Show all imported skills organized by FORGE layer.

Examples:
  forge skill list
  forge skill list --dir /path/to/project`,
		Run: func(cmd *cobra.Command, args []string) {
			runSkillList(addDir)
		},
	}

	listCmd.Flags().StringVarP(&addDir, "dir", "d", ".", "Target project directory")
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

	skillBase := detectSkillBase(target)

	if from == "" {
		from = gatherSource()
	}
	if layer == "" {
		layer = gatherLayer()
	}

	fmt.Println()
	ui.PrintHeader("Import Skill")
	ui.PrintConfig("Source", from)
	ui.PrintConfig("Layer", layer)
	ui.PrintConfig("Target", skillBase+"/"+layer+"/")
	fmt.Println()
	ui.Dim.Println("  Press Ctrl+C at any prompt to quit.")
	fmt.Println()

	if !force {
		ok := false
		ui.Ask(&survey.Confirm{Message: "Proceed?", Default: true}, &ok)
		if !ok {
			ui.Dim.Println("  Aborted.")
			return
		}
	}

	fmt.Println()
	spinner := ui.NewSpinner("Fetching skill")
	spinner.Start()

	skill, err := skillpkg.Add(from, layer, skillBase, force)
	if err != nil {
		spinner.StopWith("✗ %v", err)
		return
	}

	spinner.StopWith("✓ %s imported to %s/%s", skill.Name, layer, skill.Name)

	entries, _ := os.ReadDir(skill.Path)
	if len(entries) > 0 {
		ui.Dim.Println("\n  Contents:")
		for _, e := range entries {
			if e.IsDir() {
				ui.Dim.Printf("    %s/\n", e.Name())
			} else {
				ui.Dim.Printf("    %s\n", e.Name())
			}
		}
	}

	fmt.Println()
	ui.Dim.Println("The skill is now available in the FORGE workflow.")
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
	ui.Bold.Printf("  FORGE Skills\n")
	ui.Dim.Printf("  %s\n\n", skillBase)

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
			ui.Dim.Printf("  %-14s (empty)\n", layer)
			continue
		}

		fmt.Printf("  %-14s ", layer)
		for i, s := range layerSkills {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(s.Name)
		}
		fmt.Println()
	}

	fmt.Println()
	ui.Dim.Printf("  Total: %d skill(s)\n", totalCount)
	fmt.Println()
}

func detectSkillBase(target string) string {
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
	return filepath.Join(target, "skills/forge")
}

func gatherSource() string {
	var answer string
	ui.Ask(&survey.Input{
		Message: "Source (local path, URL, or github.com/user/repo):",
	}, &answer)
	return answer
}

func gatherLayer() string {
	var answer string
	opts := make([]string, len(skillpkg.ValidLayers))
	for i, l := range skillpkg.ValidLayers {
		opts[i] = l
	}
	ui.Ask(&survey.Select{
		Message: "Target layer:",
		Options: opts,
	}, &answer)
	return answer
}
