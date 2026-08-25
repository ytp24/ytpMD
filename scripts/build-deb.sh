#!/usr/bin/env bash
set -e

VERSION="3.1.0"
BUILD_DIR="/tmp/ytp24-deb-build"
PKG_DIR="${BUILD_DIR}/ytp24_${VERSION}_amd64"

echo "==> Building Debian package for ytp24 (v${VERSION})..."

rm -rf "${BUILD_DIR}"
mkdir -p "${PKG_DIR}/DEBIAN"
mkdir -p "${PKG_DIR}/usr/local/bin"
mkdir -p "${PKG_DIR}/usr/share/doc/ytp24"

# Copy debian control
cp packaging/debian/DEBIAN/control "${PKG_DIR}/DEBIAN/"

# Build binary for Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${PKG_DIR}/usr/local/bin/ytp24" ./cmd/ytp24
ln -sf ytp24 "${PKG_DIR}/usr/local/bin/pdf2md"

# Copy documentation
cp README.md LICENSE LEGAL.md "${PKG_DIR}/usr/share/doc/ytp24/"

# Build .deb package
mkdir -p dist
dpkg-deb --build "${PKG_DIR}" "dist/ytp24_${VERSION}_amd64.deb"

echo "==> Successfully created dist/ytp24_${VERSION}_amd64.deb"
