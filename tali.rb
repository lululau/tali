class Tali < Formula
  desc "Terminal User Interface (TUI) application for managing Alibaba Cloud resources"
  homepage "https://github.com/lululau/tali"
  url "https://github.com/lululau/tali/archive/v1.0.0.tar.gz"
  sha256 "526f073af51d91ce86c5d129aa1aec7ec6e42d199745ea9966779d8642df950b"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end

  test do
    assert_match "tali", shell_output("#{bin}/tali --help 2>&1", 1)
  end
end
