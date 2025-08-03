package reader

import (
	"github.com/scorix/grib/grib2/spec"
)

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
