package template

// GridTemplate contains grid definition template specific fields
type GridTemplate struct {
	TemplateNumber             int // Grid definition template number
	SourceOfGridDefinition     int // Source of grid definition
	NumberOfDataPoints         int // Number of data points
	NumberOfOctectsForOptional int // Number of octets for optional list of numbers
	InterpretationOfOptional   int // Interpretation of list of numbers

	// Template-specific fields (populated based on template number)
	LatLon                            *LatLonGrid                            // For template 0: Latitude/longitude
	RotatedLatLon                     *RotatedLatLonGrid                     // For template 1: Rotated latitude/longitude
	StretchedLatLon                   *StretchedLatLonGrid                   // For template 2: Stretched latitude/longitude
	StretchedRotatedLatLon            *StretchedRotatedLatLonGrid            // For template 3: Stretched and rotated latitude/longitude
	VariableResLatLon                 *VariableResLatLonGrid                 // For template 4: Variable resolution latitude/longitude
	VariableResRotatedLatLon          *VariableResRotatedLatLonGrid          // For template 5: Variable resolution rotated latitude/longitude
	Mercator                          *MercatorGrid                          // For template 10: Mercator
	TransverseMercator                *TransverseMercatorGrid                // For template 12: Transverse Mercator
	PolarStereo                       *PolarStereoGrid                       // For template 20: Polar stereographic
	Lambert                           *LambertGrid                           // For template 30: Lambert conformal
	Albers                            *AlbersGrid                            // For template 31: Albers equal-area
	Gaussian                          *GaussianGrid                          // For template 40: Gaussian latitude/longitude
	RotatedGaussian                   *RotatedGaussianGrid                   // For template 41: Rotated Gaussian latitude/longitude
	StretchedGaussian                 *StretchedGaussianGrid                 // For template 42: Stretched Gaussian latitude/longitude
	StretchedRotatedGaussian          *StretchedRotatedGaussianGrid          // For template 43: Stretched and rotated Gaussian latitude/longitude
	SphericalHarmonic                 *SphericalHarmonicGrid                 // For template 50: Spherical harmonic coefficients
	RotatedSphericalHarmonic          *RotatedSphericalHarmonicGrid          // For template 51: Rotated spherical harmonic coefficients
	StretchedSphericalHarmonic        *StretchedSphericalHarmonicGrid        // For template 52: Stretched spherical harmonic coefficients
	StretchedRotatedSphericalHarmonic *StretchedRotatedSphericalHarmonicGrid // For template 53: Stretched and rotated spherical harmonic coefficients
	SpaceView                         *SpaceViewGrid                         // For template 90: Space view perspective or orthographic
	Triangular                        *TriangularGrid                        // For template 100: Triangular grid based on an icosahedron
	Equatorial                        *EquatorialGrid                        // For template 110: Equatorial azimuthal equidistant projection
	AzimuthRange                      *AzimuthRangeGrid                      // For template 120: Azimuth-range projection
	Irregular                         *IrregularGrid                         // For template 1000: Cross-section grid with points equally spaced on the horizontal
	Curvilinear                       *CurvilinearGrid                       // For template 204: Curvilinear orthogonal grids
}

// LatLonGrid contains latitude/longitude grid specific fields (template 0)
type LatLonGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along a parallel
	NumberOfGridPointsAlongY   uint32 // Number of points along a meridian
	BasicAngleOfInitialDomain  uint32 // Basic angle of the initial production domain
	SubdivisionOfBasicAngle    uint32 // Subdivisions of basic angle used to define extreme longitudes and latitudes
	LatitudeOfFirstGridPoint   int32  // Latitude of first grid point (microdegrees)
	LongitudeOfFirstGridPoint  uint32 // Longitude of first grid point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	LatitudeOfLastGridPoint    int32  // Latitude of last grid point (microdegrees)
	LongitudeOfLastGridPoint   uint32 // Longitude of last grid point (microdegrees)
	XDirectionIncrement        uint32 // i direction increment (microdegrees)
	YDirectionIncrement        uint32 // j direction increment (microdegrees)
	ScanningMode               uint8  // Scanning mode
}

// RotatedLatLonGrid contains rotated latitude/longitude grid specific fields (template 1)
type RotatedLatLonGrid struct {
	LatLonGrid                    // Embedded lat/lon grid
	LatitudeOfSouthernPole  int32 // Latitude of the southern pole of projection (microdegrees)
	LongitudeOfSouthernPole int32 // Longitude of the southern pole of projection (microdegrees)
	AngleOfRotation         int32 // Angle of rotation of projection (microdegrees)
}

// StretchedLatLonGrid contains stretched latitude/longitude grid specific fields (template 2)
type StretchedLatLonGrid struct {
	LatLonGrid                      // Embedded lat/lon grid
	LatitudeOfPoleStretching  int32 // Latitude of the pole of stretching (microdegrees)
	LongitudeOfPoleStretching int32 // Longitude of the pole of stretching (microdegrees)
	StretchingFactor          int32 // Stretching factor
}

