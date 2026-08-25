#!/usr/bin/env bash
set -e

INSTALL_DIR="${HOME}/.local/bin"
VERSION="3.2.0"

echo "[*] Installing ytpMD (v${VERSION})..."

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
    go build -ldflags="-s -w" -o "${INSTALL_DIR}/ytpMD" ./cmd/ytpMD
    ln -sf "${INSTALL_DIR}/ytpMD" "${INSTALL_DIR}/ytp24"
    ln -sf "${INSTALL_DIR}/ytpMD" "${INSTALL_DIR}/pdf2md"
    echo "[+] Successfully installed ytpMD, ytp24, and pdf2md to ${INSTALL_DIR}/"
else
    echo "[x] Error: Go toolchain not found. Please install Go or use prebuilt .deb package."
    exit 1
fi

echo ""
echo "Installation complete! Run 'ytpMD' from your terminal."
