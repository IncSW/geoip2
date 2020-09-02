package geoip2

import (
	"errors"
	"strconv"
)

func readCityResult(buffer []byte, offset uint) (*CityResult, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, err
	}
	if dataType != dataTypeMap {
		return nil, errors.New("invalid city type: " + strconv.Itoa(int(dataType)))
	}
	var key []byte
	result := &CityResult{}
	for i := uint(0); i < size; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return nil, err
		}
		switch b2s(key) {
		case "city":
			offset, err = readCity(&result.City, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "continent":
			offset, err = readContinent(&result.Continent, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "country":
			offset, err = readCountry(&result.Country, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "location":
			offset, err = readLocation(&result.Location, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "postal":
			offset, err = readPostal(&result.Postal, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "registered_country":
			offset, err = readCountry(&result.RegisteredCountry, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "represented_country":
			offset, err = readCountry(&result.RepresentedCountry, buffer, offset)
			if err != nil {
				return nil, err
			}
		case "subdivisions":
			result.Subdivisions, offset, err = readSubdivisions(buffer, offset)
			if err != nil {
				return nil, err
			}
		case "traits":
			offset, err = readTraits(&result.Traits, buffer, offset)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown city response key: " + string(key) + ", type: " + strconv.Itoa(int(dataType)))
		}
	}
	return result, nil
}

func readCity(city *City, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readCityMap(city, buffer, size, offset)
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return 0, err
		}
		if dataType != dataTypeMap {
			return 0, errors.New("invalid city pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readCityMap(city, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid city type: " + strconv.Itoa(int(dataType)))
	}
}

func readCityMap(city *City, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "geoname_id":
			city.GeoNameID, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "names":
			city.Names, offset, err = readStringMap(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown city key: " + string(key))
		}
	}
	return offset, nil
}
