package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// TableStyle defines the box-drawing characters for a table border.
type TableStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	MiddleLeft  string
	MiddleRight string
	MiddleCross string
	TopMiddle   string
	BottomMiddle string
}

// Pre-built table styles.
var (
	TableStyleRounded = TableStyle{
		TopLeft: "╭", TopRight: "╮",
		BottomLeft: "╰", BottomRight: "╯",
		Horizontal: "─", Vertical: "│",
		MiddleLeft: "├", MiddleRight: "┤",
		MiddleCross: "┼", TopMiddle: "┬", BottomMiddle: "┴",
	}
	TableStyleBox = TableStyle{
		TopLeft: "┌", TopRight: "┐",
		BottomLeft: "└", BottomRight: "┘",
		Horizontal: "─", Vertical: "│",
		MiddleLeft: "├", MiddleRight: "┤",
		MiddleCross: "┼", TopMiddle: "┬", BottomMiddle: "┴",
	}
	TableStyleDouble = TableStyle{
		TopLeft: "╔", TopRight: "╗",
		BottomLeft: "╚", BottomRight: "╝",
		Horizontal: "═", Vertical: "║",
		MiddleLeft: "╠", MiddleRight: "╣",
		MiddleCross: "╬", TopMiddle: "╦", BottomMiddle: "╩",
	}
	TableStyleSimple = TableStyle{
		TopLeft: "+", TopRight: "+",
		BottomLeft: "+", BottomRight: "+",
		Horizontal: "-", Vertical: "|",
		MiddleLeft: "+", MiddleRight: "+",
		MiddleCross: "+", TopMiddle: "+", BottomMiddle: "+",
	}
	TableStyleMinimal = TableStyle{
		Horizontal: "─", Vertical: " ",
		TopLeft: " ", TopRight: " ", BottomLeft: " ", BottomRight: " ",
		MiddleLeft: " ", MiddleRight: " ", MiddleCross: "─", TopMiddle: "─", BottomMiddle: "─",
	}
)

// Alignment for table columns.
type Alignment int

const (
	AlignLeft   Alignment = iota
	AlignCenter Alignment = iota
	AlignRight  Alignment = iota
)

// Column defines a table column.
type Column struct {
	Header    string
	Style     Style
	Align     Alignment
	HeaderStyle Style
	MinWidth  int
	MaxWidth  int
}

// Table is a rich terminal table.
type Table struct {
	title       string
	columns     []Column
	rows        [][]string
	style       TableStyle
	headerStyle Style
	showHeader  bool
	showLines   bool
	padding     int
	w           io.Writer
	caption     string
}

