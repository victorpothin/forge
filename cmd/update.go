package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

const githubLatestURL = "https://api.github.com/repos/victorpothin/forge/releases/latest"

type githubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
}

func init() {
	var updateYes bool

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update FORGE CLI to the latest version",
		Long: `Check for updates and upgrade the FORGE CLI binary.

Examples:
  forge update          # Check and prompt for update
  forge update -y       # Auto-update if available
  forge update --check  # Only check, don't download`,
		Run: func(cmd *cobra.Command, args []string) {
			runUpdate(updateYes)
		},
	}

	updateCmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "Skip confirmation, auto-update")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(yes bool) {
	ui.Banner(version)

	spinner := ui.NewSpinner("Checking for updates...")
	spinner.Start()

	latest, err := fetchLatestRelease()
	if err != nil {
		spinner.StopWith("✗ %s", color.New(color.FgRed).Sprintf("Failed to check for updates: %v", err))
		return
	}

	latestVer := strings.TrimPrefix(latest.TagName, "v")
	currentVer := version

	spinner.Stop()

	if versionGreaterOrEqual(currentVer, latestVer) {
		ui.PrintBox(
			color.New(color.FgGreen, color.Bold).Sprint("  ✓ Up to date"),
			fmt.Sprintf("  You're running the latest version (v%s)", currentVer),
		)
		fmt.Println()
		return
	}

	fmt.Printf("  New version available: %s → %s\n",
		color.New(color.FgYellow).Sprint("v"+currentVer),
		color.New(color.FgGreen, color.Bold).Sprint("v"+latestVer))
	fmt.Printf("  Released: %s\n", formatTime(latest.PublishedAt))
	fmt.Println()

	if !yes {
		ok := false
		survey.AskOne(&survey.Confirm{
			Message: "Update to v" + latestVer + "?",
			Default: true,
		}, &ok)
		if !ok {
			ui.Dim.Println("  Update cancelled.")
			return
		}
	}

	// Perform update
	spinner2 := ui.NewSpinner("Updating forge CLI...")
	spinner2.Start()

	cmd := exec.Command("go", "install", "github.com/victorpothin/forge@latest")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	err = cmd.Run()
	if err != nil {
		spinner2.StopWith("✗ %s", color.New(color.FgRed).Sprintf("Update failed: %v", err))
		return
	}

	// Verify new version
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = os.Getenv("HOME") + "/go"
	}
	binPath := gopath + "/bin/forge"

	checkCmd := exec.Command(binPath, "--version")
	checkOut, err := checkCmd.CombinedOutput()
	if err != nil {
		spinner2.StopWith("✓ Updated to v%s", color.New(color.FgGreen).Sprint(latestVer))
	} else {
		newVer := strings.TrimSpace(strings.TrimPrefix(string(checkOut), "forge-cli "))
		spinner2.StopWith("✓ Updated to %s", color.New(color.FgGreen, color.Bold).Sprint(newVer))
	}

	fmt.Println()
	ui.Dim.Println("The CLI has been updated. Restart your terminal if forge was already loaded.")
	fmt.Println()
}

func fetchLatestRelease() (*githubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", githubLatestURL, nil)
	if err != nil {
		return nil, err
	}

	// GitHub API requires User-Agent
	req.Header.Set("User-Agent", "forge-cli/"+version)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &release, nil
}

// versionGreaterOrEqual compares semver-ish version strings.
// Handles "v0.1.0", "0.1.0", etc.
func versionGreaterOrEqual(current, latest string) bool {
	if current == latest {
		return true
	}
	// Simple comparison: if current starts with latest, treat as equal
	if strings.HasPrefix(current, latest) || strings.HasPrefix(latest, current) {
		return current == latest
	}
	// Fallback: string comparison (works for most semver cases)
	return current > latest
}

func formatTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.Format("Jan 2, 2006")
}
