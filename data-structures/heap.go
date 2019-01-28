package datastructures

import (
	"fmt"
	"strconv"
)

// SexyHeap ...
type SexyHeap struct {
	Data    []int
	Size    int
	Compare func(int, int) bool
}

// NewSexyHeap ...
func NewSexyHeap(Compare func(int, int) bool) SexyHeap {
	return SexyHeap{
		Data:    make([]int, 0),
		Size:    0,
		Compare: Compare,
	}
}

func (h *SexyHeap) siftUp(idx int) {
	var parentIndex int
	for idx != 0 {
		parentIndex = (idx - 1) / 2
		if h.Compare(h.Data[parentIndex], h.Data[idx]) {
			return
		}
		h.Data[idx], h.Data[parentIndex] = h.Data[parentIndex], h.Data[idx]
		idx = parentIndex
	}
}

func (h *SexyHeap) siftDown(idx int) {
	headCandidate := idx
	for _, childIdx := range [...]int{idx*2 + 1, idx*2 + 2} {
		if childIdx < len(h.Data) && h.Compare(h.Data[childIdx], h.Data[headCandidate]) {
			headCandidate = childIdx
		}
	}
	if headCandidate != idx {
		h.Data[idx], h.Data[headCandidate] = h.Data[headCandidate], h.Data[idx]
		h.siftDown(headCandidate)
	}
}

func (h *SexyHeap) buildHeap(idx int) {
	for ; idx >= 0; idx-- {
		h.siftDown(idx)
	}
}

// Insert ...
func (h *SexyHeap) Insert(elem ...int) {
	h.Data = append(h.Data, elem...)
	newElementsCount := len(elem)
	h.Size += newElementsCount
	switch newElementsCount {
	case 1:
		h.siftUp(len(h.Data) - 1)
	default:
		h.buildHeap(len(h.Data)/2 - 1)
	}
}

// Extract ...
func (h *SexyHeap) Extract(idx int) (int, error) {
	if idx >= h.Size || idx < 0 {
		return 0, fmt.Errorf("Invalid idx: %d", idx)
	}

	returnValue := h.Data[idx]

	h.Data[h.Size-1], h.Data[idx] = h.Data[idx], h.Data[h.Size-1]
	h.Data = h.Data[:h.Size-1]
	h.Size--
	if idx < h.Size {
		h.siftDown(idx)
	}
	return returnValue, nil
}

// GetHead ...
func (h *SexyHeap) GetHead() (int, error) {
	return h.Extract(0)
}

func (h *SexyHeap) String() string {
	result := ""
	binCount := 1
	binCountHistory := make(map[int]bool)
	for _, elem := range h.Data {
		result += strconv.Itoa(elem)
		if binCount&(binCount-1) == 0 && !binCountHistory[binCount] {
			binCountHistory[binCount] = true
			binCount = 1
			result += "\n"
		} else {
			binCount++
			result += " "
		}
	}
	return result
}
