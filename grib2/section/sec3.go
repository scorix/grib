package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type section3 struct {
	length                       uint32
	sectionNumber                uint8
	gridDefinitionSource         uint8
	numberOfDataPoints           uint32
	optionalListOctets           uint8
	optionalListInterpretation   uint8
	gridDefinitionTemplateNumber uint16
	gridDefinitionTemplate       []byte
	optionalList                 []uint32
}

var _ Section3 = (*section3)(nil)

func (s *section3) Length() uint32 {
	return s.length
}

func (s *section3) SectionNumber() uint8 {
	return s.sectionNumber
}

func (s *section3) GridDefinitionSource() uint8 {
	return s.gridDefinitionSource
}

func (s *section3) NumberOfDataPoints() uint32 {
	return s.numberOfDataPoints
}

func (s *section3) GridDefinitionTemplateNumber() uint8 {
	return uint8(s.gridDefinitionTemplateNumber)
}

func (s *section3) OptionalListOctets() uint32 {
	return uint32(s.optionalListOctets)
}

func (s *section3) OptionalListInterpretation() uint8 {
	return s.optionalListInterpretation
}

func (s *section3) OptionalList() []uint32 {
	return s.optionalList
}

func NewSection3FromBytes(data []byte) (Section3, error) {
	if len(data) < 14 {
		return nil, fmt.Errorf("section3: data too short")
	}

	br := bytes.NewReader(data)

	var s section3
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.length))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.sectionNumber))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.gridDefinitionSource))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.numberOfDataPoints))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.optionalListOctets))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.optionalListInterpretation))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.gridDefinitionTemplateNumber))

	if err != nil {
		return nil, err
	}

	// Calculate template size
	templateSize := int(s.length) - 14 - int(s.optionalListOctets)
	if templateSize > 0 {
		s.gridDefinitionTemplate = make([]byte, templateSize)
		if _, err := br.Read(s.gridDefinitionTemplate); err != nil {
			return nil, fmt.Errorf("section3: failed to read grid definition template: %w", err)
		}
	}

	// Read optional list if present
	if s.optionalListOctets > 0 {
		numEntries := int(s.optionalListOctets) / 4
		s.optionalList = make([]uint32, numEntries)
		for i := 0; i < numEntries; i++ {
			if err := binary.Read(br, binary.BigEndian, &s.optionalList[i]); err != nil {
				return nil, fmt.Errorf("section3: failed to read optional list: %w", err)
			}
		}
	}

	return &s, nil
}
