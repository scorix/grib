package reader

import "github.com/scorix/grib/grib2/section"

// SectionInfo contains metadata about a section's location
type SectionInfo struct {
	Number uint8
	Offset int64
	Length uint32
}

// MessageInfo contains metadata about a GRIB2 message's location
type MessageInfo struct {
	Index      int           // Message index in file (0-based)
	Offset     int64         // Start offset of the message (Section 0 start)
	Length     uint64        // Total length of the message (from Section 0)
	Discipline uint8         // Discipline code from Section 0
	Edition    uint8         // GRIB edition from Section 0
	Sections   []SectionInfo // Sections within this message
}

// Message represents a complete GRIB2 message
type Message struct {
	Info     MessageInfo
	Sections []section.Section
}
