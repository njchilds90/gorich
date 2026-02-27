package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// TableStyle defines the box-drawing characters used for a table border.
//
// Each field is a single character string, but a longer string also works if
// a custom visual effect is desired.
type TableStyle struct {
	TopLeft      string
	TopRight     string
	BottomLeft   string
	BottomRight  string
	Horizontal   string
	Vertical     string
	MiddleLeft   string
	MiddleRight  string
	MiddleCross  string
	TopMiddle    string
	BottomMiddle string
}

// Pre-built table styles.
var (
	// TableStyleRounded uses rounded corners.
	TableStyleRounded = TableStyle{
		TopLeft: "╭", TopRight: "╮",
		BottomLeft: "╰", BottomRight: "╯",
		Horizontal: "─", Vertical: "│",
		MiddleLeft: "├", MiddleRight: "┤",
		MiddleCross: "┼", TopMiddle: "┬", BottomMiddle: "┴",
	}
	// TableStyleBox uses box corners.
	TableStyleBox = TableStyle{
		TopLeft: "┌", TopRight: "┐",
		BottomLeft: "└", BottomRight: "┘",
		Horizontal: "─", Vertical: "│",
		MiddleLeft: "├", MiddleRight: "┤",
		MiddleCross: "┼", TopMiddle: "┬", BottomMiddle: "┴",
	}
	// TableStyleDouble uses double-line borders.
	TableStyleDouble = TableStyle{
		TopLeft: "╔", TopRight: "╗",
		BottomLeft: "╚", BottomRight: "╝",
		Horizontal: "═", Vertical: "║",
		MiddleLeft: "╠", MiddleRight: "╣",
		MiddleCross: "╬", TopMiddle: "╦", BottomMiddle: "╩",
	}
	// TableStyleSimple uses ASCII-only borders.
	TableStyleSimple = TableStyle{
		TopLeft: "+", TopRight: "+",
		BottomLeft: "+", BottomRight: "+",
		Horizontal: "-", Vertical: "|",
		MiddleLeft: "+", MiddleRight: "+",
		MiddleCross: "+", TopMiddle: "+", BottomMiddle: "+",
	}
	// TableStyleMinimal uses minimal separators suitable for compact output.
	TableStyleMinimal = TableStyle{
		Horizontal: "─", Vertical: " ",
		TopLeft: " ", TopRight: " ", BottomLeft: " ", BottomRight: " ",
		MiddleLeft: " ", MiddleRight: " ", MiddleCross: "─", TopMiddle: "─", BottomMiddle: "─",
	}
)

// Alignment describes how text is aligned within a table column.
//
// Alignment is applied after any truncation and before any per-column styling.
type Alignment int

const (
	// AlignLeft aligns text to the left.
	AlignLeft Alignment = iota
	// AlignCenter centers text.
	AlignCenter
	// AlignRight aligns text to the right.
	AlignRight
)

// Column defines a table column.
//
// The Header field is displayed in the header row when Table.showHeader is true.
// MinWidth and MaxWidth constrain the visible content width for the column.
type Column struct {
	Header      string
	Style       Style
	Align       Alignment
	HeaderStyle Style
	MinWidth    int
	MaxWidth    int
}

// Table is a rich terminal table renderer.
//
// A Table is configured using functional options passed to NewTable and then
// populated by calling AddColumn and AddRow.
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

// NewTable creates a new Table writing to standard output.
//
// Use WithWriter to render the output into a buffer for testing or for use in a
// larger application.
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

// WithTitle sets the title that is rendered into the top border.
func WithTitle(title string) func(*Table) {
	return func(t *Table) { t.title = title }
}

// WithCaption sets a dimmed caption that is rendered below the table.
func WithCaption(caption string) func(*Table) {
	return func(t *Table) { t.caption = caption }
}

// WithTableStyle sets the border drawing style.
func WithTableStyle(style TableStyle) func(*Table) {
	return func(t *Table) { t.style = style }
}

// WithHeaderStyle sets the default header style.
func WithHeaderStyle(style Style) func(*Table) {
	return func(t *Table) { t.headerStyle = style }
}

// WithShowLines controls whether separator lines are rendered between data rows.
func WithShowLines(show bool) func(*Table) {
	return func(t *Table) { t.showLines = show }
}

// WithShowHeader controls whether the header row is rendered.
func WithShowHeader(show bool) func(*Table) {
	return func(t *Table) { t.showHeader = show }
}

// WithPadding sets the number of spaces around each cell.
func WithPadding(n int) func(*Table) {
	return func(t *Table) { t.padding = n }
}

// WithWriter sets the writer used for rendering.
func WithWriter(w io.Writer) func(*Table) {
	return func(t *Table) { t.w = w }
}

