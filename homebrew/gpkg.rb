# Homebrew Formula for gpkg
# This formula is automatically updated by the release workflow
class Gpkg < Formula
  desc "Simple, user-focused package manager for GitHub releases and source builds"
  homepage "https://github.com/grave0x/gpkg"
  version "{{ VERSION }}"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/grave0x/gpkg/releases/download/v{{ VERSION }}/gpkg-v{{ VERSION }}-darwin-arm64.tar.gz"
      sha256 "{{ SHA256_DARWIN_ARM64 }}"
    else
      url "https://github.com/grave0x/gpkg/releases/download/v{{ VERSION }}/gpkg-v{{ VERSION }}-darwin-amd64.tar.gz"
      sha256 "{{ SHA256_DARWIN_AMD64 }}"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/grave0x/gpkg/releases/download/v{{ VERSION }}/gpkg-v{{ VERSION }}-linux-arm64.tar.gz"
      sha256 "{{ SHA256_LINUX_ARM64 }}"
    elsif Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/grave0x/gpkg/releases/download/v{{ VERSION }}/gpkg-v{{ VERSION }}-linux-amd64.tar.gz"
        sha256 "{{ SHA256_LINUX_AMD64 }}"
      else
        url "https://github.com/grave0x/gpkg/releases/download/v{{ VERSION }}/gpkg-v{{ VERSION }}-linux-386.tar.gz"
        sha256 "{{ SHA256_LINUX_386 }}"
      end
    end
  end

  depends_on "sqlite"

  def install
    bin.install "gpkg-v#{version}-#{OS.kernel_name.downcase}-#{Hardware::CPU.arch}" => "gpkg"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/gpkg --version")
  end
end
