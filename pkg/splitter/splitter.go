package splitter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/ytp24/ytpMD/pkg/core"
	"github.com/ytp24/ytpMD/pkg/transformer"
)

// Splitter implements the core.Splitter interface with robust TOC vs Body discrimination.
type Splitter struct {
	config      core.Config
	transformer *transformer.Transformer
}

// NewSplitter initializes a new Splitter using Factory Pattern.
func NewSplitter(cfg core.Config) *Splitter {
	return &Splitter{
		config:      cfg,
		transformer: transformer.NewTransformer(cfg),
	}
}

// SplitIntoChapters groups extracted page lines into distinct, verified chapter structures.
func (s *Splitter) SplitIntoChapters(pages []core.PDFPage, pdfTitle string) []core.Chapter {
	var rawChapters []core.Chapter
	var currentChapter *core.Chapter

	chapterHeaderPattern := regexp.MustCompile(`(?i)^(CHAPTER\s+\d+|PART\s+[IVXLCDM]+|MODULE\s+\d+)[:.]?\s*(.*)$`)
	numberedHeadingPattern := regexp.MustCompile(`^(\d+)\.\s+([A-Z][A-Za-z0-9\s,-]{3,50})$`)
	tocLeaderPattern := regexp.MustCompile(`(?i)^(chapter\s+\d+|part\s+[ivxlcdm]+).*[\.\s_–-]{3,}\s*\d+$`)

	for _, page := range pages {
		if page.IsFilteredOut {
			continue
		}

		// 1. Pre-check if page is a Table of Contents / Outline page
		chapterHeadersOnPage := 0
		for _, l := range page.Lines {
			t := strings.TrimSpace(l)
			if chapterHeaderPattern.MatchString(t) || tocLeaderPattern.MatchString(t) {
				chapterHeadersOnPage++
			}
		}
		isTOCPage := chapterHeadersOnPage >= 2

		for _, line := range page.Lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				if currentChapter != nil {
					currentChapter.Lines = append(currentChapter.Lines, "")
				}
				continue
			}

			isNewChapter := false
			var chapterTitle string

			if !isTOCPage && !tocLeaderPattern.MatchString(trimmed) {
				if m := chapterHeaderPattern.FindStringSubmatch(trimmed); m != nil {
					isNewChapter = true
					chapterTitle = m[1]
					if len(m) > 2 && m[2] != "" {
						chapterTitle += ": " + m[2]
					}
				} else if m := numberedHeadingPattern.FindStringSubmatch(trimmed); m != nil && currentChapter == nil {
					isNewChapter = true
					chapterTitle = m[1] + ". " + m[2]
				}
			}

			if isNewChapter {
				if currentChapter != nil && len(currentChapter.Lines) > 0 {
					currentChapter.Content = s.transformer.Transform([][]string{currentChapter.Lines})
					currentChapter.WordCount = len(strings.Fields(currentChapter.Content))
					currentChapter.TokenEstimate = int(float64(len(currentChapter.Content)) / 3.8)
					rawChapters = append(rawChapters, *currentChapter)
				}

				slug := sanitizeSlug(chapterTitle)
				currentChapter = &core.Chapter{
					Title:     chapterTitle,
					Slug:      slug,
					StartPage: page.PageNumber,
					Lines:     []string{line},
				}
			} else {
				if currentChapter == nil {
					currentChapter = &core.Chapter{
						Title:     "Introduction & Overview",
						Slug:      "introduction_and_overview",
						StartPage: page.PageNumber,
						Lines:     []string{},
					}
				}
				currentChapter.Lines = append(currentChapter.Lines, line)
			}
		}
	}

	if currentChapter != nil && len(currentChapter.Lines) > 0 {
		currentChapter.Content = s.transformer.Transform([][]string{currentChapter.Lines})
		currentChapter.WordCount = len(strings.Fields(currentChapter.Content))
		currentChapter.TokenEstimate = int(float64(len(currentChapter.Content)) / 3.8)
		rawChapters = append(rawChapters, *currentChapter)
	}

	// 2. Post-Processing Filter: Eliminate front-matter stub chapters and merge non-content fragments
	var finalChapters []core.Chapter

	for _, ch := range rawChapters {
		// If this is an intro stub with fewer than 30 words and there are other chapters, skip it
		if ch.Slug == "introduction_and_overview" && ch.WordCount < 30 && len(rawChapters) > 1 {
			continue
		}

		// If chapter has fewer than 20 words, merge into previous chapter
		if ch.WordCount < 20 && len(finalChapters) > 0 {
			last := &finalChapters[len(finalChapters)-1]
			last.Lines = append(last.Lines, ch.Lines...)
			last.Content = s.transformer.Transform([][]string{last.Lines})
			last.WordCount = len(strings.Fields(last.Content))
			last.TokenEstimate = int(float64(len(last.Content)) / 3.8)
			continue
		}

		finalChapters = append(finalChapters, ch)
	}

	// If no valid chapters extracted, fallback to single document chapter
	if len(finalChapters) == 0 {
		var allLines []string
		for _, p := range pages {
			if !p.IsFilteredOut {
				allLines = append(allLines, p.Lines...)
			}
		}
		content := s.transformer.Transform([][]string{allLines})
		finalChapters = append(finalChapters, core.Chapter{
			Index:         1,
			Title:         pdfTitle,
			Slug:          sanitizeSlug(pdfTitle),
			Filename:      "01_content.md",
			StartPage:     1,
			Lines:         allLines,
			Content:       content,
			WordCount:     len(strings.Fields(content)),
			TokenEstimate: int(float64(len(content)) / 3.8),
		})
	}

	// 3. Re-index and assign filenames
	for i := range finalChapters {
		finalChapters[i].Index = i + 1
		finalChapters[i].Filename = fmt.Sprintf("%02d_%s.md", i+1, finalChapters[i].Slug)
	}

	// 4. Link Previous and Next navigation pointers
	for i := range finalChapters {
		if i > 0 {
			finalChapters[i].PrevFilename = finalChapters[i-1].Filename
		}
		if i < len(finalChapters)-1 {
			finalChapters[i].NextFilename = finalChapters[i+1].Filename
		}
	}

	return finalChapters
}

