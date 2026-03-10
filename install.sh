#!/bin/sh
set -e

# PRM installer script
# Usage: curl -sSL https://raw.githubusercontent.com/cstsortan/prm/main/install.sh | sh

REPO="cstsortan/prm"
INSTALL_DIR="/usr/local/bin"

get_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    *)       echo "unsupported" ;;
  esac
}

get_arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "amd64" ;;
    arm64|aarch64)  echo "arm64" ;;
    *)              echo "unsupported" ;;
  esac
}

OS=$(get_os)
ARCH=$(get_arch)

if [ "$OS" = "unsupported" ] || [ "$ARCH" = "unsupported" ]; then
  echo "Error: unsupported platform $(uname -s)/$(uname -m)"
  exit 1
fi

# Get latest release tag
VERSION=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest version"
  exit 1
fi

FILENAME="prm_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

echo "Downloading prm ${VERSION} for ${OS}/${ARCH}..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -sSL "$URL" -o "${TMPDIR}/${FILENAME}"
tar -xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"

if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/prm" "${INSTALL_DIR}/prm"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMPDIR}/prm" "${INSTALL_DIR}/prm"
fi

echo "prm ${VERSION} installed to ${INSTALL_DIR}/prm"