// AddColumn appends a column definition.
//
// Column options can be used to set alignment, styles, and width constraints.
func (t *Table) AddColumn(header string, opts ...func(*Column)) {
	col := Column{Header: header, Align: AlignLeft}
	for _, o := range opts {
		o(&col)
	}
	t.columns = append(t.columns, col)
}

// ColStyle sets the style applied to data cells in the column.
func ColStyle(s Style) func(*Column) { return func(c *Column) { c.Style = s } }

// ColAlign sets the text alignment for the column.
func ColAlign(a Alignment) func(*Column) { return func(c *Column) { c.Align = a } }

// ColMinWidth sets the minimum visible width of the column.
func ColMinWidth(n int) func(*Column) { return func(c *Column) { c.MinWidth = n } }

// ColMaxWidth sets the maximum visible width of the column.
//
// If a cell's visible content exceeds this width, it is truncated and an
// ellipsis is appended.
func ColMaxWidth(n int) func(*Column) { return func(c *Column) { c.MaxWidth = n } }

// ColHeaderStyle sets a custom style for the header cell of the column.
func ColHeaderStyle(s Style) func(*Column) { return func(c *Column) { c.HeaderStyle = s } }

// AddRow appends a data row.
//
// Values are matched to columns by position. Missing values are rendered as
// empty cells. Extra values are ignored.
func (t *Table) AddRow(values ...string) {
	t.rows = append(t.rows, values)
}

// Render prints the table to the configured writer.
func (t *Table) Render() {
	pad := strings.Repeat(" ", t.padding)

	// Compute column widths.
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

	// Apply maximum widths.
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
	totalWidth-- // last column does not add a trailing separator

	// Title.
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

	// Header.
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

	// Rows.
	for ri, row := range t.rows {
		fmt.Fprint(t.w, s.Vertical)
		for i, col := range t.columns {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			if col.MaxWidth > 0 && visibleLen(val) > col.MaxWidth {
				// Keep (MaxWidth-1) visible runes and append an ellipsis. This avoids
				// breaking UTF-8 and also keeps ANSI escape sequences intact.
				val = truncateVisibleRunesPreservingANSI(val, col.MaxWidth-1) + "…"
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

// borderLine builds a border line segment using a left cap, middle separator, and right cap.
//
// This helper exists to keep the three border line functions consistent.
func (t *Table) borderLine(widths []int, left, mid, right string) string {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat(t.style.Horizontal, w+t.padding*2)
	}
	return left + strings.Join(parts, mid) + right
}

// topLine builds the table top border line.
func (t *Table) topLine(widths []int) string {
	s := t.style
	return t.borderLine(widths, s.TopLeft, s.TopMiddle, s.TopRight)
}

// midLine builds the table middle separator line.
func (t *Table) midLine(widths []int) string {
	s := t.style
	return t.borderLine(widths, s.MiddleLeft, s.MiddleCross, s.MiddleRight)
}

// bottomLine builds the table bottom border line.
func (t *Table) bottomLine(widths []int) string {
	s := t.style
	return t.borderLine(widths, s.BottomLeft, s.BottomMiddle, s.BottomRight)
}

// alignText pads text to the specified width based on alignment.
//
// Width is measured using visibleLen, which ignores ANSI escape sequences.
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

// truncateVisibleRunesPreservingANSI truncates a string to at most max visible runes.
//
// This helper treats bytes that are part of ANSI escape sequences as non-visible
// and copies them through without counting them toward max.
//
// The ANSI handling is intentionally minimal and matches visibleLen: it skips from
// an ESC byte through a terminating 'm' byte.
func truncateVisibleRunesPreservingANSI(s string, max int) string {
	if max <= 0 {
		return ""
	}

	inEscapeSequence := false
	visibleCount := 0

	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); {
		if s[i] == '\033' {
			inEscapeSequence = true
		}

		if inEscapeSequence {
			b.WriteByte(s[i])
			if s[i] == 'm' {
				inEscapeSequence = false
			}
			i++
			continue
		}

		if visibleCount >= max {
			break
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			// Treat invalid bytes as a single visible character.
			b.WriteByte(s[i])
			i++
			visibleCount++
			continue
		}
		b.WriteString(s[i : i+size])
		i += size
		visibleCount++
	}

	return b.String()
}

// visibleLen returns the display length of a string, ignoring ANSI escape sequences.
//
// This function counts runes, not terminal cell width. It is designed to be
// predictable and dependency-free rather than perfectly accurate for wide
// characters and combining marks.
func visibleLen(s string) int {
	inEscapeSequence := false
	count := 0
	for i := 0; i < len(s); {
		if s[i] == '\033' {
			inEscapeSequence = true
		}
		if inEscapeSequence {
			if s[i] == 'm' {
				inEscapeSequence = false
			}
			i++
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		if size == 1 && s[i] >= 0x80 {
			// This is likely invalid UTF-8. Count it as one visible character.
			count++
			i++
			continue
		}
		count++
		i += size
	}
	return count
}
