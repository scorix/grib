package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type section2 struct {
	length        uint32
	sectionNumber uint8
	localUse      []byte
}

var _ Section2 = (*section2)(nil)

func (s *section2) Length() uint32 {
	return s.length
}

func (s *section2) SectionNumber() uint8 {
	return s.sectionNumber
}

func (s *section2) LocalUseData() []byte {
	return s.localUse
}

func NewSection2FromBytes(data []byte) (Section2, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("section2: data too short")
	}

	br := bytes.NewReader(data)

	var s section2
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.length))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.sectionNumber))

	localUseN := s.length - 5
	s.localUse = make([]byte, localUseN)
	if _, err := io.ReadFull(br, s.localUse); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &s, nil
}