// StretchedRotatedLatLonGrid contains stretched and rotated latitude/longitude grid specific fields (template 3)
type StretchedRotatedLatLonGrid struct {
	RotatedLatLonGrid               // Embedded rotated lat/lon grid
	LatitudeOfPoleStretching  int32 // Latitude of the pole of stretching (microdegrees)
	LongitudeOfPoleStretching int32 // Longitude of the pole of stretching (microdegrees)
	StretchingFactor          int32 // Stretching factor
}

// VariableResLatLonGrid contains variable resolution latitude/longitude grid specific fields (template 4)
type VariableResLatLonGrid struct {
	LatLonGrid                // Embedded lat/lon grid
	NumberOfParallels uint16  // Number of parallels between a pole and the equator
	CoordinateValues  []int32 // Coordinate values (n_p × (n_p + 1)) / 2 for triangular arrays
}

// VariableResRotatedLatLonGrid contains variable resolution rotated latitude/longitude grid specific fields (template 5)
type VariableResRotatedLatLonGrid struct {
	VariableResLatLonGrid         // Embedded variable resolution lat/lon grid
	LatitudeOfSouthernPole  int32 // Latitude of the southern pole of projection (microdegrees)
	LongitudeOfSouthernPole int32 // Longitude of the southern pole of projection (microdegrees)
	AngleOfRotation         int32 // Angle of rotation of projection (microdegrees)
}

// MercatorGrid contains Mercator projection grid specific fields (template 10)
type MercatorGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY   uint32 // Number of points along y-axis
	LatitudeOfFirstGridPoint   int32  // Latitude of first grid point (microdegrees)
	LongitudeOfFirstGridPoint  uint32 // Longitude of first grid point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	LatitudeOfLastGridPoint    int32  // Latitude of last grid point (microdegrees)
	LongitudeOfLastGridPoint   uint32 // Longitude of last grid point (microdegrees)
	ScanningMode               uint8  // Scanning mode
	OrientationOfGrid          uint32 // Orientation of the grid (microdegrees)
	XDirectionIncrement        uint32 // X-direction grid length (metres)
	YDirectionIncrement        uint32 // Y-direction grid length (metres)
	LatitudeOfIntersection     int32  // Latitude at which the Mercator projection intersects the Earth (microdegrees)
}

// TransverseMercatorGrid contains transverse Mercator projection grid specific fields (template 12)
type TransverseMercatorGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY   uint32 // Number of points along y-axis
	LatitudeOfFirstGridPoint   int32  // Latitude of first grid point (microdegrees)
	LongitudeOfFirstGridPoint  uint32 // Longitude of first grid point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	LatitudeOfLastGridPoint    int32  // Latitude of last grid point (microdegrees)
	LongitudeOfLastGridPoint   uint32 // Longitude of last grid point (microdegrees)
	ScanningMode               uint8  // Scanning mode
	LatitudeOfOrigin           int32  // Latitude of origin (microdegrees)
	LongitudeOfOrigin          uint32 // Longitude of origin (microdegrees)
	XDirectionIncrement        uint32 // X-direction grid length (metres)
	YDirectionIncrement        uint32 // Y-direction grid length (metres)
	ScaleFactorAtOrigin        uint32 // Scale factor at central meridian
	XOfOrigin                  int32  // X coordinate of origin (grid lengths)
	YOfOrigin                  int32  // Y coordinate of origin (grid lengths)
}

// PolarStereoGrid contains polar stereographic projection grid specific fields (template 20)
type PolarStereoGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY   uint32 // Number of points along y-axis
	LatitudeOfFirstGridPoint   int32  // Latitude of first grid point (microdegrees)
	LongitudeOfFirstGridPoint  uint32 // Longitude of first grid point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	OrientationOfGrid          uint32 // Orientation of the grid (microdegrees)
	XDirectionIncrement        uint32 // X-direction grid length (metres)
	YDirectionIncrement        uint32 // Y-direction grid length (metres)
	ProjectionCenterFlag       uint8  // Projection center flag
	ScanningMode               uint8  // Scanning mode
}

// LambertGrid contains Lambert conformal projection grid specific fields (template 30)
type LambertGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY   uint32 // Number of points along y-axis
	LatitudeOfFirstGridPoint   int32  // Latitude of first grid point (microdegrees)
	LongitudeOfFirstGridPoint  uint32 // Longitude of first grid point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	OrientationOfGrid          uint32 // Orientation of the grid (microdegrees)
	XDirectionIncrement        uint32 // X-direction grid length (metres)
	YDirectionIncrement        uint32 // Y-direction grid length (metres)
	ProjectionCenterFlag       uint8  // Projection center flag
	ScanningMode               uint8  // Scanning mode
	LatitudeOfIntersection1    int32  // Latitude of first standard parallel (microdegrees)
	LatitudeOfIntersection2    int32  // Latitude of second standard parallel (microdegrees)
	LatitudeOfSouthernPole     int32  // Latitude of the southern pole (microdegrees)
	LongitudeOfSouthernPole    uint32 // Longitude of the southern pole (microdegrees)
}

