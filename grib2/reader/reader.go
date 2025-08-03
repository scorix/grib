package reader

import (
	"io"

	"github.com/scorix/grib/grib2/section"
	"github.com/scorix/grib/grib2/spec"
)

// Reader implements sequential reading of GRIB files using io.Reader
type Reader struct {
	io.Reader
	sections []section.Section
	messages []Message
}

// NewReader creates a new Reader from an io.Reader
func NewReader(reader io.Reader) *Reader {
	return &Reader{
		Reader: reader,
	}
}

// ReadSection reads the next section from the GRIB file
func (r *Reader) ReadSection() (section.Section, error) {
	sec, err := section.NewReader(r.Reader).ReadSection()
	if err == nil {
		r.sections = append(r.sections, sec)

		// For Section 7, we need to consume the data to advance the reader
		if sec.SectionNumber() == 7 {
			if sec7, ok := sec.(section.Section7); ok {
				_ = sec7.Data() // Force reading all data to advance the underlying reader
			}
		}
	}
	return sec, err
}

// EachMessage iterates through messages in the GRIB file
// The callback function receives the message index and MessageInfo
// Return true to continue iteration, false to stop
func (r *Reader) EachMessage(fn func(int, MessageInfo) bool) error {
	// First ensure we have read all sections
	err := r.readAllSections()
	if err != nil {
		return err
	}

	// Build messages if not already done
	if len(r.messages) == 0 {
		err = r.buildMessages()
		if err != nil {
			return err
		}
	}

	// Iterate through messages
	for i, msg := range r.messages {
		if !fn(i, msg.Info) {
			break // Stop iteration if callback returns false
		}
	}

	return nil
}

// readAllSections reads all sections sequentially from the entire file
func (r *Reader) readAllSections() error {
	// Continue reading until EOF
	for {
		_, err := r.ReadSection()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Continue reading after Section 8 to capture multiple messages
		// Each message ends with Section 8, but there can be multiple messages
	}

	return nil
}

// buildMessages constructs Message objects from cached sections according to GRIB2 specification
// Supports three levels of repetition: sections 2-7, sections 3-7, and sections 4-7
func (r *Reader) buildMessages() error {
	if len(r.sections) == 0 {
		return nil
	}

	var currentMessage *Message
	var currentLocalBlock *LocalBlock
	var currentGridBlock *GridBlock
	var currentDataField *DataField
	offset := int64(0)

	for i, sec := range r.sections {
		switch sec.SectionNumber() {
		case 0:
			// Section 0: Start new message
			if currentMessage != nil {
				r.finalizeMessage(currentMessage)
			}

			if sec0, ok := sec.(section.Section0); ok {
				currentMessage = &Message{
					Info: MessageInfo{
						Index:      len(r.messages),
						Offset:     offset,
						Length:     sec0.TotalLength(),
						Discipline: sec0.Discipline(),
						Edition:    sec0.Edition(),
					},
					Message: spec.Message{
						Indicator: sec0,
					},
				}
				currentLocalBlock = nil
				currentGridBlock = nil
				currentDataField = nil
			}

		case 1:
			// Section 1: Identification (required, once per message)
			if currentMessage != nil {
				if sec1, ok := sec.(section.Section1); ok {
					currentMessage.Message.Identification = sec1
				}
			}

		case 2:
			// Section 2: Local Use (optional, starts new local block)
			if currentMessage != nil {
				// Finalize previous local block
				if currentLocalBlock != nil {
					r.finalizeLocalBlock(currentMessage, currentLocalBlock)
				}

				// Start new local block
				currentLocalBlock = &LocalBlock{}
				if sec2, ok := sec.(section.Section2); ok {
					currentLocalBlock.LocalUse = sec2
				}
				currentGridBlock = nil
				currentDataField = nil
			}

		case 3:
			// Section 3: Grid Definition (starts new grid block)
			if currentMessage != nil {
				// Ensure we have a local block (create empty one if needed)
				if currentLocalBlock == nil {
					currentLocalBlock = &LocalBlock{}
				}

				// Finalize previous grid block
				if currentGridBlock != nil {
					r.finalizeGridBlock(currentLocalBlock, currentGridBlock)
				}

				// Start new grid block
				if sec3, ok := sec.(section.Section3); ok {
					currentGridBlock = &GridBlock{
						GridDef: sec3,
					}
				}
				currentDataField = nil
			}

		case 4:
			// Section 4: Product Definition (starts new data field)
			if currentGridBlock != nil {
				// Finalize previous data field
				if currentDataField != nil {
					r.finalizeDataField(currentGridBlock, currentDataField)
				}

				// Start new data field
				if sec4, ok := sec.(section.Section4); ok {
					currentDataField = &DataField{
						ProductDef: sec4,
					}
				}
			}

		case 5:
			// Section 5: Data Representation (required for data field)
			if currentDataField != nil {
				if sec5, ok := sec.(section.Section5); ok {
					currentDataField.DataRep = sec5
				}
			}

		case 6:
			// Section 6: Bit-map (optional for data field)
			if currentDataField != nil {
				if sec6, ok := sec.(section.Section6); ok {
					currentDataField.Bitmap = sec6
				}
			}

		case 7:
			// Section 7: Data (required for data field)
			if currentDataField != nil {
				if sec7, ok := sec.(section.Section7); ok {
					currentDataField.Data = sec7
				}
			}

		case 8:
			// Section 8: End (finalize all structures)
			if currentMessage != nil {
				if sec8, ok := sec.(section.Section8); ok {
					// Finalize current data field
					if currentDataField != nil {
						r.finalizeDataField(currentGridBlock, currentDataField)
						currentDataField = nil
					}

					// Finalize current grid block
					if currentGridBlock != nil {
						r.finalizeGridBlock(currentLocalBlock, currentGridBlock)
						currentGridBlock = nil
					}

					// Finalize current local block
					if currentLocalBlock != nil {
						r.finalizeLocalBlock(currentMessage, currentLocalBlock)
						currentLocalBlock = nil
					}

					currentMessage.Message.End = sec8
					r.buildSectionInfo(currentMessage, i)
					r.messages = append(r.messages, *currentMessage)
					currentMessage = nil
				}
			}
		}

		// Update offset
		offset += int64(r.getSectionLength(sec))
	}

	return nil
}

