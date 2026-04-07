// Package skillpkg handles importing and managing external skills.
package skillpkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidLayers lists the FORGE skill layers.
var ValidLayers = []string{
	"context", "problem", "locked-path", "planning",
	"execution", "testing", "docs",
}

// Skill represents an imported skill package.
type Skill struct {
	Name  string // e.g. "my-react-skill"
	Layer string // e.g. "execution"
	Path  string // full path on disk
}

// SourceType indicates where the skill is coming from.
type SourceType string

const (
	SourceLocal SourceType = "local"
	SourceURL   SourceType = "url"
	SourceGit   SourceType = "git"
)

// Add imports a skill from the given source into the specified layer.
func Add(source, layer, targetDir string, force bool) (*Skill, error) {
	sourceType := detectSourceType(source)

	tmpDir, err := os.MkdirTemp("", "forge-skill-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	skillDir := tmpDir

	switch sourceType {
	case SourceLocal:
		abs, err := filepath.Abs(source)
		if err != nil {
			return nil, fmt.Errorf("resolve local path: %w", err)
		}
		if _, err := os.Stat(abs); err != nil {
			return nil, fmt.Errorf("local source not found: %s", abs)
		}
		skillDir = abs

	case SourceURL:
		skillDir, err = downloadURL(source, tmpDir)
		if err != nil {
			return nil, fmt.Errorf("download: %w", err)
		}

	case SourceGit:
		skillDir, err = cloneGit(source, tmpDir)
		if err != nil {
			return nil, fmt.Errorf("git clone: %w", err)
		}
	}

	// Normalize: if the extracted dir has a single top-level dir, use that
	skillDir = findSkillRoot(skillDir)

	// Validate
	name := filepath.Base(skillDir)
	if err := validate(skillDir, name); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Validate layer
	if !isValidLayer(layer) {
		return nil, fmt.Errorf("invalid layer: %q\nValid layers: %s", layer, strings.Join(ValidLayers, ", "))
	}

	// Copy to destination
	destDir := filepath.Join(targetDir, layer, name)
	if _, err := os.Stat(destDir); err == nil && !force {
		return nil, fmt.Errorf("skill already exists at %s (use --force to overwrite)", destDir)
	}
	if _, err := os.Stat(destDir); err == nil && force {
		os.RemoveAll(destDir)
	}

	if err := copyDir(skillDir, destDir); err != nil {
		return nil, fmt.Errorf("copy skill: %w", err)
	}

	return &Skill{
		Name:  name,
		Layer: layer,
		Path:  destDir,
	}, nil
}

// List returns all skills organized by layer.
func List(targetDir string) (map[string][]Skill, error) {
	result := make(map[string][]Skill)

	for _, layer := range ValidLayers {
		layerDir := filepath.Join(targetDir, layer)

		// Check if the layer directory itself is a skill (built-in)
		if hasSkillMD(layerDir) {
			result[layer] = append(result[layer], Skill{
				Name:  layer,
				Layer: layer,
				Path:  layerDir,
			})
		}

		// Also look for sub-skills inside the layer directory
		entries, err := os.ReadDir(layerDir)
		if err != nil {
			continue
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			// Skip if this is the built-in skill itself (already added)
			skillPath := filepath.Join(layerDir, e.Name())
			if hasSkillMD(skillPath) {
				// This subdirectory is a skill (user-added)
				// Check if we already added this as the built-in
				already := false
				for _, s := range result[layer] {
					if s.Name == e.Name() {
						already = true
						break
					}
				}
				if !already {
					result[layer] = append(result[layer], Skill{
						Name:  e.Name(),
						Layer: layer,
						Path:  skillPath,
					})
				}
			}
		}
	}

	return result, nil
}

func hasSkillMD(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "SKILL.md"))
	return err == nil
}

// -- Source detection -------------------------------------------------------

func detectSourceType(source string) SourceType {
	// Git: github.com/... or git@... or .git
	if strings.Contains(source, "github.com") || strings.Contains(source, "git@") || strings.HasSuffix(source, ".git") {
		return SourceGit
	}
	// URL: starts with http:// or https://
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return SourceURL
	}
	return SourceLocal
}

