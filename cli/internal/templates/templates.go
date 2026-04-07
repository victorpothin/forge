// Package templates handles embedding and copying of FORGE templates.
package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed FORGE.md
var forgeMD []byte

//go:embed all:skills
var skillFS embed.FS

// ForgeMD exports the embedded FORGE.md content.
var ForgeMD = forgeMD

// ForgeConfig represents the .forgerc.json structure.
type ForgeConfig struct {
	AI               string            `json:"ai"`
	Model            string            `json:"model"`
	GateMode         string            `json:"gate_mode"`
	SkillsDir        string            `json:"skills_dir"`
	Layers           map[string]bool   `json:"layers"`
	LayerAI          map[string]string `json:"layer_ai"`
	ResponseLanguage string            `json:"response_language"`
}

// GenerateForgeConfig creates a .forgerc.json content string.
func GenerateForgeConfig(ai, model, gateMode string) (string, error) {
	config := ForgeConfig{
		AI:       ai,
		Model:    model,
		GateMode: gateMode,
		SkillsDir: "skills/",
		Layers: map[string]bool{
			"context":       true,
			"problem":       true,
			"locked_path":   true,
			"planning":      true,
			"execution":     true,
			"testing":       true,
			"documentation": true,
		},
		LayerAI:          map[string]string{},
		ResponseLanguage: "en",
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal config: %w", err)
	}

	return string(data) + "\n", nil
}

// SkillTargets maps AI model names to skill destination paths.
var SkillTargets = map[string]string{
	"qwen":   ".qwen/skills/forge",
	"claude": ".claude/skills/forge",
	"gemini": ".gemini/skills/forge",
	"gpt":    ".gpt/skills/forge",
	"custom": "skills/forge",
}

// AIDescriptions maps AI model names to human-readable descriptions.
var AIDescriptions = map[string]string{
	"qwen":   "Qwen Code — best for code execution & file ops",
	"claude": "Claude (Code/CLI) — best for reasoning & planning",
	"gpt":    "GPT (ChatGPT/CLI) — general purpose",
	"gemini": "Gemini (CLI) — general purpose",
	"custom": "Custom / other AI model",
}

// AIDefaults maps AI model names to default model identifiers.
var AIDefaults = map[string]string{
	"qwen":   "qwen-code",
	"claude": "claude-sonnet-4-20250514",
	"gpt":    "gpt-4o",
	"gemini": "gemini-2.5-pro",
	"custom": "",
}

// SupportedAI returns the list of supported AI models.
func SupportedAI() []string {
	return []string{"qwen", "claude", "gemini", "gpt", "custom"}
}

// SupportedGateModes returns valid gate mode values.
func SupportedGateModes() []string {
	return []string{"strict", "auto"}
}

// IsValidAI checks if the given AI model name is supported.
func IsValidAI(ai string) bool {
	_, ok := SkillTargets[ai]
	return ok
}

// IsValidGateMode checks if the given gate mode is valid.
func IsValidGateMode(mode string) bool {
	return mode == "strict" || mode == "auto"
}

// CopyForgeMD copies FORGE.md to the target directory.
// Returns true if the file was copied, false if it was skipped.
func CopyForgeMD(target string, force bool) (bool, error) {
	dest := filepath.Join(target, "FORGE.md")

	if _, err := os.Stat(dest); err == nil && !force {
		return false, nil // exists, skipped
	}

	if len(ForgeMD) == 0 {
		return false, fmt.Errorf("FORGE.md content not available")
	}

	if err := os.WriteFile(dest, ForgeMD, 0644); err != nil {
		return false, fmt.Errorf("write FORGE.md: %w", err)
	}

	return true, nil
}

// CopySkills copies all skill directories to the target.
// Returns the list of copied skill paths.
func CopySkills(target, ai string, force bool) ([]string, error) {
	destBase, ok := SkillTargets[ai]
	if !ok {
		destBase = SkillTargets["custom"]
	}

	destDir := filepath.Join(target, destBase)
	var copied []string

	entries, err := fs.ReadDir(skillFS, "skills")
	if err != nil {
		return nil, fmt.Errorf("read embedded skills: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()
		destPath := filepath.Join(destDir, skillName)

		if _, err := os.Stat(destPath); err == nil && !force {
			continue // exists, skipped
		}

		if err := copySkillDir(skillName, destPath); err != nil {
			return copied, fmt.Errorf("copy skill %s: %w", skillName, err)
		}

		copied = append(copied, filepath.Join(destBase, skillName))
	}

	return copied, nil
}

func copySkillDir(skillName, destPath string) error {
	srcPrefix := filepath.Join("skills", skillName)

	return fs.WalkDir(skillFS, srcPrefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcPrefix, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(destPath, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		data, err := skillFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded file %s: %w", path, err)
		}

		return os.WriteFile(dest, data, 0644)
	})
}

// HasForgeMD checks if FORGE.md exists in the target directory.
func HasForgeMD(target string) bool {
	_, err := os.Stat(filepath.Join(target, "FORGE.md"))
	return err == nil
}

// HasForgeConfig checks if .forgerc.json exists in the target directory.
func HasForgeConfig(target string) bool {
	_, err := os.Stat(filepath.Join(target, ".forgerc.json"))
	return err == nil
}

// HasSkills checks if skills exist in the target directory for the given AI model.
func HasSkills(target, ai string) bool {
	destBase, ok := SkillTargets[ai]
	if !ok {
		destBase = SkillTargets["custom"]
	}
	dest := filepath.Join(target, destBase)
	info, err := os.Stat(dest)
	return err == nil && info.IsDir()
}

// SkillPathFor returns the skill directory path for a given AI model.
func SkillPathFor(ai string) string {
	dest, ok := SkillTargets[ai]
	if !ok {
		return SkillTargets["custom"]
	}
	return dest
}

// SkillLayerCount returns the number of skill layers in the target.
func SkillLayerCount(target, ai string) int {
	dest := filepath.Join(target, SkillPathFor(ai))

	entries, err := os.ReadDir(dest)
	if err != nil {
		return 0
	}

	count := 0
	for _, e := range entries {
		if e.IsDir() {
			count++
		}
	}
	return count
}

// LoadForgeConfig reads and parses .forgerc.json from the target directory.
func LoadForgeConfig(target string) (*ForgeConfig, error) {
	path := filepath.Join(target, ".forgerc.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read .forgerc.json: %w", err)
	}

	var config ForgeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse .forgerc.json: %w", err)
	}

	return &config, nil
}

// IsForgeFile checks if a filename is one of the FORGE-managed files.
func IsForgeFile(name string) bool {
	name = strings.TrimPrefix(name, ".")
	name = strings.ToLower(name)

	known := map[string]bool{
		"forge.md":     true,
		"forgerc.json": true,
	}
	return known[name]
}
