COLOR_RESET = \033[0m
GREEN    	= \033[32m

.PHONY: help
help: # Show help for each of the Makefile recipes
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "${GREEN}$$(echo $$l | cut -f 1 -d':')${COLOR_RESET}:$$(echo $$l | cut -f 2- -d'#')\n"; done

build-mac: # Build on MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/go-sort .
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/go-sort .

build-windows: # Build on Windows
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/go-sort.exe .
	GOOS=windows GOARCH=arm64 go build -o bin/windows-arm64/go-sort.exe .
	GOOS=windows GOARCH=386 go build -o bin/windows-386/go-sort.exe .

build-linux: # Build on Linux
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/go-sort .
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/go-sort .
	GOOS=linux GOARCH=386 go build -o bin/linux-386/go-sort .

build: build-mac build-linux build-windows # build application on MacOS, Linux and Windows

generate: # Generate example file
	go run main.go generate -n 20000000 -c

sort: # Sort with example arguments
	go run main.go sort -c -b 8 -t 8

verify: # Verify with default arguments
	go run main.go verify -c

clean-up: # Remove input.bin and output.bin
	rm input.bin output.bin

all: clean-up generate sort verify # Test the entire application