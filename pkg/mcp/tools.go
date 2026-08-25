package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ytp24/ytpMD/pkg/batch"
	"github.com/ytp24/ytpMD/pkg/core"
	"github.com/ytp24/ytpMD/pkg/extractor"
	"github.com/ytp24/ytpMD/pkg/validator"
)

// GetRegisteredTools returns the list of all available MCP tools with their JSON Schema definitions.
func GetRegisteredTools() []Tool {
	return []Tool{
		{
			Name:        "convert_pdf",
			Description: "Converts a PDF file into a clean, chapter-segmented Markdown note directory with YAML frontmatter, breadcrumbs, and an AGENTS.md ingestion manifest, or a single concatenated Markdown file.",
			InputSchema: ToolSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"path": {
						Type:        "string",
						Description: "Absolute or relative path to the PDF document to convert.",
					},
					"output_dir": {
						Type:        "string",
						Description: "Base destination directory where extracted notes will be saved (defaults to '~/Documents/ytpMD').",
						Default:     "~/Documents/ytpMD",
					},
					"single_file": {
						Type:        "boolean",
						Description: "If true, outputs a single monolithic Markdown file instead of a chapter-segmented folder.",
						Default:     false,
					},
					"skip_front_matter": {
						Type:        "integer",
						Description: "Number of initial front-matter pages to skip (e.g. covers, copyright, dedication).",
						Default:     0,
					},
					"start_page": {
						Type:        "integer",
						Description: "Starting page number for extraction (1-indexed).",
						Default:     1,
					},
					"end_page": {
						Type:        "integer",
						Description: "Ending page number for extraction (0 = process until end of document).",
						Default:     0,
					},
					"exclude_appendix": {
						Type:        "boolean",
						Description: "Automatically detect and exclude Appendix, Index, Bibliography, and Glossary sections.",
						Default:     true,
					},
					"reflow_paragraphs": {
						Type:        "boolean",
						Description: "Join hard line-breaks and de-hyphenate split words mid-sentence into clean paragraphs.",
						Default:     true,
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "batch_convert",
			Description: "Converts an entire directory of PDF files concurrently using a managed Go worker pool, outputting structured note folders and a master batch library index.",
			InputSchema: ToolSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"directory": {
						Type:        "string",
						Description: "Path to the directory containing PDF files to convert.",
					},
					"output_dir": {
						Type:        "string",
						Description: "Base destination root directory for batch storage (defaults to '~/Documents/ytpMD').",
						Default:     "~/Documents/ytpMD",
					},
					"batch_name": {
						Type:        "string",
						Description: "Subfolder name for storing this batch library (defaults to input directory name).",
					},
					"concurrency": {
						Type:        "integer",
						Description: "Number of parallel worker goroutines to utilize (defaults to 4).",
						Default:     4,
					},
					"recursive": {
						Type:        "boolean",
						Description: "Whether to recursively search subdirectories for PDF files.",
						Default:     false,
					},
					"exclude_appendix": {
						Type:        "boolean",
						Description: "Automatically detect and exclude Appendix, Index, and Bibliography sections.",
						Default:     true,
					},
				},
				Required: []string{"directory"},
			},
		},
		{
			Name:        "inspect_pdf",
			Description: "Performs pre-flight inspection on a PDF file: validates '%PDF-' header bytes, detects total page count, checks encryption/corruption, and verifies readability without writing to disk.",
			InputSchema: ToolSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"path": {
						Type:        "string",
						Description: "Path to the PDF file to inspect.",
					},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "get_manifest",
			Description: "Reads and returns the machine-readable AGENTS.md ingestion manifest and chapter navigation map from an already converted document folder.",
			InputSchema: ToolSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"folder_path": {
						Type:        "string",
						Description: "Path to the converted document folder (containing AGENTS.md or README.md).",
					},
				},
				Required: []string{"folder_path"},
			},
		},
	}
}

// ExecuteTool handles execution of MCP tools and returns a standardized CallToolResult.
func ExecuteTool(ctx context.Context, name string, args map[string]interface{}) CallToolResult {
	switch name {
	case "convert_pdf":
		return handleConvertPDF(ctx, args)
	case "batch_convert":
		return handleBatchConvert(ctx, args)
	case "inspect_pdf":
		return handleInspectPDF(ctx, args)
	case "get_manifest":
		return handleGetManifest(ctx, args)
	default:
		return CallToolResult{
			Content: []TextContent{{Type: "text", Text: fmt.Sprintf("Unknown tool: '%s'", name)}},
			IsError: true,
		}
	}
}

