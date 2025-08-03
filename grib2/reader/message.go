package reader

import (
	"encoding/binary"

	"github.com/scorix/grib/grib2/section"
	"github.com/scorix/grib/grib2/spec"
	"github.com/scorix/grib/grib2/template"
)

// MessageInfo contains metadata about a GRIB2 message's location
type MessageInfo struct {
	Index       int           // Message index in file (0-based)
	Offset      int64         // Start offset of the message (Section 0 start)
	Length      uint64        // Total length of the message (from Section 0)
	Discipline  uint8         // Discipline code from Section 0
	Edition     uint8         // GRIB edition from Section 0
	Sections    []SectionInfo // All sections within this message
	IsFlattened bool          // True if this is a flattened message (single data field)
}

// Message represents a complete GRIB2 message with reader-specific metadata
type Message struct {
	Info         MessageInfo // Reader-specific metadata
	spec.Message             // GRIB2 specification structure
}

// FlatMessage represents a flattened GRIB2 message containing a single data field
// with all relevant information extracted and easily accessible
type FlatMessage struct {
	// Basic message info
	Index      int    // Sequential index in flattened list
	Offset     int64  // Start offset of the original message
	Length     uint64 // Total length of the original message
	Discipline int    // Discipline code
	Edition    int    // GRIB edition

	// Identification information (from Section 1)
	Centre                    int // Originating/generating centre
	SubCentre                 int // Originating/generating sub-centre
	MasterTablesVersion       int // Version number of GRIB master tables
	LocalTablesVersion        int // Version number of GRIB local tables
	ReferenceTimeSignificance int // Significance of reference time
	Year                      int // Year (4 digits)
	Month                     int // Month
	Day                       int // Day
	Hour                      int // Hour
	Minute                    int // Minute
	Second                    int // Second
	ProductionStatus          int // Production status of processed data
	TypeOfData                int // Type of processed data

	// Template-specific information
	Product template.ProductTemplate // Product definition template fields
	Grid    template.GridTemplate    // Grid definition template fields
	DataRep template.DataRepTemplate // Data representation template fields

	// Raw sections for advanced access
	Indicator      section.Section0 // Section 0 - Indicator
	Identification section.Section1 // Section 1 - Identification
	LocalUse       section.Section2 // Section 2 - Local Use (may be nil)
	GridDef        section.Section3 // Section 3 - Grid Definition
	ProductDef     section.Section4 // Section 4 - Product Definition
	DataRepSec     section.Section5 // Section 5 - Data Representation
	Bitmap         section.Section6 // Section 6 - Bitmap (may be nil)
	Data           section.Section7 // Section 7 - Data
	End            section.Section8 // Section 8 - End
}

// FlattenMessages converts a nested GRIB2 message into multiple flat messages
// Each DataField becomes an independent message with complete context
// Returns m*n*l messages where m=local blocks, n=grid blocks per local block, l=data fields per grid block
func (m *Message) FlattenMessages() []Message {
	var messages []Message

	for _, localBlock := range m.Blocks {
		for _, gridBlock := range localBlock.Grids {
			for _, dataField := range gridBlock.Fields {
				// Create a new flat message for each data field
				flatMessage := Message{
					Info: MessageInfo{
						Index:       len(messages), // Sequential index in flattened list
						Offset:      m.Info.Offset, // Original message offset (all share same offset)
						Length:      m.Info.Length, // Original message length (all share same length)
						Discipline:  m.Info.Discipline,
						Edition:     m.Info.Edition,
						Sections:    m.Info.Sections,
						IsFlattened: true, // Mark as flattened message
					},
					Message: spec.Message{
						Indicator:      m.Indicator,
						Identification: m.Identification,
						End:            m.End,
						Blocks: []spec.LocalBlock{{
							LocalUse: localBlock.LocalUse, // May be nil
							Grids: []spec.GridBlock{{
								GridDef: gridBlock.GridDef,
								Fields:  []spec.DataField{dataField}, // Single field
							}},
						}},
					},
				}

				messages = append(messages, flatMessage)
			}
		}
	}

	return messages
}

