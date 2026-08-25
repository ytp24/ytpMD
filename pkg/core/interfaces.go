package core

import (
	"context"
)

// Extractor defines the abstraction for extracting content from document sources.
type Extractor interface {
	ExtractFile(ctx context.Context, path string) (*ProcessedDocument, error)
	ExtractToDirectory(ctx context.Context, path string, targetBaseDir string) (*SplitResult, error)
}

// Filter defines the interface for eliminating noise, back-matter, and non-usable artifacts.
type Filter interface {
	ShouldStop(pageText string) (bool, string)
	CleanPageLines(rawText string) []string
}

// Transformer defines the abstraction for converting raw lines into structured Markdown.
type Transformer interface {
	Transform(pagesLines [][]string) string
}

// Splitter defines the abstraction for dividing documents into chapters/TOC entities.
type Splitter interface {
	SplitIntoChapters(pages []PDFPage, documentTitle string) []Chapter
	GenerateTOCIndex(documentTitle string, chapters []Chapter, totalPages int) string
}

// ProgressReporter defines the abstraction for progress feedback during batch execution.
type ProgressReporter interface {
	Increment(label string)
	Finish()
}

// BatchProcessor defines the interface for running concurrent multi-document conversions.
type BatchProcessor interface {
	ProcessBatch(
		ctx context.Context,
		pdfFiles []string,
		targetBaseDir string,
		batchName string,
		concurrency int,
		reporter ProgressReporter,
	) (*BatchResult, error)
}
