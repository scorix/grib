package reader

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib/grib2/section"
)

// ReaderAt implements random-access reading of GRIB files using io.ReaderAt
type ReaderAt struct {
	reader io.ReaderAt
}

// NewReaderAt creates a new ReaderAt from an io.ReaderAt
func NewReaderAt(reader io.ReaderAt) *ReaderAt {
	return &ReaderAt{
		reader: reader,
	}
}

// ReadSectionAt reads a specific section at the given offset
func (r *ReaderAt) ReadSectionAt(offset int64) (section.Section, error) {
	// First read the section header to determine the exact length
	first4 := make([]byte, 4)
	_, err := r.reader.ReadAt(first4, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to read section header at offset %d: %w", offset, err)
	}

	var sectionLength int64

	switch {
	case string(first4) == "GRIB":
		// Section 0 is always 16 bytes
		sectionLength = 16
	case string(first4) == "7777":
		// Section 8 is always 4 bytes
		sectionLength = 4
	default:
		// For other sections, first 4 bytes contain the length
		sectionLength = int64(binary.BigEndian.Uint32(first4))
	}

	// Create a section reader with the exact length
	sectionReader := io.NewSectionReader(r.reader, offset, sectionLength)

	// Use the existing section reader logic
	return section.NewReader(sectionReader).ReadSection()
}

// EachMessage iterates through messages in the GRIB file
// The callback function receives the message index and MessageInfo
// Return true to continue iteration, false to stop
func (r *ReaderAt) EachMessage(fn func(int, MessageInfo) bool) error {
	offset := int64(0)
	messageIndex := 0

	for {
		// Read first 4 bytes to check for GRIB marker
		first4 := make([]byte, 4)
		_, err := r.reader.ReadAt(first4, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read at offset %d: %w", offset, err)
		}

		// Must start with GRIB marker
		if string(first4) != "GRIB" {
			return fmt.Errorf("invalid GRIB marker at offset %d", offset)
		}

		// Read Section 0 header (16 bytes total)
		header := make([]byte, 16)
		_, err = r.reader.ReadAt(header, offset)
		if err != nil {
			return fmt.Errorf("failed to read Section 0 header at offset %d: %w", offset, err)
		}

		// Parse Section 0 data
		discipline := header[6]
		edition := header[7]
		totalLength := binary.BigEndian.Uint64(header[8:16])

		// Scan sections within this message
		sections, err := r.scanSectionsInRange(offset, offset+int64(totalLength))
		if err != nil {
			return fmt.Errorf("failed to scan sections in message %d: %w", messageIndex, err)
		}

		messageInfo := MessageInfo{
			Index:      messageIndex,
			Offset:     offset,
			Length:     totalLength,
			Discipline: discipline,
			Edition:    edition,
			Sections:   sections,
		}

		// Call the callback function
		if !fn(messageIndex, messageInfo) {
			break // Stop iteration if callback returns false
		}

		offset += int64(totalLength)
		messageIndex++
	}

	return nil
}

// scanSectionsInRange scans sections within a specific byte range
func (r *ReaderAt) scanSectionsInRange(startOffset, endOffset int64) ([]SectionInfo, error) {
	var sections []SectionInfo
	offset := startOffset

	for offset < endOffset {
		// Read first 4 bytes to determine section type
		first4 := make([]byte, 4)
		_, err := r.reader.ReadAt(first4, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read at offset %d: %w", offset, err)
		}

		var sectionNumber uint8
		var sectionLength uint32

		switch {
		case string(first4) == "GRIB":
			// Section 0
			sectionNumber = 0
			sectionLength = 16
		case string(first4) == "7777":
			// Section 8
			sectionNumber = 8
			sectionLength = 4
		default:
			// Other sections have length in first 4 bytes
			sectionLength = binary.BigEndian.Uint32(first4)
			if sectionLength < 5 {
				return nil, fmt.Errorf("invalid section length %d at offset %d", sectionLength, offset)
			}

			// Read section number (5th byte)
			sectionNumberByte := make([]byte, 1)
			_, err = r.reader.ReadAt(sectionNumberByte, offset+4)
			if err != nil {
				return nil, fmt.Errorf("failed to read section number at offset %d: %w", offset+4, err)
			}
			sectionNumber = sectionNumberByte[0]
		}

		sections = append(sections, SectionInfo{
			Number: sectionNumber,
			Offset: offset,
			Length: sectionLength,
		})

		offset += int64(sectionLength)

		// Stop at end of message
		if sectionNumber == 8 {
			break
		}
	}

	return sections, nil
}
