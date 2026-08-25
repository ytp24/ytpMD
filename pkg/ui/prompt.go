package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/devops/pdf2md/pkg/validator"
)

// ANSI Styling & Gradient Palette
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Teal-to-Green Gradient Shades (24-bit TrueColor)
	TealDark   = "\033[38;2;14;116;144m" // Deep Teal
	TealMed    = "\033[38;2;13;148;136m" // Rich Teal
	TealBright = "\033[38;2;20;184;166m" // Bright Teal
	MintLight  = "\033[38;2;45;212;191m" // Mint Green
	GreenLight = "\033[38;2;52;211;153m" // Emerald Green
	GreenNeon  = "\033[38;2;110;231;183m" // Spring Green

	// Status Colors
	ColorYellow = "\033[38;2;251;191;36m"
	ColorRed    = "\033[38;2;248;113;113m"
	ColorGray   = "\033[38;2;148;163;184m"
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

// PrintBanner renders the big stylized gradient ytpMD logo with [pdf2md] underneath.
func PrintBanner(version string) {
	fmt.Println()
	// Big Painted ASCII Art with Teal-to-Green Gradient
	fmt.Printf("%s%s  ██╗   ██╗████████╗██████╗ ███╗   ███╗██████╗ %s\n", TealDark, Bold, Reset)
	fmt.Printf("%s%s  ╚██╗ ██╔╝╚══██╔══╝██╔══██╗████╗ ████║██╔══██╗%s\n", TealMed, Bold, Reset)
	fmt.Printf("%s%s   ╚████╔╝    ██║   ██████╔╝██╔████╔██║██║  ██║%s\n", TealBright, Bold, Reset)
	fmt.Printf("%s%s    ╚██╔╝     ██║   ██╔═══╝ ██║╚██╔╝██║██║  ██║%s\n", MintLight, Bold, Reset)
	fmt.Printf("%s%s     ██║      ██║   ██║     ██║ ╚═╝ ██║██████╔╝%s\n", GreenLight, Bold, Reset)
	fmt.Printf("%s%s     ╚═╝      ╚═╝   ╚═╝     ╚═╝     ╚═╝╚═════╝ %s\n", GreenNeon, Bold, Reset)
	
	// Centered Subtitle [pdf2md] and Italicized Description
	fmt.Printf("               %s%s%s[ pdf2md ]%s  %s%sv%s%s\n", MintLight, Bold, Italic, Reset, ColorGray, Italic, version, Reset)
	fmt.Printf("   %s%s⚡ High-Performance PDF to Markdown Engine ⚡%s\n", GreenNeon, Italic, Reset)
	fmt.Printf("   %s%sTransforms PDFs into Chapter-Aware Markdown Notes with Zero Noise%s\n", ColorGray, Italic, Reset)
	fmt.Println()
}

// PromptPDFFile interactively requests and validates a PDF file path.
func PromptPDFFile(reader *bufio.Reader) (string, error) {
	for {
		fmt.Printf("%s%s[?]%s %sEnter PDF file path:%s ", MintLight, Bold, Reset, Italic, Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input stream closed")
		}

		trimmed := strings.TrimSpace(input)
		trimmed = strings.Trim(trimmed, "\"'")

		if trimmed == "" {
			fmt.Printf("%s⚠️  Please provide a valid PDF file path.%s\n", ColorYellow, Reset)
			continue
		}

		cleanPath := validator.ExpandPath(trimmed)
		if err := validator.ValidatePDFFile(cleanPath); err != nil {
			fmt.Printf("%s❌ %v%s\n", ColorRed, err, Reset)
			fmt.Printf("%sPlease try again.%s\n", ColorGray, Reset)
			continue
		}

		return cleanPath, nil
	}
}

// PromptDestinationDir requests destination folder with smart default.
func PromptDestinationDir(reader *bufio.Reader, defaultDir string) (string, error) {
	defaultExpanded := validator.ExpandPath(defaultDir)
	for {
		fmt.Printf("%s%s[?]%s %sEnter destination directory%s %s[%s]%s: ", MintLight, Bold, Reset, Italic, Reset, ColorGray, defaultDir, Reset)
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
			fmt.Printf("%s❌ %v%s\n", ColorRed, err, Reset)
			fmt.Printf("%sPlease try again.%s\n", ColorGray, Reset)
			continue
		}

		return cleanPath, nil
	}
}

// PromptInt prompts for an integer with a default fallback.
func PromptInt(reader *bufio.Reader, label string, defaultValue int, minVal int) int {
	for {
		fmt.Printf("%s%s[?]%s %s%s%s %s[%d]%s: ", MintLight, Bold, Reset, Italic, label, Reset, ColorGray, defaultValue, Reset)
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
			fmt.Printf("%s⚠️  Please enter a valid number (>= %d).%s\n", ColorYellow, minVal, Reset)
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
		fmt.Printf("%s%s[?]%s %s%s%s %s(%s)%s: ", MintLight, Bold, Reset, Italic, label, Reset, ColorGray, options, Reset)
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

		fmt.Printf("%s⚠️  Please answer with 'y' or 'n'.%s\n", ColorYellow, Reset)
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

	// 3. Ask whether to apply standard production defaults
	useDefaults := PromptBool(reader, "Use standard production defaults (TOC chapters -> Appendix cutoff)?", true)

	if useDefaults {
		fmt.Println()
		fmt.Printf("%s%s[✓] Applying production defaults (TOC chapters extracted, Appendix/Index excluded).%s\n\n", GreenLight, Bold, Reset)
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

	// Advanced Custom Settings
	fmt.Println()
	fmt.Printf("%s%s--- Advanced Custom Settings ---%s\n", TealBright, Italic, Reset)
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
	fmt.Printf("%s%s[✓] All custom settings configured. Starting extraction...%s\n\n", GreenLight, Bold, Reset)

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
