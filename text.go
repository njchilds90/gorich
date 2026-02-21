package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Style holds a set of ANSI codes to apply together.
type Style struct {
	codes []string
}

// NewStyle creates a new Style with the given ANSI codes.
func NewStyle(codes ...string) Style {
	return Style{codes: codes}
}

// Apply returns the text wrapped in this style's ANSI codes.
func (s Style) Apply(text string) string {
	return Colorize(text, s.codes...)
}

// Common pre-built styles matching Python's rich defaults.
var (
	StyleBold      = NewStyle(Bold)
	StyleDim       = NewStyle(Dim)
	StyleItalic    = NewStyle(Italic)
	StyleUnderline = NewStyle(Underline)
	StyleStrike    = NewStyle(Strike)

	StyleSuccess = NewStyle(Bold, BrightGreen)
	StyleWarning = NewStyle(Bold, BrightYellow)
	StyleError   = NewStyle(Bold, BrightRed)
	StyleInfo    = NewStyle(Bold, BrightCyan)
	StyleMuted   = NewStyle(Dim, White)
)

// Print prints styled text to stdout (like Python's rich.print).
// Supports markup like [bold red]text[/] — see Markup().
func Print(text string, style ...Style) {
	Fprint(os.Stdout, text, style...)
}

// Println prints styled text followed by a newline.
func Println(text string, style ...Style) {
	Fprintln(os.Stdout, text, style...)
}

// Fprint writes styled text to a writer.
func Fprint(w io.Writer, text string, style ...Style) {
	if len(style) > 0 {
		fmt.Fprint(w, style[0].Apply(text))
	} else {
		fmt.Fprint(w, Markup(text))
	}
}

// Fprintln writes styled text + newline to a writer.
func Fprintln(w io.Writer, text string, style ...Style) {
	if len(style) > 0 {
		fmt.Fprintln(w, style[0].Apply(text))
	} else {
		fmt.Fprintln(w, Markup(text))
	}
}

// Markup parses simple markup tags in text and applies ANSI codes.
// Supported tags: [bold], [dim], [italic], [underline], [strike],
// [red], [green], [yellow], [blue], [magenta], [cyan], [white],
// [bright_red], [bright_green], [bright_yellow], [bright_blue],
// [bright_magenta], [bright_cyan], [bright_white],
// [success], [warning], [error], [info], [muted],
// Closing tag: [/] resets all.
// Tags can be combined: [bold red]text[/]
func Markup(text string) string {
	tagMap := map[string]string{
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

	result := text
	for {
		start := strings.Index(result, "[")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "]")
		if end == -1 {
			break
		}
		end += start
		tag := result[start+1 : end]

		if tag == "/" {
			result = result[:start] + Reset + result[end+1:]
			continue
		}

		// split compound tags like "bold red"
		parts := strings.Fields(tag)
		code := ""
		matched := false
		for _, p := range parts {
			if c, ok := tagMap[p]; ok {
				code += c
				matched = true
			}
		}
		if matched {
			result = result[:start] + code + result[end+1:]
		} else {
			// not a known tag, skip it to avoid infinite loop
			result = result[:start] + result[start+1:]
		}
	}
	return result
}

// Sprint returns styled string without printing.
func Sprint(text string, style ...Style) string {
	if len(style) > 0 {
		return style[0].Apply(text)
	}
	return Markup(text)
}
