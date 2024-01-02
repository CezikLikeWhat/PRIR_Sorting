package Utils

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"fmt"
	"github.com/CezikLikeWhat/PRIR_Sorting/logger"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"runtime"
	"sync"
)

const THRESHOLD int64 = 10_000_000

type item struct {
	value     int32
	fileIndex int
}

type minHeap []item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x any) {
	*h = append(*h, x.(item))
}

func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Sort(inputFileName string, outputFileName string, numberOfThreads int32) {
	runtime.GOMAXPROCS(int(numberOfThreads))

	inputFile, _ := os.Open(inputFileName)

	reader := bufio.NewReader(inputFile)

	var numberOfElements int64
	err := binary.Read(reader, binary.NativeEndian, &numberOfElements)
	if err != nil {
		logger.Fatal("Cannot read number of generated numbers from file | Error: %s", err)
	}

	inputFile.Close()

	if numberOfElements < THRESHOLD {
		fragment, _ := readFragment(inputFileName, 0, numberOfElements)
		slices.Sort(fragment)
		WriteDataToOutputFile(outputFileName, numberOfElements, fragment)
		return
	}

	var (
		fragmentSize = THRESHOLD / (int64(numberOfThreads))
		wg           sync.WaitGroup
		tempFiles    []string
	)

	for i := int64(0); i < numberOfElements; i += fragmentSize {
		wg.Add(1)
		go func(startIndex int64) {
			defer wg.Done()
			fragment, _ := readFragment(inputFileName, startIndex, fragmentSize)
			//mergeSort(fragment)
			//quickSort(fragment, 0, len(fragment)-1)
			sampleSort(fragment, 8) // TODO: Dodać możliwość sterowania ilością kubełków i zdecydować czy wyjebać te dodatkowe algorytmy sortowania
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
		err = binary.Read(file, binary.NativeEndian, &currentValue)
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

	h := &minHeap{}
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

		if nextValue, err := readNext(readBuffers[minItem.fileIndex]); err == nil {
			heap.Push(h, item{nextValue, minItem.fileIndex})
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

func mergeSort(slice []int32) {
	if len(slice) <= 1 {
		return
	}

	middle := len(slice) / 2
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mergeSort(slice[:middle])
	}()

	go func() {
		defer wg.Done()
		mergeSort(slice[middle:])
	}()

	wg.Wait()
	merge(slice, middle)
}

func merge(slice []int32, middle int) {
	left := append([]int32(nil), slice[:middle]...)
	right := append([]int32(nil), slice[middle:]...)

	l, r := 0, 0
	for i := 0; i < len(slice); i++ {
		if l >= len(left) {
			slice[i] = right[r]
			r++
			continue
		}
		if r >= len(right) {
			slice[i] = left[l]
			l++
			continue
		}
		if left[l] < right[r] {
			slice[i] = left[l]
			l++
		} else {
			slice[i] = right[r]
			r++
		}
	}
}

func quickSort(slice []int32, left, right int) {
	if left < right {
		// Partition the slice and get the pivot index
		pivotIndex := partition(slice, left, right)

		var wg sync.WaitGroup
		wg.Add(2)

		// Sort the left part in a new goroutine
		go func() {
			defer wg.Done()
			quickSort(slice, left, pivotIndex-1)
		}()

		// Sort the right part in a new goroutine
		go func() {
			defer wg.Done()
			quickSort(slice, pivotIndex+1, right)
		}()

		wg.Wait()
	}
}

func partition(slice []int32, left, right int) int {
	pivot := slice[right]
	i := left - 1

	for j := left; j < right; j++ {
		if slice[j] < pivot {
			i++
			slice[i], slice[j] = slice[j], slice[i]
		}
	}

	slice[i+1], slice[right] = slice[right], slice[i+1]
	return i + 1
}

func sampleSort(slice []int32, bucketCount int) {
	if len(slice) <= 1 || bucketCount <= 1 {
		slices.Sort(slice)
		return
	}

	// Wybór próbek i sortowanie ich
	samples := selectSamples(slice, bucketCount)
	slices.Sort(samples)

	// Podział na kubełki
	buckets := make([][]int32, bucketCount+1)
	for i := range buckets {
		buckets[i] = make([]int32, 0)
	}

	for _, v := range slice {
		idx := binarySearchInt32(samples, v)
		if idx == len(samples) {
			idx--
		}
		buckets[idx] = append(buckets[idx], v)
	}

	// Sortowanie kubełków równolegle
	var wg sync.WaitGroup
	for i := range buckets {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			slices.Sort(buckets[i])
		}(i)
	}
	wg.Wait()

	// Łączenie posortowanych kubełków
	idx := 0
	for _, bucket := range buckets {
		for _, v := range bucket {
			slice[idx] = v
			idx++
		}
	}
}

func selectSamples(slice []int32, count int) []int32 {
	step := len(slice) / count
	samples := make([]int32, 0, count)
	for i := 0; i < len(slice); i += step {
		samples = append(samples, slice[i])
	}
	return samples
}

func binarySearchInt32(slice []int32, target int32) int {
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

func removeTempFiles(files []string) {
	for _, filename := range files {
		os.Remove(filename)
	}
}
