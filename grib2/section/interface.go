package section

import "io"

// Section0 represents the GRIB2 Indicator Section (Section 0)
// This section serves to identify the start of the record in a human readable form,
// indicate the total length of the message, and indicate the Edition number of GRIB used
// to construct or encode the message. For GRIB2, this section is always 16 octets long.
//
// Format:
// +-------------+------------------------------------------------------------------+
// | Octet Number| Content                                                          |
// +-------------+------------------------------------------------------------------+
// | 1-4         | 'GRIB' (Coded according to the International Alphabet Number 5) |
// | 5-6         | Reserved                                                         |
// | 7           | Discipline (From Table 0.0)                                     |
// | 8           | Edition number - 2 for GRIB2                                    |
// | 9-16        | Total length of GRIB message in octets (All sections)           |
// +-------------+------------------------------------------------------------------+
type Section0 interface {
	Discipline() uint8
	Edition() uint8
	TotalLength() uint64
}

// Section1 represents the GRIB2 Identification Section (Section 1)
// This section contains information about the originating center, reference time,
// and data processing status. The section length is typically 21 octets.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | Length of the section in octets (21 or N)                                   |
// | 5           | Number of the section (1)                                                   |
// | 6-7         | Identification of originating/generating center (See Table 0)               |
// | 8-9         | Identification of originating/generating subcenter (See Table C)            |
// | 10          | GRIB master tables version number (currently 2) (See Table 1.0)            |
// | 11          | Version number of GRIB local tables (see Table 1.1)                        |
// | 12          | Significance of reference time (See Table 1.2)                              |
// | 13-14       | Year (4 digits)                                                             |
// | 15          | Month                                                                        |
// | 16          | Day                                                                          |
// | 17          | Hour                                                                         |
// | 18          | Minute                                                                       |
// | 19          | Second                                                                       |
// | 20          | Production Status of Processed data (See Table 1.3)                         |
// | 21          | Type of processed data in this GRIB message (See Table 1.4)                 |
// | 22-N        | Reserved                                                                     |
// +-------------+------------------------------------------------------------------------------+
type Section1 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Originating center information
	OriginatingCenter() uint16
	OriginatingSubcenter() uint16

	// Table versions
	MasterTablesVersion() uint8
	LocalTablesVersion() uint8

	// Reference time
	ReferenceTimeSignificance() uint8
	Year() uint16
	Month() uint8
	Day() uint8
	Hour() uint8
	Minute() uint8
	Second() uint8

	// Production status
	ProductionStatus() uint8
	DataType() uint8
}

// Section2 represents the GRIB2 Local Use Section (Section 2)
// This section is used for data that is specific to the originating center.
// It is optional and may be omitted if no local use data is needed.
//
// Format:
// +-------------+-------------------------------------+
// | Octet Number| Content                             |
// +-------------+-------------------------------------+
// | 1-4         | Length of the section in octets (N) |
// | 5           | Number of the section (2)           |
// | 6-N         | Local Use                           |
// +-------------+-------------------------------------+
//
// Note: Different centers may use different formats for the local use data.
// For example, NCEP subcenter 14 (MDL) uses octet 6 to indicate which
// local use table to use.
type Section2 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Local use data
	LocalUseData() []byte
}

// Section3 represents the GRIB2 Grid Definition Section (Section 3)
// This section defines the grid geometry and geographical information
// for the data points in the GRIB message.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | Length of the section in octets (nn)                                        |
// | 5           | Number of the section (3)                                                   |
// | 6           | Source of grid definition (See Table 3.0)                                   |
// | 7-10        | Number of data points                                                        |
// | 11          | Number of octets for optional list of numbers defining number of points     |
// | 12          | Interpretation of list of numbers defining number of points (See Table 3.11)|
// | 13-14       | Grid definition template number (= N) (See Table 3.1)                      |
// | 15-xx       | Grid definition template (See Template 3.N)                                 |
// | [xx+1]-nn   | Optional list of numbers defining number of points                          |
// +-------------+------------------------------------------------------------------------------+
//
// Note: The grid definition template varies based on the template number.
// Common templates include lat/lon grids, Lambert conformal, and polar stereographic.
type Section3 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Grid definition
	GridDefinitionSource() uint8
	NumberOfDataPoints() uint32
	GridDefinitionTemplateNumber() uint8

	// Optional list information
	OptionalListOctets() uint32
	OptionalListInterpretation() uint8
	OptionalList() []uint32
}

