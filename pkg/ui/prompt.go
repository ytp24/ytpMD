package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/devops/pdf2md/pkg/validator"
)

// ANSI Styling & Teal Palette
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Teal Palette Shades (24-bit TrueColor)
	TealDeep   = "\033[38;2;15;118;110m" // Deep Teal
	TealMed    = "\033[38;2;13;148;136m" // Classic Teal
	TealBright = "\033[38;2;20;184;166m" // Bright Teal
	TealLight  = "\033[38;2;45;212;191m" // Mint Teal
	TealPale   = "\033[38;2;94;234;212m" // Light Teal

	// Status Colors (No emojis)
	ColorYellow = "\033[38;2;251;191;36m"
	ColorRed    = "\033[38;2;248;113;113m"
	ColorGray   = "\033[38;2;148;163;184m"
)

// DefaultDestinationRoot is the standard base output folder for all extracted notes.
const DefaultDestinationRoot = "~/Documents/ytp24"

// ProcessingMode defines single file or batch processing.
type ProcessingMode int

const (
	ModeSingle ProcessingMode = iota
	ModeBatch
)

// InteractiveOptions holds user-selected parameters from the interactive wizard.
type InteractiveOptions struct {
	Mode            ProcessingMode
	PDFPath         string
	BatchDir        string
	BatchName       string
	DestinationDir  string
	SkipFrontMatter int
	ExcludeAppendix bool
	SplitByChapters bool
	StartPage       int
	EndPage         int
	Concurrency     int
}

// PrintBanner renders the big italic painted teal ASCII ytpMD logo with [pdf2md] underneath.
func PrintBanner(version string) {
	fmt.Println()
	fmt.Printf("%s%s%s  ██╗   ██╗████████╗██████╗ ███╗   ███╗██████╗ %s\n", TealDeep, Bold, Italic, Reset)
	fmt.Printf("%s%s%s  ╚██╗ ██╔╝╚══██╔══╝██╔══██╗████╗ ████║██╔══██╗%s\n", TealMed, Bold, Italic, Reset)
	fmt.Printf("%s%s%s   ╚████╔╝    ██║   ██████╔╝██╔████╔██║██║  ██║%s\n", TealMed, Bold, Italic, Reset)
	fmt.Printf("%s%s%s    ╚██╔╝     ██║   ██╔═══╝ ██║╚██╔╝██║██║  ██║%s\n", TealBright, Bold, Italic, Reset)
	fmt.Printf("%s%s%s     ██║      ██║   ██║     ██║ ╚═╝ ██║██████╔╝%s\n", TealLight, Bold, Italic, Reset)
	fmt.Printf("%s%s%s     ╚═╝      ╚═╝   ╚═╝     ╚═╝     ╚═╝╚═════╝ %s\n", TealPale, Bold, Italic, Reset)

	fmt.Printf("               %s%s%s[ pdf2md ]%s  %s%sv%s%s\n", TealLight, Bold, Italic, Reset, ColorGray, Italic, version, Reset)
	fmt.Printf("   %s%s-- High-Performance Concurrent PDF to Markdown Engine --%s\n", TealBright, Italic, Reset)
	fmt.Printf("   %s%sTransforms PDFs into Chapter-Aware Markdown Notes with Zero Noise%s\n", ColorGray, Italic, Reset)
	fmt.Println()
}

// PromptPDFFile interactively requests PDF path, opening a GUI file selection window if Enter is pressed on empty input.
func PromptPDFFile(reader *bufio.Reader) (string, error) {
	for {
		fmt.Printf("%s%s%s[?]%s %s%sEnter PDF path%s %s[or press Enter to open file chooser window]%s: ",
			TealLight, Bold, Italic, Reset, TealBright, Italic, Reset, ColorGray, Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input stream closed")
		}

		trimmed := strings.TrimSpace(input)
		trimmed = strings.Trim(trimmed, "\"'")

		// If user pressed Enter with empty input -> Launch GUI File Chooser Window
		if trimmed == "" {
			fmt.Printf("%s[*] Opening file selection window...%s\n", TealLight, Reset)
			selected, err := OpenFilePicker("Select PDF Document", "PDF Files (*.pdf) | *.pdf *.PDF")
			if err == nil && selected != "" {
				cleanSelected := validator.ExpandPath(selected)
				if err := validator.ValidatePDFFile(cleanSelected); err == nil {
					fmt.Printf("%s[+] Selected file: %s%s\n", TealBright, cleanSelected, Reset)
					return cleanSelected, nil
				} else {
					fmt.Printf("%s[x] %v%s\n", ColorRed, err, Reset)
				}
			} else {
				fmt.Printf("%s[!] Window selection cancelled or not available in current terminal.%s\n", ColorYellow, Reset)
			}
			fmt.Printf("%sPlease type or paste the PDF file path:%s ", ColorGray, Reset)
			continue
		}

		cleanPath := validator.ExpandPath(trimmed)
		if err := validator.ValidatePDFFile(cleanPath); err != nil {
			fmt.Printf("%s[x] %v%s\n", ColorRed, err, Reset)
			fmt.Printf("%sPlease try again.%s\n", ColorGray, Reset)
			continue
		}

		return cleanPath, nil
	}
}

