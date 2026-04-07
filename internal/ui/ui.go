// Package ui provides clean, minimal console output for the FORGE CLI.
package ui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/AlecAivazis/survey/v2"
)

func init() {
	color.NoColor = false
}

// -- Colors ------------------------------------------------------------------

var (
	Bold    = color.New(color.Bold)
	Dim     = color.New(color.Faint)
	Green   = color.New(color.FgGreen, color.Bold)
	Red     = color.New(color.FgRed, color.Bold)
	Yellow  = color.New(color.FgYellow)
	Cyan    = color.New(color.FgCyan)
)

// -- Survey helper ------------------------------------------------------------

// Ask wraps survey.AskOne with Ctrl+C exit message.
func Ask(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	// Catch Ctrl+C globally
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		Dim.Println("  Quit (Ctrl+C)")
		os.Exit(0)
	}()
	return survey.AskOne(prompt, response, opts...)
}

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

func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:  spinnerFrames,
		message: message,
		stopCh:  make(chan struct{}),
	}
}

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
				fmt.Printf("\r  %s  %s", frame, s.message)
				s.mu.Unlock()
			}
		}
	}()
}

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

func Banner(ver string) {
	fmt.Println()
	fmt.Println(Bold.Sprint("  FORGE CLI  " + "v" + ver))
	fmt.Println()
}

// -- Helpers -----------------------------------------------------------------

func Section(title string) {
	fmt.Println()
	fmt.Println(Bold.Sprint("  " + title))
	fmt.Println()
}

func PrintHeader(title string) {
	Bold.Printf("  %s:\n", title)
}

func PrintConfig(key, value string) {
	Dim.Printf("    %-12s %s\n", key+":", value)
}

func PrintFileListPending(files []string) {
	for _, f := range files {
		fmt.Printf("  ✅ %s\n", f)
	}
}

func PrintSuccess(message string) {
	Green.Printf("  ✓ %s\n", message)
}

func PrintInfo(message string) {
	Dim.Printf("  %s\n", message)
}

func PrintError(format string, args ...interface{}) {
	Red.Fprintf(os.Stderr, "\n✗ Error: "+format+"\n", args...)
}

func PrintWarning(message string) {
	Yellow.Printf("  ⚠  %s\n", message)
}

func PrintBox(lines ...string) {
	maxW := 0
	for _, l := range lines {
		w := visibleLen(l)
		if w > maxW {
			maxW = w
		}
	}

	border := strings.Repeat("─", maxW+4)
	fmt.Println("  ┌" + border + "┐")
	for _, l := range lines {
		padding := maxW - visibleLen(l)
		fmt.Printf("  │ %s%s │\n", l, strings.Repeat(" ", padding))
	}
	fmt.Println("  └" + border + "┘")
}

func PrintSuccessBox(target string) {
	PrintBox(
		Green.Sprint("FORGE initialized"),
		Dim.Sprint("  "+target),
	)
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
