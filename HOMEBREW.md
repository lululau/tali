# Homebrew Installation Guide

This document provides detailed instructions for installing tali using Homebrew on macOS.

## Quick Installation

```bash
# Install directly from the formula URL
brew install https://raw.githubusercontent.com/lululau/tali/main/tali.rb
```

## Setting up a Homebrew Tap

For easier installation and updates, you can create a Homebrew tap:

### 1. Create a Homebrew Tap Repository

Create a new repository named `homebrew-tali` on GitHub (the `homebrew-` prefix is required).

### 2. Add the Formula

Copy the `tali.rb` file to the root of your `homebrew-tali` repository.

### 3. Update the Formula

Before publishing, update the following fields in `tali.rb`:

- **homepage**: Update to your actual repository URL
- **url**: Update to point to your release tarball
- **sha256**: Calculate the SHA256 hash of your release tarball
- **license**: Update if using a different license

### 4. Calculate SHA256 Hash

To get the SHA256 hash for your release:

```bash
# Download your release tarball
curl -L https://github.com/lululau/tali/archive/v1.0.0.tar.gz -o tali-1.0.0.tar.gz

# Calculate SHA256
shasum -a 256 tali-1.0.0.tar.gz
```

### 5. Install from Your Tap

Once your tap is set up:

```bash
# Add your tap
brew tap lululau/tali

# Install tali
brew install tali
```

## Updating the Formula

When you release a new version:

1. Update the `url` field to point to the new release
2. Update the `sha256` hash
3. Commit and push the changes to your tap repository
4. Users can update with: `brew upgrade tali`

## Testing the Formula

Test your formula locally:

```bash
# Install from local formula
brew install --build-from-source ./tali.rb

# Test the installation
tali --version

# Uninstall for testing
brew uninstall tali
```

## Formula Template

Here's the complete formula template with placeholders:

```ruby
class Tali < Formula
  desc "Terminal User Interface (TUI) application for managing Alibaba Cloud resources"
  homepage "https://github.com/lululau/tali"
  url "https://github.com/lululau/tali/archive/vVERSION.tar.gz"
  sha256 "YOUR_SHA256_HASH"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end

  test do
    assert_match "tali", shell_output("#{bin}/tali --help 2>&1", 1)
  end
end
```

## Troubleshooting

### Common Issues

1. **Formula not found**: Ensure the tap is added correctly and the repository is public
2. **Build failures**: Check that all dependencies are correctly specified
3. **SHA256 mismatch**: Recalculate the hash after any changes to the source code

### Getting Help

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Creating Homebrew Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap) 