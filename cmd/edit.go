package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/victorpothin/forge/internal/templates"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	var editDir, addLayers, removeLayers string
	var editYes bool

	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Toggle FORGE layers on/off",
		Long: `Interactively manage which FORGE layers are active in a project.

Use Space to toggle each layer on/off. Press Enter when done.

Examples:
  forge edit                 # Interactive toggle
  forge edit --add testing,docs
  forge edit --remove context,documentation`,
		Run: func(cmd *cobra.Command, args []string) {
			runEdit(editDir, addLayers, removeLayers, editYes)
		},
	}

	editCmd.Flags().StringVarP(&editDir, "dir", "d", ".", "Target project directory")
	editCmd.Flags().StringVar(&addLayers, "add", "", "Comma-separated layers to enable")
	editCmd.Flags().StringVar(&removeLayers, "remove", "", "Comma-separated layers to disable")
	editCmd.Flags().BoolVarP(&editYes, "yes", "y", false, "Skip confirmation")
	rootCmd.AddCommand(editCmd)
}

func runEdit(dir, addStr, removeStr string, yes bool) {
	target, err := filepath.Abs(dir)
	if err != nil {
		ui.PrintError("Invalid directory: %v", err)
		return
	}

	cfg, err := templates.LoadForgeConfig(target)
	if err != nil {
		ui.PrintError("No FORGE config found. Run 'forge init' first: %v", err)
		return
	}

	// Non-interactive mode
	if addStr != "" || removeStr != "" {
		layers := make(map[string]bool)
		for k, v := range cfg.Layers {
			layers[k] = v
		}

		if addStr != "" {
			for _, l := range strings.Split(addStr, ",") {
				l = strings.TrimSpace(l)
				if isValidLayerName(l) {
					layers[l] = true
				}
			}
		}
		if removeStr != "" {
			for _, l := range strings.Split(removeStr, ",") {
				l = strings.TrimSpace(l)
				layers[l] = false
			}
		}

		if err := templates.UpdateForgeLayers(target, layers); err != nil {
			ui.PrintError("Failed to update layers: %v", err)
			return
		}

		fmt.Println()
		ui.PrintHeader("Layers updated")
		printLayerStatus(layers)
		fmt.Println()
		return
	}

	// Interactive mode — toggle style
	fmt.Println()
	ui.Bold.Printf("  FORGE Edit\n")
	ui.Dim.Printf("  %s\n\n", target)

	ui.PrintHeader("Current layers")
	printLayerStatus(cfg.Layers)
	fmt.Println()

	// Build options
	opts := make([]string, len(allLayers))
	def := make([]string, 0, len(allLayers))
	for i, l := range allLayers {
		opts[i] = fmt.Sprintf("%s — %s", l, layerDescriptions[l])
		if cfg.Layers[l] {
			def = append(def, opts[i])
		}
	}

	ui.Dim.Println("  Space: toggle on/off  |  Enter: confirm  |  Ctrl+C: quit")
	fmt.Println()

	var selected []string
	ui.Ask(&survey.MultiSelect{
		Message: "Active layers:",
		Options: opts,
		Default: def,
	}, &selected)

	if len(selected) == 0 {
		ui.PrintWarning("No layers selected. At least one layer should be active.")
		fmt.Println()
		return
	}

	// Build new layers map
	newLayers := make(map[string]bool)
	for _, l := range allLayers {
		newLayers[l] = false
	}
	for _, s := range selected {
		idx := strings.Index(s, " — ")
		if idx > 0 {
			newLayers[s[:idx]] = true
		} else {
			newLayers[s] = true
		}
	}

	// Show diff
	changes := diffLayers(cfg.Layers, newLayers)
	if len(changes) == 0 {
		ui.PrintSuccess("No changes needed")
		fmt.Println()
		return
	}

	if !yes {
		fmt.Println()
		ui.PrintHeader("Changes")
		for _, c := range changes {
			fmt.Println(c)
		}
		fmt.Println()

		ok := false
		ui.Ask(&survey.Confirm{Message: "Apply changes?", Default: true}, &ok)
		if !ok {
			ui.Dim.Println("  Cancelled.")
			fmt.Println()
			return
		}
	}

	if err := templates.UpdateForgeLayers(target, newLayers); err != nil {
		ui.PrintError("Failed to update layers: %v", err)
		return
	}

	fmt.Println()
	ui.PrintHeader("New configuration")
	printLayerStatus(newLayers)
	fmt.Println()
	ui.PrintSuccess("Layers updated successfully")
	fmt.Println()
}

func diffLayers(old, new map[string]bool) []string {
	var changes []string
	for _, l := range allLayers {
		wasOn := old[l]
		isOn := new[l]
		if wasOn != isOn {
			if isOn {
				changes = append(changes, fmt.Sprintf("  + %s (enabled)", l))
			} else {
				changes = append(changes, fmt.Sprintf("  - %s (disabled)", l))
			}
		}
	}
	return changes
}

func printLayerStatus(layers map[string]bool) {
	for _, l := range allLayers {
		if layers[l] {
			fmt.Printf("  ✓ %s\n", l)
		} else {
			ui.Dim.Printf("  ✗ %s\n", l)
		}
	}
}

func isValidLayerName(name string) bool {
	for _, l := range allLayers {
		if l == name {
			return true
		}
	}
	return false
}
