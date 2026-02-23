#!/usr/bin/env bash
set -euo pipefail

REPO="cgons/hdinfo"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="hdinfo"

# Ensure running as root
if [ "$(id -u)" -ne 0 ]; then
    echo "Error: This install script must be run as root (e.g. curl ... | sudo bash)"
    exit 1
fi

echo "----------------------------"
echo "AUTO-INSTALL SCRIPT - hdinfo"
echo "----------------------------"

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH_SUFFIX="amd64" ;;
    aarch64) ARCH_SUFFIX="arm64" ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        echo "hdinfo only supports amd64 (x86_64) and arm64 (aarch64)."
        exit 1
        ;;
esac

# Detect package manager
if command -v apt-get &>/dev/null; then
    PKG_MANAGER="apt"
    INSTALL_CMD="apt-get install -y util-linux hdparm smartmontools"
elif command -v dnf &>/dev/null; then
    PKG_MANAGER="dnf"
    INSTALL_CMD="dnf install -y util-linux hdparm smartmontools"
elif command -v pacman &>/dev/null; then
    PKG_MANAGER="pacman"
    INSTALL_CMD="pacman -S --noconfirm util-linux hdparm smartmontools"
else
    echo "Error: Could not detect a supported package manager (apt, dnf, pacman)."
    echo "Please install the following packages manually: util-linux hdparm smartmontools"
    exit 1
fi

echo "Detected architecture: ${ARCH} (${ARCH_SUFFIX})"
echo "Detected package manager: ${PKG_MANAGER}"
echo ""

echo "Install Dependencies"
echo "--------------------"

# Prompt user to install dependencies
echo "The following packages are required: util-linux, hdparm, smartmontools"
echo ""
read -rp "Install dependencies using ${PKG_MANAGER}? [y/n] " CONFIRM < /dev/tty
if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "Installing dependencies..."
    $INSTALL_CMD
else
    echo "Skipping dependency installation."
    echo "Make sure dependencies are installed before using hdinfo."
fi

echo ""

# Fetch latest release download URL
ASSET_NAME="hdinfo-linux-${ARCH_SUFFIX}"
DOWNLOAD_URL="$(
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -o "\"browser_download_url\": *\"[^\"]*${ASSET_NAME}\"" \
    | cut -d'"' -f4
)"

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find a release asset matching '${ASSET_NAME}'."
    echo "Check https://github.com/${REPO}/releases for available downloads."
    exit 1
fi

echo -ne "Downloading ${ASSET_NAME} @ latest..."
curl -fsSL -o "${INSTALL_DIR}/${BINARY_NAME}" "$DOWNLOAD_URL"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
echo -ne "\rDownloading ${ASSET_NAME} @ latest... Done!"
echo ""

echo ""
cat << 'EOF'
  _         _ _        __        _____           _        _ _          _   _ 
 | |       | (_)      / _|      |_   _|         | |      | | |        | | | |
 | |__   __| |_ _ __ | |_ ___     | |  _ __  ___| |_ __ _| | | ___  __| | | |
 | '_ \ / _` | | '_ \|  _/ _ \    | | | '_ \/ __| __/ _` | | |/ _ \/ _` | | |
 | | | | (_| | | | | | || (_) |  _| |_| | | \__ \ || (_| | | |  __/ (_| | |_|
 |_| |_|\__,_|_|_| |_|_| \___/  |_____|_| |_|___/\__\__,_|_|_|\___|\__,_| (_)
                                                                             
EOF

echo "Installed successfully to: ${INSTALL_DIR}/${BINARY_NAME}"
echo "---"
echo "Run 'sudo hdinfo' to get started."
echo ""
