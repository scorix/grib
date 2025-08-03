package section

import (
	"fmt"
	"io"
)

type section8 struct {
	endMarker [4]byte
}

var _ Section8 = (*section8)(nil)

func (s *section8) Length() uint32 {
	return 4
}

func (s *section8) EndMarker() [4]byte {
	return s.endMarker
}

func (s *section8) SectionNumber() uint8 {
	return 8
}

func (s *section8) IsValid() bool {
	expected := [4]byte{'7', '7', '7', '7'}
	return s.endMarker == expected
}

func (s *section8) ReadSection(reader io.Reader) (Section, error) {
	return NewSection8FromReader(reader)
}

func NewSection8FromReader(reader io.Reader) (Section, error) {
	data := make([]byte, 4) // Section 8 is always 4 bytes
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	return NewSection8FromBytes(data)
}

func NewSection8FromBytes(data []byte) (Section8, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("section8: data too short")
	}

	var s section8
	copy(s.endMarker[:], data[:4])

	if !s.IsValid() {
		return nil, fmt.Errorf("section8: invalid end marker, expected '7777', got '%s'", string(s.endMarker[:]))
	}

	return &s, nil
}
