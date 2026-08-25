# Contributing to ytpMD [pdf2md]

Thank you for your interest in contributing to **ytpMD**! We welcome contributions from the open-source community.

---

## Development Setup

### Prerequisites
- **Go**: Version `>= 1.22`
- **Poppler Utilities**: `sudo apt install poppler-utils` (Linux) or `brew install poppler` (macOS) or `choco install poppler` (Windows).

### Build & Test Workflow
```bash
# 1. Clone the repository
git clone https://github.com/ytp24/ytpMD.git
cd ytpMD

# 2. Run formatting and tests
make all

# 3. Install binary locally
make install
```

---

## Pull Request Guidelines

1. **Branch Naming**:
   - `feat/<feature-name>` for new capabilities.
   - `fix/<bug-name>` for bug fixes.
   - `docs/<doc-name>` for documentation improvements.
2. **Code Style**:
   - Run `go fmt ./...` before committing.
   - Ensure all unit tests pass with `go test -v ./...`.
   - Maintain Clean Architecture and SOLID principles.
3. **Commit Messages**: Follow Conventional Commits:
   - `feat: add OCR fallback extraction`
   - `fix: handle edge case in hyphenated bullet points`
   - `docs: update Windows installation instructions`

---

## Code of Conduct
All contributors and participants agree to abide by our [`CODE_OF_CONDUCT.md`](./CODE_OF_CONDUCT.md).
