package transformer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/ytp24/ytpMD/pkg/core"
)

// Transformer implements the core.Transformer interface with agentic optimizations.
type Transformer struct {
	config core.Config
}

// NewTransformer initializes a new Transformer using Factory Pattern.
func NewTransformer(cfg core.Config) *Transformer {
	return &Transformer{config: cfg}
}

// Transform converts extracted page line arrays into clean Markdown.
func (t *Transformer) Transform(pagesLines [][]string) string {
	var allLines []string
	for _, pLines := range pagesLines {
		allLines = append(allLines, pLines...)
	}

	// 1. Format headings & lists first
	allLines = t.formatHeadings(allLines)
	allLines = t.formatLists(allLines)

	// 2. Detect genuine code blocks (strict syntax matching, no indentation false-positives)
	if t.config.DetectCodeBlocks {
		allLines = t.detectCodeBlocks(allLines)
	}

	var output string
	if t.config.ReflowParagraphs {
		output = t.reflowParagraphs(allLines)
	} else {
		output = strings.Join(allLines, "\n")
	}

	multiNewline := regexp.MustCompile(`\n{3,}`)
	output = multiNewline.ReplaceAllString(output, "\n\n")
	return strings.TrimSpace(output) + "\n"
}

// GenerateChapterMarkdown builds an agent-ready markdown document with YAML frontmatter & navigation breadcrumbs.
func (t *Transformer) GenerateChapterMarkdown(ch core.Chapter, sourcePDF string, totalChapters int) string {
	var sb strings.Builder

	cleanSource := strings.TrimSuffix(sourcePDF, ".pdf")
	cleanSource = strings.TrimSuffix(cleanSource, ".PDF") + ".pdf"

	words := countWords(ch.Content)
	tokens := estimateTokens(ch.Content)

	// 1. Agentic YAML Frontmatter
	if t.config.AddYAMLFrontmatter {
		sb.WriteString("---\n")
		sb.WriteString(fmt.Sprintf("title: %q\n", ch.Title))
		sb.WriteString(fmt.Sprintf("chapter: %d\n", ch.Index))
		sb.WriteString(fmt.Sprintf("total_chapters: %d\n", totalChapters))
		sb.WriteString(fmt.Sprintf("source_document: %q\n", cleanSource))
		sb.WriteString(fmt.Sprintf("start_page: %d\n", ch.StartPage))
		sb.WriteString(fmt.Sprintf("word_count: %d\n", words))
		sb.WriteString(fmt.Sprintf("estimated_tokens: %d\n", tokens))
		sb.WriteString("agent_instructions: \"Cite section headers and use code snippets directly when referencing this document.\"\n")
		sb.WriteString("---\n\n")
	}

	// 2. Navigation Breadcrumbs
	sb.WriteString("<div align=\"center\">\n\n")
	var navItems []string
	if ch.PrevFilename != "" {
		navItems = append(navItems, fmt.Sprintf("[« Previous Chapter](./%s)", ch.PrevFilename))
	}
	navItems = append(navItems, "[Table of Contents](./README.md)", "[Agent Manifest](./AGENTS.md)")
	if ch.NextFilename != "" {
		navItems = append(navItems, fmt.Sprintf("[Next Chapter »](./%s)", ch.NextFilename))
	}
	sb.WriteString(strings.Join(navItems, " • ") + "\n\n")
	sb.WriteString("</div>\n\n---\n\n")

	// 3. Main Content
	sb.WriteString(ch.Content)
	sb.WriteString("\n\n---\n\n")

	// 4. Footer Breadcrumbs
	sb.WriteString("<div align=\"center\">\n\n")
	sb.WriteString(strings.Join(navItems, " • ") + "\n\n")
	sb.WriteString("</div>\n")

	return sb.String()
}

func (t *Transformer) detectCodeBlocks(lines []string) []string {
	var result []string
	inCode := false
	codeLang := "bash"

	goPattern := regexp.MustCompile(`^\s*(package\s+\w+|import\s+\(|func\s+\w+|type\s+\w+\s+struct)`)
	pythonPattern := regexp.MustCompile(`^\s*(def\s+\w+\(|import\s+[\w.]+|class\s+\w+[:(]|from\s+\w+\s+import)`)
	yamlPattern := regexp.MustCompile(`^\s*(apiVersion:\s*[\w./]+|kind:\s+[A-Z]\w+|spec:\s*$)`)
	jsonPattern := regexp.MustCompile(`^\s*(\{\s*"|"[\w_-]+":\s*["{\[\d])`)
	bashPattern := regexp.MustCompile(`^\s*(\$\s+|#!\/bin\/|sudo\s+|kubectl\s+|docker\s+|terraform\s+|ansible\s+|helm\s+|git\s+|aws\s+|az\s+|curl\s+|npm\s+|pip\s+)`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Never treat Markdown headings as code blocks
		if strings.HasPrefix(trimmed, "#") {
			if inCode {
				inCode = false
				result = append(result, "```")
			}
			result = append(result, line)
			continue
		}

		var detectedLang string
		switch {
		case goPattern.MatchString(line):
			detectedLang = "go"
		case pythonPattern.MatchString(line):
			detectedLang = "python"
		case yamlPattern.MatchString(line):
			detectedLang = "yaml"
		case jsonPattern.MatchString(line):
			detectedLang = "json"
		case bashPattern.MatchString(line):
			detectedLang = "bash"
		}

		if detectedLang != "" && trimmed != "" {
			if !inCode {
				inCode = true
				codeLang = detectedLang
				result = append(result, "```"+codeLang)
			}
			result = append(result, line)
		} else {
			if inCode && trimmed == "" {
				result = append(result, "")
			} else if inCode && !isCodeContinuation(trimmed) {
				inCode = false
				result = append(result, "```")
				result = append(result, line)
			} else {
				result = append(result, line)
			}
		}
	}
	if inCode {
		result = append(result, "```")
	}
	return result
}

