package datastructures

// Heap ...
type Heap interface {
	Insert(...int)
	Extract(int)
	GetHead() int
	siftUp(int)
	siftDown(int)
	String() string
}
