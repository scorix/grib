package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type section5 struct {
	length                           uint32
	sectionNumber                    uint8
	numberOfDataPoints               uint32
	dataRepresentationTemplateNumber uint16
	dataRepresentationTemplate       []byte
}

var _ Section5 = (*section5)(nil)

func (s *section5) Length() uint32 {
	return s.length
}

func (s *section5) SectionNumber() uint8 {
	return s.sectionNumber
}

func (s *section5) NumberOfDataPoints() uint32 {
	return s.numberOfDataPoints
}

func (s *section5) DataRepresentationTemplateNumber() uint8 {
	return uint8(s.dataRepresentationTemplateNumber)
}

func (s *section5) ReadSection(reader io.Reader) (Section, error) {
	return NewSection5FromReader(reader)
}

func NewSection5FromReader(reader io.Reader) (Section, error) {
	var length uint32
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(reader, lengthBytes)
	if err != nil {
		return nil, err
	}
	length = binary.BigEndian.Uint32(lengthBytes)

	data := make([]byte, length)
	copy(data[:4], lengthBytes)
	_, err = io.ReadFull(reader, data[4:])
	if err != nil {
		return nil, err
	}

	return NewSection5FromBytes(data)
}

func NewSection5FromBytes(data []byte) (Section5, error) {
	if len(data) < 11 {
		return nil, fmt.Errorf("section5: data too short")
	}

	br := bytes.NewReader(data)

	var s section5
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.length))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.sectionNumber))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.numberOfDataPoints))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.dataRepresentationTemplateNumber))

	if err != nil {
		return nil, err
	}

	// Calculate template size
	templateSize := int(s.length) - 11
	if templateSize > 0 {
		s.dataRepresentationTemplate = make([]byte, templateSize)
		if _, err := br.Read(s.dataRepresentationTemplate); err != nil {
			return nil, fmt.Errorf("section5: failed to read data representation template: %w", err)
		}
	}

	return &s, nil
}
