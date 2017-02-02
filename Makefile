all: dependencies build

clean:
	rm -rf bin

dependencies:
	go install -v

build:
	go build -o bin/planrockr-cli cmd/main.go
