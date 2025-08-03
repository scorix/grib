package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type section0 struct {
	identifier  [4]byte
	reserved    [2]byte
	discipline  uint8
	edition     uint8
	totalLength uint64
}

func (s *section0) Length() uint32 {
	return 16
}

func (s *section0) StartMarker() [4]byte {
	return s.identifier
}

func (s *section0) SectionNumber() uint8 {
	return 0
}

func (s *section0) Discipline() uint8 {
	return s.discipline
}

func (s *section0) Edition() uint8 {
	return s.edition
}

func (s *section0) TotalLength() uint64 {
	return s.totalLength
}

func (s *section0) ReadSection(reader io.Reader) (Section, error) {
	return NewSection0FromReader(reader)
}

func NewSection0FromReader(reader io.Reader) (Section, error) {
	data := make([]byte, 16) // Section 0 is always 16 bytes
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	return NewSection0FromBytes(data)
}

func NewSection0FromBytes(data []byte) (Section0, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf("section0: data too short")
	}

	if string(data[:4]) != "GRIB" {
		return nil, fmt.Errorf("section0: invalid GRIB identifier")
	}

	br := bytes.NewReader(data)

	var s section0
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.identifier))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.reserved))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.discipline))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.edition))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.totalLength))

	if err != nil {
		return nil, err
	}

	return &s, nil
}
