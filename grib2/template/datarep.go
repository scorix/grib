package template

// DataRepTemplate contains data representation template specific fields
type DataRepTemplate struct {
	TemplateNumber            int     // Data representation template number
	ReferenceValue            float64 // Reference value (R)
	BinaryScaleFactor         int16   // Binary scale factor (E)
	DecimalScaleFactor        int16   // Decimal scale factor (D)
	NumberOfBitsUsedForData   uint8   // Number of bits used for each packed value
	TypeOfOriginalFieldValues uint8   // Type of original field values

	// Template-specific fields
	Simple    *SimplePackingInfo    // For template 0: Simple packing
	Complex   *ComplexPackingInfo   // For templates 2, 3: Complex packing
	IEEE      *IEEEPackingInfo      // For template 4: IEEE floating point
	RunLength *RunLengthPackingInfo // For template 200: Run length packing
	PNG       *PNGPackingInfo       // For template 41: PNG packing
	JPEG2000  *JPEG2000PackingInfo  // For template 40: JPEG2000 packing
	CCSDS     *CCSDSPackingInfo     // For template 42: CCSDS recommended lossless compression
}

// SimplePackingInfo contains simple packing specific fields (template 0)
type SimplePackingInfo struct {
	// Simple packing uses only the common fields
	// No additional template-specific fields
}

// ComplexPackingInfo contains complex packing specific fields (templates 2, 3)
type ComplexPackingInfo struct {
	GroupSplittingMethod            int     // Group splitting method
	MissingValueManagement          uint8   // Missing value management
	PrimaryMissingValueSubstitute   float32 // Primary missing value substitute
	SecondaryMissingValueSubstitute float32 // Secondary missing value substitute
	NumberOfGroupsOfDataValues      uint32  // Number of groups of data values
	ReferenceForGroupWidths         uint8   // Reference for group widths
	NumberOfBitsUsedForGroupWidths  uint8   // Number of bits used for group widths
	ReferenceForGroupLengths        uint32  // Reference for group lengths
	LengthIncrementForGroupLengths  uint8   // Length increment for group lengths
	TrueLengthOfLastGroup           uint32  // True length of last group
	NumberOfBitsUsedForGroupLengths uint8   // Number of bits used for group lengths

	// For template 3 only (complex packing and spatial differencing)
	OrderOfSpatialDifferencing     *uint8 // Order of spatial differencing (1 or 2)
	NumberOfOctetsExtraDescriptors *uint8 // Number of octets required in the data section to specify extra descriptors
}

// IEEEPackingInfo contains IEEE floating point packing specific fields (template 4)
type IEEEPackingInfo struct {
	PrecisionOfFloatingPointNumbers uint8 // Precision of floating point numbers (1=32-bit, 2=64-bit)
}

// RunLengthPackingInfo contains run length packing specific fields (template 200)
type RunLengthPackingInfo struct {
	LevelValues                      []uint8 // Level values
	NumberOfLevels                   uint8   // Number of levels
	MissingValueManagement           uint8   // Missing value management
	PrimaryMissingValueSubstitute    uint8   // Primary missing value substitute
	SecondaryMissingValueSubstitute  uint8   // Secondary missing value substitute
	NumberOfBitsForLevelValues       uint8   // Number of bits for level values
	NumberOfBitsForRunLengths        uint8   // Number of bits for run lengths
	MaximumValueOfLevelValues        uint8   // Maximum value of level values
	MaximumNumberOfBitsForRunLengths uint8   // Maximum number of bits for run lengths
}

// PNGPackingInfo contains PNG packing specific fields (template 41)
type PNGPackingInfo struct {
	// PNG packing uses the common fields plus PNG-specific encoding
	// The actual PNG data is in the data section
}

