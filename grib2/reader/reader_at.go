package reader

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib/grib2/section"
	"github.com/scorix/grib/grib2/spec"
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

// EachFlatMessage iterates through all flattened messages in the GRIB2 file
// Each nested message is flattened into multiple FlatMessage structs, one per data field
// Return true to continue iteration, false to stop
func (r *ReaderAt) EachFlatMessage(fn func(int, FlatMessage) bool) error {
	flatIndex := 0

	return r.EachMessage(func(msgIndex int, info MessageInfo) bool {
		// Build complete message from MessageInfo
		message, err := r.buildMessageFromInfo(info)
		if err != nil {
			// Log error and continue with next message
			// In a production system, you might want to handle this differently
			return true
		}

		// Flatten the message into multiple FlatMessage structs
		flatMessages := message.FlattenToFlatMessages()

		// Call callback for each flattened message
		for _, flatMsg := range flatMessages {
			// Update the index to be sequential across all flat messages
			flatMsg.Index = flatIndex
			if !fn(flatIndex, flatMsg) {
				return false // Stop iteration if callback returns false
			}
			flatIndex++
		}

		return true // Continue with next message
	})
}

// buildMessageFromInfo constructs a complete Message from MessageInfo
func (r *ReaderAt) buildMessageFromInfo(info MessageInfo) (*Message, error) {
	// Read all sections for this message
	var sections []section.Section
	for _, secInfo := range info.Sections {
		sec, err := r.ReadSectionAt(secInfo.Offset)
		if err != nil {
			return nil, fmt.Errorf("failed to read section %d at offset %d: %w", secInfo.Number, secInfo.Offset, err)
		}
		sections = append(sections, sec)
	}

	// Build the Message structure similar to Reader.buildMessages
	message := &Message{
		Info: info,
	}

	// Extract sections by type
	var sec0 section.Section0
	var sec1 section.Section1
	var sec8 section.Section8
	var currentLocalBlocks []spec.LocalBlock

	currentLocalBlock := spec.LocalBlock{}
	currentGridBlocks := []spec.GridBlock{}
	currentGridBlock := spec.GridBlock{}
	currentDataFields := []spec.DataField{}

	for _, sec := range sections {
		switch sec.SectionNumber() {
		case 0:
			sec0 = sec.(section.Section0)
			message.Message.Indicator = sec0
		case 1:
			sec1 = sec.(section.Section1)
			message.Message.Identification = sec1
		case 2:
			// Local use section - finalize previous local block and start new one
			if len(currentDataFields) > 0 {
				currentGridBlock.Fields = currentDataFields
				currentGridBlocks = append(currentGridBlocks, currentGridBlock)
				currentDataFields = []spec.DataField{}
				currentGridBlock = spec.GridBlock{}
			}
			if len(currentGridBlocks) > 0 {
				currentLocalBlock.Grids = currentGridBlocks
				currentLocalBlocks = append(currentLocalBlocks, currentLocalBlock)
				currentGridBlocks = []spec.GridBlock{}
				currentLocalBlock = spec.LocalBlock{}
			}
			currentLocalBlock.LocalUse = sec.(section.Section2)
		case 3:
			// Grid definition section - finalize previous grid block and start new one
			if len(currentDataFields) > 0 {
				currentGridBlock.Fields = currentDataFields
				currentGridBlocks = append(currentGridBlocks, currentGridBlock)
				currentDataFields = []spec.DataField{}
			}
			currentGridBlock.GridDef = sec.(section.Section3)
		case 4:
			// Product definition section - start new data field
			currentDataFields = append(currentDataFields, spec.DataField{
				ProductDef: sec.(section.Section4),
			})
		case 5:
			// Data representation section - complete current data field
			if len(currentDataFields) > 0 {
				currentDataFields[len(currentDataFields)-1].DataRep = sec.(section.Section5)
			}
		case 6:
			// Bitmap section - add to current data field
			if len(currentDataFields) > 0 {
				currentDataFields[len(currentDataFields)-1].Bitmap = sec.(section.Section6)
			}
		case 7:
			// Data section - complete current data field
			if len(currentDataFields) > 0 {
				currentDataFields[len(currentDataFields)-1].Data = sec.(section.Section7)
			}
		case 8:
			sec8 = sec.(section.Section8)
			message.Message.End = sec8
		}
	}

	// Finalize remaining blocks
	if len(currentDataFields) > 0 {
		currentGridBlock.Fields = currentDataFields
		currentGridBlocks = append(currentGridBlocks, currentGridBlock)
	}
	if len(currentGridBlocks) > 0 {
		currentLocalBlock.Grids = currentGridBlocks
		currentLocalBlocks = append(currentLocalBlocks, currentLocalBlock)
	}

	message.Message.Blocks = currentLocalBlocks
	return message, nil
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
