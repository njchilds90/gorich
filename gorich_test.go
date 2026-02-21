package gorich_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/njchilds90/gorich"
)

func TestColorize(t *testing.T) {
	result := gorich.Colorize("hello", gorich.Red, gorich.Bold)
	if !strings.Contains(result, "hello") {
		t.Error("Colorize: expected 'hello' in result")
	}
	if !strings.HasSuffix(result, gorich.Reset) {
		t.Error("Colorize: expected Reset at end")
	}
}

func TestMarkup(t *testing.T) {
	result := gorich.Markup("[bold]test[/]")
	if !strings.Contains(result, "test") {
		t.Error("Markup: expected 'test' in result")
	}
	if !strings.Contains(result, gorich.Bold) {
		t.Error("Markup: expected Bold code")
	}
}

func TestMarkupUnknownTag(t *testing.T) {
	result := gorich.Markup("[unknown]test[/]")
	if !strings.Contains(result, "test") {
		t.Error("Markup with unknown tag: expected 'test'")
	}
}

func TestStyle(t *testing.T) {
	s := gorich.NewStyle(gorich.Green, gorich.Bold)
	result := s.Apply("hi")
	if !strings.Contains(result, "hi") {
		t.Error("Style.Apply: expected 'hi'")
	}
}

func TestTableRender(t *testing.T) {
	var buf bytes.Buffer
	table := gorich.NewTable(
		gorich.WithTitle("Test Table"),
		gorich.WithWriter(&buf),
	)
	table.AddColumn("Name")
	table.AddColumn("Value", gorich.ColAlign(gorich.AlignRight))
	table.AddRow("foo", "123")
	table.AddRow("bar", "456")
	table.Render()

	out := buf.String()
	if !strings.Contains(out, "foo") {
		t.Error("Table: expected 'foo' in output")
	}
	if !strings.Contains(out, "456") {
		t.Error("Table: expected '456' in output")
	}
}

func TestPanelRender(t *testing.T) {
	var buf bytes.Buffer
	gorich.NewPanel("Hello Panel",
		gorich.PanelTitle("Test"),
		gorich.PanelWriter(&buf),
	).Render()
	out := buf.String()
	if !strings.Contains(out, "Hello Panel") {
		t.Error("Panel: expected content in output")
	}
}

func TestRuleRender(t *testing.T) {
	var buf bytes.Buffer
	gorich.NewRule(
		gorich.RuleTitle("Section"),
		gorich.RuleWidth(40),
		gorich.RuleWriter(&buf),
	).Render()
	out := buf.String()
	if !strings.Contains(out, "Section") {
		t.Error("Rule: expected title in output")
	}
}

func TestTreeRender(t *testing.T) {
	var buf bytes.Buffer
	tr := gorich.NewTree("Root", gorich.TreeWriter(&buf))
	tr.Root().Add("child1")
	child2 := tr.Root().Add("child2")
	child2.Add("grandchild")
	tr.Render()
	out := buf.String()
	if !strings.Contains(out, "Root") {
		t.Error("Tree: expected 'Root'")
	}
	if !strings.Contains(out, "grandchild") {
		t.Error("Tree: expected 'grandchild'")
	}
}

func TestProgressBar(t *testing.T) {
	var buf bytes.Buffer
	bar := gorich.NewProgressBar(10, gorich.BarWriter(&buf), gorich.BarShowETA(false))
	bar.Update(5)
	bar.Finish()
	out := buf.String()
	if !strings.Contains(out, "10/10") {
		t.Error("ProgressBar: expected '10/10'")
	}
}

func TestSyntaxHighlight(t *testing.T) {
	code := `func main() { fmt.Println("hello") }`
	result := gorich.SyntaxHighlight(code, "go")
	if !strings.Contains(result, "main") {
		t.Error("SyntaxHighlight: expected 'main' in output")
	}
}

func TestTrack(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	sum := 0
	gorich.Track(items, "Summing", func(n int) { sum += n })
	if sum != 15 {
		t.Errorf("Track: expected sum=15, got %d", sum)
	}
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	gorich.Flog(&buf, "INFO", "test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Error("Log: expected message in output")
	}
}

func TestVisibleLen(t *testing.T) {
	// test through table (indirectly) — coloring shouldn't affect widths
	var buf bytes.Buffer
	table := gorich.NewTable(gorich.WithWriter(&buf))
	table.AddColumn("Hello")
	table.AddRow(gorich.Colorize("Hello", gorich.Red)) // same visual width as "Hello"
	table.Render()
	// If widths are miscomputed, borders will be crooked — just ensure no panic
}

func TestInspect(t *testing.T) {
	// Just ensure it doesn't panic
	gorich.Inspect("myVar", map[string]int{"a": 1})
}

func TestColumns(t *testing.T) {
	// Just ensure it doesn't panic
	gorich.Columns([]string{"apple", "banana", "cherry"}, 2)
}

func TestRGB256(t *testing.T) {
	result := gorich.RGB256("hello", 196)
	if !strings.Contains(result, "hello") {
		t.Error("RGB256: expected 'hello'")
	}
}

func TestTrueColor(t *testing.T) {
	result := gorich.TrueColor("hi", 255, 128, 0)
	if !strings.Contains(result, "hi") {
		t.Error("TrueColor: expected 'hi'")
	}
}
