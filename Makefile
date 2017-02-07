all: build

clean:
	rm -rf bin

dependencies:
	go install -v

build:
	go build -o bin/planrockr-cli src/github.com/planrockr/planrockr-cli/main.go
