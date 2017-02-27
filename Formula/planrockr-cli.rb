class PlanrockrCli < Formula
  desc "planrockr-cli is a command line interface for the Planrockr API."
  homepage "https://github.com/planrockr/planrockr-cli"
  url "https://github.com/planrockr/planrockr-cli/releases/download/v1.0/planrockr-cli-1.0.0-darwin-10.12-amd64.tar.gz"
  sha256 "eb1b6aebb17f8804b9045cacbd1fddfb25b3b2ca96fb067c8c5f6d573b7d423d"


  def install
    bin.install "planrockr-cli"
  end
end
