package gorich

import (
	"fmt"
	"io"
	"strings"
)

// SyntaxTheme maps token categories to styles.
//
// A SyntaxTheme is used by SyntaxHighlight and can be replaced to customize
// how code is displayed.
type SyntaxTheme struct {
	Keyword  Style
	String   Style
	Number   Style
	Comment  Style
	Function Style
	Type     Style
	Operator Style
	Default  Style
}

// Pre-built themes.
var (
	// ThemeMonokai is the default theme, tuned for dark terminals.
	ThemeMonokai = SyntaxTheme{
		Keyword:  NewStyle(BrightMagenta, Bold),
		String:   NewStyle(BrightYellow),
		Number:   NewStyle(BrightCyan),
		Comment:  NewStyle(Dim, Green),
		Function: NewStyle(BrightGreen),
		Type:     NewStyle(BrightCyan),
		Operator: NewStyle(BrightRed),
		Default:  NewStyle(White),
	}
	// ThemeDark is an alternative theme for dark terminals with different emphasis.
	ThemeDark = SyntaxTheme{
		Keyword:  NewStyle(BrightBlue, Bold),
		String:   NewStyle(BrightGreen),
		Number:   NewStyle(BrightCyan),
		Comment:  NewStyle(Dim, BrightBlack),
		Function: NewStyle(BrightYellow),
		Type:     NewStyle(BrightMagenta),
		Operator: NewStyle(BrightWhite),
		Default:  NewStyle(White),
	}
	// ThemeLight is an alternative theme intended for light terminals.
	ThemeLight = SyntaxTheme{
		Keyword:  NewStyle(Blue, Bold),
		String:   NewStyle(Green),
		Number:   NewStyle(Cyan),
		Comment:  NewStyle(Dim, BrightBlack),
		Function: NewStyle(Magenta),
		Type:     NewStyle(Blue),
		Operator: NewStyle(Red),
		Default:  NewStyle(Black),
	}
)

var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
	"true": true, "false": true, "nil": true, "iota": true,
}

var goTypes = map[string]bool{
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"float32": true, "float64": true, "complex64": true, "complex128": true,
	"bool": true, "byte": true, "rune": true, "string": true, "error": true,
	"uintptr": true, "any": true,
}

var pyKeywords = map[string]bool{
	"False": true, "None": true, "True": true, "and": true, "as": true,
	"assert": true, "async": true, "await": true, "break": true, "class": true,
	"continue": true, "def": true, "del": true, "elif": true, "else": true,
	"except": true, "finally": true, "for": true, "from": true, "global": true,
	"if": true, "import": true, "in": true, "is": true, "lambda": true,
	"nonlocal": true, "not": true, "or": true, "pass": true, "raise": true,
	"return": true, "try": true, "while": true, "with": true, "yield": true,
}

// SyntaxHighlight applies lightweight syntax highlighting to code.
//
// The lang parameter may be "go", "python" (or "py"), "json", or "".
// Any unknown language returns the input code unchanged.
//
// This implementation is intentionally simple and designed for terminal display,
// not for full language correctness.
func SyntaxHighlight(code, lang string, theme ...SyntaxTheme) string {
	th := ThemeMonokai
	if len(theme) > 0 {
		th = theme[0]
	}

	switch strings.ToLower(lang) {
	case "go":
		return highlightGo(code, th)
	case "python", "py":
		return highlightPython(code, th)
	case "json":
		return highlightJSON(code, th)
	default:
		return code
	}
}

// PrintSyntax prints syntax-highlighted code inside a panel.
//
// This is a convenience wrapper intended for quick terminal output.
func PrintSyntax(code, lang string, title string, theme ...SyntaxTheme) {
	highlighted := SyntaxHighlight(code, lang, theme...)
	NewPanel(highlighted,
		PanelTitle(fmt.Sprintf(" %s ", title)),
		PanelBorderStyle(PanelStyleBox),
		PanelPadding(1),
	).Render()
}

// FprintSyntax writes syntax-highlighted code followed by a newline to a writer.
func FprintSyntax(w io.Writer, code, lang string) {
	fmt.Fprintln(w, SyntaxHighlight(code, lang))
}

// highlightGo highlights Go code line-by-line.
func highlightGo(code string, th SyntaxTheme) string {
	lines := strings.Split(code, "\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = highlightLine(line, goKeywords, goTypes, "//", th)
	}
	return strings.Join(result, "\n")
}

// highlightPython highlights Python code line-by-line.
func highlightPython(code string, th SyntaxTheme) string {
	lines := strings.Split(code, "\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = highlightLine(line, pyKeywords, nil, "#", th)
	}
	return strings.Join(result, "\n")
}

