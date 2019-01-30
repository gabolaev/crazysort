package main

import (
	"crazysort/algorithms"
	"crazysort/src"
	"os"
)

func main() {
	crs := src.CrazySorter{
		FilePath: os.Args[1],
		SortAlgo: algorithms.QuickSorter{},
		Parts:    make([]string, 0),
		RAMSize:  3,
	}
	err := crs.StartARiot()
	if err != nil {
		panic(err)
	}
}
