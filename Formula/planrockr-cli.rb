class PlanrockrCli < Formula
  desc "planrockr-cli is a command line interface for the Planrockr API."
  homepage "https://github.com/planrockr/planrockr-cli"
  url "https://github.com/planrockr/planrockr-cli/releases/download/v1.0.2/planrockr-cli-1.0.2-darwin-10.12-amd64.tar.gz"
  sha256 "fa3e0fa77eaa7c00bc1dff016a668d88a00843f9adc39b845e03a76a5d232dfa"

  def install
    bin.install "planrockr-cli"
  end
end
