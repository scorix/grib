package section_test

import (
	"testing"

	"github.com/scorix/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection2FromBytes(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x0f, // length: 15 octets
		0x02,                                                       // section number: 2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // local use data: 10 octets

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section2, err := section.NewSection2FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section2.Length(), uint32(15))
	assert.Equal(t, section2.SectionNumber(), uint8(2))
	assert.Equal(t, section2.LocalUseData(), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}
