class PlanrockrCli < Formula
  desc "planrockr-cli is a command line interface for the Planrockr API."
  homepage "https://github.com/planrockr/planrockr-cli"
  url "https://github.com/planrockr/planrockr-cli/releases/download/v1.0.1/planrockr-cli-1.0.1-darwin-10.12-amd64.tar.gz"
  sha256 "90d9ebf22ea73ceab5a5cb9b82333f42e0fda813168b8a912311c5d10c926516"

  def install
    bin.install "planrockr-cli"
  end
end
