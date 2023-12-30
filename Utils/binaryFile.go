package Utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/logger"
	"math/rand"
	"os"
	"time"
)

// createOrOpenFile - Function responsible for opening or creating a file when one does not exist
func createOrOpenFile(path string) *os.File {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			logger.Fatal("Cannot create input fileError: %s", err)
		}

		return file
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		if err != nil {
			logger.Fatal("Cannot open input file | Error: %s", err)
		}
	}

	return file
}

// WriteToFile - Function responsible for generating the input file for sorting
func WriteToFile(filename string, numberOfElements int64) {
	var i int64

	file := createOrOpenFile(filename)
	defer file.Close()

	err := binary.Write(file, binary.NativeEndian, numberOfElements)
	if err != nil {
		logger.Fatal("Error writing NumberOfElements | Error: %s", err)
	}
	logger.Debug("Write function | Written number of elements: %d", numberOfElements)

	source := rand.NewSource(time.Now().UnixNano())
	randSource := rand.New(source)
	for i = 0; i < numberOfElements; i++ {
		randNum := randSource.Int31()
		logger.Debug("Write function | Written number: %d", randNum)

		err := binary.Write(file, binary.NativeEndian, randNum)
		if err != nil {
			logger.Fatal("Error writing Element | Error: %s", err)
		}
	}
}

// ReadFromFile - For debugging purposes
func ReadFromFile(filename string, numberOfElements int64) {
	var (
		i                 int64
		numberOfItemsRead int64
	)

	file := createOrOpenFile(filename)
	defer file.Close()

	_, err := file.Seek(0, 0)
	if err != nil {
		logger.Fatal("Cannot seek file | Error: %s", err)
	}

	err = binary.Read(file, binary.NativeEndian, &numberOfItemsRead)
	if err != nil {
		logger.Fatal("Cannot read number of generated numbers from file | Error: %s", err)
	}
	logger.Debug("Read function | Read number of elements: %d", numberOfItemsRead)

	for i = 0; i < numberOfElements; i++ {
		var item int32
		err = binary.Read(file, binary.NativeEndian, &item)
		if err != nil {
			logger.Fatal("Cannot read generated numbers from file | Error: %s", err)
		}
		logger.Debug("Read function | Read number: %d", item)
	}
}

func VerifySortedFile(filename string) bool {
	var (
		i                 int64
		numberOfItemsRead int64
		previous          int32
		current           int32
	)

	file := createOrOpenFile(filename)
	defer file.Close()

	_, err := file.Seek(0, 0)
	if err != nil {
		logger.Fatal("Cannot seek file | Error: %s", err)
	}

	err = binary.Read(file, binary.NativeEndian, &numberOfItemsRead)
	if err != nil {
		logger.Fatal("Cannot read number of generated numbers from file | Error: %s", err)
	}
	logger.Debug(fmt.Sprintf("Verify function | Read number of elements: %d", numberOfItemsRead))

	err = binary.Read(file, binary.NativeEndian, &previous)
	if err != nil {
		logger.Fatal("Cannot read first generated number from file | Error: %s", err)
	}
	logger.Debug(fmt.Sprintf("Verify function | Read first generated number: %d", previous))

	for i = 0; i < numberOfItemsRead-1; i++ {
		err = binary.Read(file, binary.NativeEndian, &current)
		if err != nil {
			logger.Fatal("Cannot read generated numbers from file | Error: %s", err)
		}
		logger.Debug(fmt.Sprintf("Verify function | Read number: %d", current))

		if previous > current {
			return false
		}
	}

	return true
}
