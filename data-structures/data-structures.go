package datastructures

// Heap ...
type Heap interface {
	Insert(...(*Pair))
	Extract(int)
	GetHead() *Pair
	siftUp(int)
	siftDown(int)
	String() string
}

// Pair ...
type Pair struct {
	Value  int
	FileID int
}
