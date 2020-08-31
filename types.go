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
	IsInEuropeanEnion bool
	Names             map[string]string
	Type              string
}

type City struct {
	GeoNameID uint32
	Names     map[string]string
}

type Location struct {
	AccuracyRadius uint16
	MetroCode      uint16
	Latitude       float64
	Longitude      float64
	TimeZone       string
}

type Postal struct {
	Code string
}

type Traits struct {
	IsAnonymousProxy    bool
	IsSatelliteProvider bool
	StaticIPScore       float64
}

type CityResponse struct {
	Continent          Continent
	City               City
	Country            Country
	Subdivisions       []Subdivision
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
}

type ConnectionType struct {
	ConnectionType string
}
