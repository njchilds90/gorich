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
	"sync"
	"time"
)

// Log prints a styled log-level message to standard output.
//
// Common levels include: INFO, SUCCESS, WARNING, ERROR, and DEBUG.
func Log(level, message string) {
	Flog(os.Stdout, level, message)
}

// logLevelStyle maps log levels to their corresponding output Style.
//
// This is package-scoped to avoid rebuilding the map on every Flog call.
var logLevelStyle = map[string]Style{
	"INFO":    StyleInfo,
	"SUCCESS": StyleSuccess,
	"WARNING": StyleWarning,
	"ERROR":   StyleError,
	"DEBUG":   StyleMuted,
}

// Flog writes a log message to a writer.
//
// The level controls the style of the bracketed prefix. If the level is unknown,
// StyleMuted is used.
func Flog(w io.Writer, level, message string) {
	s, ok := logLevelStyle[level]
	if !ok {
		s = StyleMuted
	}
	fmt.Fprintf(w, "%s %s\n", s.Apply("["+level+"]"), message)
}

// Inspect pretty-prints any value in a labeled panel.
//
// This function is intended for quick debugging and development output.
func Inspect(label string, value interface{}) {
	content := fmt.Sprintf("Type: %T\nValue: %+v", value, value)
	NewPanel(content,
		PanelTitle(fmt.Sprintf(" %s ", label)),
		PanelBorderStyle(PanelStyleRounded),
		PanelBorderColor(NewStyle(BrightMagenta)),
	).Render()
}

// Spinner is a terminal spinner for indeterminate tasks.
//
// Start begins the animation in a background goroutine and Stop halts it.
// Stop is safe to call multiple times.
type Spinner struct {
	frames      []string
	description string
	style       Style
	w           io.Writer

	done     chan struct{}
	stopOnce sync.Once
	startOnce sync.Once
	wg       sync.WaitGroup
}

// NewSpinner creates a new Spinner.
//
// If SpinnerFrames is used to provide an empty list, NewSpinner substitutes a
// single "-" frame to avoid a division by zero in the animation loop.
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
	if len(s.frames) == 0 {
		s.frames = []string{"-"}
	}
	return s
}

// SpinnerFrames sets the animation frames used by the Spinner.
func SpinnerFrames(frames []string) func(*Spinner) { return func(s *Spinner) { s.frames = frames } }

// SpinnerStyle sets the Style applied to each frame.
func SpinnerStyle(st Style) func(*Spinner) { return func(s *Spinner) { s.style = st } }

// SpinnerWriter sets the writer used for rendering.
func SpinnerWriter(w io.Writer) func(*Spinner) { return func(s *Spinner) { s.w = w } }

// Start begins the spinner animation in a goroutine.
//
// Start is safe to call multiple times, but only the first call starts an animation.
func (s *Spinner) Start() {
	s.startOnce.Do(func() {
		// Ensure frames are valid even if the caller mutated them after construction.
		if len(s.frames) == 0 {
			s.frames = []string{"-"}
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

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
	})
}

// Stop halts the spinner and prints a completion checkmark.
//
// Stop is idempotent so that calling it multiple times does not panic.
func (s *Spinner) Stop() {
	s.stopOnce.Do(func() {
		close(s.done)
	})

	// Wait for the animation goroutine to exit so that completion output does not
	// interleave with an in-flight frame update.
	s.wg.Wait()
	fmt.Fprintf(s.w, "\r%s %s   \n", StyleSuccess.Apply("✓"), s.description)
}

// WithSpinner runs fn with a spinner and stops it when fn returns.
//
// This helper is useful for short, synchronous operations.
func WithSpinner(description string, fn func()) {
	s := NewSpinner(description)
	s.Start()
	fn()
	s.Stop()
}

// Fcolumns writes multiple strings side by side in equal-width columns.
//
// This function exists to support the library's pattern of writing to an io.Writer.
// If columnCount is less than or equal to zero, it is treated as 1.
func Fcolumns(w io.Writer, items []string, columnCount int, style ...Style) {
	if columnCount <= 0 {
		columnCount = 1
	}

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
		fmt.Fprintf(w, "%-*s", colWidth, st.Apply(item))
		if (i+1)%columnCount == 0 {
			fmt.Fprintln(w)
		}
	}
	if len(items)%columnCount != 0 {
		fmt.Fprintln(w)
	}
}

// Columns prints multiple strings side by side in equal-width columns to standard output.
//
// For non-standard output streams or tests, use Fcolumns.
func Columns(items []string, columnCount int, style ...Style) {
	Fcolumns(os.Stdout, items, columnCount, style...)
}
