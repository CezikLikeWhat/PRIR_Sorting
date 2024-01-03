package Utils

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// createOrOpenFile - Function responsible for opening or creating a file when one does not exist
func createOrOpenFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot create file %s (Error: %s)", path, err))
		}

		return file, nil
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot create file %s (Error: %s)", path, err))
		}
	}

	return file, nil
}

// GenerateInputBinaryFile - Function responsible for generating the input file for sorting
func GenerateInputBinaryFile(filename string, numberOfElements int64) error {
	file, err := createOrOpenFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = binary.Write(writer, binary.NativeEndian, numberOfElements)
	if err != nil {
		return errors.New(fmt.Sprintf("Error during writing number of elements (Error: %s)", err))
	}

	source := rand.NewSource(time.Now().UnixNano())
	randSource := rand.New(source)
	for i := int64(0); i < numberOfElements; i++ {
		randNum := randSource.Int31()

		err := binary.Write(writer, binary.NativeEndian, randNum)
		if err != nil {
			return errors.New(fmt.Sprintf("Error during writing element (Error: %s)", err))
		}
	}
	return nil
}

// VerifySortedFile - Function responsible for validating the sorted data in the file
func VerifySortedFile(filename string) (bool, error) {
	var (
		numberOfItemsRead int64
		previous          int32
		current           int32
	)

	file, err := createOrOpenFile(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	err = binary.Read(reader, binary.NativeEndian, &numberOfItemsRead)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Cannot read number of generated numbers from file %s (Error: %s)", filename, err))
	}

	err = binary.Read(reader, binary.NativeEndian, &previous)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Cannot read first generated number from file %s (Error: %s)", filename, err))
	}

	for i := int64(0); i < numberOfItemsRead-1; i++ {
		if err = binary.Read(reader, binary.NativeEndian, &current); err != nil {
			return false, errors.New(fmt.Sprintf("Cannot read generated numbers from file %s (Error: %s)", filename, err))
		}

		if previous > current {
			return false, nil
		}

		previous = current
	}

	return true, nil
}

// WriteDataToOutputFile - Function responsible for writing a small amount of data (see threshold in Utils/sort.go) to a file
func WriteDataToOutputFile(filename string, numberOfElements int64, data []int32) error {
	file, err := createOrOpenFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = binary.Write(writer, binary.NativeEndian, numberOfElements)
	if err != nil {
		return errors.New(fmt.Sprintf("Error during writing number of elements (Error: %s)", err))
	}

	for _, element := range data {
		err := binary.Write(writer, binary.NativeEndian, element)
		if err != nil {
			return errors.New(fmt.Sprintf("Error during writing element (Error: %s)", err))
		}
	}

	return nil
}
