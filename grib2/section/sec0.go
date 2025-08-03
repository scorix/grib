package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type section0 struct {
	identifier  [4]byte
	reserved    [2]byte
	discipline  uint8
	edition     uint8
	totalLength uint64
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
