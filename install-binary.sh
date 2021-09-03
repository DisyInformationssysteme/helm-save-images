#!/usr/bin/env sh

# Shamelessly copied from https://github.com/technosophos/helm-template

PROJECT_NAME="helm-save-images"
PROJECT_GH="DisyInformationssysteme/$PROJECT_NAME"
export GREP_COLOR="never"

HELM_MAJOR_VERSION=$("${HELM_BIN}" version --client --short | awk -F '.' '{print $1}')

: ${HELM_PLUGIN_DIR:="$("${HELM_BIN}" home --debug=false)/plugins/helm-save-images"}

# Convert the HELM_PLUGIN_DIR to unix if cygpath is
# available. This is the case when using MSYS2 or Cygwin
# on Windows where helm returns a Windows path but we
# need a Unix path

if type cygpath >/dev/null 2>&1; then
  HELM_PLUGIN_DIR=$(cygpath -u $HELM_PLUGIN_DIR)
fi

if [ "$SKIP_BIN_INSTALL" = "1" ]; then
  echo "Skipping binary install"
  exit
fi

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
  armv5*) ARCH="armv5" ;;
  armv6*) ARCH="armv6" ;;
  armv7*) ARCH="armv7" ;;
  aarch64) ARCH="arm64" ;;
  x86) ARCH="x86" ;;
  x86_64) ARCH="x86_64" ;;
  i686) ARCH="i386" ;;
  i386) ARCH="i386" ;;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(uname | tr '[:upper:]' '[:lower:]')

  case "$OS" in
  # Msys support
  msys*) OS='Windows' ;;
  # Minimalist GNU for Windows
  mingw*) OS='Windows' ;;
  darwin) OS='Darwin' ;;
  esac
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  supported="Darwin_arm64\nDarwin_x86_64\nLinux_arm64\nLinux_i386\nLinux_x86_64\nWindows_arm64\nWindows_i386\nWindows_x86_64\n"
  if ! echo "${supported}" | grep -q "${OS}_${ARCH}"; then
    echo "No prebuild binary for ${OS}_${ARCH}."
    exit 1
  fi

  if ! type "curl" >/dev/null && ! type "wget" >/dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# getDownloadURL checks the latest available version.
getDownloadURL() {
  version=$(git -C "$HELM_PLUGIN_DIR" describe --tags --exact-match 2>/dev/null || :)
  if [ -n "$version" ]; then
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$version/helm-save-images_$version_$OS_$ARCH.tar.gz"
  else
    # Use the GitHub API to find the download url for this project.
    url="https://api.github.com/repos/$PROJECT_GH/releases/latest"
    if type "curl" >/dev/null; then
      DOWNLOAD_URL=$(curl -s $url | grep $OS | grep $ARCH | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
    elif type "wget" >/dev/null; then
      DOWNLOAD_URL=$(wget -q -O - $url | grep $OS | grep $ARCH | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
    fi
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  PLUGIN_TMP_FILE="/tmp/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" >/dev/null; then
    curl -L "$DOWNLOAD_URL" -o "$PLUGIN_TMP_FILE"
  elif type "wget" >/dev/null; then
    wget -q -O "$PLUGIN_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  HELM_TMP="/tmp/$PROJECT_NAME"
  mkdir -p "$HELM_TMP"
  tar xf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  HELM_TMP_BIN="$HELM_TMP/helm-save-images"
  echo "Preparing to install into ${HELM_PLUGIN_DIR}"
  mkdir -p "$HELM_PLUGIN_DIR/bin"
  cp "$HELM_TMP_BIN" "$HELM_PLUGIN_DIR/bin"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    printf '\tFor support, go to https://github.com/DisyInformationssysteme/helm-save-images.\n'
  fi
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGIN_DIR/$PROJECT_NAME"
  "${HELM_PLUGIN_DIR}/bin/save-images" -h
  set -e
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e
initArch
initOS
verifySupported
getDownloadURL
downloadFile
installFile
testVersion
