package splitter

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/devops/pdf2md/pkg/core"
	"github.com/devops/pdf2md/pkg/transformer"
)

// Splitter implements the core.Splitter interface.
type Splitter struct {
	config      core.Config
	transformer core.Transformer
}

// NewSplitter initializes a new Splitter using Factory Pattern.
func NewSplitter(cfg core.Config) *Splitter {
	return &Splitter{
		config:      cfg,
		transformer: transformer.NewTransformer(cfg),
	}
}

// SplitIntoChapters groups extracted page lines into distinct chapter structures.
func (s *Splitter) SplitIntoChapters(pages []core.PDFPage, pdfTitle string) []core.Chapter {
	var chapters []core.Chapter
	var currentChapter *core.Chapter

	chapterHeaderPattern := regexp.MustCompile(`(?i)^(CHAPTER\s+\d+|PART\s+[IVXLCDM]+|MODULE\s+\d+|SECTION\s+\d+)[:.]?\s*(.*)$`)
	numberedHeadingPattern := regexp.MustCompile(`^(\d+)\.\s+([A-Z][A-Za-z0-9\s,-]{3,50})$`)

	chapterIndex := 1

	for _, page := range pages {
		if page.IsFilteredOut {
			continue
		}

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

			if isNewChapter {
				if currentChapter != nil && len(currentChapter.Lines) > 0 {
					currentChapter.Content = s.transformer.Transform([][]string{currentChapter.Lines})
					chapters = append(chapters, *currentChapter)
				}

				slug := sanitizeSlug(chapterTitle)
				filename := fmt.Sprintf("%02d_%s.md", chapterIndex, slug)

				currentChapter = &core.Chapter{
					Index:     chapterIndex,
					Title:     chapterTitle,
					Slug:      slug,
					Filename:  filename,
					StartPage: page.PageNumber,
					Lines:     []string{line},
				}
				chapterIndex++
			} else {
				if currentChapter == nil {
					currentChapter = &core.Chapter{
						Index:     chapterIndex,
						Title:     "Introduction & Overview",
						Slug:      "introduction_and_overview",
						Filename:  fmt.Sprintf("%02d_introduction_and_overview.md", chapterIndex),
						StartPage: page.PageNumber,
						Lines:     []string{},
					}
					chapterIndex++
				}
				currentChapter.Lines = append(currentChapter.Lines, line)
			}
		}
	}

	if currentChapter != nil && len(currentChapter.Lines) > 0 {
		currentChapter.Content = s.transformer.Transform([][]string{currentChapter.Lines})
		chapters = append(chapters, *currentChapter)
	}

	if len(chapters) == 0 {
		var allLines []string
		for _, p := range pages {
			if !p.IsFilteredOut {
				allLines = append(allLines, p.Lines...)
			}
		}
		content := s.transformer.Transform([][]string{allLines})
		chapters = append(chapters, core.Chapter{
			Index:     1,
			Title:     pdfTitle,
			Slug:      sanitizeSlug(pdfTitle),
			Filename:  "01_content.md",
			StartPage: 1,
			Lines:     allLines,
			Content:   content,
		})
	}

	return chapters
}

// GenerateTOCIndex creates a master README.md index linking to each chapter file.
func (s *Splitter) GenerateTOCIndex(pdfName string, chapters []core.Chapter, totalPages int) string {
	var sb strings.Builder

	cleanName := strings.ReplaceAll(pdfName, "_", " ")
	cleanName = strings.ReplaceAll(cleanName, "-", " ")

	sb.WriteString(fmt.Sprintf("# %s\n\n", cleanName))
	sb.WriteString("Extracted and transformed from PDF into chapter-based Markdown notes.\n\n")
	sb.WriteString("---\n\n")
	sb.WriteString("## Table of Contents\n\n")
	sb.WriteString("| # | Chapter Title | Start Page | Notes File |\n")
	sb.WriteString("| :- | :--- | :--- | :--- |\n")

	for _, ch := range chapters {
		sb.WriteString(fmt.Sprintf("| **%02d** | %s | Page %d | [`%s`](./%s) |\n",
			ch.Index, ch.Title, ch.StartPage, ch.Filename, ch.Filename))
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString("### Document Summary\n")
	sb.WriteString(fmt.Sprintf("- **Total Chapters Extracted:** %d\n", len(chapters)))
	sb.WriteString(fmt.Sprintf("- **Total PDF Pages:** %d\n", totalPages))
	sb.WriteString("- **Non-usable Assets (Images / Appendix / Index):** Automatically Filtered & Excluded.\n")

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
