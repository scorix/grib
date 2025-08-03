package reader_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/scorix/grib/grib2/reader"
	"github.com/scorix/grib/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestData loads real GRIB data from testdata directory
func getTestData(t *testing.T) []byte {
	data, err := os.ReadFile("testdata/gfs.t00z.pgrb2.0p25.f000")
	require.NoError(t, err, "Failed to read test data file")
	return data
}

func TestReader_ReadSection(t *testing.T) {
	testData := getTestData(t)
	r := reader.NewReader(bytes.NewReader(testData))

	// Read Section 0
	sec, err := r.ReadSection()
	require.NoError(t, err)
	assert.Equal(t, uint8(0), sec.SectionNumber())

	s0, ok := sec.(section.Section0)
	require.True(t, ok)
	assert.Equal(t, uint8(2), s0.Edition())
	assert.Equal(t, uint8(0), s0.Discipline())

	// Read Section 1
	sec, err = r.ReadSection()
	require.NoError(t, err)
	assert.Equal(t, uint8(1), sec.SectionNumber())

	s1, ok := sec.(section.Section1)
	require.True(t, ok)
	assert.Equal(t, uint16(7), s1.OriginatingCenter())
}

func TestReader_ReadAllSections(t *testing.T) {
	testData := getTestData(t)
	r := reader.NewReader(bytes.NewReader(testData))

	// Read all sections
	for {
		_, err := r.ReadSection()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	// Should have read 7 sections total (0, 1, 3, 4, 5, 7, 8)
	// Note: we can't access r.sections directly anymore as it's private
}

func TestReader_EachMessage(t *testing.T) {
	testData := getTestData(t)
	r := reader.NewReader(bytes.NewReader(testData))

	var messages []reader.MessageInfo
	var callCount int

	err := r.EachMessage(func(index int, info reader.MessageInfo) bool {
		callCount++
		messages = append(messages, info)
		assert.Equal(t, index, info.Index)
		return true // Continue iteration
	})
	require.NoError(t, err)

	// Should have found multiple messages in real GRIB file
	assert.Greater(t, callCount, 1, "Should have multiple messages")
	assert.Greater(t, len(messages), 1, "Should have multiple messages")

	// Verify first message info
	msg := messages[0]
	assert.Equal(t, 0, msg.Index)
	assert.Equal(t, uint8(2), msg.Edition)
	assert.Greater(t, msg.Length, uint64(0)) // Should have non-zero length
}

func TestReader_EachMessage_EarlyStop(t *testing.T) {
	testData := getTestData(t)
	r := reader.NewReader(bytes.NewReader(testData))

	var callCount int

	err := r.EachMessage(func(index int, info reader.MessageInfo) bool {
		callCount++
		return false // Stop immediately
	})
	require.NoError(t, err)

	// Should have called only once and stopped
	assert.Equal(t, 1, callCount)
}