// finalizeDataField adds the current data field to the grid block
func (r *Reader) finalizeDataField(gridBlock *GridBlock, dataField *DataField) {
	if gridBlock != nil && dataField != nil {
		gridBlock.Fields = append(gridBlock.Fields, *dataField)
	}
}

// finalizeGridBlock adds the current grid block to the local block
func (r *Reader) finalizeGridBlock(localBlock *LocalBlock, gridBlock *GridBlock) {
	if localBlock != nil && gridBlock != nil {
		localBlock.Grids = append(localBlock.Grids, *gridBlock)
	}
}

// finalizeLocalBlock adds the current local block to the message
func (r *Reader) finalizeLocalBlock(message *Message, localBlock *LocalBlock) {
	if message != nil && localBlock != nil {
		message.Message.Blocks = append(message.Message.Blocks, *localBlock)
	}
}

// finalizeMessage adds the message to the messages list (fallback for incomplete messages)
func (r *Reader) finalizeMessage(message *Message) {
	if message != nil {
		r.messages = append(r.messages, *message)
	}
}

// buildSectionInfo builds section info for a message
func (r *Reader) buildSectionInfo(msg *Message, endIndex int) {
	var sections []SectionInfo
	secOffset := msg.Info.Offset

	// Find start of this message
	startIndex := 0
	for i := endIndex; i >= 0; i-- {
		if r.sections[i].SectionNumber() == 0 {
			startIndex = i
			break
		}
	}

	// Build section info
	for i := startIndex; i <= endIndex; i++ {
		sec := r.sections[i]
		length := r.getSectionLength(sec)

		sections = append(sections, SectionInfo{
			Number: sec.SectionNumber(),
			Offset: secOffset,
			Length: length,
		})
		secOffset += int64(length)
	}

	msg.Info.Sections = sections
}

// getSectionLength returns the length of a section
func (r *Reader) getSectionLength(sec section.Section) uint32 {
	switch st := sec.(type) {
	case section.Section0:
		return 16
	case section.Section1:
		return st.Length()
	case section.Section2:
		return st.Length()
	case section.Section3:
		return st.Length()
	case section.Section4:
		return st.Length()
	case section.Section5:
		return st.Length()
	case section.Section6:
		return st.Length()
	case section.Section7:
		return st.Length()
	case section.Section8:
		return 4
	default:
		return 0
	}
}
