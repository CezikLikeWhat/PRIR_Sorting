package Utils

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/configuration"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"runtime"
	"sync"
)

// Sort - Function that is the main function dealing with sorting numbers from the given file
func Sort(inputFileName string, outputFileName string, numberOfCPUs int32, numberOfBuckets int32) error {
	runtime.GOMAXPROCS(int(numberOfCPUs))

	inputFile, _ := os.Open(inputFileName)

	reader := bufio.NewReader(inputFile)

	var numberOfElements int64
	err := binary.Read(reader, binary.NativeEndian, &numberOfElements)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot read number of generated numbers from file | Error: %s", err))
	}

	inputFile.Close()

	if numberOfElements < configuration.THRESHOLD {
		fragment, err := readFragment(inputFileName, 0, numberOfElements)
		if err != nil {
			return err
		}

		slices.Sort(fragment)

		err = WriteDataToOutputFile(outputFileName, numberOfElements, fragment)
		if err != nil {
			return err
		}

		return nil
	}

	var (
		fragmentSize = configuration.THRESHOLD
		wg           sync.WaitGroup
	)
	tempFiles := make(chan string, numberOfElements/fragmentSize)

	for i := int64(0); i < numberOfElements; i += fragmentSize {
		wg.Add(1)
		go func(startIndex int64) {
			defer wg.Done()
			fragment, _ := readFragment(inputFileName, startIndex, fragmentSize)

			sampleSort(fragment, numberOfBuckets)

			filename, _ := writeSortedFragmentToTempFile(fragment, startIndex)

			tempFiles <- filename
		}(i)
	}

	go func() {
		wg.Wait()
		close(tempFiles)
	}()

	var filesNames []string
	for filename := range tempFiles {
		filesNames = append(filesNames, filename)
	}

	err = mergeSortedFiles(filesNames, outputFileName, numberOfElements)
	if err != nil {
		return err
	}

	return nil
}

func SortInMemory(inputFileName string, outputFileName string, numberOfCPUs int32, numberOfBuckets int32) error {
	runtime.GOMAXPROCS(int(numberOfCPUs))

	inputFile, _ := os.Open(inputFileName)
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	var numberOfElements int64
	err := binary.Read(reader, binary.NativeEndian, &numberOfElements)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot read number of generated numbers from file | Error: %s", err))
	}

	fragment := make([]int32, 0, numberOfElements)
	var currentValue int32
	for i := int64(0); i < numberOfElements; i++ {
		err = binary.Read(reader, binary.NativeEndian, &currentValue)
		if err != nil {
			return errors.New(fmt.Sprintf("Unknown error during reading file | Error: %s", err))
		}
		fragment = append(fragment, currentValue)
	}

	sampleSort(fragment, numberOfBuckets)

	err = WriteDataToOutputFile(outputFileName, numberOfElements, fragment)
	if err != nil {
		return err
	}

	return nil
}

// readFragment - Function designed to read a specific section of a file
func readFragment(filename string, startIndex int64, fragmentSize int64) ([]int32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	sizeOfInt32 := int64(binary.Size(int32(0)))
	sizeOfInt64 := int64(binary.Size(int64(0)))

	offset := sizeOfInt64 + startIndex*sizeOfInt32
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	fragment := make([]int32, 0, fragmentSize)
	var currentValue int32
	for i := int64(0); i < fragmentSize; i++ {
		err = binary.Read(reader, binary.NativeEndian, &currentValue)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		fragment = append(fragment, currentValue)
	}

	return fragment, nil
}

// writeSortedFragmentToTempFile - Function to save slice to temporary file
func writeSortedFragmentToTempFile(fragment []int32, numberOfFile int64) (string, error) {
	tempFileName := fmt.Sprintf("temp_files/temp_%d.bin", numberOfFile)
	file, err := os.Create(tempFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, num := range fragment {
		err := binary.Write(writer, binary.NativeEndian, num)
		if err != nil {
			return "", err
		}
	}

	return tempFileName, nil
}

// mergeSortedFiles - Function designed to merge data from multiple temporary files into an output file using a min-heap
func mergeSortedFiles(tempFiles []string, outputFileName string, numberOfElements int64) error {
	files := make([]*os.File, len(tempFiles))
	readBuffers := make([]*bufio.Reader, len(files))

	for i, fileName := range tempFiles {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		files[i] = file
		readBuffers[i] = bufio.NewReader(file)
	}

	h := &configuration.MinHeap{}
	heap.Init(h)
	for i, reader := range readBuffers {
		var value int32
		if err := binary.Read(reader, binary.NativeEndian, &value); err == nil {
			heap.Push(h, configuration.Item{Value: value, FileIndex: i})
		}
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	err = binary.Write(writer, binary.NativeEndian, numberOfElements)
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing number of elements to file | Error: %s", err))
	}

	for h.Len() > 0 {
		minItem := heap.Pop(h).(configuration.Item)
		binary.Write(writer, binary.NativeEndian, minItem.Value)

		if nextValue, err := readNext(readBuffers[minItem.FileIndex]); err == nil {
			heap.Push(h, configuration.Item{Value: nextValue, FileIndex: minItem.FileIndex})
		}
	}

	for _, file := range files {
		file.Close()
		err := os.Remove(file.Name())
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot close file %s | Error: %s", file.Name(), err))
		}
	}

	return nil
}

// readNext - Function designed to read another value from a file
func readNext(reader *bufio.Reader) (int32, error) {
	var value int32
	err := binary.Read(reader, binary.NativeEndian, &value)
	return value, err
}

// sampleSort - Function designed to perform a sample sort
func sampleSort(slice []int32, bucketCount int32) {
	if len(slice) <= 1 || bucketCount <= 1 {
		slices.Sort(slice)
		return
	}

	samples := selectSamples(slice, bucketCount)
	slices.Sort(samples)

	buckets := make([][]int32, bucketCount+1)
	for i := range buckets {
		buckets[i] = make([]int32, 0)
	}

	for _, v := range slice {
		idx := binarySearch(samples, v)
		if idx == len(samples) {
			idx--
		}
		buckets[idx] = append(buckets[idx], v)
	}

	var wg sync.WaitGroup
	for i := range buckets {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			slices.Sort(buckets[i])
		}(i)
	}
	wg.Wait()

	idx := 0
	for _, bucket := range buckets {
		for _, v := range bucket {
			slice[idx] = v
			idx++
		}
	}
}

// selectSamples - Function that has the task of selecting samples for the sample sort algorithm
func selectSamples(slice []int32, count int32) []int32 {
	sliceLength := int32(len(slice))
	step := sliceLength / count
	samples := make([]int32, 0, count)
	for i := int32(0); i < sliceLength; i += step {
		samples = append(samples, slice[i])
	}
	return samples
}

// binarySearch - Function designed to find a target in a sorted slice (binary search algorithm)
func binarySearch(slice []int32, target int32) int {
	left, right := 0, len(slice)
	for left < right {
		mid := left + (right-left)/2
		if slice[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}
