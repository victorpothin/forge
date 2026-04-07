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

// -- Colors ------------------------------------------------------------------

var (
	Red     = color.New(color.FgRed, color.Bold)
	Yellow  = color.New(color.FgYellow)
	Green   = color.New(color.FgGreen)
	Cyan    = color.New(color.FgCyan)
	White   = color.New(color.FgWhite)
	Dim     = color.New(color.Faint)
	Bold    = color.New(color.Bold)
	RedBold = color.New(color.FgRed, color.Bold)
)

// -- Spinner -----------------------------------------------------------------

// Spinner represents an animated loading spinner.
type Spinner struct {
	mu      sync.Mutex
	stopCh  chan struct{}
	frames  []string
	message string
	current int
	running bool
}

var spinnerFrames = []string{
	"◜", "◝", "◞", "◟",
}

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
	// Give a tick for the goroutine to exit
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

	msg := fmt.Sprintf(format, args...)
	// Clear the spinner line completely
	fmt.Print("\r\033[K")
	fmt.Printf("  %s\n", msg)
}

// -- Progress bar ------------------------------------------------------------

func ProgressBar(done, total int, width int) string {
	if total == 0 {
		total = 1
	}
	if width == 0 {
		width = 30
	}

	filled := (done * width) / total
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += RedBold.Sprint("█")
	}
	emptyBar := strings.Repeat("·", empty)

	pct := (done * 100) / total
	return fmt.Sprintf("[%s%s] %d%%", bar, emptyBar, pct)
}

// -- Banner ------------------------------------------------------------------

const banner = `
┌──────────────────────────────────────┐
│       FORGE CLI                      │
│  Focused, Ordered, Restricted,       │
│  Guided Execution                    │`

func Banner(version string) {
	fmt.Println()
	fmt.Println(banner)
	fmt.Printf("│  %-36s│\n", "v"+version)
	fmt.Println("└──────────────────────────────────────┘")
	fmt.Println()
}

// -- Helpers -----------------------------------------------------------------

func PrintHeader(title string) {
	Bold.Printf("📋 %s:\n", title)
}

func PrintConfig(key, value string) {
	Dim.Printf("  %-12s ", key+":")
	White.Printf("%s\n", value)
}

func PrintFileList(files []string) {
	for _, f := range files {
		Green.Printf("  ✓ %s\n", f)
	}
}

func PrintFileListPending(files []string) {
	for _, f := range files {
		White.Printf("  ✅ %s\n", f)
	}
}

func PrintSuccess(message string) {
	Green.Printf("  ✓ %s\n", message)
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

func PrintBox(lines ...string) {
	maxLen := 0
	for _, l := range lines {
		// Strip ANSI color codes for width calc
		clean := l
		for i := 0; i < len(clean); {
			if clean[i] == '\033' {
				j := i
				for j < len(clean) && clean[j] != 'm' {
					j++
				}
				clean = clean[:i] + clean[j+1:]
			} else {
				i++
			}
		}
		if len(clean) > maxLen {
			maxLen = len(clean)
		}
	}

	border := strings.Repeat("─", maxLen+4)
	fmt.Println("  ┌" + border + "┐")
	for _, l := range lines {
		// Calculate visible length (strip ANSI codes)
		visible := l
		for i := 0; i < len(visible); {
			if visible[i] == '\033' {
				j := i
				for j < len(visible) && visible[j] != 'm' {
					j++
				}
				visible = visible[:i] + visible[j+1:]
			} else {
				i++
			}
		}
		padding := maxLen - len(visible)
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("  │ %s%s │\n", l, strings.Repeat(" ", padding))
	}
	fmt.Println("  └" + border + "┘")
}

func Section(title string) {
	fmt.Println()
	RedBold.Printf("── %s ──\n", title)
	fmt.Println()
}
