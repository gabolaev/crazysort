package main

import (
	"bufio"
	"bytes"
	"crazysort/algorithms"
	datastructures "crazysort/data-structures"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

var (
	// NewLineDelim ...
	NewLineDelim = []byte("\n")
	// SubPartsCount ...
	SubPartsCount = 5
)

// SortType describes sorting strategy
type SortType int

// CrazySorter ...
type CrazySorter struct {
	FilePath string
	SortAlgo algorithms.SortingAlgorithm
	RAMSize  float64
	Parts    []string
}

// NewCrazySorter ...
func NewCrazySorter(filePath string, sa algorithms.SortingAlgorithm, ramsize float64) (crs CrazySorter) {
	return CrazySorter{
		FilePath: filePath,
		SortAlgo: sa,
		RAMSize:  ramsize,
	}
}

// PartsCounter ...
func (crs *CrazySorter) PartsCounter() (int, int, error) {
	file, err := os.Open(crs.FilePath)
	if err != nil {
		return 0, 0, err
	}

	crs.RAMSize *= math.Pow(1024, 3)
	fstat, _ := file.Stat()

	fileSize, _ := strconv.Atoi(strconv.FormatInt(fstat.Size(), 10))
	log.Printf("File size: %d bytes", fileSize)
	log.Printf("RAM size: %.0f bytes", crs.RAMSize)

	partsCount := int(math.Ceil(float64(fileSize) / crs.RAMSize))
	log.Printf("Parts count: %d", partsCount)
	return partsCount, int(math.Ceil(float64(fileSize / partsCount))), nil
}

// Divide ...
func (crs *CrazySorter) Divide(partsCount, partSize int) (err error) {
	file, err := os.Open(crs.FilePath)
	defer file.Close()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	for partID := 0; partID < partsCount; partID++ {
		partData := make([]int, 0, partSize)
		for subPart := 0; subPart < SubPartsCount; subPart++ {
			buffer := make([]byte, partSize/SubPartsCount)
			n, err := io.ReadFull(reader, buffer)
			if err != nil {
				log.Println(err)
			}

			appendix, _, err := reader.ReadLine()
			if err != nil {
				log.Println(err)
			}
			buffer = append(buffer, appendix...)

			log.Printf(
				"%d bytes readed at sub-part #%d (part #%d)",
				n*(subPart+1),
				subPart,
				partID,
			)

			for _, elem := range bytes.Split(buffer, NewLineDelim) {
				value, _ := strconv.Atoi(string(elem))
				partData = append(partData, value)
			}
		}
		log.Printf("Sorting part #%d : %v ...", partID, partData[0:10])
		partData = crs.SortAlgo.Sort(partData, func(a, b int) bool {
			return a < b
		})
		log.Printf("Sorted: %v ...", partData[0:30])

		partFileName := fmt.Sprintf("%s_%d", crs.FilePath, partID)
		log.Printf("Creating part_file: %s", partFileName)

		partFile, err := os.Create(partFileName)
		defer partFile.Close()
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(partFile)

		log.Printf("Writing %d lines", len(partData))
		for _, value := range partData {
			_, err := io.WriteString(writer, strconv.Itoa(value)+"\n")
			if err != nil {
				log.Println(err)
			}
		}
		writer.Flush()
		partFile.Close()
		crs.Parts = append(crs.Parts, partFileName)
	}
	return
}

// MergeParts ...
func (crs *CrazySorter) MergeParts() error {
	minExtractorHeap := datastructures.NewSexyHeap(func(a, b int) bool { return a < b })
	resultFile, err := os.Create(fmt.Sprintf("%s_sorted", crs.FilePath))
	defer resultFile.Close()
	if err != nil {
		return err
	}
	// resultWriter := bufio.NewWriter(resultFile)

	partsFiles := make([]*os.File, 0, len(crs.Parts))
	partsReaders := make([]*bufio.Reader, 0, len(crs.Parts))
	defer func() {
		for _, file := range partsFiles {
			file.Close()
		}
	}()

	for _, partPath := range crs.Parts {
		part, err := os.Open(partPath)
		if err != nil {
			return err
		}
		partsFiles = append(partsFiles, part)
		partsReaders = append(partsReaders, bufio.NewReader(part))
	}

	subPartSize := int(crs.RAMSize) / len(crs.Parts) / SubPartsCount
	log.Printf("Reading %d bytes per each part", subPartSize)
	for idx, partReader := range partsReaders {
		buffer := make([]byte, subPartSize)
		n, err := partReader.Read(buffer)
		if err != nil {
			return err
		}
		log.Printf(
			"%d bytes readed from part #%d",
			n,
			idx,
		)

		for _, elem := range bytes.Split(buffer, NewLineDelim) {
			value, _ := strconv.Atoi(string(elem))
			minExtractorHeap.Insert(value)
		}
	}
	return nil
}

// StartARiot ...
func (crs *CrazySorter) StartARiot() error {
	partsCount, partSize, err := crs.PartsCounter()
	if err != nil {
		return err
	}
	err = crs.Divide(partsCount, partSize)
	if err != nil {
		return err
	}
	log.Println(crs.Parts)
	// crs.MergeParts()
	return nil
}
