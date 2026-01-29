#!/bin/bash
#
# mmdc installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/ageha734/mmdc/master/install.sh | bash
#

set -euo pipefail

REPO="ageha734/mmdc"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
TMPDIR="${TMPDIR:-/tmp}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

detect_platform() {
    local os arch

    case "$(uname -s)" in
    Linux*) os="linux" ;;
    Darwin*) os="darwin" ;;
    MINGW* | MSYS* | CYGWIN*) os="windows" ;;
    *) error "Unsupported operating system: $(uname -s)" ;;
    esac

    case "$(uname -m)" in
    x86_64 | amd64) arch="amd64" ;;
    arm64 | aarch64) arch="arm64" ;;
    armv7*) arch="arm" ;;
    i386 | i686) arch="386" ;;
    *) error "Unsupported architecture: $(uname -m)" ;;
    esac

    echo "${os}_${arch}"
}

get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" |
        grep '"tag_name":' |
        sed -E 's/.*"([^"]+)".*/\1/'
}

install() {
    local platform version download_url filename install_path

    platform=$(detect_platform)
    info "Detected platform: ${platform}"

    version=$(get_latest_version)
    if [ -z "$version" ]; then
        error "Failed to get latest version"
    fi
    info "Latest version: ${version}"

    filename="mmdc_${platform}.tar.gz"
    if [[ "$platform" == windows_* ]]; then
        filename="mmdc_${platform}.zip"
    fi

    download_url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    info "Downloading from: ${download_url}"

    cd "$TMPDIR"

    if ! curl -fsSLO "$download_url"; then
        error "Failed to download ${filename}"
    fi

    checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
    if curl -fsSLO "$checksums_url" 2>/dev/null; then
        info "Verifying checksum..."
        if command -v sha256sum &>/dev/null; then
            if ! grep "$filename" checksums.txt | sha256sum -c - &>/dev/null; then
                error "Checksum verification failed"
            fi
            info "Checksum verified"
        else
            warn "sha256sum not found, skipping checksum verification"
        fi
        rm -f checksums.txt
    fi

    info "Extracting..."
    if [[ "$filename" == *.zip ]]; then
        unzip -q "$filename"
    else
        tar -xzf "$filename"
    fi

    install_path="${INSTALL_DIR}/mmdc"

    if [ -w "$INSTALL_DIR" ]; then
        mv mmdc "$install_path"
        chmod +x "$install_path"
    else
        info "Requesting sudo for installation to ${INSTALL_DIR}"
        sudo mv mmdc "$install_path"
        sudo chmod +x "$install_path"
    fi

    rm -f "$filename"

    info "Successfully installed mmdc to ${install_path}"
    echo ""
    info "Run 'mmdc --version' to verify the installation"
}

check_dependencies() {
    if ! command -v curl &>/dev/null; then
        error "curl is required but not installed"
    fi

    if ! command -v tar &>/dev/null; then
        error "tar is required but not installed"
    fi
}

main() {
    echo ""
    echo "mmdc installer"
    echo "=============="
    echo ""

    check_dependencies
    install

    echo ""
    echo "Installation complete!"
    echo ""
}

main "$@"