// Section4 represents the GRIB2 Product Definition Section (Section 4)
// This section defines the meteorological product being encoded,
// including parameter information, time ranges, and vertical coordinates.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | Length of the section in octets (nn)                                        |
// | 5           | Number of the section (4)                                                   |
// | 6-7         | Number of coordinate values after template                                  |
// | 8-9         | Product definition template number (See Table 4.0)                          |
// | 10-xx       | Product definition template (See product template 4.X)                      |
// | [xx+1]-nn   | Optional list of coordinate values                                           |
// +-------------+------------------------------------------------------------------------------+
//
// Note: Coordinate values are used for hybrid coordinate vertical levels.
// They are encoded as IEEE 32-bit floating point pairs when present.
type Section4 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Product definition
	NumberOfCoordinateValues() uint32
	ProductDefinitionTemplateNumber() uint8

	// Optional coordinate values
	CoordinateValues() []float32
}

// Section5 represents the GRIB2 Data Representation Section (Section 5)
// This section describes how the data values are represented and packed,
// including compression methods and scaling factors.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | Length of the section in octets (nn)                                        |
// | 5           | Number of the section (5)                                                   |
// | 6-9         | Number of data points where one or more values are specified in Section 7   |
// | 10-11       | Data representation template number (See Table 5.0)                         |
// | 12-nn       | Data representation template (See Template 5.X)                             |
// +-------------+------------------------------------------------------------------------------+
//
// Note: The data representation template defines packing methods like
// simple packing, complex packing, or IEEE floating point representation.
type Section5 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Data representation
	NumberOfDataPoints() uint32
	DataRepresentationTemplateNumber() uint8
}

// Section6 represents the GRIB2 Bit-map Section (Section 6)
// This section defines which grid points contain valid data values.
// It is optional and used when some grid points have missing data.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | Length of the section in octets (nn)                                        |
// | 5           | Number of the section (6)                                                   |
// | 6           | Bit-map indicator (See Table 6.0)                                           |
// | 7-nn        | Bit-map                                                                     |
// +-------------+------------------------------------------------------------------------------+
//
// Note: If octet 6 is not zero, the length of this section is 6 and
// octets 7-nn are not present (no bit-map data).
type Section6 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Bit-map information
	BitMapIndicator() uint8
	BitMap() []byte
	HasBitMap() bool
}

// Section7 represents the GRIB2 Data Section (Section 7)
// This section contains the actual meteorological data values,
// packed according to the method specified in Section 5.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | Length of the section in octets (nn)                                        |
// | 5           | Number of the section (7)                                                   |
// | 6-nn        | Data in a format described by data template 5.X                             |
// +-------------+------------------------------------------------------------------------------+
//
// Note: The actual data format depends on the data representation template
// specified in Section 5. Data may be packed using various compression methods.
// For large data sections, consider using DataReader() for memory-efficient streaming.
//
// Concurrency Safety: All methods are safe for concurrent use.
type Section7 interface {
	// Section information
	Length() uint32
	SectionNumber() uint8

	// Data access
	Data() []byte          // Returns all data (may consume large memory for big sections)
	DataReader() io.Reader // Returns a reader for streaming data access
	DataSize() uint32      // Returns the size of data payload in bytes

	// Error handling
	LoadError() error // Returns any error encountered during data loading
}

// Section8 represents the GRIB2 End Section (Section 8)
// This section marks the end of the GRIB message with a fixed 4-character string.
// It is always exactly 4 octets long.
//
// Format:
// +-------------+------------------------------------------------------------------------------+
// | Octet Number| Content                                                                      |
// +-------------+------------------------------------------------------------------------------+
// | 1-4         | "7777" - Coded according to the International Alphabet Number 5            |
// +-------------+------------------------------------------------------------------------------+
//
// Note: This section serves as a definitive end marker for the GRIB2 message
// and helps validate that the message was read completely.
type Section8 interface {
	// End section verification
	EndMarker() [4]byte
	IsValid() bool
}
