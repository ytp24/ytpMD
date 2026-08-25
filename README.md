# 📄 `pdf2md` — High-Performance PDF to Markdown CLI in Go

`pdf2md` is a standalone, fast, and production-grade command-line tool written in Go that extracts text from PDF documents and transforms them into clean, structured GitHub-Flavored Markdown (`.md`).

---

## ✨ Key Features

- 🚫 **Excludes Non-Usable Sections (Appendix / Index / Glossary / References):**
  Automatically detects when an Appendix (e.g. `Appendix A`, `Appendix B`), Subject Index, Bibliography, or Glossary begins and cleanly truncates processing so your final Markdown contains only core usable content.
- 🖼️ **Excludes Images & Binary Assets:**
  Strips out binary figure references, image placeholders, and messy vector artifacts.
- 📑 **Excludes Front-Matter Noise:**
  Skip cover pages, publisher notices, copyright pages, and blank pages using `-skip-front <pages>`.
- 🔄 **Smart Paragraph Reflow & De-Hyphenation:**
  Combines split lines into continuous paragraphs and automatically rejoins broken hyphenated words (`archi-\ntecture` -> `architecture`).
- 🏷️ **Heading & Code Block Detection:**
  Automatically formats chapter headers (`# Chapter 1`), section headers (`## 1.2 Architecture`), and indented code blocks (```` ``` ````).
- 🧹 **Header & Footer Stripping:**
  Eliminates repeated running headers, "Page X of Y", and isolated page numbers.
- ⚡ **Zero CGo / Native Speed:**
  Compiled directly into a self-contained Go binary with minimal memory footprint.

---

## 🛠️ Build & Installation

### 1. Prerequisites
Ensure `pdftotext` (Poppler utilities) is installed on your Linux system:
```bash
sudo apt-get install poppler-utils
```

### 2. Build the Binary
```bash
cd ~/SHARED/Projects/pdf2md
go build -o bin/pdf2md ./cmd/pdf2md
```

### 3. (Optional) Install Globally to `$GOPATH/bin` or `/usr/local/bin`
```bash
go install ./cmd/pdf2md
# or
sudo cp bin/pdf2md /usr/local/bin/
```

---

## 🚀 CLI Usage & Commands

### 1. Convert a Single PDF
```bash
# Convert with automatic appendix and noise removal
./bin/pdf2md convert input.pdf -o output.md

# Skip first 4 pages (e.g. Cover, Copyright, TOC)
./bin/pdf2md convert book.pdf -o book.md -skip-front 4

# Extract specific page range
./bin/pdf2md convert guide.pdf -o guide.md -start-page 10 -end-page 85
```

### 2. Batch Convert a Directory of PDFs
```bash
# Convert all PDFs in a folder to a target markdown folder
./bin/pdf2md batch /path/to/pdfs/ -o /path/to/markdown/ -r
```

### 3. Check Version
```bash
./bin/pdf2md version
```

---

## ⚙️ Options & Flags Reference

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-o`, `-output <path>` | Destination `.md` file path | `<input_name>.md` |
| `-skip-front <N>` | Number of initial pages to skip (covers/TOC) | `0` |
| `-start-page <N>` | Start page number (1-indexed) | `1` |
| `-end-page <N>` | Stop converting at this page number | `0` (End of file) |
| `-keep-appendix` | Include appendix, index, and bibliography pages | `false` (Excluded by default) |
| `-no-reflow` | Disable paragraph line reflow | `false` |
| `-r`, `-recursive` | Recursively process subdirectories in batch mode | `false` |

---

## 🧪 Running Unit Tests
```bash
cd ~/SHARED/Projects/pdf2md
go test -v ./...
```
