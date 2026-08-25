package core

import (
	"path/filepath"
	"time"
)

// Config defines extraction, transformation, agentic metadata, and batch settings.
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
	GenerateAgentManifest bool
	AddYAMLFrontmatter    bool
	Concurrency           int
	StopPatterns          []string
}

// DefaultConfig returns optimal, agentic-ready production defaults.
func DefaultConfig() Config {
	return Config{
		StartPage:             1,
		EndPage:               0,
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
		GenerateAgentManifest: true,
		AddYAMLFrontmatter:    true,
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

// Chapter represents an individual section or chapter extracted from a document.
type Chapter struct {
	Index         int
	Title         string
	Slug          string
	Filename      string
	StartPage     int
	WordCount     int
	TokenEstimate int
	PrevFilename  string
	NextFilename  string
	Lines         []string
	Content       string
}

// PDFPage represents an individual page extracted from the PDF.
type PDFPage struct {
	PageNumber    int
	RawText       string
	IsFilteredOut bool
	FilterReason  string
	Lines         []string
}

// ProcessedDocument represents the single-file output.
type ProcessedDocument struct {
	SourcePath      string
	TotalPages      int
	ProcessedPages  int
	SkippedPages    int
	WordCount       int
	TokenEstimate   int
	MarkdownContent string
}

// GetFilename returns the base name of the source PDF.
func (d *ProcessedDocument) GetFilename() string {
	return filepath.Base(d.SourcePath)
}

// SplitResult represents the chapter-based directory output with agent manifests.
type SplitResult struct {
	SourcePDF       string
	PDFName         string
	TargetDirectory string
	TOCContent      string
	AgentManifest   string
	Chapters        []Chapter
	TotalPages      int
	ProcessedPages  int
	SkippedPages    int
	TotalWords      int
	TotalTokens     int
}

// FileResult holds individual file execution status in a concurrent batch.
type FileResult struct {
	PDFPath       string
	PDFName       string
	Success       bool
	Error         error
	ChaptersCount int
	TotalPages    int
	TotalTokens   int
}

// BatchResult holds aggregated metrics for multi-document batch runs.
type BatchResult struct {
	BatchName       string
	TargetDirectory string
	TotalFiles      int
	ProcessedFiles  int
	FailedFiles     int
	Results         []FileResult
	Duration        time.Duration
}
