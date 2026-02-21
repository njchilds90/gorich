// Package gorich provides beautiful, expressive terminal output for Go,
// inspired by Python's rich library.
//
// # Quick Start
//
//	import "github.com/njchilds90/gorich"
//
//	// Print styled text with markup
//	gorich.Println("[bold green]Hello, World![/]")
//
//	// Print a panel
//	gorich.PrintPanel("Task complete!", gorich.PanelTitle("Status"))
//
//	// Print a table
//	t := gorich.NewTable(gorich.WithTitle("Users"))
//	t.AddColumn("Name")
//	t.AddColumn("Role")
//	t.AddRow("Alice", "Admin")
//	t.Render()
//
//	// Print a rule
//	gorich.PrintRule("Section Title")
//
//	// Syntax highlight
//	gorich.PrintSyntax(code, "go", "main.go")
package gorich

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Log prints a timestamped, styled log-level message.
// Levels: INFO, SUCCESS, WARNING, ERROR, DEBUG
func Log(level, message string) {
	Flog(os.Stdout, level, message)
}

// Flog writes a log message to a writer.
func Flog(w io.Writer, level, message string) {
	styleMap := map[string]Style{
		"INFO":    StyleInfo,
		"SUCCESS": StyleSuccess,
		"WARNING": StyleWarning,
		"ERROR":   StyleError,
		"DEBUG":   StyleMuted,
	}
	s, ok := styleMap[level]
	if !ok {
		s = StyleMuted
	}
	fmt.Fprintf(w, "%s %s\n", s.Apply("["+level+"]"), message)
}

// Inspect pretty-prints any value in a labeled panel.
func Inspect(label string, value interface{}) {
	content := fmt.Sprintf("Type: %T\nValue: %+v", value, value)
	NewPanel(content,
		PanelTitle(fmt.Sprintf(" %s ", label)),
		PanelBorderStyle(PanelStyleRounded),
		PanelBorderColor(NewStyle(BrightMagenta)),
	).Render()
}

// Spinner is a terminal spinner for indeterminate tasks.
type Spinner struct {
	frames      []string
	description string
	style       Style
	w           io.Writer
	done        chan struct{}
}

// NewSpinner creates a new Spinner.
func NewSpinner(description string, opts ...func(*Spinner)) *Spinner {
	s := &Spinner{
		frames:      []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		description: description,
		style:       NewStyle(BrightCyan),
		w:           os.Stdout,
		done:        make(chan struct{}),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Spinner option setters.

func SpinnerFrames(frames []string) func(*Spinner) { return func(s *Spinner) { s.frames = frames } }
func SpinnerStyle(st Style) func(*Spinner)         { return func(s *Spinner) { s.style = st } }
func SpinnerWriter(w io.Writer) func(*Spinner)     { return func(s *Spinner) { s.w = w } }

// Start begins the spinner animation in a goroutine. Call Stop() when done.
func (s *Spinner) Start() {
	go func() {
		frame := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				f := s.frames[frame%len(s.frames)]
				fmt.Fprintf(s.w, "\r%s %s", s.style.Apply(f), s.description)
				frame++
			}
		}
	}()
}

// Stop halts the spinner and prints a completion checkmark.
func (s *Spinner) Stop() {
	close(s.done)
	time.Sleep(150 * time.Millisecond)
	fmt.Fprintf(s.w, "\r%s %s   \n", StyleSuccess.Apply("✓"), s.description)
}

// WithSpinner runs fn with a spinner and stops it when fn returns.
func WithSpinner(description string, fn func()) {
	s := NewSpinner(description)
	s.Start()
	fn()
	s.Stop()
}

// Columns prints multiple strings side by side in equal-width columns.
func Columns(items []string, columnCount int, style ...Style) {
	st := StyleMuted
	if len(style) > 0 {
		st = style[0]
	}
	colWidth := 0
	for _, item := range items {
		if l := visibleLen(item); l > colWidth {
			colWidth = l
		}
	}
	colWidth += 2
	for i, item := range items {
		fmt.Printf("%-*s", colWidth, st.Apply(item))
		if (i+1)%columnCount == 0 {
			fmt.Println()
		}
	}
	if len(items)%columnCount != 0 {
		fmt.Println()
	}
}