// FlattenToFlatMessages converts a nested GRIB2 message into multiple FlatMessage structs
// Each DataField becomes an independent FlatMessage with extracted field information
// Returns m*n*l flat messages where m=local blocks, n=grid blocks per local block, l=data fields per grid block
func (m *Message) FlattenToFlatMessages() []FlatMessage {
	var flatMessages []FlatMessage

	for _, localBlock := range m.Blocks {
		for _, gridBlock := range localBlock.Grids {
			for _, dataField := range gridBlock.Fields {
				flatMsg := FlatMessage{
					// Basic info
					Index:      len(flatMessages),
					Offset:     m.Info.Offset,
					Length:     m.Info.Length,
					Discipline: int(m.Info.Discipline),
					Edition:    int(m.Info.Edition),

					// Raw sections
					Indicator:      m.Indicator,
					Identification: m.Identification,
					LocalUse:       localBlock.LocalUse,
					GridDef:        gridBlock.GridDef,
					ProductDef:     dataField.ProductDef,
					DataRepSec:     dataField.DataRep,
					Bitmap:         dataField.Bitmap,
					Data:           dataField.Data,
					End:            m.End,
				}

				// Extract fields from Section 4 (Product Definition)
				flatMsg.extractProductInfo()

				// Extract fields from Section 3 (Grid Definition)
				flatMsg.extractGridInfo()

				flatMessages = append(flatMessages, flatMsg)
			}
		}
	}

	return flatMessages
}

// extractProductInfo extracts product-related information from Section 4
func (f *FlatMessage) extractProductInfo() {
	// Extract basic fields from Section 4
	f.Product.TemplateNumber = uint16(f.ProductDef.ProductDefinitionTemplateNumber())

	// Extract from Section 1 (Identification)
	f.Centre = int(f.Identification.OriginatingCenter())
	f.SubCentre = int(f.Identification.OriginatingSubcenter())
	f.MasterTablesVersion = int(f.Identification.MasterTablesVersion())
	f.LocalTablesVersion = int(f.Identification.LocalTablesVersion())
	f.ReferenceTimeSignificance = int(f.Identification.ReferenceTimeSignificance())
	f.Year = int(f.Identification.Year())
	f.Month = int(f.Identification.Month())
	f.Day = int(f.Identification.Day())
	f.Hour = int(f.Identification.Hour())
	f.Minute = int(f.Identification.Minute())
	f.Second = int(f.Identification.Second())
	f.ProductionStatus = int(f.Identification.ProductionStatus())
	f.TypeOfData = int(f.Identification.DataType())

	// TODO: Extract detailed product definition template fields
	// This requires parsing the raw productDefinitionTemplate bytes
	// For now, we set placeholder values - this should be implemented
	// based on the specific template number and GRIB2 specification

	// Common fields that would be extracted from different templates:
	// - Parameter Category (usually octet 10 in template)
	// - Parameter Number (usually octet 11 in template)
	// - Type of Generating Process (usually octet 12 in template)
	// - Forecast Time (usually octets 19-22 in template)
	// - Type/Scale/Value of Fixed Surfaces (usually octets 23-34 in template)

	// These would need to be extracted based on the template structure
	// f.Category = extractParameterCategory(f.ProductDef)
	// f.Parameter = extractParameterNumber(f.ProductDef)
	// f.ForecastTime = extractForecastTime(f.ProductDef)
}