// AlbersGrid contains Albers equal-area projection grid specific fields (template 31)
type AlbersGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY   uint32 // Number of points along y-axis
	LatitudeOfFirstGridPoint   int32  // Latitude of first grid point (microdegrees)
	LongitudeOfFirstGridPoint  uint32 // Longitude of first grid point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	OrientationOfGrid          uint32 // Orientation of the grid (microdegrees)
	XDirectionIncrement        uint32 // X-direction grid length (metres)
	YDirectionIncrement        uint32 // Y-direction grid length (metres)
	ProjectionCenterFlag       uint8  // Projection center flag
	ScanningMode               uint8  // Scanning mode
	LatitudeOfIntersection1    int32  // Latitude of first standard parallel (microdegrees)
	LatitudeOfIntersection2    int32  // Latitude of second standard parallel (microdegrees)
	LatitudeOfSouthernPole     int32  // Latitude of the southern pole (microdegrees)
	LongitudeOfSouthernPole    uint32 // Longitude of the southern pole (microdegrees)
}

// GaussianGrid contains Gaussian latitude/longitude grid specific fields (template 40)
type GaussianGrid struct {
	LatLonGrid               // Embedded lat/lon grid
	NumberOfParallels uint32 // Number of parallels between a pole and the equator
}

// RotatedGaussianGrid contains rotated Gaussian latitude/longitude grid specific fields (template 41)
type RotatedGaussianGrid struct {
	GaussianGrid                  // Embedded Gaussian grid
	LatitudeOfSouthernPole  int32 // Latitude of the southern pole of projection (microdegrees)
	LongitudeOfSouthernPole int32 // Longitude of the southern pole of projection (microdegrees)
	AngleOfRotation         int32 // Angle of rotation of projection (microdegrees)
}

// StretchedGaussianGrid contains stretched Gaussian latitude/longitude grid specific fields (template 42)
type StretchedGaussianGrid struct {
	GaussianGrid                    // Embedded Gaussian grid
	LatitudeOfPoleStretching  int32 // Latitude of the pole of stretching (microdegrees)
	LongitudeOfPoleStretching int32 // Longitude of the pole of stretching (microdegrees)
	StretchingFactor          int32 // Stretching factor
}

// StretchedRotatedGaussianGrid contains stretched and rotated Gaussian latitude/longitude grid specific fields (template 43)
type StretchedRotatedGaussianGrid struct {
	RotatedGaussianGrid             // Embedded rotated Gaussian grid
	LatitudeOfPoleStretching  int32 // Latitude of the pole of stretching (microdegrees)
	LongitudeOfPoleStretching int32 // Longitude of the pole of stretching (microdegrees)
	StretchingFactor          int32 // Stretching factor
}

// SphericalHarmonicGrid contains spherical harmonic coefficients grid specific fields (template 50)
type SphericalHarmonicGrid struct {
	J                  int32 // J - pentagonal resolution parameter
	K                  int32 // K - pentagonal resolution parameter
	M                  int32 // M - pentagonal resolution parameter
	RepresentationType uint8 // Representation type
	RepresentationMode uint8 // Representation mode
}

// RotatedSphericalHarmonicGrid contains rotated spherical harmonic coefficients grid specific fields (template 51)
type RotatedSphericalHarmonicGrid struct {
	SphericalHarmonicGrid         // Embedded spherical harmonic grid
	LatitudeOfSouthernPole  int32 // Latitude of the southern pole of projection (microdegrees)
	LongitudeOfSouthernPole int32 // Longitude of the southern pole of projection (microdegrees)
	AngleOfRotation         int32 // Angle of rotation of projection (microdegrees)
}

// StretchedSphericalHarmonicGrid contains stretched spherical harmonic coefficients grid specific fields (template 52)
type StretchedSphericalHarmonicGrid struct {
	SphericalHarmonicGrid           // Embedded spherical harmonic grid
	LatitudeOfPoleStretching  int32 // Latitude of the pole of stretching (microdegrees)
	LongitudeOfPoleStretching int32 // Longitude of the pole of stretching (microdegrees)
	StretchingFactor          int32 // Stretching factor
}

