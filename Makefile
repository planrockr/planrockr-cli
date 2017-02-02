all: dependencies build

clean:
	rm -rf bin

dependencies:
	if [ ! -d "bin" ]; then mkdir bin; fi
	curl https://glide.sh/get | sh
	cd ./src/planrockr && glide install
	cd ./src/planrockr/cmd && go install -v

build:
	go build -o bin/planrockr-cli cmd/main.go