// highlightJSON highlights JSON by styling strings, numbers, keywords, and punctuation.
//
// This function is byte-based and intentionally minimal. It includes special handling
// so that escaped quotes (\") do not terminate strings early.
func highlightJSON(code string, th SyntaxTheme) string {
	var out strings.Builder
	out.Grow(len(code))

	inString := false
	for i := 0; i < len(code); i++ {
		b := code[i]
		ch := string(b)

		if b == '"' {
			// Determine whether this quote is escaped by counting preceding backslashes.
			backslashes := 0
			for j := i - 1; j >= 0 && code[j] == '\\'; j-- {
				backslashes++
			}
			isEscaped := backslashes%2 == 1
			if !isEscaped {
				inString = !inString
			}
			out.WriteString(th.String.Apply("\""))
			continue
		}

		if inString {
			out.WriteString(th.String.Apply(ch))
			continue
		}

		switch ch {
		case ":", ",", "{", "}", "[", "]":
			out.WriteString(th.Operator.Apply(ch))
		case "t", "f", "n":
			// Check for true/false/null.
			for _, kw := range []string{"true", "false", "null"} {
				if strings.HasPrefix(code[i:], kw) {
					out.WriteString(th.Keyword.Apply(kw))
					i += len(kw) - 1
					goto nextJSON
				}
			}
			out.WriteString(th.Default.Apply(ch))
		default:
			if (b >= '0' && b <= '9') || b == '-' {
				j := i
				for j < len(code) {
					c := code[j]
					if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
						j++
						continue
					}
					break
				}
				out.WriteString(th.Number.Apply(code[i:j]))
				i = j - 1
			} else {
				out.WriteString(th.Default.Apply(ch))
			}
		}

	nextJSON:
	}
	return out.String()
}

func highlightLine(line string, keywords, types map[string]bool, commentMarker string, th SyntaxTheme) string {
	// Check for a whole-line comment.
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, commentMarker) {
		return th.Comment.Apply(line)
	}

	// Split at the first comment marker that is not within a quoted string.
	commentIdx := strings.Index(line, commentMarker)
	mainPart := line
	commentPart := ""
	if commentIdx >= 0 && !isInString(line, commentIdx) {
		mainPart = line[:commentIdx]
		commentPart = th.Comment.Apply(line[commentIdx:])
	}

	var out strings.Builder
	out.Grow(len(line))

	i := 0
	for i < len(mainPart) {
		// String literal.
		if mainPart[i] == '"' || mainPart[i] == '\'' || mainPart[i] == '`' {
			quote := mainPart[i]
			j := i + 1

			for j < len(mainPart) {
				if mainPart[j] == quote {
					j++
					break
				}

				// Raw string literals do not use escapes.
				if quote != '`' && mainPart[j] == '\\' {
					// If there is an escaped character available, consume both.
					if j+1 < len(mainPart) {
						j += 2
						continue
					}
					// Trailing backslash. Consume it and stop so slicing stays in bounds.
					j++
					break
				}

				j++
			}

			out.WriteString(th.String.Apply(mainPart[i:j]))
			i = j
			continue
		}

		// Number.
		if mainPart[i] >= '0' && mainPart[i] <= '9' {
			j := i
			for j < len(mainPart) && ((mainPart[j] >= '0' && mainPart[j] <= '9') || mainPart[j] == '.' || mainPart[j] == 'x' || mainPart[j] == 'X' || mainPart[j] == '_' ) {
				j++
			}
			out.WriteString(th.Number.Apply(mainPart[i:j]))
			i = j
			continue
		}

		// Identifier, keyword, or type.
		if isLetter(mainPart[i]) {
			j := i
			for j < len(mainPart) && (isLetter(mainPart[j]) || (mainPart[j] >= '0' && mainPart[j] <= '9')) {
				j++
			}
			word := mainPart[i:j]

			if keywords[word] {
				out.WriteString(th.Keyword.Apply(word))
			} else if types != nil && types[word] {
				out.WriteString(th.Type.Apply(word))
			} else {
				// If followed by '(', treat it as a function call.
				rest := strings.TrimLeft(mainPart[j:], " ")
				if len(rest) > 0 && rest[0] == '(' {
					out.WriteString(th.Function.Apply(word))
				} else {
					out.WriteString(th.Default.Apply(word))
				}
			}

			i = j
			continue
		}

		// Operator.
		switch mainPart[i] {
		case '+', '-', '*', '/', '=', '<', '>', '!', '&', '|', '^', '%', '~':
			out.WriteString(th.Operator.Apply(string(mainPart[i])))
		default:
			// Preserve whitespace and punctuation as-is.
			out.WriteByte(mainPart[i])
		}
		i++
	}

	return out.String() + commentPart
}

// isLetter reports whether c is a valid identifier starting character for the tokenizers.
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isInString reports whether the given byte index is within a double-quoted string.
//
// This helper is used to avoid splitting on comment markers that appear inside strings.
func isInString(line string, idx int) bool {
	inString := false
	for i := 0; i < idx && i < len(line); i++ {
		if line[i] != '"' {
			continue
		}

		// Determine whether this quote is escaped by counting preceding backslashes.
		backslashes := 0
		for j := i - 1; j >= 0 && line[j] == '\\'; j-- {
			backslashes++
		}
		if backslashes%2 == 1 {
			continue
		}
		inString = !inString
	}
	return inString
}
