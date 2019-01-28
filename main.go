package main

import (
	"crazysort/algorithms"
	"os"
)

func main() {
	crs := CrazySorter{
		FilePath: os.Args[1],
		SortAlgo: algorithms.QuickSorter{},
		Parts:    make([]string, 0),
		RAMSize:  0.1,
	}
	err := crs.StartARiot()
	if err != nil {
		panic(err)
	}
}
