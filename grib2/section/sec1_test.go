package section_test

import (
	"testing"

	"github.com/scorix/grib/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection1FromBytes(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x15, // length: 21 octets
		0x01,       // section number: 1
		0x00, 0x07, // originating center: NCEP (7)
		0x00, 0x00, // originating subcenter: none (0)
		0x02,       // master tables version: 2
		0x00,       // local tables version: 0
		0x00,       // reference time significance: analysis (0)
		0x07, 0xe8, // year: 2024 (0x07e8)
		0x03,                                                             // month: March (3)
		0x0f,                                                             // day: 15
		0x0c,                                                             // hour: 12 UTC
		0x00,                                                             // minute: 0
		0x00,                                                             // second: 0
		0x00,                                                             // production status: operational products (0)
		0x00,                                                             // product type: analysis products (0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved

		// next section
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	section1, err := section.NewSection1FromBytes(data, false)
	require.NoError(t, err)

	assert.Equal(t, section1.Length(), uint32(21))
	assert.Equal(t, section1.SectionNumber(), uint8(1))
	assert.Equal(t, section1.OriginatingCenter(), uint16(7))    // NCEP
	assert.Equal(t, section1.OriginatingSubcenter(), uint16(0)) // none
	assert.Equal(t, section1.MasterTablesVersion(), uint8(2))
	assert.Equal(t, section1.LocalTablesVersion(), uint8(0))
	assert.Equal(t, section1.ReferenceTimeSignificance(), uint8(0)) // analysis
	assert.Equal(t, section1.Year(), uint16(2024))
	assert.Equal(t, section1.Month(), uint8(3)) // March
	assert.Equal(t, section1.Day(), uint8(15))
	assert.Equal(t, section1.Hour(), uint8(12)) // 12 UTC
	assert.Equal(t, section1.Minute(), uint8(0))
	assert.Equal(t, section1.Second(), uint8(0))
	assert.Equal(t, section1.ProductionStatus(), uint8(0)) // operational
	assert.Equal(t, section1.DataType(), uint8(0))         // analysis
}