// GenerateTOCIndex creates a human-readable master README.md index linking to each chapter.
func (s *Splitter) GenerateTOCIndex(pdfName string, chapters []core.Chapter, totalPages int) string {
	var sb strings.Builder

	cleanName := strings.ReplaceAll(pdfName, "_", " ")
	cleanName = strings.ReplaceAll(cleanName, "-", " ")

	totalWords := 0
	totalTokens := 0
	for _, ch := range chapters {
		totalWords += ch.WordCount
		totalTokens += ch.TokenEstimate
	}

	sb.WriteString(fmt.Sprintf("# %s\n\n", cleanName))
	sb.WriteString("Extracted and transformed from PDF into chapter-based Markdown notes.\n\n")
	sb.WriteString("> [!TIP]\n")
	sb.WriteString("> **AI Agent Manifest Available**: For LLM chunking, vector ingestion, and multi-file agent workflows, see [`AGENTS.md`](./AGENTS.md).\n\n")
	sb.WriteString("---\n\n")
	sb.WriteString("## Table of Contents\n\n")
	sb.WriteString("| # | Chapter Title | Start Page | Words | Est. Tokens | Notes File |\n")
	sb.WriteString("| :- | :--- | :--- | :--- | :--- | :--- |\n")

	for _, ch := range chapters {
		sb.WriteString(fmt.Sprintf("| **%02d** | %s | Page %d | %d | ~%d | [`%s`](./%s) |\n",
			ch.Index, ch.Title, ch.StartPage, ch.WordCount, ch.TokenEstimate, ch.Filename, ch.Filename))
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString("### Summary Statistics\n")
	sb.WriteString(fmt.Sprintf("- **Total Chapters Extracted:** %d\n", len(chapters)))
	sb.WriteString(fmt.Sprintf("- **Total PDF Pages:** %d\n", totalPages))
	sb.WriteString(fmt.Sprintf("- **Total Word Count:** %d words (~%d tokens)\n", totalWords, totalTokens))
	sb.WriteString("- **Noise Filtered:** Image placeholders, header/footer numbers, and back-matter indexes automatically excluded.\n")

	return sb.String()
}

// GenerateAgentManifest builds an LLM/RAG-friendly AGENTS.md document with structured context maps.
func (s *Splitter) GenerateAgentManifest(pdfName string, chapters []core.Chapter, totalPages int) string {
	var sb strings.Builder

	cleanName := strings.ReplaceAll(pdfName, "_", " ")
	cleanName = strings.ReplaceAll(cleanName, "-", " ")

	type AgentChapterEntry struct {
		Index         int    `json:"chapter_index"`
		Title         string `json:"title"`
		Filename      string `json:"file_path"`
		StartPage     int    `json:"start_page"`
		WordCount     int    `json:"word_count"`
		TokenEstimate int    `json:"estimated_tokens"`
	}

	var entries []AgentChapterEntry
	totalWords := 0
	totalTokens := 0
	for _, ch := range chapters {
		entries = append(entries, AgentChapterEntry{
			Index:         ch.Index,
			Title:         ch.Title,
			Filename:      ch.Filename,
			StartPage:     ch.StartPage,
			WordCount:     ch.WordCount,
			TokenEstimate: ch.TokenEstimate,
		})
		totalWords += ch.WordCount
		totalTokens += ch.TokenEstimate
	}

	sb.WriteString(fmt.Sprintf("# AI Agent Ingestion Manifest: %s\n\n", cleanName))
	sb.WriteString("This manifest is optimized for autonomous coding agents, LLM RAG pipelines, and vector indexers.\n\n")

	sb.WriteString("## System Prompt & Usage Instructions for Agents\n\n")
	sb.WriteString("```markdown\n")
	sb.WriteString("You are navigating a structured, chapter-split technical knowledge base.\n")
	sb.WriteString("- Retrieve content from individual chapter files listed below.\n")
	sb.WriteString("- Each file begins with structured YAML frontmatter containing precise metadata.\n")
	sb.WriteString("- When citing technical instructions, reference the chapter number and section headings.\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Structured JSON Manifest (Machine-Readable)\n\n")
	sb.WriteString("```json\n")
	manifestJSON, _ := json.MarshalIndent(map[string]interface{}{
		"document_title":   cleanName,
		"total_chapters":   len(chapters),
		"total_pages":      totalPages,
		"total_words":      totalWords,
		"estimated_tokens": totalTokens,
		"chapters":         entries,
	}, "", "  ")
	sb.WriteString(string(manifestJSON))
	sb.WriteString("\n```\n\n")

	sb.WriteString("## Sequential Chapter Navigation Map\n\n")
	sb.WriteString("| Index | File | Start Page | Est. Tokens | Scope |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("| `%02d` | [`%s`](./%s) | Page %d | ~%d | %s |\n",
			e.Index, e.Filename, e.Filename, e.StartPage, e.TokenEstimate, e.Title))
	}

	return sb.String()
}

func sanitizeSlug(s string) string {
	prefixRegex := regexp.MustCompile(`(?i)^(chapter|part|module|section)\s+\d+[:.]?\s*`)
	cleaned := prefixRegex.ReplaceAllString(s, "")

	cleaned = strings.ToLower(cleaned)
	cleaned = strings.ReplaceAll(cleaned, ":", " ")
	cleaned = strings.ReplaceAll(cleaned, ".", " ")
	cleaned = strings.ReplaceAll(cleaned, "/", " ")
	cleaned = strings.ReplaceAll(cleaned, "\\", " ")
	cleaned = strings.ReplaceAll(cleaned, "-", " ")

	var words []string
	for _, word := range strings.Fields(cleaned) {
		cleanWord := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return r
			}
			return -1
		}, word)
		if cleanWord != "" {
			words = append(words, cleanWord)
		}
	}

	if len(words) == 0 {
		return "chapter"
	}
	if len(words) > 6 {
		words = words[:6]
	}
	return strings.Join(words, "_")
}