func isCodeContinuation(trimmed string) bool {
	// If line ends with typical code punctuation or starts with symbols/indentation
	if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "}") || strings.HasPrefix(trimmed, "]") {
		return true
	}
	if strings.Contains(trimmed, "=") || strings.Contains(trimmed, ":") || strings.Contains(trimmed, "{") || strings.Contains(trimmed, "}") {
		return true
	}
	return false
}

func (t *Transformer) formatHeadings(lines []string) []string {
	var result []string

	h1Pattern := regexp.MustCompile(`(?i)^(CHAPTER\s+\d+|PART\s+[IVXLCDM]+|MODULE\s+\d+)[:.]?\s*(.*)$`)
	h2Pattern := regexp.MustCompile(`^(\d+\.\d+)\s+([A-Z].*)$`)
	h3Pattern := regexp.MustCompile(`^(\d+\.\d+\.\d+)\s+([A-Z].*)$`)
	h1NumPattern := regexp.MustCompile(`^(\d+)\.\s+([A-Z][A-Za-z0-9\s,-]{3,50})$`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "#") {
			result = append(result, line)
			continue
		}

		if m := h1Pattern.FindStringSubmatch(trimmed); m != nil {
			title := m[1]
			if len(m) > 2 && m[2] != "" {
				title += ": " + m[2]
			}
			result = append(result, "\n# "+strings.TrimSpace(title)+"\n")
			continue
		}

		if m := h3Pattern.FindStringSubmatch(trimmed); m != nil {
			result = append(result, "\n### "+m[1]+" "+strings.TrimSpace(m[2])+"\n")
			continue
		}

		if m := h2Pattern.FindStringSubmatch(trimmed); m != nil {
			result = append(result, "\n## "+m[1]+" "+strings.TrimSpace(m[2])+"\n")
			continue
		}

		if m := h1NumPattern.FindStringSubmatch(trimmed); m != nil && len(trimmed) < 60 {
			result = append(result, "\n# "+m[1]+". "+strings.TrimSpace(m[2])+"\n")
			continue
		}

		if isAllUpper(trimmed) && len(trimmed) > 4 && len(trimmed) < 50 && !strings.HasSuffix(trimmed, ".") {
			result = append(result, "\n## "+toTitle(trimmed)+"\n")
			continue
		}

		result = append(result, line)
	}

	return result
}

func (t *Transformer) formatLists(lines []string) []string {
	var result []string
	bulletPattern := regexp.MustCompile(`^[•▪‣⁃–]\s*(.*)$`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if m := bulletPattern.FindStringSubmatch(trimmed); m != nil {
			indentLen := len(line) - len(strings.TrimLeft(line, " "))
			indent := strings.Repeat(" ", indentLen)
			result = append(result, indent+"- "+m[1])
		} else {
			result = append(result, line)
		}
	}
	return result
}

func (t *Transformer) reflowParagraphs(lines []string) string {
	var paragraphs []string
	var currentP []string
	inCode := false

	listPattern := regexp.MustCompile(`^(\s*[-*]\s+|\s*\d+\.\s+)`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if len(currentP) > 0 {
				paragraphs = append(paragraphs, t.joinLines(currentP))
				currentP = nil
			}
			inCode = !inCode
			paragraphs = append(paragraphs, line)
			continue
		}

		if inCode {
			paragraphs = append(paragraphs, line)
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			if len(currentP) > 0 {
				paragraphs = append(paragraphs, t.joinLines(currentP))
				currentP = nil
			}
			paragraphs = append(paragraphs, line)
			continue
		}

		if listPattern.MatchString(trimmed) {
			if len(currentP) > 0 {
				paragraphs = append(paragraphs, t.joinLines(currentP))
				currentP = nil
			}
			paragraphs = append(paragraphs, line)
			continue
		}

		if trimmed == "" {
			if len(currentP) > 0 {
				paragraphs = append(paragraphs, t.joinLines(currentP))
				currentP = nil
			}
			continue
		}

		currentP = append(currentP, trimmed)
	}

	if len(currentP) > 0 {
		paragraphs = append(paragraphs, t.joinLines(currentP))
	}

	return strings.Join(paragraphs, "\n\n")
}

func (t *Transformer) joinLines(lines []string) string {
	var sb strings.Builder
	for i, l := range lines {
		if i == 0 {
			sb.WriteString(l)
		} else {
			curr := sb.String()
			if strings.HasSuffix(curr, "-") && !strings.HasSuffix(curr, " -") {
				trimmed := strings.TrimSuffix(curr, "-")
				sb.Reset()
				sb.WriteString(trimmed + l)
			} else {
				sb.WriteString(" " + l)
			}
		}
	}
	return sb.String()
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func estimateTokens(s string) int {
	return int(float64(len(s)) / 3.8)
}

func isAllUpper(s string) bool {
	hasLetters := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetters = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetters
}

func toTitle(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {
		if len(w) > 0 {
			r := []rune(w)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}
