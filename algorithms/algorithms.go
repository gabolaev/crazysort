package algorithms

// SortingAlgorithm ...
type SortingAlgorithm interface {
	Sort([]int, func(a, b int) bool) []int
}
