package geoip2

const (
	dataTypeExtended           = 0
	dataTypePointer            = 1
	dataTypeString             = 2
	dataTypeFloat64            = 3
	dataTypeBytes              = 4
	dataTypeUint16             = 5
	dataTypeUint32             = 6
	dataTypeMap                = 7
	dataTypeInt32              = 8
	dataTypeUint64             = 9
	dataTypeUint128            = 10
	dataTypeSlice              = 11
	dataTypeDataCacheContainer = 12
	dataTypeEndMarker          = 13
	dataTypeBool               = 14
	dataTypeFloat32            = 15

	dataSectionSeparatorSize = 16
)

type Continent struct {
	GeoNameID uint32
	Code      string
	Names     map[string]string
}

type Country struct {
	GeoNameID         uint32
	ISOCode           string
	Names             map[string]string
	Confidence        uint8 // Enterprise [Country, RegisteredCountry]
	IsInEuropeanEnion bool
	Type              string // [RepresentedCountry]
}

type Subdivision struct {
	GeoNameID  uint32
	ISOCode    string
	Names      map[string]string
	Confidence uint8 // Enterprise
}

type City struct {
	GeoNameID  uint32
	Names      map[string]string
	Confidence uint8 // Enterprise
}

type Location struct {
	AccuracyRadius uint16
	MetroCode      uint16
	Latitude       float64
	Longitude      float64
	TimeZone       string
}

type Postal struct {
	Code       string
	Confidence uint8 // Enterprise
}

type Traits struct {
	IsAnonymousProxy             bool
	IsSatelliteProvider          bool
	IsLegitimateProxy            bool    // Enterprise
	StaticIPScore                float64 // Enterprise
	AutonomousSystemNumber       uint32  // Enterprise
	AutonomousSystemOrganization string  // Enterprise
	ISP                          string  // Enterprise
	Organization                 string  // Enterprise
	ConnectionType               string  // Enterprise
	Domain                       string  // Enterprise
	UserType                     string  // Enterprise
}

type CityResult struct {
	Continent          Continent
	Country            Country
	Subdivisions       []Subdivision
	City               City
	Location           Location
	Postal             Postal
	RegisteredCountry  Country
	RepresentedCountry Country
	Traits             Traits
}

type ISPResult struct {
	AutonomousSystemNumber       uint32
	AutonomousSystemOrganization string
	ISP                          string
	Organization                 string
}

type ConnectionTypeResult struct {
	ConnectionType string
}

type AnonymousIPResult struct {
	IsAnonymous       bool
	IsAnonymousVPN    bool
	IsHostingProvider bool
	IsPublicProxy     bool
	IsTorExitNode     bool
}

type ASNResult struct {
	AutonomousSystemNumber       uint32
	AutonomousSystemOrganization string
}

type DomainResult struct {
	Domain string
}
