package template

// ProductTemplate contains product definition template specific fields
type ProductTemplate struct {
	TemplateNumber              uint16 // Product definition template number (2 bytes)
	Category                    uint8  // Parameter category (1 byte)
	Parameter                   uint8  // Parameter number (1 byte)
	TypeOfGeneratingProcess     uint8  // Type of generating process (1 byte)
	BackgroundProcess           uint8  // Background generating process identifier (1 byte)
	GeneratingProcessIdentifier uint8  // Generating process or model identifier (1 byte)
	HoursAfterDataCutoff        uint16 // Hours after reference time of data cutoff (2 bytes)
	MinutesAfterDataCutoff      uint8  // Minutes after reference time of data cutoff (1 byte)
	IndicatorOfUnitOfTimeRange  uint8  // Indicator of unit of time range (1 byte)
	ForecastTime                uint32 // Forecast time in units defined by previous octet (4 bytes)

	// Fixed surface information
	TypeOfFirstFixedSurface         uint8  // Type of first fixed surface (1 byte)
	ScaleFactorOfFirstFixedSurface  int8   // Scale factor of first fixed surface (1 byte, signed)
	ScaledValueOfFirstFixedSurface  uint32 // Scaled value of first fixed surface (4 bytes)
	TypeOfSecondFixedSurface        uint8  // Type of second fixed surface (1 byte)
	ScaleFactorOfSecondFixedSurface int8   // Scale factor of second fixed surface (1 byte, signed)
	ScaledValueOfSecondFixedSurface uint32 // Scaled value of second fixed surface (4 bytes)

	// Template-specific fields (populated based on template number)
	TimeRange   *TimeRangeInfo   // For templates with time ranges (8, 9, 10, 11, 12, 13, 14)
	Ensemble    *EnsembleInfo    // For ensemble templates (1, 2, 3, 4, 11, 12, 13, 14)
	Probability *ProbabilityInfo // For probability templates (5, 9)
	Percentile  *PercentileInfo  // For percentile templates (6, 10)
	Derived     *DerivedInfo     // For derived templates (7, 12, 13, 14)
}

// TimeRangeInfo contains time range specific information
type TimeRangeInfo struct {
	TypeOfTimeIncrement             uint8  // Type of time increment (1 byte)
	IndicatorOfUnitForTimeRange     uint8  // Indicator of unit for time range (1 byte)
	LengthOfTimeRange               uint32 // Length of time range (4 bytes)
	IndicatorOfUnitForTimeIncrement uint8  // Indicator of unit for time increment (1 byte)
	TimeIncrement                   uint32 // Time increment (4 bytes)
	TypeOfStatisticalProcessing     uint8  // Type of statistical processing (1 byte)
	NumberOfTimeRanges              uint16 // Number of time ranges (2 bytes)

	// For multiple time ranges
	TimeRanges []TimeRangeSpec // List of time range specifications
}

// TimeRangeSpec represents a single time range specification
type TimeRangeSpec struct {
	StatisticalProcessType uint8  // Type of statistical processing (1 byte)
	TimeIncrementType      uint8  // Type of time increment (1 byte)
	TimeRangeLength        uint32 // Length of time range (4 bytes)
	TimeIncrement          uint32 // Time increment (4 bytes)
}

// EnsembleInfo contains ensemble forecast specific information
type EnsembleInfo struct {
	TypeOfEnsembleForecast      uint8 // Type of ensemble forecast (1 byte)
	PerturbationNumber          uint8 // Perturbation number (1 byte)
	NumberOfForecastsInEnsemble uint8 // Number of forecasts in ensemble (1 byte)
}

// ProbabilityInfo contains probability forecast specific information
type ProbabilityInfo struct {
	ForecastProbabilityNumber          uint8   // Forecast probability number (1 byte)
	TotalNumberOfForecastProbabilities uint8   // Total number of forecast probabilities (1 byte)
	ProbabilityType                    uint8   // Probability type (1 byte)
	ScaleFactorOfLowerLimit            int8    // Scale factor of lower limit (1 byte, signed)
	ScaledValueOfLowerLimit            uint32  // Scaled value of lower limit (4 bytes)
	ScaleFactorOfUpperLimit            int8    // Scale factor of upper limit (1 byte, signed)
	ScaledValueOfUpperLimit            uint32  // Scaled value of upper limit (4 bytes)
	LowerLimitValue                    float64 // Calculated lower limit value
	UpperLimitValue                    float64 // Calculated upper limit value
}

// PercentileInfo contains percentile forecast specific information
type PercentileInfo struct {
	PercentileValue uint8 // Percentile value (0-100, 1 byte)
}

// DerivedInfo contains derived forecast specific information
type DerivedInfo struct {
	DerivedForecastType            uint8  // Type of derived forecast (1 byte)
	NumberOfForecastsInEnsemble    uint8  // Number of forecasts used to create derived forecast (1 byte)
	ClusterIdentifier              uint8  // Cluster identifier (1 byte)
	NumberOfClustersOfEnsemble     uint8  // Number of clusters (1 byte)
	ClusteringMethod               uint8  // Clustering method (1 byte)
	NorthernLatitudeOfCluster      int32  // Northern latitude of cluster domain (4 bytes, signed, microdegrees)
	SouthernLatitudeOfCluster      int32  // Southern latitude of cluster domain (4 bytes, signed, microdegrees)
	EasternLongitudeOfCluster      uint32 // Eastern longitude of cluster domain (4 bytes, microdegrees)
	WesternLongitudeOfCluster      uint32 // Western longitude of cluster domain (4 bytes, microdegrees)
	NumberOfForecastsInCluster     uint8  // Number of forecasts in cluster (1 byte)
	ScaleFactorOfCentralWaveNumber int8   // Scale factor of central wave number (1 byte, signed)
	ScaledValueOfCentralWaveNumber uint32 // Scaled value of central wave number (4 bytes)
}