// StretchedRotatedSphericalHarmonicGrid contains stretched and rotated spherical harmonic coefficients grid specific fields (template 53)
type StretchedRotatedSphericalHarmonicGrid struct {
	RotatedSphericalHarmonicGrid       // Embedded rotated spherical harmonic grid
	LatitudeOfPoleStretching     int32 // Latitude of the pole of stretching (microdegrees)
	LongitudeOfPoleStretching    int32 // Longitude of the pole of stretching (microdegrees)
	StretchingFactor             int32 // Stretching factor
}

// SpaceViewGrid contains space view perspective or orthographic grid specific fields (template 90)
type SpaceViewGrid struct {
	ShapeOfEarth                uint8  // Shape of the Earth
	ScaleFactorRadiusEarth      uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth      uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis        uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis        uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis        uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis        uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX    uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY    uint32 // Number of points along y-axis
	LapValue                    int32  // Lap value (angle in microdegrees)
	LopValue                    uint32 // Lop value (angle in microdegrees)
	ResolutionAndComponentFlag  uint8  // Resolution and component flags
	XDirectionIncrement         uint32 // Apparent diameter of Earth in grid lengths
	YDirectionIncrement         uint32 // Apparent diameter of Earth in grid lengths
	XCoordinateOfOrigin         int32  // X-coordinate of sub-satellite point
	YCoordinateOfOrigin         int32  // Y-coordinate of sub-satellite point
	ScanningMode                uint8  // Scanning mode
	OrientationOfGrid           uint32 // Orientation of the grid (microdegrees)
	NrValue                     uint32 // Nr value (height of view point in Earth radii × 10^6)
	XCoordinateOfOriginOfSector int32  // X-coordinate of origin of sector image
	YCoordinateOfOriginOfSector int32  // Y-coordinate of origin of sector image
}

// TriangularGrid contains triangular grid based on an icosahedron specific fields (template 100)
type TriangularGrid struct {
	N2Value                    uint8  // n2 - exponent of 2 for the number of intervals on main triangle sides
	N3Value                    uint8  // n3 - exponent of 3 for the number of intervals on main triangle sides
	NumberOfIntsInDiamond      uint8  // n_i - number of intervals along each side of main icosahedral diamond
	ScanningModeForOneTriangle uint8  // Scanning mode for one triangle
	NumberOfPointsInTriangle   uint32 // Number of points in each triangle
	ScanningModeForDiamond     uint8  // Scanning mode for sequence of diamonds
	NumberOfDiamondsInGrid     uint32 // Total number of diamonds in grid
}

// EquatorialGrid contains equatorial azimuthal equidistant projection grid specific fields (template 110)
type EquatorialGrid struct {
	ShapeOfEarth               uint8  // Shape of the Earth
	ScaleFactorRadiusEarth     uint8  // Scale factor of radius of spherical Earth
	ScaledValueRadiusEarth     uint32 // Scaled value of radius of spherical Earth
	ScaleFactorMajorAxis       uint8  // Scale factor of major axis of oblate spheroid Earth
	ScaledValueMajorAxis       uint32 // Scaled value of major axis of oblate spheroid Earth
	ScaleFactorMinorAxis       uint8  // Scale factor of minor axis of oblate spheroid Earth
	ScaledValueMinorAxis       uint32 // Scaled value of minor axis of oblate spheroid Earth
	NumberOfGridPointsAlongX   uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY   uint32 // Number of points along y-axis
	LatitudeOfTangencyPoint    int32  // Latitude of tangency point (microdegrees)
	LongitudeOfTangencyPoint   uint32 // Longitude of tangency point (microdegrees)
	ResolutionAndComponentFlag uint8  // Resolution and component flags
	XDirectionIncrement        uint32 // X-direction grid length (metres)
	YDirectionIncrement        uint32 // Y-direction grid length (metres)
	ProjectionCenterFlag       uint8  // Projection center flag
	ScanningMode               uint8  // Scanning mode
}

// AzimuthRangeGrid contains azimuth-range projection grid specific fields (template 120)
type AzimuthRangeGrid struct {
	NumberOfDataPoints     uint32  // Number of data points
	CoordinateValuesList   []int32 // List of coordinate values
	CoordinateValuesNumber uint32  // Number of coordinate values
}

// IrregularGrid contains cross-section grid with points equally spaced on the horizontal specific fields (template 1000)
type IrregularGrid struct {
	NumberOfHorizontalPoints  uint32  // Number of points on horizontal level
	HorizontalPointDefinition uint8   // Basic angle of the initial production domain
	PhysicalMeaning1          uint8   // Physical meaning of vertical coordinate
	VerticalCoordinateValues  []int32 // Vertical coordinate values
}

// CurvilinearGrid contains curvilinear orthogonal grids specific fields (template 204)
type CurvilinearGrid struct {
	NumberOfGridPointsAlongX uint32 // Number of points along x-axis
	NumberOfGridPointsAlongY uint32 // Number of points along y-axis
	ScanningMode             uint8  // Scanning mode
}
