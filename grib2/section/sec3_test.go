package section_test

import (
	"testing"

	"github.com/scorix/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection3FromBytes(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x48, // length: 72 octets
		0x03,                   // section number: 3
		0x00,                   // grid definition source: specified in code table 3.0 (0)
		0x00, 0x00, 0x27, 0x10, // number of data points: 10000
		0x00,       // number of octets for optional list: 0
		0x00,       // interpretation of optional list: none (0)
		0x00, 0x00, // grid definition template number: lat/lon grid (0)
		// Grid definition template (lat/lon grid template 3.0, simplified)
		0x06,                   // shape of the earth: spherical with radius 6,367,470 m
		0xff,                   // scale factor of radius
		0xff, 0xff, 0xff, 0xff, // scaled value of radius
		0xff,                   // scale factor of major axis
		0xff, 0xff, 0xff, 0xff, // scaled value of major axis
		0xff,                   // scale factor of minor axis
		0xff, 0xff, 0xff, 0xff, // scaled value of minor axis
		0x00, 0x00, 0x00, 0x64, // Ni - number of points along parallel: 100
		0x00, 0x00, 0x00, 0x64, // Nj - number of points along meridian: 100
		0x00, 0x00, 0x00, 0x00, // basic angle of initial production domain: 0
		0x00, 0x00, 0x00, 0x00, // subdivisions of basic angle: missing
		0x00, 0x00, 0x00, 0x00, // latitude of first grid point: 0°
		0x00, 0x00, 0x00, 0x00, // longitude of first grid point: 0°
		0x00,                   // resolution and component flags
		0x05, 0x9d, 0x80, 0x00, // latitude of last grid point: 90°
		0x15, 0x7c, 0x00, 0x00, // longitude of last grid point: 359°
		0x00, 0x98, 0x96, 0x80, // i direction increment: 1°
		0x00, 0x98, 0x96, 0x80, // j direction increment: 1°
		0x00, // scanning mode flags

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section3, err := section.NewSection3FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section3.Length(), uint32(72))
	assert.Equal(t, section3.SectionNumber(), uint8(3))
	assert.Equal(t, section3.GridDefinitionSource(), uint8(0))         // specified in code table
	assert.Equal(t, section3.NumberOfDataPoints(), uint32(10000))      // 100x100 grid
	assert.Equal(t, section3.GridDefinitionTemplateNumber(), uint8(0)) // lat/lon grid
	assert.Equal(t, section3.OptionalListOctets(), uint32(0))          // no optional list
	assert.Equal(t, section3.OptionalListInterpretation(), uint8(0))   // none
	assert.Empty(t, section3.OptionalList())                           // no optional list
}
