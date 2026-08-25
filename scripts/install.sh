#!/usr/bin/env bash
# ==============================================================================
# ytpMD [pdf2md] - Official Unix/Linux/macOS Installer
# ==============================================================================
set -e

INSTALL_DIR="${HOME}/.local/bin"
VERSION="3.1.0"

echo "[*] Installing ytpMD (ytp24 v${VERSION})..."

# Check dependencies
if ! command -v pdftotext >/dev/null 2>&1; then
    echo "[!] Warning: 'pdftotext' not found. Please install poppler-utils:"
    echo "    Ubuntu/Debian: sudo apt install poppler-utils"
    echo "    Fedora/RHEL:   sudo dnf install poppler-utils"
    echo "    Arch Linux:    sudo pacman -S poppler"
    echo "    macOS:         brew install poppler"
fi

mkdir -p "${INSTALL_DIR}"

if command -v go >/dev/null 2>&1; then
    echo "[*] Building binary from source..."
    go build -ldflags="-s -w" -o "${INSTALL_DIR}/ytp24" ./cmd/ytp24
    ln -sf "${INSTALL_DIR}/ytp24" "${INSTALL_DIR}/pdf2md"
    echo "[+] Successfully installed ytp24 and pdf2md to ${INSTALL_DIR}/"
else
    echo "[x] Error: Go toolchain not found. Please install Go or use prebuilt .deb package."
    exit 1
fi

echo ""
echo "Installation complete! Run 'ytp24' from your terminal."
