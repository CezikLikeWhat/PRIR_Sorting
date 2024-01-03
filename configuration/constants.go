package configuration

import "time"

const THRESHOLD int64 = 10_000_000

var (
	IsDebugMode   = false
	TimeMeasuring = false
	AppTimer      = Timer{}
)

type Timer struct {
	Start   time.Time
	Elapsed time.Duration
}

type Item struct {
	Value     int32
	FileIndex int
}

type MinHeap []Item

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Value < h[j].Value }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(Item))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