// PromptBatchDir interactively requests batch directory, opening a GUI folder chooser if Enter is pressed on empty input.
func PromptBatchDir(reader *bufio.Reader) (string, error) {
	for {
		fmt.Printf("%s%s%s[?]%s %s%sEnter Batch PDF directory path%s %s[or press Enter to open folder picker]%s: ",
			TealLight, Bold, Italic, Reset, TealBright, Italic, Reset, ColorGray, Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input stream closed")
		}

		trimmed := strings.TrimSpace(input)
		trimmed = strings.Trim(trimmed, "\"'")

		if trimmed == "" {
			fmt.Printf("%s[*] Opening directory selection window...%s\n", TealLight, Reset)
			selected, err := OpenDirectoryPicker("Select PDF Batch Directory", "")
			if err == nil && selected != "" {
				cleanSelected := validator.ExpandPath(selected)
				if err := validator.ValidateDirectory(cleanSelected); err == nil {
					fmt.Printf("%s[+] Selected directory: %s%s\n", TealBright, cleanSelected, Reset)
					return cleanSelected, nil
				} else {
					fmt.Printf("%s[x] %v%s\n", ColorRed, err, Reset)
				}
			} else {
				fmt.Printf("%s[!] Directory selection cancelled.%s\n", ColorYellow, Reset)
			}
			fmt.Printf("%sPlease type or paste the directory path:%s ", ColorGray, Reset)
			continue
		}

		cleanPath := validator.ExpandPath(trimmed)
		if err := validator.ValidateDirectory(cleanPath); err != nil {
			fmt.Printf("%s[x] %v%s\n", ColorRed, err, Reset)
			fmt.Printf("%sPlease try again.%s\n", ColorGray, Reset)
			continue
		}

		return cleanPath, nil
	}
}

// PromptDestinationDir requests destination folder, defaulting to ~/Documents/ytp24.
func PromptDestinationDir(reader *bufio.Reader, defaultDir string) (string, error) {
	defaultExpanded := validator.ExpandPath(defaultDir)
	for {
		fmt.Printf("%s%s%s[?]%s %s%sEnter destination root%s %s[%s, or 'b' to browse]%s: ",
			TealLight, Bold, Italic, Reset, TealBright, Italic, Reset, ColorGray, defaultDir, Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			_ = os.MkdirAll(defaultExpanded, 0755)
			return defaultExpanded, nil
		}

		trimmed := strings.TrimSpace(input)
		trimmed = strings.Trim(trimmed, "\"'")

		// Default on Enter
		if trimmed == "" {
			_ = os.MkdirAll(defaultExpanded, 0755)
			return defaultExpanded, nil
		}

		if strings.ToLower(trimmed) == "b" || strings.ToLower(trimmed) == "browse" {
			fmt.Printf("%s[*] Opening directory selection window...%s\n", TealLight, Reset)
			selected, err := OpenDirectoryPicker("Select Destination Directory", defaultExpanded)
			if err == nil && selected != "" {
				cleanSelected := validator.ExpandPath(selected)
				if err := validator.ValidateDirectory(cleanSelected); err == nil {
					fmt.Printf("%s[+] Selected destination: %s%s\n", TealBright, cleanSelected, Reset)
					return cleanSelected, nil
				} else {
					fmt.Printf("%s[x] %v%s\n", ColorRed, err, Reset)
				}
			} else {
				fmt.Printf("%s[!] Directory selection cancelled. Using default: %s%s\n", ColorYellow, defaultDir, Reset)
				_ = os.MkdirAll(defaultExpanded, 0755)
				return defaultExpanded, nil
			}
			continue
		}

		cleanPath := validator.ExpandPath(trimmed)
		if err := validator.ValidateDirectory(cleanPath); err != nil {
			fmt.Printf("%s[x] %v%s\n", ColorRed, err, Reset)
			fmt.Printf("%sPlease try again.%s\n", ColorGray, Reset)
			continue
		}

		return cleanPath, nil
	}
}

// PromptString prompts for a generic string with a default value.
func PromptString(reader *bufio.Reader, label string, defaultVal string) string {
	fmt.Printf("%s%s%s[?]%s %s%s%s%s %s[%s]%s: ", TealLight, Bold, Italic, Reset, TealBright, Italic, label, Reset, ColorGray, defaultVal, Reset)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultVal
	}
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return defaultVal
	}
	return trimmed
}

