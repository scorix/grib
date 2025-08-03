package section

import (
	"bytes"
	"fmt"
	"io"
)

type Reader struct {
	io.Reader
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{Reader: reader}
}

var readFunc = map[uint8]func(io.Reader) (Section, error){
	0: NewSection0FromReader,
	1: NewSection1FromReader,
	2: NewSection2FromReader,
	3: NewSection3FromReader,
	4: NewSection4FromReader,
	5: NewSection5FromReader,
	6: NewSection6FromReader,
	7: NewSection7FromReader,
	8: NewSection8FromReader,
}

func (r *Reader) ReadSection() (Section, error) {
	first4Bytes := make([]byte, 4)
	_, err := r.Read(first4Bytes)
	if err != nil {
		return nil, err
	}

	switch {
	case first4Bytes[0] == 'G' && first4Bytes[1] == 'R' && first4Bytes[2] == 'I' && first4Bytes[3] == 'B':
		return NewSection0FromReader(io.MultiReader(bytes.NewReader(first4Bytes), r.Reader))
	case first4Bytes[0] == '7' && first4Bytes[1] == '7' && first4Bytes[2] == '7' && first4Bytes[3] == '7':
		return NewSection8FromReader(io.MultiReader(bytes.NewReader(first4Bytes), r.Reader))
	default:
		nextByte := make([]byte, 1)
		_, err := r.Read(nextByte)
		if err != nil {
			return nil, err
		}

		reader := io.MultiReader(bytes.NewReader(first4Bytes), bytes.NewReader(nextByte), r.Reader)

		f, ok := readFunc[nextByte[0]]
		if !ok {
			return nil, fmt.Errorf("invalid section marker")
		}

		return f(reader)
	}
}
