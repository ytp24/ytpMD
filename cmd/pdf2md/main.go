package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/devops/pdf2md/pkg/extractor"
	"github.com/devops/pdf2md/pkg/models"
	"github.com/devops/pdf2md/pkg/ui"
	"github.com/devops/pdf2md/pkg/validator"
)

const version = "2.1.0"

func printUsage() {
	ui.PrintBanner(version)
	fmt.Printf(`%s%sUSAGE:%s
   %spdf2md%s                          Launch interactive wizard (prompts for file, dest & settings)
   %spdf2md convert <input.pdf>%s      Convert a single PDF file
   %spdf2md batch <directory>%s        Batch convert all PDFs in a directory
   %spdf2md help%s                     Show this help screen
   %spdf2md version%s                  Show version

%s%sOPTIONS (for non-interactive CLI flags):%s
   -o, -output <path>              Destination folder (default: alongside PDF)
   -skip-front <N>                 Skip first N pages (covers, copyright, TOC) (default: 0)
   -start-page <N>                 Start page number (default: 1)
   -end-page <N>                   End page number (default: 0 / until end)
   -keep-appendix                  Do NOT exclude appendix, index, and bibliography (default: false)
   -single-file                    Output single monolithic .md instead of chapter folder
   -r, -recursive                  Recursively search subdirectories in batch mode

%s%sEXAMPLES:%s
   # Interactive mode (asks questions with smart defaults):
   pdf2md

   # Direct conversion into a named chapter folder:
   pdf2md convert DevOps_Handbook.pdf -skip-front 3

   # Batch process an entire directory:
   pdf2md batch ~/Downloads/PDFs/ -o ~/Projects/Notes/ -r
`, ui.Bold, ui.MintLight, ui.Reset,
		ui.GreenLight, ui.Reset,
		ui.GreenLight, ui.Reset,
		ui.GreenLight, ui.Reset,
		ui.GreenLight, ui.Reset,
		ui.GreenLight, ui.Reset,
		ui.Bold, ui.MintLight, ui.Reset,
		ui.Bold, ui.MintLight, ui.Reset,
	)
}

