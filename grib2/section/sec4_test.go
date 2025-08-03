package section_test

import (
	"testing"

	"github.com/scorix/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection4FromBytes(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x22, // length: 34 octets
		0x04,       // section number: 4
		0x00, 0x00, // number of coordinate values: 0
		0x00, 0x00, // product definition template number: analysis/forecast at horizontal level (0)
		// Product definition template 4.0 (analysis/forecast at horizontal level)
		0x00,       // parameter category: temperature (0)
		0x00,       // parameter number: temperature (0)
		0x02,       // type of generating process: forecast (2)
		0x00,       // background generating process identifier
		0x00,       // analysis/forecast generating process identifier
		0x00, 0x00, // hours after reference time
		0x00,                   // minutes after reference time
		0x01,                   // indicator of unit of time range: hour (1)
		0x00, 0x00, 0x00, 0x0c, // forecast time in units: 12 hours
		0x01,                   // type of first fixed surface: surface (1)
		0x00,                   // scale factor of first fixed surface
		0x00, 0x00, 0x00, 0x00, // scaled value of first fixed surface
		0xff,                   // type of second fixed surface: missing (255)
		0xff,                   // scale factor of second fixed surface: missing
		0xff, 0xff, 0xff, 0xff, // scaled value of second fixed surface: missing

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section4, err := section.NewSection4FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section4.Length(), uint32(34))
	assert.Equal(t, section4.SectionNumber(), uint8(4))
	assert.Equal(t, section4.NumberOfCoordinateValues(), uint32(0))       // no coordinate values
	assert.Equal(t, section4.ProductDefinitionTemplateNumber(), uint8(0)) // template 4.0
	assert.Empty(t, section4.CoordinateValues())                          // no coordinate values
}
