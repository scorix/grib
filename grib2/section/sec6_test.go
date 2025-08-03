package section_test

import (
	"testing"

	"github.com/scorix/grib/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection6FromBytes_WithBitMap(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x0e, // length: 14 octets
		0x06, // section number: 6
		0x00, // bit-map indicator: bit-map applies (0)
		// Bit-map data (8 octets)
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section6, err := section.NewSection6FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section6.Length(), uint32(14))
	assert.Equal(t, section6.SectionNumber(), uint8(6))
	assert.Equal(t, section6.BitMapIndicator(), uint8(0))                                      // bit-map applies
	assert.True(t, section6.HasBitMap())                                                       // has bit-map data
	assert.Equal(t, section6.BitMap(), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}) // all bits set
}

func TestNewSection6FromBytes_WithoutBitMap(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x06, // length: 6 octets
		0x06, // section number: 6
		0xff, // bit-map indicator: bit-map does not apply (255)

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section6, err := section.NewSection6FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section6.Length(), uint32(6))
	assert.Equal(t, section6.SectionNumber(), uint8(6))
	assert.Equal(t, section6.BitMapIndicator(), uint8(255)) // bit-map does not apply
	assert.False(t, section6.HasBitMap())                   // no bit-map data
	assert.Empty(t, section6.BitMap())                      // empty bit-map
}
