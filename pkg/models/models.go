package models

import (
	"path/filepath"
	"time"
)

// Config defines the extraction and transformation settings.
type Config struct {
	StartPage             int
	EndPage               int
	SkipFrontMatterPages  int
	ExcludeAppendix       bool
	ExcludeIndex          bool
	ExcludeBibliography   bool
	ExcludeHeadersFooters bool
	StripAssets           bool
	ReflowParagraphs      bool
	DetectCodeBlocks      bool
	SplitByChapter        bool
	GenerateTOCIndex      bool
	Concurrency           int
	StopPatterns          []string
}

// DefaultConfig returns optimal production defaults.
func DefaultConfig() Config {
	return Config{
		StartPage:             1,
		EndPage:               0, // 0 = all pages
		SkipFrontMatterPages:  0,
		ExcludeAppendix:       true,
		ExcludeIndex:          true,
		ExcludeBibliography:   true,
		ExcludeHeadersFooters: true,
		StripAssets:           true,
		ReflowParagraphs:      true,
		DetectCodeBlocks:      true,
		SplitByChapter:        true,
		GenerateTOCIndex:      true,
		Concurrency:           4,
		StopPatterns: []string{
			`(?i)^appendix\s+[a-z0-9]`,
			`(?i)^appendices\b`,
			`(?i)^index\b`,
			`(?i)^subject\s+index\b`,
			`(?i)^author\s+index\b`,
			`(?i)^bibliography\b`,
			`(?i)^references\b`,
			`(?i)^glossary\b`,
		},
	}
}

// Chapter represents an individual section or chapter extracted from the PDF.
type Chapter struct {
	Index     int
	Title     string
	Slug      string
	Filename  string
	StartPage int
	Lines     []string
	Content   string
}

// SplitResult represents the multi-file TOC output.
type SplitResult struct {
	SourcePDF       string
	PDFName         string
	TargetDirectory string
	TOCContent      string
	Chapters        []Chapter
	TotalPages      int
	ProcessedPages  int
	SkippedPages    int
}

// FileResult holds individual file execution status in a concurrent batch.
type FileResult struct {
	PDFPath       string
	PDFName       string
	Success       bool
	Error         error
	ChaptersCount int
	TotalPages    int
}

// BatchResult holds aggregated batch execution metrics.
type BatchResult struct {
	BatchName       string
	TargetDirectory string
	TotalFiles      int
	ProcessedFiles  int
	FailedFiles     int
	Results         []FileResult
	Duration        time.Duration
}

// PDFPage represents an individual page extracted from the PDF.
type PDFPage struct {
	PageNumber    int
	RawText       string
	IsFilteredOut bool
	FilterReason  string
	Lines         []string
}

// ProcessedDocument represents the final transformed output.
type ProcessedDocument struct {
	SourcePath      string
	TotalPages      int
	ProcessedPages  int
	SkippedPages    int
	MarkdownContent string
}

// GetFilename returns the base name of the source PDF.
func (d *ProcessedDocument) GetFilename() string {
	return filepath.Base(d.SourcePath)
}
