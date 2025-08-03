package spec

import "github.com/scorix/grib/grib2/section"

// DataField represents the innermost repeatable sequence (sections 4-7)
// This is the atomic unit of data in GRIB2 - a single data field with its metadata
type DataField struct {
	ProductDef section.Section4 // Section 4 - Product Definition (required)
	DataRep    section.Section5 // Section 5 - Data Representation (required)
	Bitmap     section.Section6 // Section 6 - Bitmap (optional, nil if not present)
	Data       section.Section7 // Section 7 - Data (required)
}

// GridBlock represents the middle repeatable sequence (sections 3-7)
// Contains a grid definition followed by one or more data fields using that grid
type GridBlock struct {
	GridDef section.Section3 // Section 3 - Grid Definition (required for this block)
	Fields  []DataField      // Data fields using this grid (sections 4-7 repeated)
}

// LocalBlock represents the outermost repeatable sequence (sections 2-7)
// Contains optional local use section followed by one or more grid blocks
type LocalBlock struct {
	LocalUse section.Section2 // Section 2 - Local Use (optional, nil if not present)
	Grids    []GridBlock      // Grid blocks (sections 3-7 repeated)
}

// Message represents a complete GRIB2 message according to WMO specification
// Supports the full three-level nesting structure defined in the standard
//
// Structure pattern:
// Section 0 (Indicator) - appears once at start
// Section 1 (Identification) - appears once after Section 0
// [Repeated blocks containing sections 2-7, 3-7, or 4-7]
// Section 8 (End) - appears once at end
//
// Three levels of repetition are supported:
// 1. Local blocks (sections 2-7 repeated)
// 2. Grid blocks (sections 3-7 repeated within a local block)
// 3. Data fields (sections 4-7 repeated within a grid block)
type Message struct {
	// Fixed sections - appear exactly once per message
	Indicator      section.Section0 // Section 0 - Indicator (required, appears once)
	Identification section.Section1 // Section 1 - Identification (required, appears once)

	// Variable sections - can be repeated according to the specification
	Blocks []LocalBlock // Local blocks (sections 2-7 repeated)

	// Terminator section - appears exactly once per message
	End section.Section8 // Section 8 - End (required, appears once)
}
