package geoip2

import (
	"errors"
	"io/ioutil"
	"net"
	"strconv"
)

type CityReader struct {
	*reader
}

func (r *CityReader) Lookup(ip net.IP) (*CityResult, error) {
	offset, err := r.getOffset(ip)
	if err != nil {
		return nil, err
	}
	dataType, size, offset, err := readControl(r.decoderBuffer, offset)
	if err != nil {
		return nil, err
	}
	if dataType != dataTypeMap {
		return nil, errors.New("invalid City type: " + strconv.Itoa(int(dataType)))
	}
	var key []byte
	result := &CityResult{}
	for i := uint(0); i < size; i++ {
		key, offset, err = readMapKey(r.decoderBuffer, offset)
		if err != nil {
			return nil, err
		}
		switch b2s(key) {
		case "city":
			offset, err = readCity(&result.City, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "continent":
			offset, err = readContinent(&result.Continent, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "country":
			offset, err = readCountry(&result.Country, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "location":
			offset, err = readLocation(&result.Location, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "postal":
			offset, err = readPostal(&result.Postal, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "registered_country":
			offset, err = readCountry(&result.RegisteredCountry, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "represented_country":
			offset, err = readCountry(&result.RepresentedCountry, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "subdivisions":
			result.Subdivisions, offset, err = readSubdivisions(r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		case "traits":
			offset, err = readTraits(&result.Traits, r.decoderBuffer, offset)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown City response key: " + string(key) + ", type: " + strconv.Itoa(int(dataType)))
		}
	}
	return result, nil
}

// NewCityReaderWithType creates a new CityReader that accepts MMDB files with a custom database type. Note that
// CityReader only implements the fields provided by MaxMind Geo*-City and GeoIP2-Enterprise databases, and will ignore
// other fields. It is up to the developer to ensure that the database provides a compatible selection of fields.
func NewCityReaderWithType(buffer []byte, expectedTypes ...string) (*CityReader, error) {
	reader, err := newReader(buffer)
	if err != nil {
		return nil, err
	}
	if !isExpectedDatabaseType(reader.metadata.DatabaseType, expectedTypes...) {
		return nil, errors.New("wrong MaxMind DB City type: " + reader.metadata.DatabaseType)
	}
	return &CityReader{
		reader: reader,
	}, nil
}

func NewCityReader(buffer []byte) (*CityReader, error) {
	return NewCityReaderWithType(buffer, "GeoIP2-City", "GeoLite2-City", "GeoIP2-Enterprise")
}

func NewCityReaderFromFile(filename string) (*CityReader, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewCityReader(buffer)
}

func NewEnterpriseReader(buffer []byte) (*CityReader, error) {
	return NewCityReader(buffer)
}

func NewEnterpriseReaderFromFile(filename string) (*CityReader, error) {
	return NewCityReaderFromFile(filename)
}
