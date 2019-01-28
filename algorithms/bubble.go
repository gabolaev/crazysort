package algorithms

// BubbleSorter ...
type BubbleSorter struct {
}

// Sort ...
func (bs BubbleSorter) Sort(slice []int, compare func(a, b int) bool) []int {
	swapped := false
	for !swapped {
		swapped = true
		for index, element := range slice[:len(slice)-1] {
			if compare(slice[index+1], element) {
				swapped = false
				slice[index+1], slice[index] = element, slice[index+1]
			}
		}
	}
	return slice
}
