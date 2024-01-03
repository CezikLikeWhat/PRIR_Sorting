
generate:
	go run main.go generate -n 1000000 -c

sort:
	go run main.go sort -c

verify:
	go run main.go verify -c

clean-up:
	rm input.bin output.bin

all: clean-up generate sort verify