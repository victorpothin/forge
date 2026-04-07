package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/victorpothin/forge/internal/templates"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	var doctorDir string

	doctorCmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check if FORGE is properly set up in a project",
		Long: `Examine the current project and report on the health of its FORGE setup.

Checks:
  - FORGE.md exists and is valid
  - .forgerc.json exists and is parseable
  - AI model configuration is valid
  - Gate mode configuration is valid
  - Skill files are present in the correct location

Examples:
  forge doctor            # Check current directory
  forge doctor --dir ../my-project`,
		Run: func(cmd *cobra.Command, args []string) {
			runDoctor(doctorDir)
		},
	}

	doctorCmd.Flags().StringVarP(&doctorDir, "dir", "d", ".", "Target directory")
	rootCmd.AddCommand(doctorCmd)
}

type checkResult struct {
	name    string
	ok      bool
	status  string
	details string
}

var (
	redBold   = color.New(color.FgRed, color.Bold)
	greenBold = color.New(color.FgGreen, color.Bold)
)

func runDoctor(dir string) {
	target, err := filepath.Abs(dir)
	if err != nil {
		ui.PrintError("Invalid directory: %v", err)
		return
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		ui.PrintError("Directory does not exist: %s", target)
		return
	}

	// Animated check
	spinner := ui.NewSpinner("Scanning project...")
	spinner.Start()

	fmt.Printf("\n  FORGE Doctor — Checking project setup\n")
	fmt.Printf("  Target: %s\n\n", target)

	var checks []checkResult

	// Check FORGE.md
	checks = append(checks, checkFile(
		filepath.Join(target, "FORGE.md"),
		"FORGE.md",
		"Master template",
	))

	// Check .forgerc.json
	checks = append(checks, checkFile(
		filepath.Join(target, ".forgerc.json"),
		".forgerc.json",
		"CLI configuration",
	))

	// Validate .forgerc.json
	forgeCfg, err := templates.LoadForgeConfig(target)
	if err != nil {
		checks = append(checks, checkResult{
			name:    ".forgerc.json",
			ok:      false,
			status:  "INVALID",
			details: err.Error(),
		})
	} else {
		var issues []string
		if !templates.IsValidAI(forgeCfg.AI) {
			issues = append(issues, fmt.Sprintf("invalid ai: %s", forgeCfg.AI))
		}
		if !templates.IsValidGateMode(forgeCfg.GateMode) {
			issues = append(issues, fmt.Sprintf("invalid gate mode: %s", forgeCfg.GateMode))
		}

		if len(issues) > 0 {
			checks = append(checks, checkResult{
				name:    ".forgerc.json",
				ok:      false,
				status:  "INVALID",
				details: strings.Join(issues, ", "),
			})
		} else {
			checks = append(checks, checkResult{
				name:    ".forgerc.json",
				ok:      true,
				status:  "OK",
				details: fmt.Sprintf("ai=%s, gate_mode=%s", forgeCfg.AI, forgeCfg.GateMode),
			})
		}
	}

	// Check skills
	if forgeCfg != nil {
		skillDir := templates.SkillPathFor(forgeCfg.AI)
		fullPath := filepath.Join(target, skillDir)
		info, statErr := os.Stat(fullPath)

		if statErr == nil && info.IsDir() {
			layers := templates.SkillLayerCount(target, forgeCfg.AI)
			checks = append(checks, checkResult{
				name:    "skills",
				ok:      true,
				status:  "OK",
				details: fmt.Sprintf("%d layers in %s", layers, skillDir),
			})
		} else {
			checks = append(checks, checkResult{
				name:    "skills",
				ok:      false,
				status:  "MISSING",
				details: fmt.Sprintf("expected in %s", skillDir),
			})
		}
	} else {
		checks = append(checks, checkResult{
			name:    "skills",
			ok:      false,
			status:  "UNKNOWN",
			details: "cannot check skills without .forgerc.json",
		})
	}

	spinner.Stop()

	// Print table
	printTable(checks)
	fmt.Println()

	// Overall status
	allOK := true
	for _, c := range checks {
		if !c.ok {
			allOK = false
			break
		}
	}

	if allOK {
		greenBold.Println("  ✓ All checks passed.")
	} else {
		ui.PrintWarning("Some checks failed. Run 'forge init' to set up missing components.")
	}
	fmt.Println()
}

func checkFile(filePath, name, description string) checkResult {
	if _, err := os.Stat(filePath); err == nil {
		return checkResult{
			name:    name,
			ok:      true,
			status:  "OK",
			details: description,
		}
	}
	return checkResult{
		name:    name,
		ok:      false,
		status:  "MISSING",
		details: description,
	}
}

func printTable(checks []checkResult) {
	nameW := 16

	border := strings.Repeat("─", 70)
	fmt.Println("  " + border)
	fmt.Printf("  %-*s │ %-14s │ %s\n", nameW, "Check", "Status", "Details")
	fmt.Println("  " + border)

	for _, c := range checks {
		statusStr := c.status
		if c.ok {
			statusStr = greenBold.Sprintf("✓ %s", c.status)
		} else {
			statusStr = redBold.Sprintf("✗ %s", c.status)
		}
		fmt.Printf("  %-*s │ %-20s │ %s\n", nameW, c.name, statusStr, c.details)
	}
	fmt.Println("  " + border)
}
