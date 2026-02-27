package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Style holds a set of ANSI escape codes that are applied together.
//
// A Style is intended to be created with NewStyle and then reused for many
// strings to keep calling code readable and consistent.
type Style struct {
	codes []string
}

// NewStyle creates a new Style using the provided ANSI escape codes.
//
// The provided codes are concatenated in the order they are passed.
func NewStyle(codes ...string) Style {
	return Style{codes: codes}
}

// Apply wraps text with this Style's ANSI escape codes and a trailing reset.
//
// If the Style has no codes, Apply still returns a reset-wrapped string because
// Colorize always appends Reset.
func (s Style) Apply(text string) string {
	return Colorize(text, s.codes...)
}

// Common pre-built styles matching Python's rich defaults.
var (
	// StyleBold renders text in bold.
	StyleBold = NewStyle(Bold)
	// StyleDim renders text in a dimmed intensity.
	StyleDim = NewStyle(Dim)
	// StyleItalic renders text in italic.
	StyleItalic = NewStyle(Italic)
	// StyleUnderline renders text with an underline.
	StyleUnderline = NewStyle(Underline)
	// StyleStrike renders text with a strike-through.
	StyleStrike = NewStyle(Strike)

	// StyleSuccess renders text in a bold green style suitable for successful status messages.
	StyleSuccess = NewStyle(Bold, BrightGreen)
	// StyleWarning renders text in a bold yellow style suitable for warnings.
	StyleWarning = NewStyle(Bold, BrightYellow)
	// StyleError renders text in a bold red style suitable for errors.
	StyleError = NewStyle(Bold, BrightRed)
	// StyleInfo renders text in a bold cyan style suitable for informational messages.
	StyleInfo = NewStyle(Bold, BrightCyan)
	// StyleMuted renders text in a dim white style suitable for de-emphasized output.
	StyleMuted = NewStyle(Dim, White)
)

// Print prints styled text to standard output.
//
// If a Style is provided, the Style is applied directly.
// If no Style is provided, the text is interpreted as markup via Markup.
func Print(text string, style ...Style) {
	Fprint(os.Stdout, text, style...)
}

// Println prints styled text followed by a newline to standard output.
//
// If a Style is provided, the Style is applied directly.
// If no Style is provided, the text is interpreted as markup via Markup.
func Println(text string, style ...Style) {
	Fprintln(os.Stdout, text, style...)
}

// Fprint writes styled text to a writer.
//
// If a Style is provided, the Style is applied directly.
// If no Style is provided, the text is interpreted as markup via Markup.
func Fprint(w io.Writer, text string, style ...Style) {
	if len(style) > 0 {
		fmt.Fprint(w, style[0].Apply(text))
		return
	}
	fmt.Fprint(w, Markup(text))
}

// Fprintln writes styled text followed by a newline to a writer.
//
// If a Style is provided, the Style is applied directly.
// If no Style is provided, the text is interpreted as markup via Markup.
func Fprintln(w io.Writer, text string, style ...Style) {
	if len(style) > 0 {
		fmt.Fprintln(w, style[0].Apply(text))
		return
	}
	fmt.Fprintln(w, Markup(text))
}

// markupTagMap maps markup tag names to ANSI escape sequences.
//
// This is package-scoped to avoid rebuilding the map on every Markup call.
var markupTagMap = map[string]string{
	"bold":           Bold,
	"dim":            Dim,
	"italic":         Italic,
	"underline":      Underline,
	"strike":         Strike,
	"reverse":        Reverse,
	"blink":          Blink,
	"black":          Black,
	"red":            Red,
	"green":          Green,
	"yellow":         Yellow,
	"blue":           Blue,
	"magenta":        Magenta,
	"cyan":           Cyan,
	"white":          White,
	"bright_black":   BrightBlack,
	"bright_red":     BrightRed,
	"bright_green":   BrightGreen,
	"bright_yellow":  BrightYellow,
	"bright_blue":    BrightBlue,
	"bright_magenta": BrightMagenta,
	"bright_cyan":    BrightCyan,
	"bright_white":   BrightWhite,
	"bg_black":       BgBlack,
	"bg_red":         BgRed,
	"bg_green":       BgGreen,
	"bg_yellow":      BgYellow,
	"bg_blue":        BgBlue,
	"bg_magenta":     BgMagenta,
	"bg_cyan":        BgCyan,
	"bg_white":       BgWhite,
	"success":        Bold + BrightGreen,
	"warning":        Bold + BrightYellow,
	"error":          Bold + BrightRed,
	"info":           Bold + BrightCyan,
	"muted":          Dim + White,
}

// Markup parses simple markup tags in text and converts them into ANSI escape codes.
//
// Supported tags include:
//   - Text attributes: [bold], [dim], [italic], [underline], [strike], [reverse], [blink]
//   - Foreground colors: [red], [green], [yellow], [blue], [magenta], [cyan], [white], and their bright variants
//   - Background colors: [bg_red], [bg_green], and similar
//   - Convenience tags: [success], [warning], [error], [info], [muted]
//
// The closing tag [/] resets all formatting.
//
// Tags can be combined, for example: [bold red]text[/]. Unknown tags are preserved
// as literal text so that Markup does not unexpectedly delete user content.
func Markup(text string) string {
	var out strings.Builder
	out.Grow(len(text))

	for i := 0; i < len(text); {
		if text[i] != '[' {
			out.WriteByte(text[i])
			i++
			continue
		}

		// Attempt to parse a tag. If there is no closing bracket, treat the
		// opening bracket as literal text.
		endRel := strings.IndexByte(text[i:], ']')
		if endRel < 0 {
			out.WriteByte(text[i])
			i++
			continue
		}
		end := i + endRel
		tag := text[i+1 : end]

		// Reset tag.
		if tag == "/" {
			out.WriteString(Reset)
			i = end + 1
			continue
		}

		parts := strings.Fields(tag)
		code := ""
		matched := false
		for _, part := range parts {
			if c, ok := markupTagMap[part]; ok {
				code += c
				matched = true
			}
		}

		if matched {
			out.WriteString(code)
			i = end + 1
			continue
		}

		// Unknown tag: keep it literal and advance by one byte so that the
		// parser always makes forward progress.
		out.WriteByte(text[i])
		i++
	}

	return out.String()
}

// Sprint returns a styled string without printing.
//
// If a Style is provided, the Style is applied directly.
// If no Style is provided, the text is interpreted as markup via Markup.
func Sprint(text string, style ...Style) string {
	if len(style) > 0 {
		return style[0].Apply(text)
	}
	return Markup(text)
}
