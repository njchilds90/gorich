package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const defaultTermWidth = 80

// termWidthForWriter returns the best-effort terminal width for the given writer.
//
// If the writer is not an *os.File or the terminal size cannot be determined,
// this function returns defaultTermWidth.
func termWidthForWriter(w io.Writer) int {
	f, ok := w.(*os.File)
	if !ok {
		return defaultTermWidth
	}
	width, _, err := term.GetSize(int(f.Fd()))
	if err != nil || width <= 0 {
		return defaultTermWidth
	}
	return width
}

// Rule prints a horizontal divider similar to Python's rich.rule.
//
// A Rule can be customized with a title, styles, a drawing character, and a
// fixed width. If no width is set, the terminal width is used when possible.
type Rule struct {
	title     string
	style     Style
	ruleStyle Style
	char      string
	width     int
	w         io.Writer
}

// NewRule creates a Rule with sensible defaults.
//
// Use RuleWriter to render to a buffer or to a non-standard output stream.
func NewRule(opts ...func(*Rule)) *Rule {
	r := &Rule{
		style:     NewStyle(Bold, BrightMagenta),
		ruleStyle: NewStyle(Dim),
		char:      "─",
		w:         os.Stdout,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// RuleTitle sets the title rendered in the center of the rule.
func RuleTitle(title string) func(*Rule) { return func(r *Rule) { r.title = title } }

// RuleStyle sets the Style applied to the title.
func RuleStyle(s Style) func(*Rule) { return func(r *Rule) { r.style = s } }

// RuleChar sets the drawing character used for the left and right lines.
func RuleChar(c string) func(*Rule) { return func(r *Rule) { r.char = c } }

// RuleWidth sets a fixed width. If width is less than or equal to zero, the terminal width is used.
func RuleWidth(n int) func(*Rule) { return func(r *Rule) { r.width = n } }

// RuleWriter sets the writer used for rendering.
func RuleWriter(w io.Writer) func(*Rule) { return func(r *Rule) { r.w = w } }

// Render prints the rule.
func (r *Rule) Render() {
	width := r.width
	if width <= 0 {
		width = termWidthForWriter(r.w)
	}

	if r.title == "" {
		fmt.Fprintln(r.w, r.ruleStyle.Apply(strings.Repeat(r.char, width)))
		return
	}

	titleLen := visibleLen(r.title)
	side := (width - titleLen - 2) / 2
	if side < 1 {
		side = 1
	}

	leftLine := r.ruleStyle.Apply(strings.Repeat(r.char, side))
	rightWidth := width - side - titleLen - 2
	if rightWidth < 0 {
		rightWidth = 0
	}
	rightLine := r.ruleStyle.Apply(strings.Repeat(r.char, rightWidth))
	fmt.Fprintf(r.w, "%s %s %s\n", leftLine, r.style.Apply(r.title), rightLine)
}

// PrintRule is a convenience wrapper that creates and renders a Rule.
func PrintRule(title string, opts ...func(*Rule)) {
	NewRule(append([]func(*Rule){RuleTitle(title)}, opts...)...).Render()
}
