package section_test

import (
	"testing"

	"github.com/scorix/grib/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSection8FromBytes_Valid(t *testing.T) {
	data := []byte{
		'7', '7', '7', '7', // end marker: "7777"
	}

	section8, err := section.NewSection8FromBytes(data)
	require.NoError(t, err)

	assert.Equal(t, section8.EndMarker(), [4]byte{'7', '7', '7', '7'})
	assert.True(t, section8.IsValid())
}

func TestNewSection8FromBytes_Invalid(t *testing.T) {
	data := []byte{
		'7', '7', '7', '8', // invalid end marker: "7778"
	}

	section8, err := section.NewSection8FromBytes(data)
	require.Error(t, err)
	assert.Nil(t, section8)
	assert.Contains(t, err.Error(), "invalid end marker")
}

func TestNewSection8FromBytes_TooShort(t *testing.T) {
	data := []byte{
		'7', '7', '7', // only 3 bytes
	}

	section8, err := section.NewSection8FromBytes(data)
	require.Error(t, err)
	assert.Nil(t, section8)
	assert.Contains(t, err.Error(), "data too short")
}
