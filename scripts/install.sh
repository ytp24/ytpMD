#!/usr/bin/env bash
set -e

VERSION="3.2.0"
INSTALL_DIR="${HOME}/.local/bin"
REPO="ytp24/ytpMD"

# ANSI Colors
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${CYAN}====================================================${NC}"
echo -e "${CYAN}   ytpMD [pdf2md] Installer (v${VERSION})${NC}"
echo -e "${CYAN}   High-Performance PDF to Chapter Markdown Engine${NC}"
echo -e "${CYAN}====================================================${NC}"
echo ""

# 1. Dependency check
if ! command -v pdftotext >/dev/null 2>&1; then
    echo -e "${YELLOW}[!] Warning: 'pdftotext' was not found in your PATH.${NC}"
    echo -e "    ytpMD requires poppler-utils. Please install it using:"
    echo -e "    - Ubuntu/Debian:  ${GREEN}sudo apt install poppler-utils${NC}"
    echo -e "    - Fedora/RHEL:    ${GREEN}sudo dnf install poppler-utils${NC}"
    echo -e "    - Arch Linux:     ${GREEN}sudo pacman -S poppler${NC}"
    echo -e "    - macOS:          ${GREEN}brew install poppler${NC}"
    echo ""
fi

# 2. Architecture & OS detection
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "${ARCH}" in
    x86_64|amd64)
        GOARCH="amd64"
        ;;
    arm64|aarch64)
        GOARCH="arm64"
        ;;
    *)
        echo -e "${RED}[x] Unsupported architecture: ${ARCH}${NC}"
        exit 1
        ;;
esac

case "${OS}" in
    linux)
        GOOS="linux"
        ;;
    darwin)
        GOOS="darwin"
        ;;
    *)
        echo -e "${RED}[x] Unsupported OS: ${OS}. For Windows, please use install.ps1 in PowerShell.${NC}"
        exit 1
        ;;
esac

mkdir -p "${INSTALL_DIR}"

# 3. Installation strategy: If local source directory with go.mod is present, compile; otherwise download binary
INSTALLED=false

if [ -f "./go.mod" ] && command -v go >/dev/null 2>&1; then
    echo -e "${CYAN}[*] Building binary from local repository source with Go...${NC}"
    CGO_ENABLED=0 go build -ldflags="-s -w" -o "${INSTALL_DIR}/ytpmd" ./cmd/ytpmd
    CGO_ENABLED=0 go build -ldflags="-s -w" -o "${INSTALL_DIR}/ytpmd-mcp" ./cmd/ytpmd-mcp
    INSTALLED=true
else
    # Download pre-built archive from GitHub release
    TARBALL="ytpMD-${VERSION}-${GOOS}-${GOARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}"
    TMP_DIR="/tmp/ytpmd-install-$$"
    
    echo -e "${CYAN}[*] Downloading prebuilt binary for ${GOOS}/${GOARCH} (v${VERSION})...${NC}"
    mkdir -p "${TMP_DIR}"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${TARBALL}" || true
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${DOWNLOAD_URL}" -O "${TMP_DIR}/${TARBALL}" || true
    fi

    if [ -f "${TMP_DIR}/${TARBALL}" ] && [ -s "${TMP_DIR}/${TARBALL}" ]; then
        tar -xzf "${TMP_DIR}/${TARBALL}" -C "${TMP_DIR}"
        cp -f "${TMP_DIR}/ytpMD-${VERSION}-${GOOS}-${GOARCH}/ytpmd" "${INSTALL_DIR}/ytpmd" 2>/dev/null || cp -f "${TMP_DIR}/ytpMD-${VERSION}-${GOOS}-${GOARCH}/ytpMD" "${INSTALL_DIR}/ytpmd"
        if [ -f "${TMP_DIR}/ytpMD-${VERSION}-${GOOS}-${GOARCH}/ytpmd-mcp" ]; then
            cp -f "${TMP_DIR}/ytpMD-${VERSION}-${GOOS}-${GOARCH}/ytpmd-mcp" "${INSTALL_DIR}/ytpmd-mcp"
        fi
        rm -rf "${TMP_DIR}"
        INSTALLED=true
    elif command -v go >/dev/null 2>&1; then
        echo -e "${YELLOW}[!] Pre-built release archive not reachable directly. Building via 'go install'...${NC}"
        go install github.com/${REPO}/cmd/ytpmd@latest
        if [ -f "${GOPATH:-$HOME/go}/bin/ytpmd" ]; then
            cp -f "${GOPATH:-$HOME/go}/bin/ytpmd" "${INSTALL_DIR}/ytpmd"
            INSTALLED=true
        fi
    fi
fi

if [ "${INSTALLED}" = "true" ]; then
    chmod +x "${INSTALL_DIR}/ytpmd" 2>/dev/null || true
    ln -sf "${INSTALL_DIR}/ytpmd" "${INSTALL_DIR}/ytpMD"
    ln -sf "${INSTALL_DIR}/ytpmd" "${INSTALL_DIR}/ytp24"
    ln -sf "${INSTALL_DIR}/ytpmd" "${INSTALL_DIR}/pdf2md"
    if [ -f "${INSTALL_DIR}/ytpmd-mcp" ]; then
        chmod +x "${INSTALL_DIR}/ytpmd-mcp" 2>/dev/null || true
    fi

    echo ""
    echo -e "${GREEN}[+] Successfully installed ytpMD to ${INSTALL_DIR}/${NC}"
    echo -e "    Binaries & Aliases: ytpmd, ytpMD, ytp24, pdf2md, ytpmd-mcp"
    echo ""

    # PATH Check
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        echo -e "${YELLOW}[!] Note: '${INSTALL_DIR}' is not in your current PATH.${NC}"
        echo -e "    Add it by running:"
        echo -e "    ${GREEN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
        echo -e "    And add that line to your ~/.bashrc or ~/.zshrc."
        echo ""
    fi

    echo -e "${CYAN}Ready! Launch the tool by running:${NC} ${GREEN}ytpmd${NC}"
    echo -e "${CYAN}For MCP IDE integration:${NC} ${GREEN}ytpmd mcp${NC}"
else
    echo -e "${RED}[x] Installation failed. Please check network connection or compile with Go.${NC}"
    exit 1
fi