func handleConvertPDF(ctx context.Context, args map[string]interface{}) CallToolResult {
	path, ok := args["path"].(string)
	if !ok || strings.TrimSpace(path) == "" {
		return errorResult("Missing required argument 'path'")
	}

	cleanPath := validator.ExpandPath(path)
	if err := validator.ValidatePDFFile(cleanPath); err != nil {
		return errorResult(fmt.Sprintf("PDF validation failed: %v", err))
	}

	outputDir := "~/Documents/ytpMD"
	if customOut, ok := args["output_dir"].(string); ok && strings.TrimSpace(customOut) != "" {
		outputDir = customOut
	}
	cleanOutputDir := validator.ExpandPath(outputDir)

	cfg := core.DefaultConfig()

	if val, ok := args["skip_front_matter"].(float64); ok {
		cfg.SkipFrontMatterPages = int(val)
	}
	if val, ok := args["start_page"].(float64); ok {
		cfg.StartPage = int(val)
	}
	if val, ok := args["end_page"].(float64); ok {
		cfg.EndPage = int(val)
	}
	if val, ok := args["exclude_appendix"].(bool); ok {
		cfg.ExcludeAppendix = val
	}
	if val, ok := args["reflow_paragraphs"].(bool); ok {
		cfg.ReflowParagraphs = val
	}

	singleFile := false
	if val, ok := args["single_file"].(bool); ok {
		singleFile = val
	}

	ext := extractor.NewPDFExtractor(cfg)

	if singleFile {
		doc, err := ext.ExtractFile(ctx, cleanPath)
		if err != nil {
			return errorResult(fmt.Sprintf("Extraction failed: %v", err))
		}

		baseName := strings.TrimSuffix(filepath.Base(cleanPath), filepath.Ext(cleanPath))
		_ = os.MkdirAll(cleanOutputDir, 0755)
		targetFile := filepath.Join(cleanOutputDir, baseName+".md")
		if err := os.WriteFile(targetFile, []byte(doc.MarkdownContent), 0644); err != nil {
			return errorResult(fmt.Sprintf("Failed to save Markdown file: %v", err))
		}

		resp := map[string]interface{}{
			"status":           "success",
			"mode":             "single_file",
			"source_pdf":       cleanPath,
			"target_file":      targetFile,
			"total_pages":      doc.TotalPages,
			"processed_pages":  doc.ProcessedPages,
			"skipped_pages":    doc.SkippedPages,
			"word_count":       doc.WordCount,
			"estimated_tokens": doc.TokenEstimate,
		}
		jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
		return successResult(string(jsonBytes))
	}

	result, err := ext.ExtractToDirectory(ctx, cleanPath, cleanOutputDir)
	if err != nil {
		return errorResult(fmt.Sprintf("Chapter extraction failed: %v", err))
	}

	type ChapterSummary struct {
		Index         int    `json:"chapter_index"`
		Title         string `json:"title"`
		Filename      string `json:"filename"`
		StartPage     int    `json:"start_page"`
		WordCount     int    `json:"word_count"`
		TokenEstimate int    `json:"estimated_tokens"`
		Path          string `json:"file_path"`
	}

	var chapters []ChapterSummary
	for _, ch := range result.Chapters {
		chapters = append(chapters, ChapterSummary{
			Index:         ch.Index,
			Title:         ch.Title,
			Filename:      ch.Filename,
			StartPage:     ch.StartPage,
			WordCount:     ch.WordCount,
			TokenEstimate: ch.TokenEstimate,
			Path:          filepath.Join(result.TargetDirectory, ch.Filename),
		})
	}

	resp := map[string]interface{}{
		"status":           "success",
		"mode":             "chapter_split",
		"source_pdf":       cleanPath,
		"target_directory": result.TargetDirectory,
		"toc_readme":       filepath.Join(result.TargetDirectory, "README.md"),
		"agent_manifest":   filepath.Join(result.TargetDirectory, "AGENTS.md"),
		"total_chapters":   len(chapters),
		"total_pages":      result.TotalPages,
		"processed_pages":  result.ProcessedPages,
		"skipped_pages":    result.SkippedPages,
		"total_words":      result.TotalWords,
		"total_tokens":     result.TotalTokens,
		"chapters":         chapters,
	}

	jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
	return successResult(string(jsonBytes))
}