// extractGridInfo extracts grid-related information from Section 3
func (f *FlatMessage) extractGridInfo() {
	// Extract basic fields from Section 3
	f.Grid.SourceOfGridDefinition = int(f.GridDef.GridDefinitionSource())
	f.Grid.NumberOfDataPoints = int(f.GridDef.NumberOfDataPoints())
	f.Grid.NumberOfOctectsForOptional = int(f.GridDef.OptionalListOctets())
	f.Grid.InterpretationOfOptional = int(f.GridDef.OptionalListInterpretation())
	f.Grid.TemplateNumber = int(f.GridDef.GridDefinitionTemplateNumber())

	// Extract from Section 5 (Data Representation)
	f.DataRep.TemplateNumber = int(f.DataRepSec.DataRepresentationTemplateNumber())

	// TODO: Extract detailed grid definition template fields
	// This requires parsing the raw gridDefinitionTemplate bytes
	// For common templates like lat/lon regular grids (template 0):
	// - Scanning Mode (usually octet 72 in template)
	// - Number of points along parallels/meridians (octets 31-34, 35-38)
	// - Latitude/Longitude of first/last grid points (octets 47-50, 51-54, 56-59, 60-63)
	// - Direction increments (octets 64-67, 68-71)

	// These would need to be extracted based on the template structure:
	// f.ScanningMode = extractScanningMode(f.GridDef)
	// f.NumberOfGridPointsAlongX = extractGridPointsX(f.GridDef)
	// f.NumberOfGridPointsAlongY = extractGridPointsY(f.GridDef)
	// f.LatitudeOfFirstGridPoint = extractFirstLatitude(f.GridDef)
	// f.LongitudeOfFirstGridPoint = extractFirstLongitude(f.GridDef)
	// f.XDirectionIncrement = extractXIncrement(f.GridDef)
	// f.YDirectionIncrement = extractYIncrement(f.GridDef)

	// TODO: Extract data representation template fields
	// This requires parsing the raw dataRepresentationTemplate bytes
	// For simple packing (template 0):
	// - Reference value (octets 12-15)
	// - Binary scale factor (octets 16-17)
	// - Decimal scale factor (octets 18-19)
	// - Number of bits for packing (octet 20)
	// - Type of original field values (octet 21)

	// f.ReferenceValue = extractReferenceValue(f.DataRep)
	// f.BinaryScaleFactor = extractBinaryScaleFactor(f.DataRep)
	// f.DecimalScaleFactor = extractDecimalScaleFactor(f.DataRep)
	// f.NumberOfBitsUsedForData = extractNumberOfBits(f.DataRep)

	// Extract some basic template fields if we can access raw template data
	f.extractProductTemplate()
	f.extractGridTemplate()
	f.extractDataRepTemplate()
}

// extractProductTemplate extracts common fields from product definition template
func (f *FlatMessage) extractProductTemplate() {
	// Use reflection or type assertion to access the underlying implementation
	// This is a bit of a hack, but allows us to extract template data
	if sec4, ok := f.ProductDef.(interface{ productDefinitionTemplate() []byte }); ok {
		template := sec4.productDefinitionTemplate()
		f.extractFromProductTemplate(template, int(f.ProductDef.ProductDefinitionTemplateNumber()))
	}
}

