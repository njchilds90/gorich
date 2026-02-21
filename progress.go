package gorich

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

// ProgressBar is a terminal progress bar.
type ProgressBar struct {
	total       int
	current     int
	width       int
	description string
	style       Style
	doneStyle   Style
	showPercent bool
	showCount   bool
	showETA     bool
	barChar     string
	emptyChar   string
	startTime   time.Time
	mu          sync.Mutex
	w           io.Writer
	complete    bool
}

// NewProgressBar creates a ProgressBar with the given total.
func NewProgressBar(total int, opts ...func(*ProgressBar)) *ProgressBar {
	p := &ProgressBar{
		total:       total,
		width:       40,
		style:       NewStyle(BrightGreen),
		doneStyle:   NewStyle(Bold, BrightGreen),
		showPercent: true,
		showCount:   true,
		showETA:     true,
		barChar:     "█",
		emptyChar:   "░",
		startTime:   time.Now(),
		w:           os.Stdout,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Progress bar option setters.

func BarDescription(desc string) func(*ProgressBar) { return func(p *ProgressBar) { p.description = desc } }
func BarWidth(n int) func(*ProgressBar)             { return func(p *ProgressBar) { p.width = n } }
func BarStyle(s Style) func(*ProgressBar)           { return func(p *ProgressBar) { p.style = s } }
func BarChar(c string) func(*ProgressBar)           { return func(p *ProgressBar) { p.barChar = c } }
func BarEmptyChar(c string) func(*ProgressBar)      { return func(p *ProgressBar) { p.emptyChar = c } }
func BarShowPercent(b bool) func(*ProgressBar)      { return func(p *ProgressBar) { p.showPercent = b } }
func BarShowCount(b bool) func(*ProgressBar)        { return func(p *ProgressBar) { p.showCount = b } }
func BarShowETA(b bool) func(*ProgressBar)          { return func(p *ProgressBar) { p.showETA = b } }
func BarWriter(w io.Writer) func(*ProgressBar)      { return func(p *ProgressBar) { p.w = w } }

// Update sets the current progress to n and re-renders.
func (p *ProgressBar) Update(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = n
	if p.current > p.total {
		p.current = p.total
	}
	p.render()
}

// Advance increments progress by delta and re-renders.
func (p *ProgressBar) Advance(delta int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current += delta
	if p.current > p.total {
		p.current = p.total
	}
	p.render()
}

// Finish completes the progress bar at 100%.
func (p *ProgressBar) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = p.total
	p.complete = true
	p.render()
	fmt.Fprintln(p.w)
}

func (p *ProgressBar) render() {
	var pct float64
	if p.total > 0 {
		pct = float64(p.current) / float64(p.total)
	}
	filled := int(math.Round(pct * float64(p.width)))
	if filled > p.width {
		filled = p.width
	}

	st := p.style
	if p.complete {
		st = p.doneStyle
	}
	bar := st.Apply(strings.Repeat(p.barChar, filled)) + NewStyle(Dim).Apply(strings.Repeat(p.emptyChar, p.width-filled))

	line := ""
	if p.description != "" {
		line += p.description + " "
	}
	line += "[" + bar + "]"
	if p.showPercent {
		line += fmt.Sprintf(" %3.0f%%", pct*100)
	}
	if p.showCount {
		line += fmt.Sprintf(" %d/%d", p.current, p.total)
	}
	if p.showETA && !p.complete && p.current > 0 {
		elapsed := time.Since(p.startTime)
		totalEst := time.Duration(float64(elapsed) / pct)
		eta := totalEst - elapsed
		line += fmt.Sprintf(" ETA %s", fmtDuration(eta))
	}

	fmt.Fprintf(p.w, "\r%s", line)
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// Track is a convenience function that runs fn for each item in items,
// advancing a progress bar automatically.
func Track[T any](items []T, description string, fn func(T)) {
	bar := NewProgressBar(len(items), BarDescription(description))
	for i, item := range items {
		fn(item)
		bar.Update(i + 1)
	}
	bar.Finish()
}
