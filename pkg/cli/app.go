package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ytp24/ytpMD/pkg/batch"
	"github.com/ytp24/ytpMD/pkg/core"
	"github.com/ytp24/ytpMD/pkg/extractor"
	"github.com/ytp24/ytpMD/pkg/ui"
	"github.com/ytp24/ytpMD/pkg/validator"
)

const AppVersion = "3.2.0"

// Run is the main entrypoint for the CLI application.
func Run() {
	// 1. Global Panic Recovery Handler
	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			fmt.Printf("%s[x] An unexpected internal error occurred:%s %v\n", ui.ColorRed, ui.Reset, r)
			fmt.Printf("Please ensure the PDF file is valid and poppler-utils is installed.\n")
			os.Exit(1)
		}
	}()

	// 2. Graceful Termination on Ctrl+C / SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
		fmt.Println()
		fmt.Printf("\n%s[!] Operation cancelled by user. Exiting cleanly.%s\n", ui.ColorYellow, ui.Reset)
		os.Exit(0)
	}()

	// Pre-flight check: Dependencies
	if err := validator.CheckDependencies(); err != nil {
		fmt.Printf("%s[x] Dependency Error:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	// Interactive mode when run without arguments
	if len(os.Args) == 1 {
		runInteractiveMode(ctx)
		return
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("ytpMD [pdf2md] v%s (built with Go %s)\n", AppVersion, "1.22")
		return

	case "help", "-h", "--help":
		printUsage()
		return

	case "interactive", "wizard":
		runInteractiveMode(ctx)

	case "convert":
		runConvert(ctx, os.Args[2:])

	case "batch":
		runBatch(ctx, os.Args[2:])

	default:
		if strings.HasSuffix(strings.ToLower(command), ".pdf") {
			runConvert(ctx, os.Args[1:])
			return
		}
		fmt.Printf("%s[x] Unknown command or file:%s %s\n\n", ui.ColorRed, ui.Reset, command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	ui.PrintBanner(AppVersion)
	fmt.Printf(`%s%sUSAGE:%s
   %sytpmd%s                           Launch interactive wizard (single PDF or concurrent batch)
   %sytpmd convert <input.pdf>%s       Convert a single PDF file into a chapter-based notes folder
   %sytpmd batch <directory>%s         Batch convert all PDFs using concurrent Goroutines
   %sytpmd help%s                      Show this help screen
   %sytpmd version%s                   Show version

%s%sOPTIONS (for non-interactive CLI flags):%s
   -o, -output <path>              Destination root folder (default: ~/Documents/ytpMD)
   -name <batch_name>              Subfolder name for batch storage (default: input folder name)
   -f, -force, -overwrite          Overwrite existing destination folder without prompting
   -concurrency <N>                Number of parallel worker goroutines (default: 4)
   -skip-front <N>                 Skip first N pages (covers, copyright, TOC) (default: 0)
   -start-page <N>                 Start page number (default: 1)
   -end-page <N>                   End page number (default: 0 / until end)
   -keep-appendix                  Do NOT exclude appendix, index, and bibliography (default: false)
   -single-file                    Output single monolithic .md instead of chapter folder
   -r, -recursive                  Recursively search subdirectories in batch mode

%s%sEXAMPLES:%s
   # Interactive mode (prompts with file chooser & overwrite checks):
   ytpmd

   # Convert single PDF into ~/Documents/ytpMD/DevOps_Handbook/:
   ytpmd convert DevOps_Handbook.pdf

   # Force overwrite if destination already exists:
   ytpmd convert DevOps_Handbook.pdf -force

   # Concurrent batch conversion into ~/Documents/ytpMD/CloudBooks/:
   ytpmd batch ~/Downloads/PDFs/ -name CloudBooks -concurrency 6 -force
`, ui.Bold, ui.TealLight, ui.Reset,
		ui.TealBright, ui.Reset,
		ui.TealBright, ui.Reset,
		ui.TealBright, ui.Reset,
		ui.TealBright, ui.Reset,
		ui.TealBright, ui.Reset,
		ui.Bold, ui.TealLight, ui.Reset,
		ui.Bold, ui.TealLight, ui.Reset,
	)
}

func runInteractiveMode(ctx context.Context) {
	opts, err := ui.RunInteractiveWizard(AppVersion)
	if err != nil {
		fmt.Printf("%s[x] %v%s\n", ui.ColorYellow, err, ui.Reset)
		return
	}

	cfg := core.DefaultConfig()
	cfg.SkipFrontMatterPages = opts.SkipFrontMatter
	cfg.ExcludeAppendix = opts.ExcludeAppendix
	cfg.StartPage = opts.StartPage
	cfg.EndPage = opts.EndPage
	cfg.Concurrency = opts.Concurrency

	ext := extractor.NewPDFExtractor(cfg)

	if opts.Mode == ui.ModeBatch {
		pdfFiles, err := findPDFFiles(opts.BatchDir, true)
		if err != nil || len(pdfFiles) == 0 {
			fmt.Printf("%s[!] No PDF documents found in '%s'%s\n", ui.ColorYellow, opts.BatchDir, ui.Reset)
			return
		}

		fmt.Printf("%s[*] Launching concurrent batch engine (%d workers) for %d PDF(s)...%s\n\n",
			ui.TealLight, opts.Concurrency, len(pdfFiles), ui.Reset)

		bar := ui.NewProgressBar(len(pdfFiles), "Converting Batch")
		engine := batch.NewConcurrentBatchEngine(ext)
		res, err := engine.ProcessBatch(ctx, pdfFiles, opts.DestinationDir, opts.BatchName, opts.Concurrency, bar)
		if err != nil {
			fmt.Printf("%s[x] Batch execution error:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}

		printBatchSummary(res)
		return
	}

	// Single PDF Mode
	fmt.Printf("%s%s[*] Extracting and transforming:%s %s...\n", ui.TealLight, ui.Bold, ui.Reset, filepath.Base(opts.PDFPath))

	if !opts.SplitByChapters {
		doc, err := ext.ExtractFile(ctx, opts.PDFPath)
		if err != nil {
			fmt.Printf("%s[x] Extraction Error:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		baseName := strings.TrimSuffix(filepath.Base(opts.PDFPath), filepath.Ext(opts.PDFPath))
		targetFile := filepath.Join(opts.DestinationDir, baseName+".md")
		if err := os.WriteFile(targetFile, []byte(doc.MarkdownContent), 0644); err != nil {
			fmt.Printf("%s[x] Failed to save file:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		fmt.Printf("\n%s%s[+] Success! Single file created:%s %s\n", ui.TealBright, ui.Bold, ui.Reset, targetFile)
		return
	}

	result, err := ext.ExtractToDirectory(ctx, opts.PDFPath, opts.DestinationDir)
	if err != nil {
		fmt.Printf("%s[x] Extraction Error:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	fmt.Printf("\n%s%s[+] Extraction Complete (Agent-Ready)%s\n", ui.TealBright, ui.Bold, ui.Reset)
	fmt.Printf("   [-] %sDestination Folder:%s %s%s/%s\n", ui.Bold, ui.Reset, ui.TealLight, result.TargetDirectory, ui.Reset)
	fmt.Printf("   [-] %sTOC Index:         %s %s/README.md\n", ui.Bold, ui.Reset, result.TargetDirectory)
	fmt.Printf("   [-] %sAI Agent Manifest: %s %s/AGENTS.md\n", ui.Bold, ui.Reset, result.TargetDirectory)
	fmt.Printf("   [-] %sTotal Chapters:    %s %d\n", ui.Bold, ui.Reset, len(result.Chapters))
	for _, ch := range result.Chapters {
		fmt.Printf("       + [%02d] %s %s-> %s (~%d tokens)%s\n", ch.Index, ch.Title, ui.ColorGray, ch.Filename, ch.TokenEstimate, ui.Reset)
	}
	fmt.Printf("   [-] %sStats:             %s %d total pages | %d converted | %d skipped/filtered.\n\n",
		ui.Bold, ui.Reset, result.TotalPages, result.ProcessedPages, result.SkippedPages)
}

func runConvert(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("convert", flag.ExitOnError)

	var (
		outputDir    string
		skipFront    int
		startPage    int
		endPage      int
		keepAppendix bool
		singleFile   bool
		noReflow     bool
		force        bool
	)

	fs.StringVar(&outputDir, "o", ui.DefaultDestinationRoot, "Destination root directory")
	fs.StringVar(&outputDir, "output", ui.DefaultDestinationRoot, "Destination root directory")
	fs.IntVar(&skipFront, "skip-front", 0, "Skip first N pages")
	fs.IntVar(&startPage, "start-page", 1, "Start page number")
	fs.IntVar(&endPage, "end-page", 0, "End page number")
	fs.BoolVar(&keepAppendix, "keep-appendix", false, "Keep appendix and index")
	fs.BoolVar(&singleFile, "single-file", false, "Generate single monolithic .md file")
	fs.BoolVar(&noReflow, "no-reflow", false, "Disable paragraph reflow")
	fs.BoolVar(&force, "f", false, "Overwrite existing destination folder")
	fs.BoolVar(&force, "force", false, "Overwrite existing destination folder")
	fs.BoolVar(&force, "overwrite", false, "Overwrite existing destination folder")

	_ = fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 1 {
		fmt.Printf("%s[x] Error: missing input PDF file path.%s\n", ui.ColorRed, ui.Reset)
		fmt.Println("Usage: ytpmd convert <input.pdf> [-o <output_dir>] [-force]")
		os.Exit(1)
	}

	inputPdf := validator.ExpandPath(positional[0])
	if err := validator.ValidatePDFFile(inputPdf); err != nil {
		fmt.Printf("%s[x] %v%s\n", ui.ColorRed, err, ui.Reset)
		os.Exit(1)
	}

	destDir := validator.ExpandPath(outputDir)
	if err := validator.ValidateDirectory(destDir); err != nil {
		fmt.Printf("%s[x] %v%s\n", ui.ColorRed, err, ui.Reset)
		os.Exit(1)
	}

	pdfBaseName := strings.TrimSuffix(filepath.Base(inputPdf), filepath.Ext(inputPdf))
	targetDocDir := filepath.Join(destDir, pdfBaseName)

	if !force && ui.CheckDirectoryNonEmpty(targetDocDir) {
		reader := bufio.NewReader(os.Stdin)
		if !ui.PromptOverwrite(reader, targetDocDir) {
			fmt.Printf("%s[!] Conversion cancelled: output folder '%s' already exists.%s\n", ui.ColorYellow, targetDocDir, ui.Reset)
			return
		}
	}

	cfg := core.DefaultConfig()
	cfg.SkipFrontMatterPages = skipFront
	cfg.StartPage = startPage
	cfg.EndPage = endPage
	cfg.ExcludeAppendix = !keepAppendix
	cfg.ReflowParagraphs = !noReflow

	ext := extractor.NewPDFExtractor(cfg)

	fmt.Printf("%s[*] Processing:%s %s...\n", ui.TealLight, ui.Reset, filepath.Base(inputPdf))

	if singleFile {
		targetFile := filepath.Join(destDir, pdfBaseName+".md")
		if !force && ui.CheckFileExists(targetFile) {
			reader := bufio.NewReader(os.Stdin)
			if !ui.PromptBool(reader, fmt.Sprintf("File '%s' already exists. Overwrite?", filepath.Base(targetFile)), false) {
				fmt.Printf("%s[!] Conversion cancelled.%s\n", ui.ColorYellow, ui.Reset)
				return
			}
		}

		doc, err := ext.ExtractFile(ctx, inputPdf)
		if err != nil {
			fmt.Printf("%s[x] Extraction failed:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		if err := os.WriteFile(targetFile, []byte(doc.MarkdownContent), 0644); err != nil {
			fmt.Printf("%s[x] Failed to save file:%s %v\n", ui.ColorRed, ui.Reset, err)
			os.Exit(1)
		}
		fmt.Printf("%s[+] Success! Single file saved to:%s %s\n", ui.TealBright, ui.Reset, targetFile)
		return
	}

	result, err := ext.ExtractToDirectory(ctx, inputPdf, destDir)
	if err != nil {
		fmt.Printf("%s[x] Extraction failed:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	fmt.Printf("%s[+] Success! Extracted into folder:%s %s/\n", ui.TealBright, ui.Reset, result.TargetDirectory)
	fmt.Printf("   [-] TOC Index: %s/README.md\n", result.TargetDirectory)
	fmt.Printf("   [-] AI Agent Manifest: %s/AGENTS.md\n", result.TargetDirectory)
	fmt.Printf("   [-] Total Chapters: %d\n", len(result.Chapters))
	for _, ch := range result.Chapters {
		fmt.Printf("       + [%02d] %s -> %s (~%d tokens)\n", ch.Index, ch.Title, ch.Filename, ch.TokenEstimate)
	}
	fmt.Printf("   [-] Stats: %d total pages | %d converted | %d skipped/filtered.\n",
		result.TotalPages, result.ProcessedPages, result.SkippedPages)
}

func runBatch(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("batch", flag.ExitOnError)

	var (
		outputDir    string
		batchName    string
		concurrency  int
		skipFront    int
		keepAppendix bool
		recursive    bool
		force        bool
	)

	fs.StringVar(&outputDir, "o", ui.DefaultDestinationRoot, "Base destination root")
	fs.StringVar(&outputDir, "output", ui.DefaultDestinationRoot, "Base destination root")
	fs.StringVar(&batchName, "name", "", "Batch subfolder name (default: input directory name)")
	fs.IntVar(&concurrency, "concurrency", 4, "Number of concurrent worker goroutines")
	fs.IntVar(&skipFront, "skip-front", 0, "Skip first N pages on all files")
	fs.BoolVar(&keepAppendix, "keep-appendix", false, "Keep appendix and index")
	fs.BoolVar(&recursive, "r", false, "Recursive search")
	fs.BoolVar(&recursive, "recursive", false, "Recursive search")
	fs.BoolVar(&force, "f", false, "Overwrite existing destination folder")
	fs.BoolVar(&force, "force", false, "Overwrite existing destination folder")
	fs.BoolVar(&force, "overwrite", false, "Overwrite existing destination folder")

	_ = fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 1 {
		fmt.Printf("%s[x] Error: missing input directory.%s\n", ui.ColorRed, ui.Reset)
		fmt.Println("Usage: ytpmd batch <directory> [-name <batch_name>] [-o <output_dir>] [-concurrency <N>] [-force]")
		os.Exit(1)
	}

	inputDir := validator.ExpandPath(positional[0])
	destRoot := validator.ExpandPath(outputDir)

	if batchName == "" {
		batchName = filepath.Base(inputDir)
	}

	if err := validator.ValidateDirectory(destRoot); err != nil {
		fmt.Printf("%s[x] %v%s\n", ui.ColorRed, err, ui.Reset)
		os.Exit(1)
	}

	targetBatchDir := filepath.Join(destRoot, batchName)
	if !force && ui.CheckDirectoryNonEmpty(targetBatchDir) {
		reader := bufio.NewReader(os.Stdin)
		if !ui.PromptOverwrite(reader, targetBatchDir) {
			fmt.Printf("%s[!] Batch cancelled: destination directory '%s' already exists.%s\n", ui.ColorYellow, targetBatchDir, ui.Reset)
			return
		}
	}

	pdfFiles, err := findPDFFiles(inputDir, recursive)
	if err != nil || len(pdfFiles) == 0 {
		fmt.Printf("%s[!] No PDF files found in '%s'%s\n", ui.ColorYellow, inputDir, ui.Reset)
		return
	}

	cfg := core.DefaultConfig()
	cfg.SkipFrontMatterPages = skipFront
	cfg.ExcludeAppendix = !keepAppendix
	cfg.Concurrency = concurrency

	ext := extractor.NewPDFExtractor(cfg)
	engine := batch.NewConcurrentBatchEngine(ext)

	fmt.Printf("%s[*] Launching concurrent batch engine (%d workers) for %d PDF(s)...%s\n\n",
		ui.TealLight, concurrency, len(pdfFiles), ui.Reset)

	bar := ui.NewProgressBar(len(pdfFiles), "Converting Batch")
	result, err := engine.ProcessBatch(ctx, pdfFiles, destRoot, batchName, concurrency, bar)
	if err != nil {
		fmt.Printf("%s[x] Batch execution failed:%s %v\n", ui.ColorRed, ui.Reset, err)
		os.Exit(1)
	}

	printBatchSummary(result)
}

func findPDFFiles(dir string, recursive bool) ([]string, error) {
	var pdfFiles []string
	if recursive {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".pdf") {
				pdfFiles = append(pdfFiles, path)
			}
			return nil
		})
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".pdf") {
				pdfFiles = append(pdfFiles, filepath.Join(dir, entry.Name()))
			}
		}
	}
	return pdfFiles, nil
}

func printBatchSummary(res *core.BatchResult) {
	fmt.Println()
	fmt.Printf("%s%s[+] Batch Processing Completed in %s%s\n", ui.TealBright, ui.Bold, res.Duration.Round(time.Millisecond), ui.Reset)
	fmt.Printf("   [-] %sBatch Directory:%s  %s%s/%s\n", ui.Bold, ui.Reset, ui.TealLight, res.TargetDirectory, ui.Reset)
	fmt.Printf("   [-] %sMaster Library:%s   %s/README.md\n", ui.Bold, ui.Reset, res.TargetDirectory)
	fmt.Printf("   [-] %sTotal Files:   %s   %d converted | %d failed\n", ui.Bold, ui.Reset, res.ProcessedFiles, res.FailedFiles)
	for _, r := range res.Results {
		if r.Success {
			fmt.Printf("       + %s %s(%d chapters, %d pages)%s\n", r.PDFName, ui.ColorGray, r.ChaptersCount, r.TotalPages, ui.Reset)
		} else {
			fmt.Printf("       x %s %s(failed: %v)%s\n", r.PDFName, ui.ColorRed, r.Error, ui.Reset)
		}
	}
	fmt.Println()
}
