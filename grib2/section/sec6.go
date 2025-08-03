package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type section6 struct {
	length          uint32
	sectionNumber   uint8
	bitMapIndicator uint8
	bitMap          []byte
}

var _ Section6 = (*section6)(nil)

func (s *section6) Length() uint32 {
	return s.length
}

func (s *section6) SectionNumber() uint8 {
	return s.sectionNumber
}

func (s *section6) BitMapIndicator() uint8 {
	return s.bitMapIndicator
}

func (s *section6) BitMap() []byte {
	return s.bitMap
}

func (s *section6) HasBitMap() bool {
	return s.bitMapIndicator == 0
}

func (s *section6) ReadSection(reader io.Reader) (Section, error) {
	return NewSection6FromReader(reader)
}

func NewSection6FromReader(reader io.Reader) (Section, error) {
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

	return NewSection6FromBytes(data)
}

func NewSection6FromBytes(data []byte) (Section6, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("section6: data too short")
	}

	br := bytes.NewReader(data)

	var s section6
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.length))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.sectionNumber))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.bitMapIndicator))

	if err != nil {
		return nil, err
	}

	// Read bit-map data if present (only when bit-map indicator is 0)
	if s.bitMapIndicator == 0 {
		bitMapSize := int(s.length) - 6
		if bitMapSize > 0 {
			s.bitMap = make([]byte, bitMapSize)
			if _, err := br.Read(s.bitMap); err != nil {
				return nil, fmt.Errorf("section6: failed to read bit-map data: %w", err)
			}
		}
	}

	return &s, nil
}
