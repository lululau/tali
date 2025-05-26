class Tali < Formula
  desc "Terminal User Interface (TUI) application for managing Alibaba Cloud resources"
  homepage "https://github.com/lululau/tali"
  url "https://github.com/lululau/tali/archive/v1.0.0.tar.gz"
  sha256 "your-sha256-hash-here"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end

  test do
    # Test that the binary was installed correctly
    assert_match "tali", shell_output("#{bin}/tali --help 2>&1", 1)
  end
end