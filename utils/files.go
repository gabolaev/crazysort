package utils

import (
	"bufio"
	"io"

	log "github.com/kataras/golog"
)

// SafeReadPart ...
func SafeReadPart(subPartID, partID, partSize int, reader *bufio.Reader) (buffer []byte) {
	buffer = make([]byte, partSize)
	n, err := io.ReadFull(reader, buffer)
	if err != nil {
		log.Warn(err)
	}

	appendix, _, err := reader.ReadLine()
	if err != nil {
		log.Warn(err)
	}
	buffer = append(buffer, appendix...)

	log.Infof(
		"%d bytes read from part #%d",
		n*(subPartID+1),
		partID,
	)
	return buffer
}
