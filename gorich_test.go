package gorich

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
)

// stripANSI removes common ANSI Select Graphic Rendition sequences.
//
// The production code intentionally only understands a minimal subset of ANSI,
// so the test helper mirrors that simplicity.
func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

func TestSyntaxHighlight_GoTrailingBackslashDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SyntaxHighlight panicked: %v", r)
		}
	}()

	// This input includes an unterminated string literal with a trailing backslash.
	// The highlighter must not panic even when input is invalid source code.
	_ = SyntaxHighlight("package main\nfunc main() { s := \"abc\\\" }\n", "go")
}

func TestSyntaxHighlight_JSONEscapedQuoteKeepsOutputIntact(t *testing.T) {
	input := `{"a":"b\"c","n":1}`
	out := SyntaxHighlight(input, "json")
	plain := stripANSI(out)
	if plain != input {
		t.Fatalf("expected highlighted JSON to preserve original content after stripping ANSI; got %q want %q", plain, input)
	}
}

func TestTable_Render_TruncatesUTF8Safely(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(
		WithWriter(&buf),
		WithShowHeader(true),
		WithPadding(1),
	)
	table.AddColumn("Word", ColMaxWidth(2))
	table.AddRow("你好")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Table.Render panicked: %v", r)
		}
	}()
	table.Render()

	plain := stripANSI(buf.String())
	if !strings.Contains(plain, "你…") {
		t.Fatalf("expected truncated output to contain a valid rune + ellipsis; got %q", plain)
	}
}

func TestTable_Render_TruncatesANSISafely(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(
		WithWriter(&buf),
		WithShowHeader(false),
	)
	table.AddColumn("Status", ColMaxWidth(1))
	table.AddRow(StyleSuccess.Apply("OK"))

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Table.Render panicked: %v", r)
		}
	}()
	table.Render()

	plain := stripANSI(buf.String())
	if !strings.Contains(plain, "…") {
		t.Fatalf("expected ANSI-styled cell to truncate to an ellipsis; got %q", plain)
	}
}

func TestSpinner_StopIsIdempotentAndEmptyFramesDoNotPanic(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner("Working",
		SpinnerFrames([]string{}),
		SpinnerWriter(&buf),
	)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Spinner panicked: %v", r)
		}
	}()

	s.Start()
	time.Sleep(10 * time.Millisecond)
	s.Stop()
	s.Stop()

	plain := stripANSI(buf.String())
	if !strings.Contains(plain, "Working") {
		t.Fatalf("expected spinner output to include description; got %q", plain)
	}
}

func TestFcolumns_WritesToProvidedWriterAndHandlesInvalidColumnCount(t *testing.T) {
	var buf bytes.Buffer
	Fcolumns(&buf, []string{"a", "b", "c"}, 0)
	if buf.Len() == 0 {
		t.Fatalf("expected output to be written to the provided writer")
	}
}
