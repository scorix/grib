package reader

import (
	"io"

	"github.com/scorix/grib/grib2/section"
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

// buildMessages constructs Message objects from cached sections
func (r *Reader) buildMessages() error {
	if len(r.sections) == 0 {
		return nil
	}

	var currentMessage *Message
	var currentSections []section.Section
	offset := int64(0)

	for _, sec := range r.sections {
		switch sec.SectionNumber() {
		case 0:
			// Start new message
			if currentMessage != nil {
				// Finalize previous message
				currentMessage.Sections = currentSections

				// Build section info for current message
				var secInfos []SectionInfo
				secOffset := currentMessage.Info.Offset
				for _, s := range currentSections {
					length := s.Length()

					secInfos = append(secInfos, SectionInfo{
						Number: s.SectionNumber(),
						Offset: secOffset,
						Length: length,
					})
					secOffset += int64(length)
				}
				currentMessage.Info.Sections = secInfos

				r.messages = append(r.messages, *currentMessage)
			}

			// Start new message
			if sec0, ok := sec.(section.Section0); ok {
				currentMessage = &Message{
					Info: MessageInfo{
						Index:      len(r.messages),
						Offset:     offset,
						Length:     sec0.TotalLength(),
						Discipline: sec0.Discipline(),
						Edition:    sec0.Edition(),
					},
				}
				currentSections = []section.Section{sec}
			}

		case 8:
			// End current message
			if currentMessage != nil {
				currentSections = append(currentSections, sec)
				currentMessage.Sections = currentSections

				// Build section info for current message
				var secInfos []SectionInfo
				secOffset := currentMessage.Info.Offset
				for _, s := range currentSections {
					length := s.Length()

					secInfos = append(secInfos, SectionInfo{
						Number: s.SectionNumber(),
						Offset: secOffset,
						Length: length,
					})
					secOffset += int64(length)
				}
				currentMessage.Info.Sections = secInfos

				r.messages = append(r.messages, *currentMessage)
				currentMessage = nil
				currentSections = nil
			}

		default:
			// Add section to current message
			if currentMessage != nil {
				currentSections = append(currentSections, sec)
			}
		}

		// Update offset based on section type
		length := sec.Length()
		offset += int64(length)
	}

	return nil
}
