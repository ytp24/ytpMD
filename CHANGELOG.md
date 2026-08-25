# Changelog

All notable changes to **ytpMD [pdf2md]** are documented in this file in compliance with [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [3.1.0] - 2026-08-26

### Added
- **Clean Architecture Refactoring**: Domain models, interfaces, and factory pattern instantiation.
- **Concurrent Batch Processing**: Multi-core Goroutine worker pool with context cancellation.
- **Animated Teal Gradient Progress Bar**: Shifts from deep dark teal to glowing pale teal near completion.
- **Windows & Linux Cross-Platform Dialogs**: Added native PowerShell GUI file and folder picker integration for Windows.
- **Open Source Governance**: Added Apache 2.0 `LICENSE`, `LEGAL.md`, `SECURITY.md`, and `CONTRIBUTING.md`.
- **Linux Debian Packaging**: Added `.deb` packaging generator and automated one-line install scripts.

### Changed
- Default destination folder set to `~/Documents/ytp24`.
- Standardized CLI entrypoint to global command `ytp24`.
- Removed all emojis for a clean, professional, italic teal terminal UI.

---

## [2.0.0] - 2026-08-26

### Added
- Interactive prompt wizard with smart production defaults.
- Native GUI file/folder chooser fallback via `zenity` / `kdialog`.
- Chapter-by-chapter splitting with Table of Contents `README.md` index generation.

---

## [1.0.0] - 2026-08-25

### Added
- Initial release in Go (Golang).
- Appendix, Index, and Bibliography automatic cutoff filter.
- Smart paragraph reflow and de-hyphenation engine.
