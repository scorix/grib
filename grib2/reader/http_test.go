package reader_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/scorix/grib/grib2/reader"
	"github.com/scorix/grib/grib2/section"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getCurrentGFSURL generates a URL for recent GFS data
func getCurrentGFSURL() string {
	// Use a date from a few days ago to ensure data availability
	now := time.Now().UTC().AddDate(0, 0, -1).Truncate(time.Hour * 6) // 2 days ago
	hour := now.Hour()

	return fmt.Sprintf("https://noaa-gfs-bdp-pds.s3.amazonaws.com/gfs.%04d%02d%02d/%02d/atmos/gfs.t%02dz.sfluxgrbf000.grib2",
		now.Year(), now.Month(), now.Day(), hour, hour)
}

func TestHTTPReaderAt_RealGRIB(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	url := getCurrentGFSURL()
	t.Logf("Testing with URL: %s", url)

	// Create HTTP ReaderAt
	httpReader, err := reader.NewHTTPReaderAt(url)
	if err != nil {
		t.Skipf("Failed to create HTTP reader (server may be unavailable): %v", err)
	}

	t.Logf("File size: %d bytes (%.2f MB)", httpReader.Size(), float64(httpReader.Size())/(1024*1024))

	// Create GRIB ReaderAt
	r := reader.NewReaderAt(httpReader)

	// Test message functionality using EachMessage iterator - only check first message
	var firstMessage *reader.MessageInfo
	var messageCount int

	err = r.EachMessage(func(index int, info reader.MessageInfo) bool {
		messageCount++
		if index == 0 {
			firstMessage = &info
		}
		return index == 0 // Only process first message
	})
	require.NoError(t, err)

	t.Logf("Found first message in GRIB file (total count not fully scanned)")
	require.NotNil(t, firstMessage)

	// Test first message if found
	t.Logf("Message 0 - Offset: %d, Length: %d bytes (%.2f MB), Discipline: %d, Edition: %d",
		firstMessage.Offset, firstMessage.Length, float64(firstMessage.Length)/(1024*1024), firstMessage.Discipline, firstMessage.Edition)
	t.Logf("Message 0 has %d sections", len(firstMessage.Sections))

	// Test reading individual sections
	if len(firstMessage.Sections) > 0 {
		sec0, err := r.ReadSectionAt(firstMessage.Sections[0].Offset)
		require.NoError(t, err)
		assert.Equal(t, uint8(0), sec0.SectionNumber())
	}

	// Test reading section details
	if len(firstMessage.Sections) > 1 {
		sec1, err := r.ReadSectionAt(firstMessage.Sections[1].Offset)
		if err == nil {
			if s1, ok := sec1.(section.Section1); ok {
				t.Logf("Center: %d, Year: %d, Month: %d, Day: %d, Hour: %d",
					s1.OriginatingCenter(), s1.Year(), s1.Month(), s1.Day(), s1.Hour())
			}
		}
	}
}

func TestHTTPReaderAt_Messages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	url := getCurrentGFSURL()
	httpReader, err := reader.NewHTTPReaderAt(url)
	if err != nil {
		t.Skipf("Failed to create HTTP reader: %v", err)
	}

	r := reader.NewReaderAt(httpReader)

	// Test EachMessage iterator - get message metadata without downloading content
	var messages []reader.MessageInfo
	err = r.EachMessage(func(index int, info reader.MessageInfo) bool {
		messages = append(messages, info)
		return index < 2 // Only process first 3 messages (index 0, 1, 2)
	})
	require.NoError(t, err)

	assert.Greater(t, len(messages), 0, "Should have at least one message")

	t.Logf("Processed %d messages using EachMessage iterator", len(messages))

	// Verify message metadata
	for i, msg := range messages {
		t.Logf("Message %d:", i)
		t.Logf("  Offset: %d bytes", msg.Offset)
		t.Logf("  Length: %d bytes (%.2f MB)", msg.Length, float64(msg.Length)/(1024*1024))
		t.Logf("  Discipline: %d", msg.Discipline)
		t.Logf("  Edition: %d", msg.Edition)
		t.Logf("  Sections: %d", len(msg.Sections))

		// Verify message structure
		assert.Equal(t, i, msg.Index)
		assert.Equal(t, uint8(2), msg.Edition, "Should be GRIB2")
		assert.Greater(t, msg.Length, uint64(0), "Message should have non-zero length")
		assert.Greater(t, len(msg.Sections), 0, "Message should have sections")

		// Verify sections start with Section 0 and end with Section 8
		if len(msg.Sections) > 0 {
			assert.Equal(t, uint8(0), msg.Sections[0].Number, "First section should be Section 0")
			assert.Equal(t, uint8(8), msg.Sections[len(msg.Sections)-1].Number, "Last section should be Section 8")
		}
	}

	// Test reading first few messages
	testCount := min(len(messages), 2) // Test first 2 messages only
	for msgIdx := 0; msgIdx < testCount; msgIdx++ {
		t.Logf("Testing reading message %d...", msgIdx)

		// Test reading first section from this message
		if len(messages[msgIdx].Sections) > 0 {
			firstSection := messages[msgIdx].Sections[0]
			sec, err := r.ReadSectionAt(firstSection.Offset)
			require.NoError(t, err)
			assert.Equal(t, firstSection.Number, sec.SectionNumber())

			// Log basic section information
			if msgIdx == 0 {
				t.Logf("Message %d Section 0 read successfully", msgIdx)
			}
		}

		t.Logf("Successfully read message %d structure", msgIdx)
	}
}

func TestHTTPReaderAt_GetMessageInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	url := getCurrentGFSURL()
	httpReader, err := reader.NewHTTPReaderAt(url)
	if err != nil {
		t.Skipf("Failed to create HTTP reader: %v", err)
	}

	r := reader.NewReaderAt(httpReader)

	// Test EachMessage - scan only first few messages efficiently
	var messages []reader.MessageInfo
	err = r.EachMessage(func(index int, info reader.MessageInfo) bool {
		messages = append(messages, info)
		return index < 2 // Only process first 3 messages (index 0, 1, 2)
	})
	require.NoError(t, err)

	// Verify message info structure
	assert.Greater(t, len(messages), 0, "Should have at least one message")
	assert.LessOrEqual(t, len(messages), 3, "Should not exceed limit")

	t.Logf("Scanned %d messages efficiently using EachMessage iterator", len(messages))
}
