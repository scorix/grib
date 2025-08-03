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

func TestReader_EachFlatMessage(t *testing.T) {
	testData := getTestData(t)
	r := reader.NewReader(bytes.NewReader(testData))

	var flatMessages []reader.FlatMessage
	var callCount int

	err := r.EachFlatMessage(func(index int, flatMsg reader.FlatMessage) bool {
		callCount++
		flatMessages = append(flatMessages, flatMsg)
		return true // Continue processing
	})
	require.NoError(t, err)

	// Verify we got flattened messages
	assert.Greater(t, callCount, 0, "Should have processed at least one flat message")
	assert.Equal(t, len(flatMessages), callCount, "Should have same number of messages as calls")

	t.Logf("Processed %d flattened messages", callCount)

	// Verify each flat message has expected structure
	for i, flatMsg := range flatMessages {
		assert.Equal(t, i, flatMsg.Index, "Flat message index should be sequential")
		assert.Greater(t, flatMsg.Length, uint64(0), "Flat message should have non-zero length")
		assert.Equal(t, 2, flatMsg.Edition, "Should be GRIB2")

		// Verify template information is populated
		assert.NotNil(t, flatMsg.Indicator, "Should have Section 0")
		assert.NotNil(t, flatMsg.Identification, "Should have Section 1")
		assert.NotNil(t, flatMsg.GridDef, "Should have Section 3")
		assert.NotNil(t, flatMsg.ProductDef, "Should have Section 4")
		assert.NotNil(t, flatMsg.DataRepSec, "Should have Section 5")
		assert.NotNil(t, flatMsg.Data, "Should have Section 7")
		assert.NotNil(t, flatMsg.End, "Should have Section 8")

		t.Logf("Flat message %d: Discipline=%d, Centre=%d, Product.Category=%d",
			i, flatMsg.Discipline, flatMsg.Centre, flatMsg.Product.Category)
	}
}

func TestReader_EachFlatMessage_EarlyStop(t *testing.T) {
	testData := getTestData(t)
	r := reader.NewReader(bytes.NewReader(testData))

	var callCount int

	err := r.EachFlatMessage(func(index int, flatMsg reader.FlatMessage) bool {
		callCount++
		return false // Stop immediately
	})
	require.NoError(t, err)

	// Should have called only once and stopped
	assert.Equal(t, 1, callCount)
}
