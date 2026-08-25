package transformer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/devops/pdf2md/pkg/core"
)

// Transformer implements the core.Transformer interface.
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

	if t.config.DetectCodeBlocks {
		allLines = t.detectCodeBlocks(allLines)
	}

	allLines = t.formatHeadings(allLines)
	allLines = t.formatLists(allLines)

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

func (t *Transformer) detectCodeBlocks(lines []string) []string {
	var result []string
	inCode := false

	codePattern := regexp.MustCompile(`^(\s{4,}|\t|\$\s+|#\s*!/|def\s+\w+|function\s+\w+|import\s+[\w.]+|package\s+\w+|apiVersion:|kind:)`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if codePattern.MatchString(line) && trimmed != "" {
			if !inCode {
				inCode = true
				result = append(result, "```bash")
			}
			result = append(result, strings.TrimRight(line, " \t\r\n"))
		} else {
			if inCode && trimmed == "" {
				result = append(result, "")
			} else if inCode {
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
