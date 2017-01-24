# Planrockr CLI 

```
planrockr-cli is a command line interface for the Planrockr API.

Usage:
  planrockr-cli [command]

Available Commands:
  auth        auth commands
  import      import a project to Planrockr
  version     show the current version

Flags:
  -h, --help                  help for doctl
  -v, --verbose               verbose output

Use "planrockr-cli [command] --help" for more information about a command.
```

## Installation

### Option 1 - Use a Package Manager (preferred method)

MacOS

You can use Homebrew to install planrockr-cli on Mac OS X by using the command below:

	brew install planrockr-cli

GNU/Linux

Integrations with package managers for GNU/Linux

### Option 2 – Download a Release from GitHub

Visit the Releases page for the planrockr-cli GitHub project, and find the appropriate archive for your operating system and architecture. (For MacOS systems, remember to use the darwin archive.)

MacOS and GNU/Linux

You can download the archive from your browser, or copy its URL and retrieve it to your home directory with wget or curl:	

**MacOS**

	cd ~
	curl -L https://github.com/planrockr/planrockr-cli/releases/download/v1.0.0/planrockr-cli-1.0.0-darwin-10.6-amd64.tar.gz | tar xz

**Gnu/Linux (with wget)**

	cd ~
	wget -qO- https://github.com/planrockr/planrockr-cli/releases/download/v1.0.0/planrockr-cli-1.0.0-linux-amd64.tar.gz  | tar xz

**Gnu/Linux (with curl)**

	cd ~
	curl -L https://github.com/planrockr/planrockr-cli/releases/download/v1.0.0/planrockr-cli-1.0.0-linux-amd64.tar.gz  | tar xz

Move the planrockr-cli binary to somewhere in your path. For example:

	sudo mv ./planrockr-cli /usr/local/bin


### Option 3 – Build From Source

Alternatively, if you have a Go environment configured, you can install the development version of planrockr-cli from the command line like so:

	go get github.com/digitalocean/doctl/cmd/doctl

### Option 4 – Build with Docker

If you have Docker installed, you can build with the Dockerfile a Docker image and run planrockr-cli within a Docker container.

**Build Docker image**

	docker build -t planrockr-cli .

**Usage**

	docker run planrockr-cli <followed by planrockr-cli commands>

## Examples

**Login**

	planrockr-cli login -u email@email.com -p password

The login command will create a file in ~/.planrockr/config.yml with your token. This token will be used in your future commands.

**Import a Jira projet**

	planrockr-cli import -t jira -h http://your_jira_host -u jira_user -p jira_password

A new project will be created on Planrockr and a Webhook will be created on Jira. This Jira user must have permissions to create a new Webhook.

**Import a Gitlab projet**

	planrockr-cli import -t gitlab -h http://your_gitlab_host -u gitlab_user -p gitlab_password

A new project will be created on Planrockr and a Webhook will be created on Gitlab. This Gitlab user must have permissions to create a new Webhook.