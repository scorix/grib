package section_test

import (
	"testing"

	"github.com/scorix/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection0FromBytes(t *testing.T) {
	data := []byte{
		'G', 'R', 'I', 'B',
		0x00, 0x00, // reserved
		0x00,                                           // discipline
		0x02,                                           // edition
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // total length

		// next section
		0x00, 0x00, 0x00, 0x15, // length: 21 octets
		0x01, // section number: 1
	}

	section0, err := section.NewSection0FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section0.Discipline(), uint8(0), "discipline")
	assert.Equal(t, section0.Edition(), uint8(2), "edition")
	assert.Equal(t, section0.TotalLength(), uint64(16), "total length")
}
