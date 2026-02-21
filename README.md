# gorich

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/gorich.svg)](https://pkg.go.dev/github.com/njchilds90/gorich)
[![Go Report Card](https://goreportcard.com/badge/github.com/njchilds90/gorich)](https://goreportcard.com/report/github.com/njchilds90/gorich)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> **Beautiful terminal output for Go** — inspired by Python's [rich](https://github.com/Textualize/rich) library.

`gorich` gives you rich terminal output with zero dependencies: styled text, tables, panels, progress bars, spinners, trees, rules, and syntax highlighting — all from a clean, composable API. Works great for CLIs, AI agents, and developer tooling.

---

## Install
```sh
go get github.com/njchilds90/gorich
```

Requires Go 1.21+.

---

## Quick Start
```go
package main

import "github.com/njchilds90/gorich"

func main() {
    // Markup syntax (like Python's rich)
    gorich.Println("[bold green]Hello, World![/]")
    gorich.Println("[error]Something went wrong[/]")
    gorich.Println("[success]All done![/]")

    // Log levels
    gorich.Log("INFO", "Server started on :8080")
    gorich.Log("WARNING", "High memory usage")

    // Panel
    gorich.PrintPanel("Task complete!", gorich.PanelTitle("Status"))

    // Rule
    gorich.PrintRule("Section Title")
}
```

---

## Features

### Styled Text & Markup
```go
gorich.Println("[bold red]Error![/] Something went wrong.")
gorich.Println("[bold bright_cyan]gorich[/] — beautiful terminal output")

// Direct styles
gorich.Println("hello", gorich.StyleSuccess)

// Colorize
fmt.Println(gorich.Colorize("Important", gorich.Bold, gorich.BrightYellow))

// 256-color and true-color
fmt.Println(gorich.RGB256("custom", 214))
fmt.Println(gorich.TrueColor("rgb", 255, 128, 0))
```

**Supported markup tags:** `bold`, `dim`, `italic`, `underline`, `strike`, `reverse`, `blink`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `bright_red`…`bright_white`, `bg_red`…`bg_white`, `success`, `warning`, `error`, `info`, `muted`. Close with `[/]`. Combine: `[bold red]text[/]`.

---

### Table
```go
t := gorich.NewTable(
    gorich.WithTitle("Users"),
    gorich.WithShowLines(true),
    gorich.WithTableStyle(gorich.TableStyleRounded),
)
t.AddColumn("Name")
t.AddColumn("Role",  gorich.ColStyle(gorich.StyleInfo))
t.AddColumn("Score", gorich.ColAlign(gorich.AlignRight))
t.AddRow("Alice", "Admin", "9,842")
t.AddRow("Bob",   "User",  "3,210")
t.Render()
```

**Table styles:** `TableStyleRounded`, `TableStyleBox`, `TableStyleDouble`, `TableStyleSimple`, `TableStyleMinimal`

---

### Panel
```go
gorich.PrintPanel(
    "Deployment successful!\nVersion: v2.3.1",
    gorich.PanelTitle("Deploy"),
    gorich.PanelSubtitle("prod-us-east-1"),
    gorich.PanelBorderStyle(gorich.PanelStyleDouble),
    gorich.PanelBorderColor(gorich.NewStyle(gorich.BrightGreen)),
)
```

**Panel styles:** `PanelStyleRounded`, `PanelStyleBox`, `PanelStyleDouble`, `PanelStyleHeavy`, `PanelStyleSimple`

---

### Rule (Divider)
```go
gorich.PrintRule("Section Title")
gorich.NewRule(gorich.RuleChar("═")).Render()  // custom char
gorich.NewRule().Render()                       // plain rule
```

---

### Tree
```go
gorich.PrintTree("project/", func(root *gorich.TreeNode) {
    cmd := root.Add("cmd/")
    cmd.Add("main.go")
    internal := root.Add("internal/")
    internal.Add("handler.go")
    root.Add("go.mod")
})
```

---

### Progress Bar
```go
bar := gorich.NewProgressBar(100, gorich.BarDescription("Downloading"))
for i := 1; i <= 100; i++ {
    doWork()
    bar.Update(i)
}
bar.Finish()
```

**Track slices automatically:**
```go
gorich.Track(files, "Processing", func(f string) {
    processFile(f)
})
```

---

### Spinner
```go
gorich.WithSpinner("Loading data...", func() {
    fetchData()
})

// Or manually:
s := gorich.NewSpinner("Thinking...")
s.Start()
doWork()
s.Stop()
```

---

### Syntax Highlighting
```go
gorich.PrintSyntax(code, "go", "main.go")
gorich.PrintSyntax(jsonStr, "json", "config.json")
gorich.PrintSyntax(pyCode, "python", "script.py")

// Or get the string:
highlighted := gorich.SyntaxHighlight(code, "go")
```

**Themes:** `ThemeMonokai` (default), `ThemeDark`, `ThemeLight`

---

### Log
```go
gorich.Log("INFO",    "Server started")
gorich.Log("SUCCESS", "Build passed")
gorich.Log("WARNING", "Disk usage high")
gorich.Log("ERROR",   "Connection failed")
gorich.Log("DEBUG",   "Request: GET /api")
```

---

### Inspect
```go
gorich.Inspect("Config", myStruct)
```

---

### Columns
```go
gorich.Columns([]string{"fmt", "os", "io", "net", "http", "json"}, 3)
```

---

## Writing to Any `io.Writer`

Every component accepts a writer option — great for AI agents, tests, and logging pipelines:
```go
var buf bytes.Buffer

table := gorich.NewTable(gorich.WithWriter(&buf))
gorich.NewPanel("...", gorich.PanelWriter(&buf))
gorich.NewRule(gorich.RuleWriter(&buf))
bar := gorich.NewProgressBar(10, gorich.BarWriter(&buf))
```

---

## Run the Demo
```sh
go run ./examples/main.go
```

---

## GoDoc

Full API documentation: [pkg.go.dev/github.com/njchilds90/gorich](https://pkg.go.dev/github.com/njchilds90/gorich)

---

## License

MIT © 2026 njchilds90
