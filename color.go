// Package gorich provides rich terminal output utilities for Go,
// inspired by Python's rich library. It supports styled text,
// tables, panels, progress bars, trees, rules, and syntax highlighting.
package gorich

import "fmt"

// ANSI escape codes for colors and styles.
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Strike    = "\033[9m"

	// Foreground colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright foreground colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Colorize wraps text with ANSI color/style codes and a reset.
// Multiple codes can be combined: gorich.Colorize("hello", gorich.Bold, gorich.Red)
func Colorize(text string, codes ...string) string {
	prefix := ""
	for _, c := range codes {
		prefix += c
	}
	return fmt.Sprintf("%s%s%s", prefix, text, Reset)
}

// RGB256 returns an ANSI escape for a 256-color foreground.
// n must be in range 0–255.
func RGB256(text string, n int) string {
	return fmt.Sprintf("\033[38;5;%dm%s%s", n, text, Reset)
}

// BgRGB256 returns an ANSI escape for a 256-color background.
func BgRGB256(text string, n int) string {
	return fmt.Sprintf("\033[48;5;%dm%s%s", n, text, Reset)
}

// TrueColor returns ANSI true-color (24-bit) escape for a foreground.
func TrueColor(text string, r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s%s", r, g, b, text, Reset)
}

// BgTrueColor returns ANSI true-color (24-bit) escape for a background.
func BgTrueColor(text string, r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s%s", r, g, b, text, Reset)
}
