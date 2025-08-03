package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type section4 struct {
	length                          uint32
	sectionNumber                   uint8
	numberOfCoordinateValues        uint16
	productDefinitionTemplateNumber uint16
	productDefinitionTemplate       []byte
	coordinateValues                []float32
}

var _ Section4 = (*section4)(nil)

func (s *section4) Length() uint32 {
	return s.length
}

func (s *section4) SectionNumber() uint8 {
	return s.sectionNumber
}

func (s *section4) NumberOfCoordinateValues() uint32 {
	return uint32(s.numberOfCoordinateValues)
}

func (s *section4) ProductDefinitionTemplateNumber() uint8 {
	return uint8(s.productDefinitionTemplateNumber)
}

func (s *section4) CoordinateValues() []float32 {
	return s.coordinateValues
}

func NewSection4FromBytes(data []byte) (Section4, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("section4: data too short")
	}

	br := bytes.NewReader(data)

	var s section4
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.length))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.sectionNumber))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.numberOfCoordinateValues))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.productDefinitionTemplateNumber))

	if err != nil {
		return nil, err
	}

	// Calculate template size
	coordinateSize := int(s.numberOfCoordinateValues) * 4 // 4 bytes per float32
	templateSize := int(s.length) - 9 - coordinateSize
	if templateSize > 0 {
		s.productDefinitionTemplate = make([]byte, templateSize)
		if _, err := br.Read(s.productDefinitionTemplate); err != nil {
			return nil, fmt.Errorf("section4: failed to read product definition template: %w", err)
		}
	}

	// Read coordinate values if present
	if s.numberOfCoordinateValues > 0 {
		s.coordinateValues = make([]float32, s.numberOfCoordinateValues)
		for i := 0; i < int(s.numberOfCoordinateValues); i++ {
			if err := binary.Read(br, binary.BigEndian, &s.coordinateValues[i]); err != nil {
				return nil, fmt.Errorf("section4: failed to read coordinate values: %w", err)
			}
		}
	}

	return &s, nil
}
