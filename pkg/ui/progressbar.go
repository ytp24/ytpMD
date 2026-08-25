package ui

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// ProgressBar manages a thread-safe, animated in-place terminal progress bar
// with teal shades shifting from dark to bright near completion.
type ProgressBar struct {
	mu      sync.Mutex
	total   int
	current int
	title   string
	width   int
}

// NewProgressBar creates a new progress bar instance.
func NewProgressBar(total int, title string) *ProgressBar {
	return &ProgressBar{
		total: total,
		title: title,
		width: 22, // Compact width to guarantee single-line rendering without wrap
	}
}

// Increment advances the progress bar by 1 and updates the terminal in-place.
func (p *ProgressBar) Increment(currentFile string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current++
	if p.current > p.total {
		p.current = p.total
	}
	p.render(currentFile)
}

// Set updates the progress to a specific value.
func (p *ProgressBar) Set(current int, currentFile string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = current
	if p.current > p.total {
		p.current = p.total
	}
	p.render(currentFile)
}

// Finish completes the progress bar and prints a single clean newline.
func (p *ProgressBar) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = p.total
	p.render("Complete")
	fmt.Println()
}

func (p *ProgressBar) render(statusText string) {
	if p.total <= 0 {
		return
	}

	percent := float64(p.current) / float64(p.total)
	filled := int(percent * float64(p.width))
	if filled > p.width {
		filled = p.width
	}
	empty := p.width - filled

	// Color gradient shifting from dark teal to glowing pale teal
	color := getTealShade(percent)

	barFilled := strings.Repeat("█", filled)
	barEmpty := strings.Repeat("░", empty)

	// Clean & strictly truncate filename to max 24 chars to avoid terminal wrapping
	cleanStatus := strings.TrimSpace(statusText)
	if len(cleanStatus) > 24 {
		cleanStatus = cleanStatus[:21] + "..."
	}

	// \r\033[2K returns to column 0 AND clears the entire current line
	fmt.Printf("\r\033[2K%s%s %s[%s%s%s%s] %3.0f%% (%d/%d)%s %s| %s%s",
		TealLight, p.title,
		color,
		color, barFilled, ColorGray, barEmpty,
		percent*100, p.current, p.total,
		Reset,
		ColorGray, cleanStatus, Reset)

	_ = os.Stdout.Sync()
}

// getTealShade returns the ANSI color code matching the current completion percentage.
func getTealShade(percent float64) string {
	switch {
	case percent >= 0.95:
		return "\033[38;2;94;234;212m" // Pale Glowing Teal (Final stage)
	case percent >= 0.75:
		return "\033[38;2;45;212;191m" // Mint Teal (Near completion)
	case percent >= 0.50:
		return "\033[38;2;20;184;166m" // Bright Teal (Halfway)
	case percent >= 0.25:
		return "\033[38;2;13;148;136m" // Classic Teal (In progress)
	default:
		return "\033[38;2;15;118;110m" // Deep Dark Teal (Starting)
	}
}
