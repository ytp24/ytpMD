package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/devops/pdf2md/pkg/validator"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorCyan   = "\033[36m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorBlue   = "\033[34m"
)

// InteractiveOptions holds user-selected parameters from the interactive wizard.
type InteractiveOptions struct {
	PDFPath         string
	DestinationDir  string
	SkipFrontMatter int
	ExcludeAppendix bool
	SplitByChapters bool
	StartPage       int
	EndPage         int
}

// PrintBanner prints the tool welcome banner.
func PrintBanner(version string) {
	fmt.Println()
	fmt.Printf("%s%s=======================================================%s\n", ColorCyan, ColorBold, ColorReset)
	fmt.Printf("%s%s   📄  pdf2md (v%s) — Production PDF to Markdown CLI%s\n", ColorCyan, ColorBold, version, ColorReset)
	fmt.Printf("%s%s=======================================================%s\n", ColorCyan, ColorBold, ColorReset)
	fmt.Println("Transform PDFs into clean Markdown notes by chapter with zero noise.")
	fmt.Println()
}

// PromptPDFFile interactively requests and validates a PDF file path.
func PromptPDFFile(reader *bufio.Reader) (string, error) {
	for {
		fmt.Printf("%s%s[?] Enter PDF file path: %s", ColorCyan, ColorBold, ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input stream closed")
		}

		trimmed := strings.TrimSpace(input)
		trimmed = strings.Trim(trimmed, "\"'")

		if trimmed == "" {
			fmt.Printf("%s⚠️  Please provide a valid PDF path.%s\n", ColorYellow, ColorReset)
			continue
		}

		cleanPath := validator.ExpandPath(trimmed)
		if err := validator.ValidatePDFFile(cleanPath); err != nil {
			fmt.Printf("%s❌ %v%s\n", ColorRed, err, ColorReset)
			fmt.Println("Please try again.")
			continue
		}

		return cleanPath, nil
	}
}

// PromptDestinationDir requests destination folder with smart default.
func PromptDestinationDir(reader *bufio.Reader, defaultDir string) (string, error) {
	defaultExpanded := validator.ExpandPath(defaultDir)
	for {
		fmt.Printf("%s%s[?] Enter destination directory [%s]: %s", ColorCyan, ColorBold, defaultDir, ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultExpanded, nil
		}

		trimmed := strings.TrimSpace(input)
		trimmed = strings.Trim(trimmed, "\"'")

		if trimmed == "" {
			return defaultExpanded, nil
		}

		cleanPath := validator.ExpandPath(trimmed)
		if err := validator.ValidateDirectory(cleanPath); err != nil {
			fmt.Printf("%s❌ %v%s\n", ColorRed, err, ColorReset)
			fmt.Println("Please try again.")
			continue
		}

		return cleanPath, nil
	}
}

// PromptInt prompts for an integer with a default fallback.
func PromptInt(reader *bufio.Reader, label string, defaultValue int, minVal int) int {
	for {
		fmt.Printf("%s%s[?] %s [%d]: %s", ColorCyan, ColorBold, label, defaultValue, ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultValue
		}

		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			return defaultValue
		}

		val, err := strconv.Atoi(trimmed)
		if err != nil || val < minVal {
			fmt.Printf("%s⚠️  Please enter a valid number (>= %d).%s\n", ColorYellow, minVal, ColorReset)
			continue
		}

		return val
	}
}

// PromptBool prompts for Yes/No with a default value.
func PromptBool(reader *bufio.Reader, label string, defaultVal bool) bool {
	options := "Y/n"
	if !defaultVal {
		options = "y/N"
	}

	for {
		fmt.Printf("%s%s[?] %s (%s): %s", ColorCyan, ColorBold, label, options, ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultVal
		}

		trimmed := strings.ToLower(strings.TrimSpace(input))
		if trimmed == "" {
			return defaultVal
		}

		if trimmed == "y" || trimmed == "yes" || trimmed == "true" {
			return true
		}
		if trimmed == "n" || trimmed == "no" || trimmed == "false" {
			return false
		}

		fmt.Printf("%s⚠️  Please answer with 'y' or 'n'.%s\n", ColorYellow, ColorReset)
	}
}

// RunInteractiveWizard executes the full question-answer wizard with smart defaults.
func RunInteractiveWizard(version string) (*InteractiveOptions, error) {
	PrintBanner(version)

	reader := bufio.NewReader(os.Stdin)

	// 1. Ask for input PDF file
	pdfPath, err := PromptPDFFile(reader)
	if err != nil {
		return nil, err
	}

	// 2. Ask for destination directory (default to current directory)
	currentDir, _ := os.Getwd()
	destDir, err := PromptDestinationDir(reader, currentDir)
	if err != nil {
		return nil, err
	}

	// 3. Ask whether to apply standard production defaults:
	// - Extract Table of Contents chapters into named folder
	// - Automatically cutoff at Appendix / Index / Bibliography
	// - Exclude images and running headers/footers
	useDefaults := PromptBool(reader, "Use standard production defaults (TOC chapters -> Appendix cutoff)?", true)

	if useDefaults {
		fmt.Println()
		fmt.Printf("%s%s[+] Applying production defaults (TOC chapters extracted, Appendix/Index excluded). Starting...%s\n\n", ColorGreen, ColorBold, ColorReset)
		return &InteractiveOptions{
			PDFPath:         pdfPath,
			DestinationDir:  destDir,
			SkipFrontMatter: 0,
			ExcludeAppendix: true,
			SplitByChapters: true,
			StartPage:       1,
			EndPage:         0,
		}, nil
	}

	// Advanced Custom Settings:
	fmt.Println()
	fmt.Printf("%s--- Custom Extraction Settings ---%s\n", ColorBlue, ColorReset)
	skipFront := PromptInt(reader, "Skip initial front-matter pages (covers, copyright, dedication)", 0, 0)
	excludeAppendix := PromptBool(reader, "Automatically exclude Appendix, Index, & Bibliography sections?", true)
	splitChapters := PromptBool(reader, "Split into individual chapter files inside a named folder?", true)

	customRange := PromptBool(reader, "Set custom start/end page range?", false)
	startPage := 1
	endPage := 0
	if customRange {
		startPage = PromptInt(reader, "Start page number", 1, 1)
		endPage = PromptInt(reader, "End page number (0 = until end)", 0, 0)
	}

	fmt.Println()
	fmt.Printf("%s%s[+] All custom settings configured. Starting extraction...%s\n\n", ColorGreen, ColorBold, ColorReset)

	return &InteractiveOptions{
		PDFPath:         pdfPath,
		DestinationDir:  destDir,
		SkipFrontMatter: skipFront,
		ExcludeAppendix: excludeAppendix,
		SplitByChapters: splitChapters,
		StartPage:       startPage,
		EndPage:         endPage,
	}, nil
}
