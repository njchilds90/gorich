//go:build ignore

// Run with: go run ./examples/main.go
package main

import (
	"time"

	"github.com/njchilds90/gorich"
)

func main() {
	// ── Markup & Styled Text ──────────────────────────────────────────
	gorich.PrintRule("gorich Demo", gorich.RuleStyle(gorich.NewStyle(gorich.Bold, gorich.BrightMagenta)))

	gorich.Println("[bold bright_cyan]gorich[/] — Beautiful terminal output for Go")
	gorich.Println("[success]✓ All systems operational[/]")
	gorich.Println("[warning]⚠ Disk usage above 80%[/]")
	gorich.Println("[error]✗ Connection refused[/]")
	gorich.Println("[muted]debug: verbose output here[/]")
	println()

	// ── Log ──────────────────────────────────────────────────────────
	gorich.PrintRule("Logging")
	gorich.Log("INFO", "Starting server on :8080")
	gorich.Log("SUCCESS", "Server started successfully")
	gorich.Log("WARNING", "High memory usage detected")
	gorich.Log("ERROR", "Database connection failed")
	gorich.Log("DEBUG", "Request received: GET /api/v1/users")
	println()

	// ── Panel ─────────────────────────────────────────────────────────
	gorich.PrintRule("Panel")
	gorich.PrintPanel(
		"This is a [bold]gorich[/] panel.\nIt supports [cyan]multiple lines[/] and [yellow]markup[/].",
		gorich.PanelTitle("Welcome"),
		gorich.PanelBorderStyle(gorich.PanelStyleRounded),
		gorich.PanelBorderColor(gorich.NewStyle(gorich.BrightBlue)),
	)
	println()

	// ── Table ─────────────────────────────────────────────────────────
	gorich.PrintRule("Table")
	t := gorich.NewTable(
		gorich.WithTitle("Go Modules"),
		gorich.WithShowLines(true),
	)
	t.AddColumn("Module", gorich.ColStyle(gorich.NewStyle(gorich.BrightCyan)))
	t.AddColumn("Version", gorich.ColAlign(gorich.AlignCenter))
	t.AddColumn("License", gorich.ColStyle(gorich.StyleSuccess))
	t.AddColumn("Stars", gorich.ColAlign(gorich.AlignRight), gorich.ColStyle(gorich.NewStyle(gorich.BrightYellow)))
	t.AddRow("gin-gonic/gin", "v1.9.1", "MIT", "75k")
	t.AddRow("go-chi/chi", "v5.1.0", "MIT", "17k")
	t.AddRow("rs/zerolog", "v1.32.0", "MIT", "10k")
	t.AddRow("spf13/cobra", "v1.8.0", "Apache-2", "37k")
	t.Render()
	println()

	// ── Tree ──────────────────────────────────────────────────────────
	gorich.PrintRule("Tree")
	gorich.PrintTree("myapp/", func(root *gorich.TreeNode) {
		cmd := root.Add("cmd/", gorich.NewStyle(gorich.BrightBlue))
		cmd.Add("main.go")
		cmd.Add("server.go")
		internal := root.Add("internal/", gorich.NewStyle(gorich.BrightBlue))
		handlers := internal.Add("handlers/")
		handlers.Add("user.go")
		handlers.Add("auth.go")
		internal.Add("models/").Add("user.go")
		root.Add("go.mod")
		root.Add("go.sum")
		root.Add("README.md")
	})
	println()

	// ── Rule ──────────────────────────────────────────────────────────
	gorich.PrintRule("Rules")
	gorich.PrintRule("Section One")
	gorich.NewRule().Render() // plain rule
	gorich.PrintRule("Section Two", gorich.RuleChar("═"))
	println()

	// ── Progress Bar ─────────────────────────────────────────────────
	gorich.PrintRule("Progress Bar")
	bar := gorich.NewProgressBar(20, gorich.BarDescription("Processing"))
	for i := 1; i <= 20; i++ {
		time.Sleep(50 * time.Millisecond)
		bar.Update(i)
	}
	bar.Finish()
	println()

	// ── Track ─────────────────────────────────────────────────────────
	gorich.PrintRule("Track")
	items := []string{"alpha.go", "beta.go", "gamma.go", "delta.go", "epsilon.go"}
	gorich.Track(items, "Compiling", func(f string) {
		time.Sleep(80 * time.Millisecond)
	})
	println()

	// ── Spinner ───────────────────────────────────────────────────────
	gorich.PrintRule("Spinner")
	gorich.WithSpinner("Fetching data...", func() {
		time.Sleep(1200 * time.Millisecond)
	})
	println()

	// ── Syntax Highlight ─────────────────────────────────────────────
	gorich.PrintRule("Syntax Highlighting")
	gorich.PrintSyntax(`package main

import (
	"fmt"
	"net/http"
)

// Handler responds to HTTP requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	code := 200
	fmt.Fprintf(w, "status: %d", code)
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}`, "go", "main.go")
	println()

	gorich.PrintSyntax(`{
  "name": "gorich",
  "version": "1.0.0",
  "active": true,
  "stars": 9999,
  "tags": ["go", "terminal", "rich"]
}`, "json", "config.json")
	println()

	// ── Inspect ───────────────────────────────────────────────────────
	gorich.PrintRule("Inspect")
	gorich.Inspect("Config", struct {
		Host string
		Port int
		SSL  bool
	}{"localhost", 8080, true})
	println()

	// ── Columns ───────────────────────────────────────────────────────
	gorich.PrintRule("Columns")
	gorich.Columns([]string{
		"fmt", "os", "io", "net", "http",
		"json", "time", "sync", "math", "sort",
		"strings", "bytes", "errors", "context", "reflect",
	}, 5, gorich.NewStyle(gorich.BrightCyan))
	println()

	// ── Colors ────────────────────────────────────────────────────────
	gorich.PrintRule("Colors")
	colors := []struct{ name, code string }{
		{"Red", gorich.Red}, {"Green", gorich.Green}, {"Yellow", gorich.Yellow},
		{"Blue", gorich.Blue}, {"Magenta", gorich.Magenta}, {"Cyan", gorich.Cyan},
		{"BrightRed", gorich.BrightRed}, {"BrightGreen", gorich.BrightGreen},
		{"BrightYellow", gorich.BrightYellow}, {"BrightBlue", gorich.BrightBlue},
	}
	for _, c := range colors {
		gorich.Println(gorich.Colorize("  "+c.name+"  ", c.code))
	}

	gorich.PrintRule("Done", gorich.RuleStyle(gorich.StyleSuccess))
}
