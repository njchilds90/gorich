package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
	"golang.org/x/term"
)

const defaultTermWidth = 80

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return defaultTermWidth
	}
	return w
}

// Rule prints a horizontal rule/divider (like Python's rich.rule).
type Rule struct {
	title     string
	style     Style
	ruleStyle Style
	char      string
	width     int
	w         io.Writer
}

// NewRule creates a Rule.
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

// Rule option setters.

func RuleTitle(title string) func(*Rule)      { return func(r *Rule) { r.title = title } }
func RuleStyle(s Style) func(*Rule)           { return func(r *Rule) { r.style = s } }
func RuleChar(c string) func(*Rule)           { return func(r *Rule) { r.char = c } }
func RuleWidth(n int) func(*Rule)             { return func(r *Rule) { r.width = n } }
func RuleWriter(w io.Writer) func(*Rule)      { return func(r *Rule) { r.w = w } }

// Render prints the rule.
func (r *Rule) Render() {
	width := r.width
	if width <= 0 {
		width = termWidth()
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
	rightLine := r.ruleStyle.Apply(strings.Repeat(r.char, width-side-titleLen-2))
	fmt.Fprintf(r.w, "%s %s %s\n", leftLine, r.style.Apply(r.title), rightLine)
}

// PrintRule is a convenience wrapper.
func PrintRule(title string, opts ...func(*Rule)) {
	NewRule(append([]func(*Rule){RuleTitle(title)}, opts...)...).Render()
}