// PromptInt prompts for an integer with a default fallback.
func PromptInt(reader *bufio.Reader, label string, defaultValue int, minVal int) int {
	for {
		fmt.Printf("%s%s%s[?]%s %s%s%s%s %s[%d]%s: ", TealLight, Bold, Italic, Reset, TealBright, Italic, label, Reset, ColorGray, defaultValue, Reset)
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
			fmt.Printf("%s[!] Please enter a valid number (>= %d).%s\n", ColorYellow, minVal, Reset)
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
		fmt.Printf("%s%s%s[?]%s %s%s%s%s %s(%s)%s: ", TealLight, Bold, Italic, Reset, TealBright, Italic, label, Reset, ColorGray, options, Reset)
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

		fmt.Printf("%s[!] Please answer with 'y' or 'n'.%s\n", ColorYellow, Reset)
	}
}

// RunInteractiveWizard executes the full question-answer wizard with smart defaults.
func RunInteractiveWizard(version string) (*InteractiveOptions, error) {
	PrintBanner(version)

	reader := bufio.NewReader(os.Stdin)

	// Mode Selection: Single PDF or Batch Folder
	fmt.Printf("%s%s%s[?]%s %sSelect Mode:%s\n", TealLight, Bold, Italic, Reset, TealBright, Reset)
	fmt.Printf("    %s1.%s Single PDF File %s(default)%s\n", TealBright, Reset, ColorGray, Reset)
	fmt.Printf("    %s2.%s Batch Directory of PDFs %s(Concurrent Goroutines)%s\n", TealBright, Reset, ColorGray, Reset)
	fmt.Printf("%s%sChoice [1]:%s ", TealLight, Bold, Reset)
	
	modeInput, _ := reader.ReadString('\n')
	modeChoice := strings.TrimSpace(modeInput)

	if modeChoice == "2" {
		// Batch Mode
		batchDir, err := PromptBatchDir(reader)
		if err != nil {
			return nil, err
		}

		defaultBatchName := filepath.Base(batchDir)
		batchName := PromptString(reader, "Batch name for storage subfolder", defaultBatchName)
		destDir, err := PromptDestinationDir(reader, DefaultDestinationRoot)
		if err != nil {
			return nil, err
		}

		useDefaults := PromptBool(reader, "Use standard production defaults (TOC chapters -> Appendix cutoff)?", true)

		skipFront := 0
		excludeAppendix := true
		splitChapters := true
		concurrency := 4

		if !useDefaults {
			fmt.Println()
			fmt.Printf("%s%s--- Advanced Batch Settings ---%s\n", TealBright, Italic, Reset)
			skipFront = PromptInt(reader, "Skip initial front-matter pages on each PDF", 0, 0)
			excludeAppendix = PromptBool(reader, "Automatically exclude Appendix, Index, & Bibliography?", true)
			concurrency = PromptInt(reader, "Concurrent Worker Goroutines", 4, 1)
		}

		fmt.Println()
		fmt.Printf("%s%s[+] Batch configured: ~/Documents/ytp24/%s/%s\n\n", TealBright, Bold, batchName, Reset)

		return &InteractiveOptions{
			Mode:            ModeBatch,
			BatchDir:        batchDir,
			BatchName:       batchName,
			DestinationDir:  destDir,
			SkipFrontMatter: skipFront,
			ExcludeAppendix: excludeAppendix,
			SplitByChapters: splitChapters,
			Concurrency:     concurrency,
		}, nil
	}

	// Single PDF Mode (Default)
	pdfPath, err := PromptPDFFile(reader)
	if err != nil {
		return nil, err
	}

	destDir, err := PromptDestinationDir(reader, DefaultDestinationRoot)
	if err != nil {
		return nil, err
	}

	useDefaults := PromptBool(reader, "Use standard production defaults (TOC chapters -> Appendix cutoff)?", true)

	if useDefaults {
		fmt.Println()
		fmt.Printf("%s%s[+] Applying production defaults (TOC chapters extracted, Appendix/Index excluded).%s\n\n", TealBright, Bold, Reset)
		return &InteractiveOptions{
			Mode:            ModeSingle,
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
	fmt.Printf("%s%s[+] All custom settings configured. Starting extraction...%s\n\n", TealBright, Bold, Reset)

	return &InteractiveOptions{
		Mode:            ModeSingle,
		PDFPath:         pdfPath,
		DestinationDir:  destDir,
		SkipFrontMatter: skipFront,
		ExcludeAppendix: excludeAppendix,
		SplitByChapters: splitChapters,
		StartPage:       startPage,
		EndPage:         endPage,
	}, nil
}