func handleBatchConvert(ctx context.Context, args map[string]interface{}) CallToolResult {
	dir, ok := args["directory"].(string)
	if !ok || strings.TrimSpace(dir) == "" {
		return errorResult("Missing required argument 'directory'")
	}

	cleanDir := validator.ExpandPath(dir)
	if err := validator.ValidateDirectory(cleanDir); err != nil {
		return errorResult(fmt.Sprintf("Invalid directory: %v", err))
	}

	outputDir := "~/Documents/ytpMD"
	if customOut, ok := args["output_dir"].(string); ok && strings.TrimSpace(customOut) != "" {
		outputDir = customOut
	}
	cleanOutputDir := validator.ExpandPath(outputDir)

	batchName := filepath.Base(cleanDir)
	if customName, ok := args["batch_name"].(string); ok && strings.TrimSpace(customName) != "" {
		batchName = customName
	}

	concurrency := 4
	if val, ok := args["concurrency"].(float64); ok && val >= 1 {
		concurrency = int(val)
	}

	recursive := false
	if val, ok := args["recursive"].(bool); ok {
		recursive = val
	}

	excludeAppendix := true
	if val, ok := args["exclude_appendix"].(bool); ok {
		excludeAppendix = val
	}

	var pdfFiles []string
	if recursive {
		_ = filepath.Walk(cleanDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".pdf") {
				pdfFiles = append(pdfFiles, path)
			}
			return nil
		})
	} else {
		entries, err := os.ReadDir(cleanDir)
		if err != nil {
			return errorResult(fmt.Sprintf("Failed to read directory: %v", err))
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".pdf") {
				pdfFiles = append(pdfFiles, filepath.Join(cleanDir, entry.Name()))
			}
		}
	}

	if len(pdfFiles) == 0 {
		return errorResult(fmt.Sprintf("No PDF documents found in '%s'", cleanDir))
	}

	cfg := core.DefaultConfig()
	cfg.ExcludeAppendix = excludeAppendix
	cfg.Concurrency = concurrency

	ext := extractor.NewPDFExtractor(cfg)
	engine := batch.NewConcurrentBatchEngine(ext)

	res, err := engine.ProcessBatch(ctx, pdfFiles, cleanOutputDir, batchName, concurrency, nil)
	if err != nil {
		return errorResult(fmt.Sprintf("Batch execution failed: %v", err))
	}

	resp := map[string]interface{}{
		"status":           "success",
		"batch_name":       res.BatchName,
		"target_directory": res.TargetDirectory,
		"master_readme":    filepath.Join(res.TargetDirectory, "README.md"),
		"total_files":      res.TotalFiles,
		"processed_files":  res.ProcessedFiles,
		"failed_files":     res.FailedFiles,
		"duration_seconds": res.Duration.Seconds(),
		"results":          res.Results,
	}

	jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
	return successResult(string(jsonBytes))
}

func handleInspectPDF(ctx context.Context, args map[string]interface{}) CallToolResult {
	path, ok := args["path"].(string)
	if !ok || strings.TrimSpace(path) == "" {
		return errorResult("Missing required argument 'path'")
	}

	cleanPath := validator.ExpandPath(path)
	if err := validator.ValidatePDFFile(cleanPath); err != nil {
		return errorResult(fmt.Sprintf("Validation failed: %v", err))
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		return errorResult(fmt.Sprintf("Cannot stat file: %v", err))
	}

	cfg := core.DefaultConfig()
	ext := extractor.NewPDFExtractor(cfg)

	// Attempt reading single page to verify text layers
	doc, err := ext.ExtractFile(ctx, cleanPath)
	var readableText bool
	var totalPages int
	var tokenEst int
	var words int

	if err == nil && doc != nil {
		readableText = len(doc.MarkdownContent) >= 10
		totalPages = doc.TotalPages
		tokenEst = doc.TokenEstimate
		words = doc.WordCount
	}

	resp := map[string]interface{}{
		"status":              "valid",
		"file_name":           filepath.Base(cleanPath),
		"absolute_path":       cleanPath,
		"file_size_bytes":     info.Size(),
		"file_size_mb":        fmt.Sprintf("%.2f MB", float64(info.Size())/(1024*1024)),
		"total_pages":         totalPages,
		"has_text_layers":     readableText,
		"estimated_words":     words,
		"estimated_tokens":    tokenEst,
		"ocr_required":        !readableText,
		"ready_for_converter": true,
	}

	jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
	return successResult(string(jsonBytes))
}

func handleGetManifest(_ context.Context, args map[string]interface{}) CallToolResult {
	folderPath, ok := args["folder_path"].(string)
	if !ok || strings.TrimSpace(folderPath) == "" {
		return errorResult("Missing required argument 'folder_path'")
	}

	cleanDir := validator.ExpandPath(folderPath)
	manifestPath := filepath.Join(cleanDir, "AGENTS.md")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// Fallback to README.md
		readmePath := filepath.Join(cleanDir, "README.md")
		if content, err := os.ReadFile(readmePath); err == nil {
			return successResult(string(content))
		}
		return errorResult(fmt.Sprintf("No AGENTS.md or README.md found in '%s'", cleanDir))
	}

	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return errorResult(fmt.Sprintf("Failed to read AGENTS.md: %v", err))
	}

	return successResult(string(content))
}

func successResult(text string) CallToolResult {
	return CallToolResult{
		Content: []TextContent{{Type: "text", Text: text}},
		IsError: false,
	}
}

func errorResult(msg string) CallToolResult {
	return CallToolResult{
		Content: []TextContent{{Type: "text", Text: msg}},
		IsError: true,
	}
}
