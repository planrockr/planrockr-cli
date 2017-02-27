# Planrockr CLI 

[![Build Status](https://travis-ci.org/planrockr/planrockr-cli.svg?branch=master)](https://travis-ci.org/planrockr/planrockr-cli)

[![codecov](https://codecov.io/gh/planrockr/planrockr-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/planrockr/planrockr-cli)


```
planrockr-cli is a command line interface for the Planrockr API.

Usage:
  planrockr-cli [command]

Available Commands:
  auth        auth commands
  import      import a project to Planrockr
  version     show the current version

Flags:
  -h, --help                  help for planrockr-cli
  -v, --verbose               verbose output

Use "planrockr-cli [command] --help" for more information about a command.
```

## Installation

### Option 1 - Use a Package Manager (preferred method)

*MacOS*

You can use Homebrew to install planrockr-cli on Mac OS X by using the command below:

	brew tap planrockr/planrockr-cli
	brew install planrockr-cli

*GNU/Linux*

Integrations with package managers for GNU/Linux are to come.

### Option 2 – Download a Release from GitHub

Visit the Releases page for the planrockr-cli GitHub project, and find the appropriate archive for your operating system and architecture. (For MacOS systems, remember to use the darwin archive.)

*MacOS* and *GNU/Linux*

You can download the archive from your browser, or copy its URL and retrieve it to your home directory with wget or curl:	

**MacOS**

	cd ~
	curl -L https://github.com/planrockr/planrockr-cli/releases/download/v1.0.2/planrockr-cli-1.0.2-darwin-10.12-amd64.tar.gz | tar xz

**Gnu/Linux (with wget)**

	cd ~
	wget -qO- https://github.com/planrockr/planrockr-cli/releases/download/v1.0.2/planrockr-cli-1.0.2-linux-amd64.tar.gz  | tar xz

**Gnu/Linux (with curl)**

	cd ~
	curl -L https://github.com/planrockr/planrockr-cli/releases/download/v1.0.2/planrockr-cli-1.0.2-linux-amd64.tar.gz  | tar xz

Move the planrockr-cli binary to somewhere in your path. For example:

	sudo mv ./planrockr-cli /usr/local/bin


### Option 3 – Build From Source

Alternatively, if you have a Go environment configured, you can install the development version of planrockr-cli from the command line like so:

	go get github.com/planrockr/planrockr-cli/cmd/planrockr-cli

### Option 4 – Build with Docker

If you have Docker installed, you can build with the Dockerfile a Docker image and run planrockr-cli within a Docker container.

**Build Docker image**

	docker build -t planrockr-cli .

**Usage**

	docker run planrockr-cli <followed by planrockr-cli commands>

## Examples

**Auth**

	planrockr-cli auth -u email@email.com -p password

The auth command will create a file in ~/.planrockr/config.yml with your token. This token will be used in your future commands. If you omit the -u and -p parameters you will be asked to provide the username and password right after the command is executed.

**Import a Jira project**

	planrockr-cli import -t jira -s http://your_jira_host -u jira_user -p jira_password

A new project will be created on Planrockr and a Webhook will be created on Jira. This Jira user must have permissions to create a new Webhook. If you omit the -u and -p parameters you will be asked to provide the username and password right after the command is executed.

**Import a Gitlab project**

	planrockr-cli import -t gitlab -s http://your_gitlab_host -u gitlab_user -p gitlab_password

A new project will be created on Planrockr and a *Webhook* will be created on Gitlab. This Gitlab user must have permissions to create a new *Webhook*.
