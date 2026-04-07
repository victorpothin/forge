// Package ui provides themed console output for the FORGE CLI.
// Theme: red/forge — bold, warm, intense.
package ui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func init() {
	// Force colors even when piped or in CI
	color.NoColor = false
}

// -- Colors ------------------------------------------------------------------

var (
	Red     = color.New(color.FgRed, color.Bold)
	Green   = color.New(color.FgGreen)
	Yellow  = color.New(color.FgYellow)
	Cyan    = color.New(color.FgCyan)
	White   = color.New(color.FgWhite)
	Dim     = color.New(color.Faint)
	Bold    = color.New(color.Bold)
	RedBold = color.New(color.FgRed, color.Bold)
	GreenB  = color.New(color.FgGreen, color.Bold)
)

// -- Spinner -----------------------------------------------------------------

type Spinner struct {
	mu      sync.Mutex
	stopCh  chan struct{}
	frames  []string
	message string
	current int
	running bool
}

var spinnerFrames = []string{"◜", "◝", "◞", "◟"}

// NewSpinner creates a new spinner with the given message.
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:  spinnerFrames,
		message: message,
		stopCh:  make(chan struct{}),
	}
}

// Start begins the spinner animation.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopCh:
				fmt.Print("\r\033[K")
				return
			case <-ticker.C:
				s.mu.Lock()
				frame := s.frames[s.current%len(s.frames)]
				s.current++
				fmt.Printf("\r  %s %s", Red.Sprint(frame), s.message)
				s.mu.Unlock()
			}
		}
	}()
}

// Stop halts the spinner and clears the line.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()
	close(s.stopCh)
	time.Sleep(20 * time.Millisecond)
	fmt.Print("\r\033[K")
}

// StopWith replaces the spinner line with a final message.
func (s *Spinner) StopWith(format string, args ...interface{}) {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		fmt.Printf("  "+format+"\n", args...)
		return
	}
	s.running = false
	s.mu.Unlock()
	close(s.stopCh)
	time.Sleep(20 * time.Millisecond)

	fmt.Print("\r\033[K")
	fmt.Printf("  "+format+"\n", args...)
}

// -- Banner ------------------------------------------------------------------

const bannerLines = `
┌──────────────────────────────────────┐
│       FORGE CLI                      │
│  Focused, Ordered, Restricted,       │
│  Guided Execution                    │`

func Banner(ver string) {
	fmt.Println()
	fmt.Println(RedBold.Sprint(bannerLines))
	fmt.Printf("│  %-36s│\n", "v"+ver)
	fmt.Println(RedBold.Sprint("└──────────────────────────────────────┘"))
	fmt.Println()
}

// -- Helpers -----------------------------------------------------------------

func Section(title string) {
	fmt.Println()
	RedBold.Printf("── %s ──\n", title)
	fmt.Println()
}

func PrintHeader(title string) {
	Bold.Printf("📋 %s:\n", title)
}

func PrintConfig(key, value string, extra ...string) {
	Dim.Printf("  %-12s ", key+":")
	White.Printf("%s", value)
	if len(extra) > 0 {
		fmt.Print(" " + extra[0])
	}
	fmt.Println()
}

func PrintFileListPending(files []string) {
	for _, f := range files {
		Red.Printf("  ✅ %s\n", f)
	}
}

func PrintSuccess(message string) {
	GreenB.Printf("  ✓ %s\n", message)
}

func PrintInfo(message string) {
	Dim.Printf("  %s\n", message)
}

func PrintError(format string, args ...interface{}) {
	Red.Fprintf(os.Stderr, "\n❌ Error: "+format+"\n", args...)
}

func PrintWarning(message string) {
	Yellow.Printf("  ⚠️  %s\n", message)
}

func PrintSuccessBox(target string) {
	lines := []string{
		GreenB.Sprint("  ✓ FORGE initialized"),
		fmt.Sprintf("    %s", target),
	}
	printBox(lines)
}

func PrintBox(lines ...string) {
	printBox(lines)
}

func printBox(lines []string) {
	// Calculate max visible width
	maxW := 0
	for _, l := range lines {
		w := visibleLen(l)
		if w > maxW {
			maxW = w
		}
	}

	border := "─" + strings.Repeat("─", maxW) + "─"
	top := "┌" + border + "┐"
	bot := "└" + border + "┘"

	fmt.Println("  " + RedBold.Sprint(top))
	for _, l := range lines {
		padding := maxW - visibleLen(l)
		fmt.Printf("  %s %s%s %s\n",
			RedBold.Sprint("│"),
			l,
			strings.Repeat(" ", padding),
			RedBold.Sprint("│"),
		)
	}
	fmt.Println("  " + RedBold.Sprint(bot))
}

// visibleLen returns string length ignoring ANSI escape codes.
func visibleLen(s string) int {
	count := 0
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		count++
	}
	return count
}
