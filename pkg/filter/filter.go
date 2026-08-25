package filter

import (
	"regexp"
	"strings"

	"github.com/ytp24/ytp24/pkg/core"
)

// ContentFilter implements core.Filter interface.
type ContentFilter struct {
	config               core.Config
	compiledStopPatterns []*regexp.Regexp
	headerFooterPatterns []*regexp.Regexp
}

// NewContentFilter initializes a new ContentFilter using Factory Pattern.
func NewContentFilter(cfg core.Config) *ContentFilter {
	var stopPatterns []*regexp.Regexp
	for _, p := range cfg.StopPatterns {
		if re, err := regexp.Compile(p); err == nil {
			stopPatterns = append(stopPatterns, re)
		}
	}

	hfPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)^page\s+\d+(\s+of\s+\d+)?$`),
		regexp.MustCompile(`^\d+\s*[/|]\s*\d+$`),
		regexp.MustCompile(`^[—–-]\s*\d+\s*[—–-]$`),
		regexp.MustCompile(`^\d{1,4}$`),
		regexp.MustCompile(`(?i)^(copyright|all rights reserved|confidential)\b`),
	}

	return &ContentFilter{
		config:               cfg,
		compiledStopPatterns: stopPatterns,
		headerFooterPatterns: hfPatterns,
	}
}

// ShouldStop checks if the page triggers an Appendix or Index cutoff.
func (f *ContentFilter) ShouldStop(pageText string) (bool, string) {
	if !f.config.ExcludeAppendix && !f.config.ExcludeIndex {
		return false, ""
	}

	lines := strings.Split(pageText, "\n")
	checkedLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		checkedLines++
		if checkedLines > 8 {
			break
		}

		for _, re := range f.compiledStopPatterns {
			if re.MatchString(trimmed) {
				return true, "Matched stop pattern: '" + trimmed + "'"
			}
		}
	}
	return false, ""
}

// CleanPageLines removes headers, footers, page numbers, and binary/asset placeholders.
func (f *ContentFilter) CleanPageLines(rawText string) []string {
	lines := strings.Split(rawText, "\n")
	var cleaned []string

	imgRegex := regexp.MustCompile(`(?i)\[image:[^\]]*\]`)
	figRegex := regexp.MustCompile(`(?i)\bFigure\s+\d+([.:].*)?$`)
	nonPrintable := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			cleaned = append(cleaned, "")
			continue
		}

		if f.config.ExcludeHeadersFooters {
			isNoise := false
			for _, re := range f.headerFooterPatterns {
				if re.MatchString(trimmed) {
					isNoise = true
					break
				}
			}
			if isNoise {
				continue
			}
		}

		sanitized := nonPrintable.ReplaceAllString(line, "")

		if f.config.StripAssets {
			sanitized = imgRegex.ReplaceAllString(sanitized, "")
			sanitized = figRegex.ReplaceAllString(sanitized, "")
		}

		cleaned = append(cleaned, sanitized)
	}

	return cleaned
}
