package Utils

import (
	"encoding/binary"
	"errors"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func CreateFile(path string) *os.File {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			log.Panicf("Cannot create input file | Error: %s", err)
		}

		return file
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		if err != nil {
			log.Panicf("Cannot open input file | Error: %s", err)
		}
	}

	return file
}

func WriteToFile(file *os.File, numberOfElements int64) {
	var i int64

	err := binary.Write(file, binary.NativeEndian, numberOfElements)
	if err != nil {
		log.Fatalf("Error writing NumberOfElements | Error: %s", err)
	}
	if IsDebugMode {
		log.Printf("| Write function | Written number of elements: %d", numberOfElements)
	}

	source := rand.NewSource(time.Now().UnixNano())
	randSource := rand.New(source)
	for i = 0; i < numberOfElements; i++ {
		randNum := randSource.Int31()
		if IsDebugMode {
			log.Printf("| Write function | Written number: %d", randNum)
		}
		err := binary.Write(file, binary.NativeEndian, randNum)
		if err != nil {
			log.Fatalf("Error writing Element | Error: %s", err)
		}
	}
}

func WriteToFileAsync(file *os.File, numberOfElements int64) {
	var wg sync.WaitGroup

	err := binary.Write(file, binary.NativeEndian, numberOfElements)
	if err != nil {
		log.Fatalf("Error writing NumberOfElements | Error: %s", err)
	}
	if IsDebugMode {
		log.Printf("| Write function | Written number of elements: %d", numberOfElements)
	}

	wg.Add(int(numberOfElements))
	for i := int64(0); i < numberOfElements; i++ {
		go func(index int64) {
			defer wg.Done()
			source := rand.NewSource(time.Now().UnixNano())
			randSource := rand.New(source)
			randNum := randSource.Int31()
			if IsDebugMode {
				log.Printf("| Write function | Written number: %d", randNum)
			}
			err := binary.Write(file, binary.NativeEndian, randNum)
			if err != nil {
				log.Fatalf("Error writing Element | Error: %s", err)
			}
		}(i)
	}

	wg.Wait()
}

func ReadFromFile(file *os.File, numberOfElements int64) {
	var (
		i                 int64
		numberOfItemsRead int64
	)
	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Cannot seek file | Error: %s", err)
	}

	err = binary.Read(file, binary.NativeEndian, &numberOfItemsRead)
	if IsDebugMode {
		log.Printf("Read function | Read number of elements: %d", numberOfItemsRead)
	}
	if err != nil {
		log.Fatalf("Cannot read number of generated numbers from file | Error: %s", err)
	}

	for i = 0; i < numberOfElements; i++ {
		var item int32
		err = binary.Read(file, binary.NativeEndian, &item)
		if IsDebugMode {
			log.Printf("Read function | Read number: %d", item)
		}
		if err != nil {
			log.Fatalf("Cannot read generated numbers from file | Error: %s", err)
		}
	}
}

func ReadFromFileAsync(file *os.File, numberOfElements int64) {
	var (
		wg                sync.WaitGroup
		numberOfItemsRead int64
	)

	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Cannot seek file | Error: %s", err)
	}

	err = binary.Read(file, binary.NativeEndian, &numberOfItemsRead)
	if IsDebugMode {
		log.Printf("Read function | Read number of elements: %d", numberOfItemsRead)
	}
	if err != nil {
		log.Fatalf("Cannot read number of generated numbers from file | Error: %s", err)
	}

	wg.Add(int(numberOfElements))
	for i := int64(0); i < numberOfElements; i++ {
		go func() {
			defer wg.Done()
			var item int32
			err = binary.Read(file, binary.NativeEndian, &item)
			if IsDebugMode {
				log.Printf("Read function | Read number: %d", item)
			}
			if err != nil {
				log.Fatalf("Cannot read generated numbers from file | Error: %s", err)
			}
		}()
	}

	wg.Wait()
}
