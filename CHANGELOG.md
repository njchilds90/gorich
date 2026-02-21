# Changelog

All notable changes to gorich will be documented in this file.

## [Unreleased]

## [v0.1.0] - 2026-02-21

### Added
- `Colorize`, `RGB256`, `TrueColor`, `BgTrueColor` — ANSI color helpers
- `Style`, `NewStyle` — composable style definitions
- `Markup` — Python-rich-style `[bold red]text[/]` markup parser
- `Print`, `Println`, `Sprint`, `Fprint`, `Fprintln` — styled print functions
- Pre-built styles: `StyleSuccess`, `StyleWarning`, `StyleError`, `StyleInfo`, `StyleMuted`
- `Table` — rich tables with multiple border styles, column alignment, and headers
- `Panel` — bordered panel with title, subtitle, and padding
- `Rule` — horizontal rule/divider with optional title
- `Tree` / `TreeNode` — hierarchical tree view
- `ProgressBar` — progress bar with percentage, count, ETA
- `Track[T]` — generic progress tracking for slices
- `Spinner` — animated terminal spinner
- `WithSpinner` — convenience wrapper for spinner + function
- `SyntaxHighlight` — syntax highlighting for Go, Python, JSON
- `PrintSyntax` — highlighted code in a panel
- `Log`, `Flog` — styled log levels (INFO, SUCCESS, WARNING, ERROR, DEBUG)
- `Inspect` — pretty-print any value in a labeled panel
- `Columns` — multi-column layout
