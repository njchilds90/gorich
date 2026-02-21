package gorich

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// SyntaxTheme maps token types to styles.
type SyntaxTheme struct {
	Keyword   Style
	String    Style
	Number    Style
	Comment   Style
	Function  Style
	Type      Style
	Operator  Style
	Default   Style
}

// Pre-built themes.
var (
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

// SyntaxHighlight applies syntax highlighting to code.
// lang can be "go", "python", "json", or "": returns plain text.
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

// PrintSyntax prints syntax-highlighted code with a panel.
func PrintSyntax(code, lang string, title string, theme ...SyntaxTheme) {
	highlighted := SyntaxHighlight(code, lang, theme...)
	NewPanel(highlighted,
		PanelTitle(fmt.Sprintf(" %s ", title)),
		PanelBorderStyle(PanelStyleBox),
		PanelPadding(1),
	).Render()
}

// FprintSyntax writes syntax-highlighted code to a writer.
func FprintSyntax(w io.Writer, code, lang string) {
	fmt.Fprintln(w, SyntaxHighlight(code, lang))
}

// Simple line-by-line tokenizer for Go.
func highlightGo(code string, th SyntaxTheme) string {
	lines := strings.Split(code, "\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = highlightLine(line, goKeywords, goTypes, "//", th)
	}
	return strings.Join(result, "\n")
}

func highlightPython(code string, th SyntaxTheme) string {
	lines := strings.Split(code, "\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = highlightLine(line, pyKeywords, nil, "#", th)
	}
	return strings.Join(result, "\n")
}

func highlightJSON(code string, th SyntaxTheme) string {
	var out strings.Builder
	inString := false
	for i := 0; i < len(code); i++ {
		ch := string(code[i])
		if code[i] == '"' {
			if inString {
				out.WriteString(th.String.Apply(ch))
				inString = false
			} else {
				inString = true
				out.WriteString(th.String.Apply(ch))
			}
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
			// check for true/false/null
			for _, kw := range []string{"true", "false", "null"} {
				if strings.HasPrefix(code[i:], kw) {
					out.WriteString(th.Keyword.Apply(kw))
					i += len(kw) - 1
					goto nextJSON
				}
			}
			out.WriteString(th.Default.Apply(ch))
		default:
			if code[i] >= '0' && code[i] <= '9' || code[i] == '-' {
				j := i
				for j < len(code) && (code[j] >= '0' && code[j] <= '9' || code[j] == '.' || code[j] == 'e' || code[j] == 'E' || code[j] == '+' || code[j] == '-') {
					j++
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
	// Check for whole-line comment
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, commentMarker) {
		return th.Comment.Apply(line)
	}

	// Split at comment
	commentIdx := strings.Index(line, commentMarker)
	mainPart := line
	commentPart := ""
	if commentIdx >= 0 && !isInString(line, commentIdx) {
		mainPart = line[:commentIdx]
		commentPart = th.Comment.Apply(line[commentIdx:])
	}

	var out strings.Builder
	i := 0
	for i < len(mainPart) {
		// string literal
		if mainPart[i] == '"' || mainPart[i] == '\'' || mainPart[i] == '`' {
			quote := mainPart[i]
			j := i + 1
			for j < len(mainPart) && mainPart[j] != quote {
				if mainPart[j] == '\\' {
					j++
				}
				j++
			}
			if j < len(mainPart) {
				j++
			}
			out.WriteString(th.String.Apply(mainPart[i:j]))
			i = j
			continue
		}
		// number
		if mainPart[i] >= '0' && mainPart[i] <= '9' {
			j := i
			for j < len(mainPart) && ((mainPart[j] >= '0' && mainPart[j] <= '9') || mainPart[j] == '.' || mainPart[j] == 'x' || mainPart[j] == 'X' || mainPart[j] == '_') {
				j++
			}
			out.WriteString(th.Number.Apply(mainPart[i:j]))
			i = j
			continue
		}
		// identifier / keyword
		if isLetter(mainPart[i]) {
			j := i
			for j < len(mainPart) && (isLetter(mainPart[j]) || mainPart[j] >= '0' && mainPart[j] <= '9') {
				j++
			}
			word := mainPart[i:j]
			if keywords[word] {
				out.WriteString(th.Keyword.Apply(word))
			} else if types != nil && types[word] {
				out.WriteString(th.Type.Apply(word))
			} else {
				// Check if followed by '(' — it's a function call
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
		// operator
		switch mainPart[i] {
		case '+', '-', '*', '/', '=', '<', '>', '!', '&', '|', '^', '%', '~':
			out.WriteString(th.Operator.Apply(string(mainPart[i])))
		default:
			out.WriteString(string(mainPart[i]))
		}
		i++
	}

	return out.String() + commentPart
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isInString(line string, idx int) bool {
	count := 0
	for i := 0; i < idx; i++ {
		if line[i] == '"' && (i == 0 || line[i-1] != '\\') {
			count++
		}
	}
	return count%2 != 0
}
