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
	ISOCode           string
	Names             map[string]string
	Type              string // [RepresentedCountry]
	GeoNameID         uint32
	Confidence        uint16 // Enterprise [Country, RegisteredCountry]
	IsInEuropeanUnion bool
}

type Subdivision struct {
	ISOCode    string
	Names      map[string]string
	GeoNameID  uint32
	Confidence uint16 // Enterprise
}

type City struct {
	Names      map[string]string
	GeoNameID  uint32
	Confidence uint16 // Enterprise
}

type Location struct {
	Latitude       float64
	Longitude      float64
	TimeZone       string
	AccuracyRadius uint16
	MetroCode      uint16
}

type Postal struct {
	Code       string
	Confidence uint16 // Enterprise
}

type Traits struct {
	StaticIPScore                float64 // Enterprise
	ISP                          string  // Enterprise
	Organization                 string  // Enterprise
	ConnectionType               string  // Enterprise
	Domain                       string  // Enterprise
	UserType                     string  // Enterprise
	AutonomousSystemOrganization string  // Enterprise
	AutonomousSystemNumber       uint32  // Enterprise
	IsLegitimateProxy            bool    // Enterprise
	MobileCountryCode            string  // Enterprise
	MobileNetworkCode            string  // Enterprise
	IsAnonymousProxy             bool
	IsSatelliteProvider          bool
}

type CountryResult struct {
	Continent          Continent
	Country            Country
	RegisteredCountry  Country
	RepresentedCountry Country
	Traits             Traits
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

type ISP struct {
	AutonomousSystemNumber       uint32
	AutonomousSystemOrganization string
	ISP                          string
	Organization                 string
	MobileCountryCode            string
	MobileNetworkCode            string
}

type ConnectionType struct {
	ConnectionType string
}

type AnonymousIP struct {
	IsAnonymous        bool
	IsAnonymousVPN     bool
	IsHostingProvider  bool
	IsPublicProxy      bool
	IsTorExitNode      bool
	IsResidentialProxy bool
}

type ASN struct {
	AutonomousSystemNumber       uint32
	AutonomousSystemOrganization string
	Network                      string
}

type Domain struct {
	Domain string
}
