package section_test

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/scorix/grib2/section"
	"github.com/stretchr/testify/suite"
)

// Section7ReaderTestSuite defines a test suite for different reader types with Section7
type Section7ReaderTestSuite struct {
	suite.Suite
	testData    []byte
	sectionData []byte
}

// SetupSuite runs once before all tests in the suite
func (s *Section7ReaderTestSuite) SetupSuite() {
	// Create test data
	s.testData = make([]byte, 100)
	for i := range s.testData {
		s.testData[i] = byte(i)
	}

	// Create complete section data with header
	s.sectionData = make([]byte, 0, len(s.testData)+5)
	length := uint32(len(s.testData) + 5)

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	s.sectionData = append(s.sectionData, lengthBytes...)
	s.sectionData = append(s.sectionData, 7) // section number
	s.sectionData = append(s.sectionData, s.testData...)
}

// testReaderCompatibility is a helper method that tests a reader implementation
func (s *Section7ReaderTestSuite) testReaderCompatibility(reader io.Reader) {
	section7, err := section.NewSection7FromReader(reader)
	s.Require().NoError(err)

	s.Assert().Equal(uint8(7), section7.SectionNumber())
	s.Assert().Equal(uint32(len(s.testData)), section7.DataSize())
	s.Assert().NoError(section7.LoadError())

	// Test all data access methods
	data := section7.Data()
	s.Assert().Equal(s.testData, data)

	// Test repeated reads
	data2 := section7.Data()
	s.Assert().Equal(s.testData, data2)

	// Test streaming access
	dataReader := section7.DataReader()
	streamData, err := io.ReadAll(dataReader)
	s.Assert().NoError(err)
	s.Assert().Equal(s.testData, streamData)
}

// TestSection7ReaderTypes runs the test suite for different reader types
func TestSection7ReaderTypes(t *testing.T) {
	suite.Run(t, new(Section7ReaderTestSuite))
}

// TestBytesReader tests bytes.Reader compatibility
func (s *Section7ReaderTestSuite) TestBytesReader() {
	reader := bytes.NewReader(s.sectionData)
	s.testReaderCompatibility(reader)
}

// TestGzipReader tests gzip.Reader compatibility
func (s *Section7ReaderTestSuite) TestGzipReader() {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(s.sectionData)
	s.Require().NoError(err)
	s.Require().NoError(w.Close())

	reader, err := gzip.NewReader(&buf)
	s.Require().NoError(err)
	s.testReaderCompatibility(reader)
}

// TestGzipReaderCorrupted tests corrupted gzip data handling
func (s *Section7ReaderTestSuite) TestGzipReaderCorrupted() {
	// Create invalid/corrupted gzip data
	corruptedGzip := []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff}

	reader, err := gzip.NewReader(bytes.NewReader(corruptedGzip))
	s.Require().NoError(err) // gzip.NewReader might succeed initially

	// This should fail when trying to read from corrupted gzip data
	_, err = section.NewSection7FromReader(reader)
	s.Assert().Error(err)
	// The error could be from gzip decompression or from section parsing
}

// TestHTTPReader tests HTTP response body compatibility
func (s *Section7ReaderTestSuite) TestHTTPReader() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(s.sectionData)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Test directly with HTTP response body (one-time reader)
	s.testReaderCompatibility(resp.Body)
}

// TestInvalidSectionNumber tests error handling for invalid section number
func (s *Section7ReaderTestSuite) TestInvalidSectionNumber() {
	// Invalid section number
	invalidData := make([]byte, len(s.sectionData))
	copy(invalidData, s.sectionData)
	invalidData[4] = 6 // Change section number from 7 to 6

	_, err := section.NewSection7FromReader(bytes.NewReader(invalidData))
	s.Assert().ErrorContains(err, "invalid section number")
}

// TestHTTPDisconnect tests HTTP connection disconnect and data truncation scenario
func (s *Section7ReaderTestSuite) TestHTTPDisconnect() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		// Write only partial data (header only) then disconnect
		_, _ = w.Write(s.sectionData[:3]) // Only 3 bytes, not enough for header

		// Force connection close by using Hijacker to close the underlying connection
		if hijacker, ok := w.(http.Hijacker); ok {
			conn, _, err := hijacker.Hijack()
			if err == nil {
				conn.Close() // Force disconnect
			}
		}
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// This should fail because we don't have enough data for a valid section header
	// and the connection might be interrupted during reading
	_, err = section.NewSection7FromReader(resp.Body)
	s.Assert().ErrorContains(err, "failed to read")
}

func TestSection7_ConcurrencySafety(t *testing.T) {
	// Create test data
	testData := make([]byte, 100)
	for i := range testData {
		testData[i] = byte(i)
	}

	// Create complete section data with header
	sectionData := make([]byte, 0, len(testData)+5)
	length := uint32(len(testData) + 5)

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	sectionData = append(sectionData, lengthBytes...)
	sectionData = append(sectionData, 7) // section number
	sectionData = append(sectionData, testData...)

	reader := bytes.NewReader(sectionData)
	section7, err := section.NewSection7FromReader(reader)
	if err != nil {
		t.Fatalf("Failed to create section7: %v", err)
	}

	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make([][]byte, numGoroutines)

	// Launch multiple goroutines that access data concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Mix of Data() and DataReader() calls
			if index%2 == 0 {
				results[index] = section7.Data()
			} else {
				reader := section7.DataReader()
				if reader != nil {
					data, _ := io.ReadAll(reader)
					results[index] = data
				}
			}
		}(i)
	}

	wg.Wait()

	// All results should be identical
	for i, result := range results {
		if !bytes.Equal(testData, result) {
			t.Errorf("Result %d should match expected data", i)
		}
	}

	if err := section7.LoadError(); err != nil {
		t.Errorf("LoadError should be nil, got: %v", err)
	}
}