// NewTable creates a new Table writing to stdout.
func NewTable(opts ...func(*Table)) *Table {
	t := &Table{
		style:       TableStyleRounded,
		headerStyle: NewStyle(Bold, BrightCyan),
		showHeader:  true,
		showLines:   false,
		padding:     1,
		w:           os.Stdout,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Option setters for Table.

func WithTitle(title string) func(*Table) {
	return func(t *Table) { t.title = title }
}
func WithCaption(caption string) func(*Table) {
	return func(t *Table) { t.caption = caption }
}
func WithTableStyle(style TableStyle) func(*Table) {
	return func(t *Table) { t.style = style }
}
func WithHeaderStyle(style Style) func(*Table) {
	return func(t *Table) { t.headerStyle = style }
}
func WithShowLines(show bool) func(*Table) {
	return func(t *Table) { t.showLines = show }
}
func WithShowHeader(show bool) func(*Table) {
	return func(t *Table) { t.showHeader = show }
}
func WithPadding(n int) func(*Table) {
	return func(t *Table) { t.padding = n }
}
func WithWriter(w io.Writer) func(*Table) {
	return func(t *Table) { t.w = w }
}

// AddColumn appends a column definition.
func (t *Table) AddColumn(header string, opts ...func(*Column)) {
	col := Column{Header: header, Align: AlignLeft}
	for _, o := range opts {
		o(&col)
	}
	t.columns = append(t.columns, col)
}

// Column option setters.

func ColStyle(s Style) func(*Column)        { return func(c *Column) { c.Style = s } }
func ColAlign(a Alignment) func(*Column)    { return func(c *Column) { c.Align = a } }
func ColMinWidth(n int) func(*Column)       { return func(c *Column) { c.MinWidth = n } }
func ColMaxWidth(n int) func(*Column)       { return func(c *Column) { c.MaxWidth = n } }
func ColHeaderStyle(s Style) func(*Column)  { return func(c *Column) { c.HeaderStyle = s } }

// AddRow appends a data row. Values are matched to columns by position.
func (t *Table) AddRow(values ...string) {
	t.rows = append(t.rows, values)
}

// Render prints the table to the configured writer.
func (t *Table) Render() {
	pad := strings.Repeat(" ", t.padding)

	// compute column widths
	widths := make([]int, len(t.columns))
	for i, col := range t.columns {
		widths[i] = visibleLen(col.Header)
		if col.MinWidth > widths[i] {
			widths[i] = col.MinWidth
		}
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if i >= len(widths) {
				break
			}
			l := visibleLen(cell)
			if l > widths[i] {
				widths[i] = l
			}
		}
	}
	// apply max widths
	for i, col := range t.columns {
		if col.MaxWidth > 0 && widths[i] > col.MaxWidth {
			widths[i] = col.MaxWidth
		}
	}

	s := t.style
	totalWidth := 2 // left + right border
	for _, w := range widths {
		totalWidth += w + t.padding*2 + 1
	}
	totalWidth-- // last column doesn't add trailing separator

	// title
	if t.title != "" {
		titleLine := t.headerStyle.Apply(" " + t.title + " ")
		inner := totalWidth - 2
		titlePad := (inner - visibleLen(t.title) - 2) / 2
		if titlePad < 0 {
			titlePad = 0
		}
		leftPad := strings.Repeat(s.Horizontal, titlePad)
		rightPad := strings.Repeat(s.Horizontal, inner-titlePad-visibleLen(t.title)-2)
		fmt.Fprintf(t.w, "%s%s%s%s%s\n", s.TopLeft, leftPad, titleLine, rightPad, s.TopRight)
	} else {
		fmt.Fprintln(t.w, t.topLine(widths))
	}

	// header
	if t.showHeader {
		fmt.Fprint(t.w, s.Vertical)
		for i, col := range t.columns {
			hs := t.headerStyle
			if col.HeaderStyle.codes != nil {
				hs = col.HeaderStyle
			}
			cell := hs.Apply(alignText(col.Header, widths[i], col.Align))
			fmt.Fprintf(t.w, "%s%s%s%s", pad, cell, pad, s.Vertical)
		}
		fmt.Fprintln(t.w)
		fmt.Fprintln(t.w, t.midLine(widths))
	}

	// rows
	for ri, row := range t.rows {
		fmt.Fprint(t.w, s.Vertical)
		for i, col := range t.columns {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			if col.MaxWidth > 0 && visibleLen(val) > col.MaxWidth {
				val = val[:col.MaxWidth-1] + "…"
			}
			cell := alignText(val, widths[i], col.Align)
			if col.Style.codes != nil {
				cell = col.Style.Apply(cell)
			}
			fmt.Fprintf(t.w, "%s%s%s%s", pad, cell, pad, s.Vertical)
		}
		fmt.Fprintln(t.w)
		if t.showLines && ri < len(t.rows)-1 {
			fmt.Fprintln(t.w, t.midLine(widths))
		}
	}

	fmt.Fprintln(t.w, t.bottomLine(widths))

	if t.caption != "" {
		fmt.Fprintf(t.w, "  %s\n", NewStyle(Dim).Apply(t.caption))
	}
}

func (t *Table) topLine(widths []int) string {
	s := t.style
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat(s.Horizontal, w+t.padding*2)
	}
	return s.TopLeft + strings.Join(parts, s.TopMiddle) + s.TopRight
}

func (t *Table) midLine(widths []int) string {
	s := t.style
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat(s.Horizontal, w+t.padding*2)
	}
	return s.MiddleLeft + strings.Join(parts, s.MiddleCross) + s.MiddleRight
}

func (t *Table) bottomLine(widths []int) string {
	s := t.style
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat(s.Horizontal, w+t.padding*2)
	}
	return s.BottomLeft + strings.Join(parts, s.BottomMiddle) + s.BottomRight
}

func alignText(text string, width int, align Alignment) string {
	l := visibleLen(text)
	if l >= width {
		return text
	}
	pad := width - l
	switch align {
	case AlignRight:
		return strings.Repeat(" ", pad) + text
	case AlignCenter:
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	default:
		return text + strings.Repeat(" ", pad)
	}
}

// visibleLen returns the display length of a string, ignoring ANSI escape sequences.
func visibleLen(s string) int {
	inEsc := false
	count := 0
	for i := 0; i < len(s); {
		if s[i] == '\033' {
			inEsc = true
		}
		if inEsc {
			if s[i] == 'm' {
				inEsc = false
			}
			i++
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		count++
		i += size
	}
	return count
}
