
input:
	go run main.go input --number 100000000 --time

sort:
	go run main.go sort -t

verify:
	go run main.go verify -t


all: input sort verify