func main() {
	// 1. Global Panic Recovery Handler: Never crash or panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			fmt.Printf("%s❌ An unexpected internal error occurred:%s %v\n", ui.ColorRed, ui.Reset, r)
			fmt.Printf("Please ensure the PDF file is valid and poppler-utils is installed.\n")
			os.Exit(1)
		}
	}()

	// 2. Graceful Termination on Ctrl+C / SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println()
		fmt.Printf("\n%s⚠️  Operation cancelled by user. Exiting cleanly.%s\n", ui.ColorYellow, ui.Reset)
		os.Exit(0)
	}()

	// Pre-flight check: Dependencies
	if err := validator.CheckDependencies(); err != nil {
		fmt.Printf("%s❌ Dependency Error:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	// If no arguments provided -> Run interactive prompting mode
	if len(os.Args) == 1 {
		runInteractiveMode()
		return
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("ytpMD [pdf2md] v%s (built with Go %s)\n", version, "1.22")
		return

	case "help", "-h", "--help":
		printUsage()
		return

	case "interactive", "wizard":
		runInteractiveMode()

	case "convert":
		runConvert(os.Args[2:])

	case "batch":
		runBatch(os.Args[2:])

	default:
		// If user passes a direct file path without subcommand (e.g. `pdf2md document.pdf`)
		if strings.HasSuffix(strings.ToLower(command), ".pdf") {
			runConvert(os.Args[1:])
			return
		}
		fmt.Printf("%s❌ Unknown command or file:%s %s\n\n", ui.ColorRed, ui.Reset, command)
		printUsage()
		os.Exit(1)
	}
}

func runInteractiveMode() {
	opts, err := ui.RunInteractiveWizard(version)
	if err != nil {
		fmt.Printf("%s❌ Error:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	cfg := models.DefaultConfig()
	cfg.SkipFrontMatterPages = opts.SkipFrontMatter
	cfg.ExcludeAppendix = opts.ExcludeAppendix
	cfg.StartPage = opts.StartPage
	cfg.EndPage = opts.EndPage

	ext := extractor.NewPDFExtractor(cfg)

	fmt.Printf("%s%s⏳ Extracting and transforming:%s %s...\n", ui.MintLight, ui.Bold, ui.Reset, filepath.Base(opts.PDFPath))

	if !opts.SplitByChapters {
		doc, err := ext.ExtractFile(opts.PDFPath)
		if err != nil {
			fmt.Printf("%s❌ Extraction Error:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		baseName := strings.TrimSuffix(filepath.Base(opts.PDFPath), filepath.Ext(opts.PDFPath))
		targetFile := filepath.Join(opts.DestinationDir, baseName+".md")
		if err := os.WriteFile(targetFile, []byte(doc.MarkdownContent), 0644); err != nil {
			fmt.Printf("%s❌ Failed to save file:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		fmt.Printf("\n%s%s✅ Success! Single file created:%s %s\n", ui.GreenLight, ui.Bold, ui.Reset, targetFile)
		return
	}

	result, err := ext.ExtractToDirectory(opts.PDFPath, opts.DestinationDir)
	if err != nil {
		fmt.Printf("%s❌ Extraction Error:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	fmt.Printf("\n%s%s✨ Extraction Complete! ✨%s\n", ui.GreenLight, ui.Bold, ui.Reset)
	fmt.Printf("   📂 %sDestination Folder:%s %s%s/%s\n", ui.Bold, ui.Reset, ui.MintLight, result.TargetDirectory, ui.Reset)
	fmt.Printf("   📑 %sTOC Index:%s          %s/README.md\n", ui.Bold, ui.Reset, result.TargetDirectory)
	fmt.Printf("   📚 %sTotal Chapters:%s     %d\n", ui.Bold, ui.Reset, len(result.Chapters))
	for _, ch := range result.Chapters {
		fmt.Printf("      %s✓%s [%02d] %s %s-> %s%s\n", ui.GreenLight, ui.Reset, ch.Index, ch.Title, ui.ColorGray, ch.Filename, ui.Reset)
	}
	fmt.Printf("   📊 %sStats:%s              %d total pages | %d converted | %d skipped/filtered.\n\n",
		ui.Bold, ui.Reset, result.TotalPages, result.ProcessedPages, result.SkippedPages)
}

func runConvert(args []string) {
	fs := flag.NewFlagSet("convert", flag.ExitOnError)

	var (
		outputDir    string
		skipFront    int
		startPage    int
		endPage      int
		keepAppendix bool
		singleFile   bool
		noReflow     bool
	)

	fs.StringVar(&outputDir, "o", "", "Destination base directory")
	fs.StringVar(&outputDir, "output", "", "Destination base directory")
	fs.IntVar(&skipFront, "skip-front", 0, "Skip first N pages")
	fs.IntVar(&startPage, "start-page", 1, "Start page number")
	fs.IntVar(&endPage, "end-page", 0, "End page number")
	fs.BoolVar(&keepAppendix, "keep-appendix", false, "Keep appendix and index")
	fs.BoolVar(&singleFile, "single-file", false, "Generate single monolithic .md file")
	fs.BoolVar(&noReflow, "no-reflow", false, "Disable paragraph reflow")

	_ = fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 1 {
		fmt.Printf("%s❌ Error: missing input PDF file path.%s\n", ui.ColorRed, ui.Reset)
		fmt.Println("Usage: pdf2md convert <input.pdf> [-o <output_dir>]")
		os.Exit(1)
	}

	inputPdf := validator.ExpandPath(positional[0])
	if err := validator.ValidatePDFFile(inputPdf); err != nil {
		fmt.Printf("%s❌ %v%s\n", ui.ColorRed, err, ui.Reset)
		os.Exit(1)
	}

	if outputDir == "" {
		outputDir = filepath.Dir(inputPdf)
	} else {
		outputDir = validator.ExpandPath(outputDir)
	}

	if err := validator.ValidateDirectory(outputDir); err != nil {
		fmt.Printf("%s❌ %v%s\n", ui.ColorRed, err, ui.Reset)
		os.Exit(1)
	}

	cfg := models.DefaultConfig()
	cfg.SkipFrontMatterPages = skipFront
	cfg.StartPage = startPage
	cfg.EndPage = endPage
	cfg.ExcludeAppendix = !keepAppendix
	cfg.ReflowParagraphs = !noReflow

	ext := extractor.NewPDFExtractor(cfg)

	fmt.Printf("%s📄 Processing:%s %s...\n", ui.MintLight, ui.Reset, filepath.Base(inputPdf))

	if singleFile {
		doc, err := ext.ExtractFile(inputPdf)
		if err != nil {
			fmt.Printf("%s❌ Extraction failed:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		targetFile := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(inputPdf), filepath.Ext(inputPdf))+".md")
		if err := os.WriteFile(targetFile, []byte(doc.MarkdownContent), 0644); err != nil {
			fmt.Printf("%s❌ Failed to save file:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		fmt.Printf("%s✅ Success! Single file saved to:%s %s\n", ui.GreenLight, ui.Reset, targetFile)
		return
	}

	result, err := ext.ExtractToDirectory(inputPdf, outputDir)
	if err != nil {
		fmt.Printf("%s❌ Extraction failed:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	fmt.Printf("%s✅ Success! Extracted into folder:%s %s/\n", ui.GreenLight, ui.Reset, result.TargetDirectory)
	fmt.Printf("   📑 TOC Index: %s/README.md\n", result.TargetDirectory)
	fmt.Printf("   📚 Total Chapters: %d\n", len(result.Chapters))
	for _, ch := range result.Chapters {
		fmt.Printf("      - [%02d] %s -> %s\n", ch.Index, ch.Title, ch.Filename)
	}
	fmt.Printf("   📊 Stats: %d total pages | %d converted | %d skipped/filtered.\n",
		result.TotalPages, result.ProcessedPages, result.SkippedPages)
}

func runBatch(args []string) {
	fs := flag.NewFlagSet("batch", flag.ExitOnError)

	var (
		outputDir    string
		skipFront    int
		keepAppendix bool
		recursive    bool
	)

	fs.StringVar(&outputDir, "o", "", "Destination directory")
	fs.StringVar(&outputDir, "output", "", "Destination directory")
	fs.IntVar(&skipFront, "skip-front", 0, "Skip first N pages on all files")
	fs.BoolVar(&keepAppendix, "keep-appendix", false, "Keep appendix and index")
	fs.BoolVar(&recursive, "r", false, "Recursive search")
	fs.BoolVar(&recursive, "recursive", false, "Recursive search")

	_ = fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 1 {
		fmt.Printf("%s❌ Error: missing input directory.%s\n", ui.ColorRed, ui.Reset)
		fmt.Println("Usage: pdf2md batch <directory> [-o <output_dir>] [-r]")
		os.Exit(1)
	}

	inputDir := validator.ExpandPath(positional[0])
	if outputDir == "" {
		outputDir = inputDir
	} else {
		outputDir = validator.ExpandPath(outputDir)
	}

	if err := validator.ValidateDirectory(outputDir); err != nil {
		fmt.Printf("%s❌ %v%s\n", ui.ColorRed, err, ui.Reset)
		os.Exit(1)
	}

	cfg := models.DefaultConfig()
	cfg.SkipFrontMatterPages = skipFront
	cfg.ExcludeAppendix = !keepAppendix

	ext := extractor.NewPDFExtractor(cfg)

	var pdfFiles []string
	if recursive {
		_ = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".pdf") {
				pdfFiles = append(pdfFiles, path)
			}
			return nil
		})
	} else {
		entries, err := os.ReadDir(inputDir)
		if err != nil {
			fmt.Printf("%s❌ Error reading directory:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".pdf") {
				pdfFiles = append(pdfFiles, filepath.Join(inputDir, entry.Name()))
			}
		}
	}

	if len(pdfFiles) == 0 {
		fmt.Printf("%s⚠️  No PDF files found in '%s'%s\n", ui.ColorYellow, inputDir, ui.Reset)
		return
	}

	fmt.Printf("%s🔍 Found %d PDF file(s) in %s%s\n", ui.MintLight, len(pdfFiles), inputDir, ui.Reset)

	successCount := 0
	for _, pdf := range pdfFiles {
		fmt.Printf("\n📄 Converting %s...\n", filepath.Base(pdf))
		res, err := ext.ExtractToDirectory(pdf, outputDir)
		if err != nil {
			fmt.Printf("   %s✗ Failed:%s %v\n", ui.ColorRed, ui.Reset, err)
			continue
		}

		fmt.Printf("   %s✓ Created folder:%s %s/ (%d chapters)\n", ui.GreenLight, ui.Reset, filepath.Base(res.TargetDirectory), len(res.Chapters))
		successCount++
	}

	fmt.Printf("\n%s🎉 Batch processing complete! %d/%d PDF documents converted into chapter folders.%s\nBase destination: %s\n",
		ui.GreenLight, successCount, len(pdfFiles), ui.Reset, outputDir)
}