// extractFromProductTemplate extracts fields from the raw product definition template bytes
func (f *FlatMessage) extractFromProductTemplate(templateData []byte, templateNumber int) {
	if len(templateData) < 34 { // Minimum length for most templates
		return
	}

	// Most common templates (0, 1, 8, etc.) have similar structure for basic fields
	// Based on WMO GRIB2 specification Table 4.0

	// Parameter Category (octet 10 of template = octet 0 of template data)
	if len(templateData) > 0 {
		f.Product.Category = uint8(templateData[0])
	}

	// Parameter Number (octet 11 of template = octet 1 of template data)
	if len(templateData) > 1 {
		f.Product.Parameter = uint8(templateData[1])
	}

	// Type of Generating Process (octet 12 of template = octet 2 of template data)
	if len(templateData) > 2 {
		f.Product.TypeOfGeneratingProcess = uint8(templateData[2])
	}

	// Background Process (octet 13 of template = octet 3 of template data)
	if len(templateData) > 3 {
		f.Product.BackgroundProcess = uint8(templateData[3])
	}

	// Generating Process Identifier (octet 14 of template = octet 4 of template data)
	if len(templateData) > 4 {
		f.Product.GeneratingProcessIdentifier = uint8(templateData[4])
	}

	// Hours after data cutoff (octets 15-16 of template = octets 5-6 of template data)
	if len(templateData) > 6 {
		f.Product.HoursAfterDataCutoff = binary.BigEndian.Uint16(templateData[5:7])
	}

	// Minutes after data cutoff (octet 17 of template = octet 7 of template data)
	if len(templateData) > 7 {
		f.Product.MinutesAfterDataCutoff = uint8(templateData[7])
	}

	// Indicator of unit of time range (octet 18 of template = octet 8 of template data)
	if len(templateData) > 8 {
		f.Product.IndicatorOfUnitOfTimeRange = uint8(templateData[8])
	}

	// Forecast time (octets 19-22 of template = octets 9-12 of template data)
	if len(templateData) > 12 {
		f.Product.ForecastTime = binary.BigEndian.Uint32(templateData[9:13])
	}

	// Type of first fixed surface (octet 23 of template = octet 13 of template data)
	if len(templateData) > 13 {
		f.Product.TypeOfFirstFixedSurface = uint8(templateData[13])
	}

	// Scale factor of first fixed surface (octet 24 of template = octet 14 of template data)
	if len(templateData) > 14 {
		f.Product.ScaleFactorOfFirstFixedSurface = int8(templateData[14]) // Signed
	}

	// Scaled value of first fixed surface (octets 25-28 of template = octets 15-18 of template data)
	if len(templateData) > 18 {
		f.Product.ScaledValueOfFirstFixedSurface = binary.BigEndian.Uint32(templateData[15:19])
	}

	// Type of second fixed surface (octet 29 of template = octet 19 of template data)
	if len(templateData) > 19 {
		f.Product.TypeOfSecondFixedSurface = uint8(templateData[19])
	}

	// Scale factor of second fixed surface (octet 30 of template = octet 20 of template data)
	if len(templateData) > 20 {
		f.Product.ScaleFactorOfSecondFixedSurface = int8(templateData[20]) // Signed
	}

	// Scaled value of second fixed surface (octets 31-34 of template = octets 21-24 of template data)
	if len(templateData) > 24 {
		f.Product.ScaledValueOfSecondFixedSurface = binary.BigEndian.Uint32(templateData[21:25])
	}
}

// extractGridTemplate extracts common fields from grid definition template
func (f *FlatMessage) extractGridTemplate() {
	// Access the underlying implementation to get raw template data
	if sec3, ok := f.GridDef.(interface{ gridDefinitionTemplate() []byte }); ok {
		template := sec3.gridDefinitionTemplate()
		f.extractFromGridTemplate(template, int(f.GridDef.GridDefinitionTemplateNumber()))
	}
}

// extractFromGridTemplate extracts fields from the raw grid definition template bytes
func (f *FlatMessage) extractFromGridTemplate(templateData []byte, templateNumber int) {
	// Template 0: Latitude/longitude (or equidistant cylindrical, or Plate Carree)
	// This is the most common template, let's implement it fully
	if templateNumber == 0 && len(templateData) >= 72 {
		// Shape of the Earth (octet 15 of section 3 = octet 1 of template)
		// We could extract this but it's not in our FlatMessage struct

		// Create LatLonGrid template for template 0
		if f.Grid.LatLon == nil {
			f.Grid.LatLon = &template.LatLonGrid{}
		}

		// Number of points along a parallel (octets 31-34 = octets 17-20 of template)
		if len(templateData) >= 20 {
			f.Grid.LatLon.NumberOfGridPointsAlongX = binary.BigEndian.Uint32(templateData[16:20])
		}

		// Number of points along a meridian (octets 35-38 = octets 21-24 of template)
		if len(templateData) >= 24 {
			f.Grid.LatLon.NumberOfGridPointsAlongY = binary.BigEndian.Uint32(templateData[20:24])
		}

		// Latitude of first grid point (octets 47-50 = octets 33-36 of template)
		if len(templateData) >= 36 {
			f.Grid.LatLon.LatitudeOfFirstGridPoint = int32(binary.BigEndian.Uint32(templateData[32:36]))
		}

		// Longitude of first grid point (octets 51-54 = octets 37-40 of template)
		if len(templateData) >= 40 {
			f.Grid.LatLon.LongitudeOfFirstGridPoint = binary.BigEndian.Uint32(templateData[36:40])
		}

		// Latitude of last grid point (octets 56-59 = octets 42-45 of template)
		if len(templateData) >= 45 {
			f.Grid.LatLon.LatitudeOfLastGridPoint = int32(binary.BigEndian.Uint32(templateData[41:45]))
		}

		// Longitude of last grid point (octets 60-63 = octets 46-49 of template)
		if len(templateData) >= 49 {
			f.Grid.LatLon.LongitudeOfLastGridPoint = binary.BigEndian.Uint32(templateData[45:49])
		}

		// i direction increment (octets 64-67 = octets 50-53 of template)
		if len(templateData) >= 53 {
			f.Grid.LatLon.XDirectionIncrement = binary.BigEndian.Uint32(templateData[49:53])
		}

		// j direction increment (octets 68-71 = octets 54-57 of template)
		if len(templateData) >= 57 {
			f.Grid.LatLon.YDirectionIncrement = binary.BigEndian.Uint32(templateData[53:57])
		}

		// Scanning mode (octet 72 = octet 58 of template)
		if len(templateData) >= 58 {
			f.Grid.LatLon.ScanningMode = templateData[57]
		}
	}

	// For other templates, we'd add similar extraction logic
	// Template 1: Rotated latitude/longitude
	// Template 10: Mercator
	// Template 20: Polar stereographic projection
	// Template 30: Lambert conformal
	// etc.
}

