# ytpMD [pdf2md] — Production Cloud & Document Pipeline in Go

`ytpMD` (also invocable as `pdf2md` or `ytp24`) is an enterprise-grade, concurrent CLI tool and Go library designed to transform PDF documents into chapter-aware Markdown notes with zero noise.

---

## 🏛️ Clean Architecture & Design Patterns

The codebase strictly adheres to **SOLID principles**, the **Standard Go Project Layout**, and established software design patterns:

- **Factory Pattern**: Centralized component instantiation (`NewPDFExtractor`, `NewContentFilter`, `NewTransformer`, `NewSplitter`, `NewConcurrentBatchEngine`).
- **Strategy / Pipeline Pattern**: Decoupled document processing stages:
  $$\text{Source PDF} \xrightarrow{\text{Extractor}} \text{Raw Text} \xrightarrow{\text{Filter}} \text{Sanitized Lines} \xrightarrow{\text{Transformer}} \text{Markdown} \xrightarrow{\text{Splitter}} \text{TOC Chapters}$$
- **Worker Pool (Concurrency Pattern)**: High-throughput batch processing using bounded worker channels and context cancellation (`pkg/batch/pool.go`).
- **Observer Pattern**: Dynamic progress reporting abstracted through `core.ProgressReporter`, driving the multi-shade teal terminal progress bar.
- **Dependency Inversion Principle (DIP)**: Concrete engines depend on domain abstractions defined in [`pkg/core/interfaces.go`](./pkg/core/interfaces.go).

---

## 📂 Package Layout

```
├── cmd/
│   ├── ytp24/main.go            # Primary CLI entrypoint (invokes pkg/cli.Run)
│   └── pdf2md/main.go           # Canonical tool alias
├── pkg/
│   ├── core/                    # Domain models, entities & interfaces (DIP)
│   │   ├── interfaces.go
│   │   └── models.go
│   ├── batch/                   # Concurrent worker pool engine
│   │   ├── pool.go
│   │   └── pool_test.go
│   ├── extractor/               # PDF extraction implementation
│   │   └── extractor.go
│   ├── filter/                  # Noise, appendix & asset removal
│   │   ├── filter.go
│   │   └── filter_test.go
│   ├── splitter/                # Table of Contents & Chapter segmentation
│   │   ├── splitter.go
│   │   └── splitter_test.go
│   ├── transformer/             # Markdown formatting & de-hyphenation
│   │   ├── transformer.go
│   │   └── transformer_test.go
│   ├── ui/                      # Terminal UI, dialogs & teal gradient progress bar
│   │   ├── dialog.go
│   │   ├── progressbar.go
│   │   └── prompt.go
│   ├── validator/               # Pre-flight dependency & magic byte checks
│   │   ├── validator.go
│   │   └── validator_test.go
│   └── cli/                     # Top-level CLI command routing & panic handler
│       └── app.go
├── examples/                    # Test input PDFs & sample batch
├── Makefile                     # Build, test, format & install automation
├── LICENSE                      # MIT License
└── go.mod
```

---

## 🛠️ Build & Installation

```bash
# Clone or navigate to the repository
cd ~/SHARED/Projects/pdf2md

# Format, test, and install globally to ~/.local/bin/
make install
```

---

## 🚀 Usage

### 1. Interactive Mode
Run without arguments to launch the wizard with native GUI file selection:
```bash
ytp24
```

### 2. Single Document Conversion
```bash
ytp24 convert book.pdf -skip-front 3
```

### 3. Concurrent Batch Processing
```bash
ytp24 batch ~/Downloads/PDFs/ -name CloudEngineering -concurrency 6
```
