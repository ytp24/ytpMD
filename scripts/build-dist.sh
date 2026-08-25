#!/usr/bin/env bash
set -e

VERSION="3.2.0"
DIST_DIR="dist"

echo "==> Building cross-platform distribution archives (v${VERSION})..."
mkdir -p "${DIST_DIR}"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUTPUT_NAME="ytpMD-${VERSION}-${GOOS}-${GOARCH}"
    TARGET_DIR="/tmp/${OUTPUT_NAME}"
    
    echo "  --> Compiling for ${GOOS}/${GOARCH}..."
    rm -rf "${TARGET_DIR}"
    mkdir -p "${TARGET_DIR}"

    BINARY="ytpMD"
    if [ "${GOOS}" = "windows" ]; then
        BINARY="ytpMD.exe"
    fi

    CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags="-s -w" -o "${TARGET_DIR}/${BINARY}" ./cmd/ytpMD
    cp README.md LICENSE LEGAL.md "${TARGET_DIR}/"

    if [ "${GOOS}" = "windows" ]; then
        (cd /tmp && zip -q -r "${OUTPUT_NAME}.zip" "${OUTPUT_NAME}")
        mv "/tmp/${OUTPUT_NAME}.zip" "${DIST_DIR}/"
    else
        tar -czf "${DIST_DIR}/${OUTPUT_NAME}.tar.gz" -C /tmp "${OUTPUT_NAME}"
    fi
    rm -rf "${TARGET_DIR}"
done

echo "==> All distribution binaries built in ./${DIST_DIR}/"
ls -lh "${DIST_DIR}/"
