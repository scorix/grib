package section_test

import (
	"testing"

	"github.com/scorix/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection5FromBytes(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x17, // length: 23 octets
		0x05,                   // section number: 5
		0x00, 0x00, 0x27, 0x10, // number of data points: 10000
		0x00, 0x00, // data representation template number: grid point data - simple packing (0)
		// Data representation template 5.0 (grid point data - simple packing)
		0x42, 0x50, 0x00, 0x00, // reference value (IEEE 32-bit): 52.0
		0x00, 0x10, // binary scale factor: 16
		0x00, 0x05, // decimal scale factor: 5
		0x08, // number of bits used for each packed value: 8
		0x00, // type of original field values: floating point (0)

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section5, err := section.NewSection5FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section5.Length(), uint32(23))
	assert.Equal(t, section5.SectionNumber(), uint8(5))
	assert.Equal(t, section5.NumberOfDataPoints(), uint32(10000))          // 10000 data points
	assert.Equal(t, section5.DataRepresentationTemplateNumber(), uint8(0)) // simple packing
}
