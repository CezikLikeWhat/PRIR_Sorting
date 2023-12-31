package Utils

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/logger"
	"golang.org/x/exp/slices"
	"os"
	"sync"
)

type item struct {
	value   int32
	fileIdx int
}

type itemHeap []item

func (h itemHeap) Len() int           { return len(h) }
func (h itemHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h itemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *itemHeap) Push(x any) {
	*h = append(*h, x.(item))
}

func (h *itemHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Sort(inputFileName string, outputFileName string, numberOfThreads int32, threshold int64) {
	inputFile, _ := os.Open(inputFileName)

	reader := bufio.NewReader(inputFile)

	var numberOfElements int64
	err := binary.Read(reader, binary.NativeEndian, &numberOfElements)
	if err != nil {
		logger.Fatal("Cannot read number of generated numbers from file | Error: %s", err)
	}

	inputFile.Close()

	if numberOfElements < 10_000_000 {
		fragment, _ := readFragment(inputFileName, 0, numberOfElements)
		slices.Sort(fragment)
		WriteDataToOutputFile(outputFileName, numberOfElements, fragment)
	}

	var fragmentSize int64 = 100_000_000
	//var fragmentSize int64 = numberOfElements
	var wg sync.WaitGroup
	var tempFiles []string

	for i := int64(0); i < numberOfElements; i += fragmentSize {
		wg.Add(1)
		go func(startIndex int64) {
			defer wg.Done()
			fragment, _ := readFragment(inputFileName, startIndex, fragmentSize)
			slices.Sort(fragment)
			filename, _ := writeSortedFragmentToTempFile(fragment, startIndex)
			tempFiles = append(tempFiles, filename)
		}(i)
	}

	wg.Wait()

	mergeSortedFiles(tempFiles, outputFileName, numberOfElements)

	removeTempFiles(tempFiles)
}

func readFragment(filename string, startIndex int64, fragmentSize int64) ([]int32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	offset := int64(8) + startIndex*int64(binary.Size(int32(0)))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	fragment := make([]int32, fragmentSize)
	for i := int64(0); i < fragmentSize; i++ {
		err = binary.Read(file, binary.NativeEndian, &fragment[i])
		if err != nil {
			return nil, err
		}
	}

	return fragment, nil
}

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

func mergeSortedFiles(tempFiles []string, outputFileName string, numberOfElements int64) error {
	files := make([]*os.File, len(tempFiles))
	for i, fileName := range tempFiles {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		files[i] = file
	}

	h := &itemHeap{}
	heap.Init(h)
	for i, file := range files {
		var value int32
		if err := binary.Read(file, binary.NativeEndian, &value); err == nil {
			heap.Push(h, item{value, i})
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
		logger.Fatal("Error writing NumberOfElements | Error: %s", err)
	}

	readBuffers := make([]*bufio.Reader, len(files))
	for i, file := range files {
		readBuffers[i] = bufio.NewReader(file)
	}

	for h.Len() > 0 {
		minItem := heap.Pop(h).(item)
		binary.Write(writer, binary.NativeEndian, minItem.value)

		if nextValue, err := readNext(readBuffers[minItem.fileIdx]); err == nil {
			heap.Push(h, item{nextValue, minItem.fileIdx})
		}
	}

	for _, file := range files {
		file.Close()
	}

	return nil
}

func readNext(reader *bufio.Reader) (int32, error) {
	var value int32
	err := binary.Read(reader, binary.NativeEndian, &value)
	return value, err
}

func removeTempFiles(files []string) {
	for _, filename := range files {
		os.Remove(filename)
	}
}
