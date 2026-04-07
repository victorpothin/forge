package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/victorpothin/forge/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	var memDir string

	memoryCmd := &cobra.Command{
		Use:   "memory",
		Short: "Manage persistent: rules and context that persist across sessions",
		Long: `Store important project rules and context that persist across sessions.
Memories are lightweight and only loaded when relevant to the current context.

Examples:
  forge memory add "PostgreSQL only — no SQLite"
  forge memory add --tag db "Always use migrations"
  forge memory list
  forge memory search database
  forge memory delete 1
  forge memory clear`,
	}

	memoryCmd.PersistentFlags().StringVarP(&memDir, "dir", "d", ".", "Project directory")
	memoryCmd.AddCommand(memAddCmd())
	memoryCmd.AddCommand(memListCmd())
	memoryCmd.AddCommand(memSearchCmd())
	memoryCmd.AddCommand(memDeleteCmd())
	memoryCmd.AddCommand(memClearCmd())

	rootCmd.AddCommand(memoryCmd)
}

type Memory struct {
	ID   int
	Tags string
	Text string
}

func memFile(dir string) string {
	return filepath.Join(dir, ".forge-memory")
}

func loadMemories(dir string) ([]Memory, error) {
	f, err := os.Open(memFile(dir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var mems []Memory
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}
		id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		mems = append(mems, Memory{
			ID:   id,
			Tags: strings.TrimSpace(parts[1]),
			Text: strings.TrimSpace(parts[2]),
		})
	}
	return mems, scanner.Err()
}

func saveMemories(dir string, mems []Memory) error {
	f, err := os.Create(memFile(dir))
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("# forge-memory v1\n")
	f.WriteString("# id|tags|text\n")
	for _, m := range mems {
		fmt.Fprintf(f, "%d|%s|%s\n", m.ID, m.Tags, m.Text)
	}
	return nil
}

func nextID(mems []Memory) int {
	max := 0
	for _, m := range mems {
		if m.ID > max {
			max = m.ID
		}
	}
	return max + 1
}

func plural(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

func getDir(cmd *cobra.Command) string {
	dir, _ := cmd.Flags().GetString("dir")
	return dir
}

// -- Subcommands -------------------------------------------------------------

func memAddCmd() *cobra.Command {
	var tags string
	cmd := &cobra.Command{
		Use:   "add [text]",
		Short: "Save a project memory",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			var text string
			if len(args) > 0 {
				text = args[0]
			}
			if text == "" {
				ui.Ask(&survey.Input{Message: "Memory text:"}, &text)
				if text == "" {
					ui.Dim.Println("  Aborted.")
					return
				}
			}

			mems, err := loadMemories(dir)
			if err != nil {
				ui.PrintError("Failed to load memories: %v", err)
				return
			}

			mems = append(mems, Memory{
				ID:   nextID(mems),
				Tags: tags,
				Text: text,
			})

			if err := saveMemories(dir, mems); err != nil {
				ui.PrintError("Failed to save memory: %v", err)
				return
			}

			fmt.Println()
			ui.PrintSuccess(fmt.Sprintf("Memory saved (ID: %d)", len(mems)))
			fmt.Println()
		},
	}
	cmd.Flags().StringVarP(&tags, "tag", "t", "", "Comma-separated tags")
	return cmd
}

func memListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all stored memories",
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			mems, err := loadMemories(dir)
			if err != nil {
				ui.PrintError("Failed to load memories: %v", err)
				return
			}
			if len(mems) == 0 {
				ui.Dim.Println("  No memories stored.")
				fmt.Println()
				return
			}

			fmt.Println()
			ui.Bold.Printf("  Project Memories\n")
			ui.Dim.Printf("  %s\n\n", memFile(dir))

			for _, m := range mems {
				tagStr := ""
				if m.Tags != "" {
					tagStr = fmt.Sprintf(" [%s]", m.Tags)
				}
				fmt.Printf("  %d%s: %s\n", m.ID, tagStr, m.Text)
			}
			fmt.Println()
			ui.Dim.Printf("  Total: %d %s\n", len(mems), plural(len(mems), "memory", "memories"))
			fmt.Println()
		},
	}
}

func memSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Search memories by text or tag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			query := strings.ToLower(args[0])
			mems, err := loadMemories(dir)
			if err != nil {
				ui.PrintError("Failed to load memories: %v", err)
				return
			}

			var found []Memory
			for _, m := range mems {
				if strings.Contains(strings.ToLower(m.Text), query) ||
					strings.Contains(strings.ToLower(m.Tags), query) {
					found = append(found, m)
				}
			}

			if len(found) == 0 {
				ui.Dim.Printf("  No memories matching %q\n", query)
				fmt.Println()
				return
			}

			fmt.Println()
			ui.Bold.Printf("  Matching memories (%d)\n\n", len(found))
			for _, m := range found {
				tagStr := ""
				if m.Tags != "" {
					tagStr = fmt.Sprintf(" [%s]", m.Tags)
				}
				fmt.Printf("  %d%s: %s\n", m.ID, tagStr, m.Text)
			}
			fmt.Println()
		},
	}
}

func memDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete [id...]",
		Short:   "Delete memories by ID",
		Aliases: []string{"rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			mems, err := loadMemories(dir)
			if err != nil {
				ui.PrintError("Failed to load memories: %v", err)
				return
			}

			toDelete := make(map[int]bool)
			for _, a := range args {
				id, err := strconv.Atoi(a)
				if err != nil {
					ui.PrintError("Invalid ID: %s", a)
					return
				}
				toDelete[id] = true
			}

			var kept []Memory
			for _, m := range mems {
				if !toDelete[m.ID] {
					kept = append(kept, m)
				}
			}

			deleted := len(mems) - len(kept)
			if deleted == 0 {
				ui.Dim.Println("  Nothing to delete.")
				fmt.Println()
				return
			}

			if err := saveMemories(dir, kept); err != nil {
				ui.PrintError("Failed to save memories: %v", err)
				return
			}

			fmt.Println()
			ui.PrintSuccess(fmt.Sprintf("Deleted %d %s", deleted, plural(deleted, "memory", "memories")))
			fmt.Println()
		},
	}
}

func memClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Delete all stored memories",
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			mems, err := loadMemories(dir)
			if err != nil {
				if os.IsNotExist(err) {
					ui.Dim.Println("  No memories to clear.")
					fmt.Println()
					return
				}
				ui.PrintError("Failed to load memories: %v", err)
				return
			}
			if len(mems) == 0 {
				ui.Dim.Println("  No memories to clear.")
				fmt.Println()
				return
			}

			ok := false
			ui.Ask(&survey.Confirm{
				Message: fmt.Sprintf("Delete all %d memories?", len(mems)),
				Default: false,
			}, &ok)
			if !ok {
				ui.Dim.Println("  Cancelled.")
				fmt.Println()
				return
			}

			if err := os.Remove(memFile(dir)); err != nil {
				ui.PrintError("Failed to clear memories: %v", err)
				return
			}

			fmt.Println()
			ui.PrintSuccess("All memories cleared")
			fmt.Println()
		},
	}
}
