#!/bin/bash

sizes=(1000 10000 50000 100000 250000 500000 750000 1000000 2500000 5000000 10000000 25000000 50000000)
buckets=16
cpu=8

for size in "${sizes[@]}"; do
    echo "Size: $size"
    for i in {1..5}; do
        rm -f input.bin output.bin 2> /dev/null
        ./bin/linux-amd64/go-sort generate -n 100000000
        echo "Try $i: $(./bin/linux-amd64/go-sort sort -c -t $cpu -b $buckets -s $size | cut -d" " -f3)"
        echo "$(./bin/linux-amd64/go-sort verify)"
    done
    echo "======================================================="
done

rm -f input.bin output.bin 2> /dev/null