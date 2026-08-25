class Ytpmd < Formula
  desc "High-Performance Concurrent PDF to Chapter Markdown Engine & MCP Server"
  homepage "https://github.com/ytp24/ytpMD"
  version "3.2.0"
  license "Apache-2.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/ytp24/ytpMD/releases/download/v3.2.0/ytpMD-3.2.0-darwin-arm64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    else
      url "https://github.com/ytp24/ytpMD/releases/download/v3.2.0/ytpMD-3.2.0-darwin-amd64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/ytp24/ytpMD/releases/download/v3.2.0/ytpMD-3.2.0-linux-arm64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    else
      url "https://github.com/ytp24/ytpMD/releases/download/v3.2.0/ytpMD-3.2.0-linux-amd64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    end
  end

  depends_on "poppler"

  def install
    bin.install "ytpmd"
    bin.install "ytpmd-mcp" if File.exist?("ytpmd-mcp")
    bin.install_symlink "ytpmd" => "ytpMD"
    bin.install_symlink "ytpmd" => "ytp24"
    bin.install_symlink "ytpmd" => "pdf2md"
  end

  test do
    system "#{bin}/ytpmd", "version"
  end
end
