package src

import (
	"bufio"
	"bytes"
	"crazysort/algorithms"
	ds "crazysort/data-structures"
	utils "crazysort/utils"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	log "github.com/kataras/golog"
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

// partsCounter ...
func (crs *CrazySorter) partsCounter() (int, int, error) {
	file, err := os.Open(crs.FilePath)
	if err != nil {
		return 0, 0, err
	}

	crs.RAMSize *= math.Pow(1024, 3)
	fstat, _ := file.Stat()

	fileSize, _ := strconv.Atoi(strconv.FormatInt(fstat.Size(), 10))
	log.Infof("File size: %d bytes", fileSize)
	log.Infof("RAM size: %.0f bytes", crs.RAMSize)

	partsCount := int(math.Ceil(float64(fileSize) / crs.RAMSize))
	log.Infof("Parts count: %d", partsCount)
	return partsCount, int(math.Ceil(float64(fileSize / partsCount))), nil
}

// divide ...
func (crs *CrazySorter) divide(partsCount, partSize int) (err error) {
	file, err := os.Open(crs.FilePath)
	defer file.Close()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	for partID := 0; partID < partsCount; partID++ {
		partData := make([]int, 0, partSize)
		for subPartID := 0; subPartID < SubPartsCount; subPartID++ {

			buffer := utils.SafeReadPart(subPartID, partID, partSize/SubPartsCount, reader)

			for _, elem := range bytes.Split(buffer, NewLineDelim) {
				value, _ := strconv.Atoi(string(elem))
				partData = append(partData, value)
			}
		}

		log.Infof("Sorting part #%d : %v ...", partID, partData[0:10])
		partData = crs.SortAlgo.Sort(partData, func(a, b int) bool {
			return a < b
		})
		log.Infof("Sorted: %v ...", partData[0:30])

		partFileName := fmt.Sprintf("%s_%d", crs.FilePath, partID)
		log.Infof("Creating part_file: %s", partFileName)

		partFile, err := os.Create(partFileName)
		defer partFile.Close()
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(partFile)

		log.Infof("Writing %d lines", len(partData))
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

func (crs *CrazySorter) heapInitialFilling(tree *ds.SexyHeap, partsReaders []*bufio.Reader, subPartSize int) error {
	for partID, partReader := range partsReaders {
		buffer := utils.SafeReadPart(0, partID, subPartSize, partReader)
		for _, elem := range bytes.Split(buffer, NewLineDelim) {
			value, _ := strconv.Atoi(string(elem))
			tree.Insert(
				&ds.Pair{
					Value:  value,
					FileID: partID,
				},
			)
		}
	}
	return nil
}

func (crs *CrazySorter) resultWriteQueueOrganizer(heap *ds.SexyHeap, fileWriter io.Writer, partsReaders []io.Reader) {

}

// mergeParts ...
func (crs *CrazySorter) mergeParts() error {
	minExtractorHeap := ds.NewSexyHeap(func(a, b *ds.Pair) bool { return a.Value < b.Value })
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
	log.Infof("Reading %d bytes per each part", subPartSize)
	crs.heapInitialFilling(&minExtractorHeap, partsReaders, subPartSize)
	fmt.Println(minExtractorHeap.GetHead())

	return nil
}

// StartARiot ...
func (crs *CrazySorter) StartARiot() error {
	partsCount, partSize, err := crs.partsCounter()
	if err != nil {
		return err
	}
	err = crs.divide(partsCount, partSize)
	if err != nil {
		return err
	}
	log.Println(crs.Parts)
	crs.mergeParts()
	return nil
}
