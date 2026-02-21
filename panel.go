package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PanelStyle defines the box-drawing characters for a panel border.
type PanelStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

// Pre-built panel styles.
var (
	PanelStyleRounded = PanelStyle{
		TopLeft: "╭", TopRight: "╮",
		BottomLeft: "╰", BottomRight: "╯",
		Horizontal: "─", Vertical: "│",
	}
	PanelStyleBox = PanelStyle{
		TopLeft: "┌", TopRight: "┐",
		BottomLeft: "└", BottomRight: "┘",
		Horizontal: "─", Vertical: "│",
	}
	PanelStyleDouble = PanelStyle{
		TopLeft: "╔", TopRight: "╗",
		BottomLeft: "╚", BottomRight: "╝",
		Horizontal: "═", Vertical: "║",
	}
	PanelStyleHeavy = PanelStyle{
		TopLeft: "┏", TopRight: "┓",
		BottomLeft: "┗", BottomRight: "┛",
		Horizontal: "━", Vertical: "┃",
	}
	PanelStyleSimple = PanelStyle{
		TopLeft: "+", TopRight: "+",
		BottomLeft: "+", BottomRight: "+",
		Horizontal: "-", Vertical: "|",
	}
)

// Panel displays text inside a decorative border box.
type Panel struct {
	content     string
	title       string
	subtitle    string
	style       PanelStyle
	titleStyle  Style
	borderStyle Style
	padding     int
	width       int
	w           io.Writer
}

// NewPanel creates a Panel with the given content.
func NewPanel(content string, opts ...func(*Panel)) *Panel {
	p := &Panel{
		content:    content,
		style:      PanelStyleRounded,
		titleStyle: NewStyle(Bold),
		padding:    1,
		w:          os.Stdout,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Panel option setters.

func PanelTitle(title string) func(*Panel)            { return func(p *Panel) { p.title = title } }
func PanelSubtitle(s string) func(*Panel)             { return func(p *Panel) { p.subtitle = s } }
func PanelBorderStyle(s PanelStyle) func(*Panel)      { return func(p *Panel) { p.style = s } }
func PanelTitleStyle(s Style) func(*Panel)            { return func(p *Panel) { p.titleStyle = s } }
func PanelBorderColor(s Style) func(*Panel)           { return func(p *Panel) { p.borderStyle = s } }
func PanelWidth(n int) func(*Panel)                   { return func(p *Panel) { p.width = n } }
func PanelPadding(n int) func(*Panel)                 { return func(p *Panel) { p.padding = n } }
func PanelWriter(w io.Writer) func(*Panel)            { return func(p *Panel) { p.w = w } }

// Render prints the panel to the configured writer.
func (p *Panel) Render() {
	lines := strings.Split(p.content, "\n")
	contentWidth := 0
	for _, l := range lines {
		if w := visibleLen(l); w > contentWidth {
			contentWidth = w
		}
	}

	innerWidth := contentWidth + p.padding*2
	if p.width > 0 && p.width-2 > innerWidth {
		innerWidth = p.width - 2
	}

	border := func(s string) string {
		if p.borderStyle.codes != nil {
			return p.borderStyle.Apply(s)
		}
		return s
	}

	s := p.style
	hbar := strings.Repeat(s.Horizontal, innerWidth)

	// Top line
	if p.title != "" {
		titleStr := p.titleStyle.Apply(p.title)
		titleLen := visibleLen(p.title)
		padLeft := (innerWidth - titleLen - 2) / 2
		if padLeft < 0 {
			padLeft = 0
		}
		padRight := innerWidth - padLeft - titleLen - 2
		if padRight < 0 {
			padRight = 0
		}
		topLeft := border(s.TopLeft + strings.Repeat(s.Horizontal, padLeft))
		topRight := border(strings.Repeat(s.Horizontal, padRight) + s.TopRight)
		fmt.Fprintf(p.w, "%s %s %s\n", topLeft, titleStr, topRight)
	} else {
		fmt.Fprintf(p.w, "%s\n", border(s.TopLeft+hbar+s.TopRight))
	}

	// Content lines
	pad := strings.Repeat(" ", p.padding)
	for _, line := range lines {
		llen := visibleLen(line)
		trailing := strings.Repeat(" ", innerWidth-p.padding-llen-p.padding)
		if innerWidth-p.padding*2-llen < 0 {
			trailing = ""
		}
		fmt.Fprintf(p.w, "%s%s%s%s%s\n", border(s.Vertical), pad, line, trailing+pad, border(s.Vertical))
	}

	// Bottom line
	if p.subtitle != "" {
		subStr := NewStyle(Dim).Apply(p.subtitle)
		subLen := visibleLen(p.subtitle)
		padLeft := (innerWidth - subLen - 2) / 2
		if padLeft < 0 {
			padLeft = 0
		}
		padRight := innerWidth - padLeft - subLen - 2
		if padRight < 0 {
			padRight = 0
		}
		botLeft := border(s.BottomLeft + strings.Repeat(s.Horizontal, padLeft))
		botRight := border(strings.Repeat(s.Horizontal, padRight) + s.BottomRight)
		fmt.Fprintf(p.w, "%s %s %s\n", botLeft, subStr, botRight)
	} else {
		fmt.Fprintf(p.w, "%s\n", border(s.BottomLeft+hbar+s.BottomRight))
	}
}

// PrintPanel is a convenience wrapper to print a panel directly.
func PrintPanel(content string, opts ...func(*Panel)) {
	NewPanel(content, opts...).Render()
}
