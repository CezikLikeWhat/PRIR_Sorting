
input:
	go run main.go input --number 10 --time

sort:
	go run main.go sort -t

verify:
	go run main.go verify -t


all: input sort verify