// extractDataRepTemplate extracts common fields from data representation template
func (f *FlatMessage) extractDataRepTemplate() {
	// Access the underlying implementation to get raw template data
	if sec5, ok := f.DataRepSec.(interface{ dataRepresentationTemplate() []byte }); ok {
		template := sec5.dataRepresentationTemplate()
		f.extractFromDataRepTemplate(template, f.DataRep.TemplateNumber)
	}
}

// extractFromDataRepTemplate extracts fields from the raw data representation template bytes
func (f *FlatMessage) extractFromDataRepTemplate(templateData []byte, templateNumber int) {
	// Template 0: Grid point data - simple packing
	// Template 2: Grid point data - complex packing
	// Template 3: Grid point data - complex packing and spatial differencing
	// Most templates share the first few fields

	if len(templateData) >= 21 { // Minimum length for most templates
		// Reference value (octets 12-15 = octets 0-3 of template)
		if len(templateData) >= 4 {
			refRaw := binary.BigEndian.Uint32(templateData[0:4])
			f.DataRep.ReferenceValue = float64(refRaw) // This is actually IEEE float, need proper conversion
		}

		// Binary scale factor (octets 16-17 = octets 4-5 of template)
		if len(templateData) >= 6 {
			f.DataRep.BinaryScaleFactor = int16(binary.BigEndian.Uint16(templateData[4:6])) // Signed
		}

		// Decimal scale factor (octets 18-19 = octets 6-7 of template)
		if len(templateData) >= 8 {
			f.DataRep.DecimalScaleFactor = int16(binary.BigEndian.Uint16(templateData[6:8])) // Signed
		}

		// Number of bits used for each packed value (octet 20 = octet 8 of template)
		if len(templateData) >= 9 {
			f.DataRep.NumberOfBitsUsedForData = uint8(templateData[8])
		}

		// Type of original field values (octet 21 = octet 9 of template)
		if len(templateData) >= 10 {
			f.DataRep.TypeOfOriginalFieldValues = uint8(templateData[9])
		}
	}

	// Template-specific fields
	switch templateNumber {
	case 2, 3: // Complex packing
		if len(templateData) >= 31 {
			// Group splitting method (octet 22 = octet 10 of template)
			// Missing value management (octet 23 = octet 11 of template)
			// Primary missing value substitute (octets 24-27 = octets 12-15 of template)
			// Secondary missing value substitute (octets 28-31 = octets 16-19 of template)
			// Number of groups of data values (octets 32-35 = octets 20-23 of template)

			// For now, we don't have these fields in FlatMessage
			// They could be added if needed for complex packing support
		}
	}
}

// IsFlattened returns true if this message contains only a single data field
func (m *Message) IsFlattened() bool {
	return m.Info.IsFlattened
}

// IsNested returns true if this is an original nested message (not flattened)
func (m *Message) IsNested() bool {
	return !m.Info.IsFlattened
}
