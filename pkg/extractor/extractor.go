package extractor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/devops/pdf2md/pkg/core"
	"github.com/devops/pdf2md/pkg/filter"
	"github.com/devops/pdf2md/pkg/splitter"
	"github.com/devops/pdf2md/pkg/transformer"
	"github.com/devops/pdf2md/pkg/validator"
)

// PDFExtractor implements the core.Extractor interface.
type PDFExtractor struct {
	config      core.Config
	filter      core.Filter
	transformer core.Transformer
	splitter    core.Splitter
}

// NewPDFExtractor initializes a new PDFExtractor using Factory Pattern.
func NewPDFExtractor(cfg core.Config) *PDFExtractor {
	return &PDFExtractor{
		config:      cfg,
		filter:      filter.NewContentFilter(cfg),
		transformer: transformer.NewTransformer(cfg),
		splitter:    splitter.NewSplitter(cfg),
	}
}

// ExtractToDirectory extracts all chapters as individual markdown files inside a directory named after the PDF.
func (e *PDFExtractor) ExtractToDirectory(ctx context.Context, pdfPath string, targetBaseDir string) (result *core.SplitResult, finalErr error) {
	defer func() {
		if r := recover(); r != nil {
			finalErr = fmt.Errorf("internal processing error: %v", r)
		}
	}()

	absPath := validator.ExpandPath(pdfPath)
	if err := validator.ValidatePDFFile(absPath); err != nil {
		return nil, err
	}

	cleanBaseDir := validator.ExpandPath(targetBaseDir)
	if err := validator.ValidateDirectory(cleanBaseDir); err != nil {
		return nil, err
	}

	pdfBaseName := strings.TrimSuffix(filepath.Base(absPath), filepath.Ext(absPath))
	targetDir := filepath.Join(cleanBaseDir, pdfBaseName)

	totalPages, _ := e.getPageCount(ctx, absPath)
	if totalPages <= 0 {
		totalPages = 1
	}

	pages, err := e.extractPages(ctx, absPath)
	if err != nil {
		return nil, err
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no content could be extracted from '%s'", filepath.Base(absPath))
	}

	var processedPages []core.PDFPage
	processedCount := 0
	totalCharsExtracted := 0

	for i := range pages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		page := &pages[i]

		if page.PageNumber <= e.config.SkipFrontMatterPages {
			page.IsFilteredOut = true
			page.FilterReason = fmt.Sprintf("Skipped front matter page <= %d", e.config.SkipFrontMatterPages)
			continue
		}

		if page.PageNumber < e.config.StartPage {
			page.IsFilteredOut = true
			continue
		}
		if e.config.EndPage > 0 && page.PageNumber > e.config.EndPage {
			page.IsFilteredOut = true
			break
		}

		if shouldStop, reason := e.filter.ShouldStop(page.RawText); shouldStop {
			page.IsFilteredOut = true
			page.FilterReason = reason
			break
		}

		page.Lines = e.filter.CleanPageLines(page.RawText)
		for _, l := range page.Lines {
			totalCharsExtracted += len(strings.TrimSpace(l))
		}

		processedPages = append(processedPages, *page)
		processedCount++
	}

	if totalCharsExtracted < 10 {
		return nil, fmt.Errorf("PDF contains no readable text layers (likely an image-only scanned document). OCR preprocessing is required.")
	}

	chapters := e.splitter.SplitIntoChapters(processedPages, pdfBaseName)
	tocIndex := e.splitter.GenerateTOCIndex(pdfBaseName, chapters, totalPages)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	for _, ch := range chapters {
		chPath := filepath.Join(targetDir, ch.Filename)
		if err := os.WriteFile(chPath, []byte(ch.Content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write chapter file %s: %w", chPath, err)
		}
	}

	readmePath := filepath.Join(targetDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(tocIndex), 0644); err != nil {
		return nil, fmt.Errorf("failed to write README.md: %w", err)
	}

	return &core.SplitResult{
		SourcePDF:       absPath,
		PDFName:         pdfBaseName,
		TargetDirectory: targetDir,
		TOCContent:      tocIndex,
		Chapters:        chapters,
		TotalPages:      totalPages,
		ProcessedPages:  processedCount,
		SkippedPages:    totalPages - processedCount,
	}, nil
}

// ExtractFile processes a single PDF file and returns a single concatenated markdown document.
func (e *PDFExtractor) ExtractFile(ctx context.Context, pdfPath string) (doc *core.ProcessedDocument, finalErr error) {
	defer func() {
		if r := recover(); r != nil {
			finalErr = fmt.Errorf("internal processing error: %v", r)
		}
	}()

	absPath := validator.ExpandPath(pdfPath)
	if err := validator.ValidatePDFFile(absPath); err != nil {
		return nil, err
	}

	totalPages, _ := e.getPageCount(ctx, absPath)
	if totalPages <= 0 {
		totalPages = 1
	}

	pages, err := e.extractPages(ctx, absPath)
	if err != nil {
		return nil, err
	}

	var cleanedPagesLines [][]string
	processedCount := 0

	for _, page := range pages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if page.PageNumber <= e.config.SkipFrontMatterPages {
			continue
		}
		if page.PageNumber < e.config.StartPage {
			continue
		}
		if e.config.EndPage > 0 && page.PageNumber > e.config.EndPage {
			break
		}
		if shouldStop, _ := e.filter.ShouldStop(page.RawText); shouldStop {
			break
		}

		lines := e.filter.CleanPageLines(page.RawText)
		cleanedPagesLines = append(cleanedPagesLines, lines)
		processedCount++
	}

	markdownContent := e.transformer.Transform(cleanedPagesLines)

	return &core.ProcessedDocument{
		SourcePath:      absPath,
		TotalPages:      totalPages,
		ProcessedPages:  processedCount,
		SkippedPages:    totalPages - processedCount,
		MarkdownContent: markdownContent,
	}, nil
}

func (e *PDFExtractor) getPageCount(ctx context.Context, pdfPath string) (int, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(reqCtx, "pdfinfo", pdfPath)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		lines := strings.Split(out.String(), "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.ToLower(line), "pages:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					if count, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						return count, nil
					}
				}
			}
		}
	}

	pages, err := e.extractPages(ctx, pdfPath)
	if err == nil && len(pages) > 0 {
		return len(pages), nil
	}

	return 1, nil
}

func (e *PDFExtractor) extractPages(ctx context.Context, pdfPath string) ([]core.PDFPage, error) {
	if err := validator.CheckDependencies(); err != nil {
		return nil, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(reqCtx, "pdftotext", "-layout", pdfPath, "-")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.ToLower(stderr.String())
		if strings.Contains(errMsg, "encrypted") || strings.Contains(errMsg, "password") {
			return nil, fmt.Errorf("PDF is encrypted or password-protected: %s", filepath.Base(pdfPath))
		}
		if strings.Contains(errMsg, "syntax error") || strings.Contains(errMsg, "corrupt") {
			return nil, fmt.Errorf("PDF file appears to be corrupted or invalid: %s", filepath.Base(pdfPath))
		}
		return nil, fmt.Errorf("failed to process PDF: %s (%s)", err, strings.TrimSpace(stderr.String()))
	}

	rawPages := strings.Split(stdout.String(), "\x0c")
	var pages []core.PDFPage

	for idx, text := range rawPages {
		trimmed := strings.TrimSpace(text)
		if idx == len(rawPages)-1 && trimmed == "" {
			continue
		}
		pages = append(pages, core.PDFPage{
			PageNumber: idx + 1,
			RawText:    text,
		})
	}

	return pages, nil
}