// JPEG2000PackingInfo contains JPEG2000 packing specific fields (template 40)
type JPEG2000PackingInfo struct {
	TargetCompressionRatio uint8 // Target compression ratio M:1
	CompressionType        uint8 // Type of compression
	CompressionRatio       uint8 // Compression ratio

	// JPEG2000 specific parameters
	NumberOfWaveletLevels uint8  // Number of wavelet levels
	ProgressionOrder      uint8  // Progression order
	NumberOfTilesTotalX   uint16 // Total number of tiles in x direction
	NumberOfTilesTotalY   uint16 // Total number of tiles in y direction
	TileSizeX             uint16 // Tile size in x direction
	TileSizeY             uint16 // Tile size in y direction
}

// CCSDSPackingInfo contains CCSDS recommended lossless compression specific fields (template 42)
type CCSDSPackingInfo struct {
	CCSDSFlags uint8 // CCSDS flags
	BlockSize  uint8 // Block size
	RSILength  uint8 // Reference sample interval length
	Flags      uint8 // Additional flags

	// CCSDS specific parameters
	CompressionOption    uint8   // Compression option
	ReferenceValues      []int32 // Reference values
	NumberOfBitsPerPixel uint8   // Number of bits per pixel

	// Advanced CCSDS parameters
	PreprocessorFlag         uint8  // Preprocessor flag
	DynamicRange             uint16 // Dynamic range
	SampleType               uint8  // Sample type
	NumberOfResolutionLevels uint8  // Number of resolution levels
}

// GetCommonInfo returns common data representation information regardless of template type
func (dr *DataRepTemplate) GetCommonInfo() DataRepCommonInfo {
	return DataRepCommonInfo{
		TemplateNumber:            dr.TemplateNumber,
		ReferenceValue:            dr.ReferenceValue,
		BinaryScaleFactor:         dr.BinaryScaleFactor,
		DecimalScaleFactor:        dr.DecimalScaleFactor,
		NumberOfBitsUsedForData:   dr.NumberOfBitsUsedForData,
		TypeOfOriginalFieldValues: dr.TypeOfOriginalFieldValues,
	}
}

// DataRepCommonInfo contains the common fields across all data representation templates
type DataRepCommonInfo struct {
	TemplateNumber            int     // Data representation template number
	ReferenceValue            float64 // Reference value (R)
	BinaryScaleFactor         int16   // Binary scale factor (E)
	DecimalScaleFactor        int16   // Decimal scale factor (D)
	NumberOfBitsUsedForData   uint8   // Number of bits used for each packed value
	TypeOfOriginalFieldValues uint8   // Type of original field values
}

// IsLossyCompression returns true if the template uses lossy compression
func (dr *DataRepTemplate) IsLossyCompression() bool {
	switch dr.TemplateNumber {
	case 40: // JPEG2000
		return dr.JPEG2000 != nil && dr.JPEG2000.CompressionType == 0 // 0 indicates lossy
	case 41: // PNG (always lossless)
		return false
	case 42: // CCSDS (always lossless)
		return false
	default:
		return false
	}
}

// HasMissingValues returns true if the template supports missing value management
func (dr *DataRepTemplate) HasMissingValues() bool {
	switch dr.TemplateNumber {
	case 2, 3: // Complex packing
		return dr.Complex != nil && dr.Complex.MissingValueManagement > 0
	case 200: // Run length packing
		return dr.RunLength != nil && dr.RunLength.MissingValueManagement > 0
	default:
		return false
	}
}

// GetMissingValueSubstitutes returns the missing value substitutes if available
func (dr *DataRepTemplate) GetMissingValueSubstitutes() (primary float32, secondary float32, hasValues bool) {
	switch dr.TemplateNumber {
	case 2, 3: // Complex packing
		if dr.Complex != nil && dr.Complex.MissingValueManagement > 0 {
			return dr.Complex.PrimaryMissingValueSubstitute,
				dr.Complex.SecondaryMissingValueSubstitute,
				true
		}
	case 200: // Run length packing
		if dr.RunLength != nil && dr.RunLength.MissingValueManagement > 0 {
			return float32(dr.RunLength.PrimaryMissingValueSubstitute),
				float32(dr.RunLength.SecondaryMissingValueSubstitute),
				true
		}
	}
	return 0, 0, false
}
