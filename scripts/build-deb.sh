#!/usr/bin/env bash
set -e

VERSION="3.2.0"
BUILD_DIR="/tmp/ytpmd-deb-build"
PKG_DIR="${BUILD_DIR}/ytpmd_${VERSION}_amd64"

echo "==> Building Debian package for ytpMD (v${VERSION})..."

rm -rf "${BUILD_DIR}"
mkdir -p "${PKG_DIR}/DEBIAN"
mkdir -p "${PKG_DIR}/usr/local/bin"
mkdir -p "${PKG_DIR}/usr/share/doc/ytpmd"

cp packaging/debian/DEBIAN/control "${PKG_DIR}/DEBIAN/"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${PKG_DIR}/usr/local/bin/ytpmd" ./cmd/ytpmd
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${PKG_DIR}/usr/local/bin/ytpmd-mcp" ./cmd/ytpmd-mcp
ln -sf ytpmd "${PKG_DIR}/usr/local/bin/ytpMD"
ln -sf ytpmd "${PKG_DIR}/usr/local/bin/ytp24"
ln -sf ytpmd "${PKG_DIR}/usr/local/bin/pdf2md"

cp README.md LICENSE LEGAL.md "${PKG_DIR}/usr/share/doc/ytpmd/"

mkdir -p dist
dpkg-deb --build "${PKG_DIR}" "dist/ytpmd_${VERSION}_amd64.deb"

echo "==> Successfully created dist/ytpmd_${VERSION}_amd64.deb"
