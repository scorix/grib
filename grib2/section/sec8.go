package section

import (
	"fmt"
)

type section8 struct {
	endMarker [4]byte
}

var _ Section8 = (*section8)(nil)

func (s *section8) EndMarker() [4]byte {
	return s.endMarker
}

func (s *section8) IsValid() bool {
	expected := [4]byte{'7', '7', '7', '7'}
	return s.endMarker == expected
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
