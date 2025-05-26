#!/bin/bash

# Script to generate Homebrew formula for tali
# Usage: ./scripts/generate-formula.sh <version> <github-username>

set -e

VERSION=${1:-"1.0.0"}
GITHUB_USER=${2:-"lululau"}

if [ "$GITHUB_USER" = "lululau" ]; then
    echo "Usage: $0 <version> <github-username>"
    echo "Example: $0 1.0.0 lululau"
    exit 1
fi

TARBALL_URL="https://github.com/${GITHUB_USER}/tali/archive/v${VERSION}.tar.gz"
TEMP_FILE="/tmp/tali-${VERSION}.tar.gz"

echo "Downloading tarball to calculate SHA256..."
curl -L "$TARBALL_URL" -o "$TEMP_FILE"

if [ ! -f "$TEMP_FILE" ]; then
    echo "Error: Failed to download tarball from $TARBALL_URL"
    echo "Make sure the release v${VERSION} exists on GitHub"
    exit 1
fi

SHA256=$(shasum -a 256 "$TEMP_FILE" | cut -d' ' -f1)
echo "SHA256: $SHA256"

# Generate the formula
cat > tali.rb << EOF
class Tali < Formula
  desc "Terminal User Interface (TUI) application for managing Alibaba Cloud resources"
  homepage "https://github.com/${GITHUB_USER}/tali"
  url "https://github.com/${GITHUB_USER}/tali/archive/v${VERSION}.tar.gz"
  sha256 "${SHA256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end

  test do
    assert_match "tali", shell_output("#{bin}/tali --help 2>&1", 1)
  end
end
EOF

echo "Generated tali.rb formula for version ${VERSION}"
echo "GitHub user: ${GITHUB_USER}"
echo "SHA256: ${SHA256}"

# Clean up
rm -f "$TEMP_FILE"

echo ""
echo "Next steps:"
echo "1. Review the generated tali.rb file"
echo "2. Test the formula: brew install --build-from-source ./tali.rb"
echo "3. Create a homebrew-tali repository on GitHub"
echo "4. Copy tali.rb to the root of that repository"
echo "5. Users can then install with: brew tap ${GITHUB_USER}/tali && brew install tali" 