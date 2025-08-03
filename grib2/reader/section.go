package reader

import "github.com/scorix/grib/grib2/spec"

// SectionInfo contains metadata about a section's location
type SectionInfo struct {
	Number uint8
	Offset int64
	Length uint32
}

// Use GRIB2 specification types from spec package
type DataField = spec.DataField
type GridBlock = spec.GridBlock
type LocalBlock = spec.LocalBlock

// MessageInfo contains metadata about a GRIB2 message's location
type MessageInfo struct {
	Index      int           // Message index in file (0-based)
	Offset     int64         // Start offset of the message (Section 0 start)
	Length     uint64        // Total length of the message (from Section 0)
	Discipline uint8         // Discipline code from Section 0
	Edition    uint8         // GRIB edition from Section 0
	Sections   []SectionInfo // All sections within this message
}

// Message represents a complete GRIB2 message with reader-specific metadata
type Message struct {
	Info         MessageInfo // Reader-specific metadata
	spec.Message             // GRIB2 specification structure
}