// -- Downloaders ------------------------------------------------------------

func downloadURL(url, tmpDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	// Try tar.gz first
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "gzip") || strings.Contains(contentType, "tar") || strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz") {
		return extractTarGz(resp.Body, tmpDir)
	}

	// Try zip
	if strings.Contains(contentType, "zip") || strings.HasSuffix(url, ".zip") {
		return extractZip(resp.Body, tmpDir)
	}

	// Try tar.gz anyway
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b {
		return extractTarGz(io.NopCloser(strings.NewReader(string(body))), tmpDir)
	}

	return "", fmt.Errorf("unsupported file type from URL (need .tar.gz, .tgz, or .zip)")
}

func cloneGit(repo, tmpDir string) (string, error) {
	dest := filepath.Join(tmpDir, "repo")

	url := repo
	// Normalize: github.com/user/repo -> https://github.com/user/repo.git
	if strings.HasPrefix(url, "github.com/") {
		url = "https://" + url
	} else if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "git@") {
		url = "https://github.com/" + url
	}
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}

	cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git clone failed: %s", string(output))
	}

	return dest, nil
}

// -- Extractors -------------------------------------------------------------

func extractTarGz(r io.Reader, dest string) (string, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return "", fmt.Errorf("not a gzip file: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var firstDir string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read tar: %w", err)
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", fmt.Errorf("create dir: %w", err)
			}
			if firstDir == "" {
				firstDir = target
			}
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			f, err := os.Create(target)
			if err != nil {
				return "", fmt.Errorf("create file: %w", err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return "", fmt.Errorf("write file: %w", err)
			}
			f.Close()
		}
	}

	return firstDir, nil
}

func extractZip(r io.Reader, dest string) (string, error) {
	// Write to temp file first
	tmpFile := filepath.Join(dest, "archive.zip")
	f, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return "", err
	}
	f.Close()

	// Use unzip command
	outDir := filepath.Join(dest, "unzipped")
	cmd := exec.Command("unzip", "-o", tmpFile, "-d", outDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unzip failed: %s", string(output))
	}

	return findSkillRoot(outDir), nil
}

// -- Helpers ----------------------------------------------------------------

func findSkillRoot(dir string) string {
	// Check if dir itself has SKILL.md
	if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err == nil {
		return dir
	}

	// Look for a single subdirectory
	entries, _ := os.ReadDir(dir)
	if len(entries) == 1 && entries[0].IsDir() {
		sub := filepath.Join(dir, entries[0].Name())
		if _, err := os.Stat(filepath.Join(sub, "SKILL.md")); err == nil {
			return sub
		}
	}

	return dir
}

func validate(skillDir, name string) error {
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if _, err := os.Stat(skillFile); err != nil {
		return fmt.Errorf("SKILL.md not found in %s", skillDir)
	}

	data, err := os.ReadFile(skillFile)
	if err != nil {
		return fmt.Errorf("cannot read SKILL.md: %w", err)
	}

	content := string(data)

	// Check for YAML frontmatter
	if !strings.HasPrefix(content, "---") {
		return fmt.Errorf("SKILL.md missing YAML frontmatter (must start with ---)")
	}

	// Find closing ---
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return fmt.Errorf("SKILL.md has unclosed YAML frontmatter")
	}

	frontmatter := rest[:idx]

	// Must have name and description fields
	hasName := regexp.MustCompile(`(?m)^name:\s+\S+`).MatchString(frontmatter)
	hasDesc := regexp.MustCompile(`(?m)^description:\s+\S+`).MatchString(frontmatter)

	if !hasName {
		return fmt.Errorf("SKILL.md frontmatter missing 'name' field")
	}
	if !hasDesc {
		return fmt.Errorf("SKILL.md frontmatter missing 'description' field")
	}

	return nil
}

func isValidLayer(layer string) bool {
	for _, l := range ValidLayers {
		if l == layer {
			return true
		}
	}
	return false
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

// SkillPathFor returns the skill destination base path.
func SkillPathFor(targetDir string) string {
	return targetDir
}
