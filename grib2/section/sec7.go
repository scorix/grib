package section

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

type section7 struct {
	length        uint32
	sectionNumber uint8
	dataSize      uint32

	// Smart buffering - stores data as it's read
	buffer         []byte
	originalReader io.Reader
	isFullyRead    bool
	readErr        error

	// Concurrency control
	mu sync.RWMutex
}

var _ Section7 = (*section7)(nil)

func (s *section7) Length() uint32 {
	return s.length
}

func (s *section7) SectionNumber() uint8 {
	return s.sectionNumber
}

// readChunk reads data from the original reader into buffer as needed
// Returns the number of bytes read and any error
func (s *section7) readChunk(minBytes uint32) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Already fully read or has error
	if s.isFullyRead || s.readErr != nil {
		return 0, s.readErr
	}

	// Already have enough data
	if uint32(len(s.buffer)) >= minBytes {
		return 0, nil
	}

	// Calculate how much more we need to read
	targetSize := minBytes
	if targetSize > s.dataSize {
		targetSize = s.dataSize
	}

	// Read in chunks to avoid large allocations
	const chunkSize = 64 * 1024 // 64KB chunks
	totalRead := 0

	for uint32(len(s.buffer)) < targetSize && s.originalReader != nil {
		// Calculate chunk size for this iteration
		remainingNeed := targetSize - uint32(len(s.buffer))
		currentChunkSize := chunkSize
		if remainingNeed < chunkSize {
			currentChunkSize = int(remainingNeed)
		}

		chunk := make([]byte, currentChunkSize)
		n, err := s.originalReader.Read(chunk)

		if n > 0 {
			s.buffer = append(s.buffer, chunk[:n]...)
			totalRead += n
		}

		if err != nil {
			if err == io.EOF {
				s.isFullyRead = true
				s.originalReader = nil
			} else {
				s.readErr = err
			}
			break
		}

		// Check if we've read everything we need
		if uint32(len(s.buffer)) >= s.dataSize {
			s.isFullyRead = true
			s.originalReader = nil
			break
		}
	}

	return totalRead, s.readErr
}

func (s *section7) Data() []byte {
	// Try to read all data
	_, _ = s.readChunk(s.dataSize)

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.readErr != nil {
		return nil
	}

	// Return a copy to prevent external modification
	result := make([]byte, len(s.buffer))
	copy(result, s.buffer)
	return result
}

func (s *section7) DataReader() io.Reader {
	return &section7Reader{section: s, offset: 0}
}

func (s *section7) DataSize() uint32 {
	return s.dataSize
}

func (s *section7) LoadError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readErr
}

// section7Reader implements io.Reader for section7
type section7Reader struct {
	section *section7
	offset  uint32
}

func (r *section7Reader) Read(p []byte) (n int, err error) {
	// Ensure we have enough data buffered
	needed := r.offset + uint32(len(p))
	if needed > r.section.dataSize {
		needed = r.section.dataSize
	}

	_, _ = r.section.readChunk(needed)

	r.section.mu.RLock()
	defer r.section.mu.RUnlock()

	// Check for read errors
	if r.section.readErr != nil {
		return 0, r.section.readErr
	}

	// Calculate available data
	available := uint32(len(r.section.buffer)) - r.offset
	if available == 0 {
		return 0, io.EOF
	}

	// Copy data to user buffer
	toCopy := available
	if toCopy > uint32(len(p)) {
		toCopy = uint32(len(p))
	}

	copy(p, r.section.buffer[r.offset:r.offset+toCopy])
	r.offset += toCopy

	return int(toCopy), nil
}

// NewSection7FromReader creates a Section7 from a reader positioned at the section start.
// This reads the header (length and section number) and sets up smart buffering for the data.
// Works with any io.Reader including HTTP response bodies, gzip readers, files, etc.
//
// The reader should be positioned at the beginning of the section (including header).
func NewSection7FromReader(reader io.Reader) (Section7, error) {
	var length uint32
	var sectionNumber uint8

	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("section7: failed to read length: %w", err)
	}

	if err := binary.Read(reader, binary.BigEndian, &sectionNumber); err != nil {
		return nil, fmt.Errorf("section7: failed to read section number: %w", err)
	}

	if sectionNumber != 7 {
		return nil, fmt.Errorf("section7: invalid section number, expected 7, got %d", sectionNumber)
	}

	// Use LimitReader to ensure we only read the data portion
	dataSize := length - 5
	dataReader := io.LimitReader(reader, int64(dataSize))

	// Delegate to NewSection7FromDataReader for the actual construction
	return NewSection7FromDataReader(length, sectionNumber, dataReader), nil
}

// NewSection7FromDataReader creates a Section7 with smart buffering from a data-only reader.
// This is useful when you only have the data portion without the section header.
// Works with any io.Reader including HTTP response bodies, gzip readers, files, etc.
//
// Parameters:
// - length: Total section length including header (from section header)
// - sectionNumber: Section number (should be 7)
// - dataReader: Reader positioned at the start of section data (after length and section number)
func NewSection7FromDataReader(length uint32, sectionNumber uint8, dataReader io.Reader) Section7 {
	return &section7{
		length:         length,
		sectionNumber:  sectionNumber,
		dataSize:       length - 5, // length minus header (4 bytes length + 1 byte section number)
		originalReader: dataReader,
		buffer:         make([]byte, 0), // Start with empty buffer
	}
}
