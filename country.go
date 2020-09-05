package geoip2

import (
	"errors"
	"strconv"
)

func readCountry(country *Country, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readCountryMap(country, buffer, size, offset)
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
			return 0, errors.New("invalid country pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readCountryMap(country, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid country type: " + strconv.Itoa(int(dataType)))
	}
}

func readCountryMap(country *Country, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "geoname_id":
			country.GeoNameID, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "iso_code":
			country.ISOCode, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "names":
			country.Names, offset, err = readStringMap(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_in_european_union":
			country.IsInEuropeanEnion, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "type":
			country.Type, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "confidence":
			country.Confidence, offset, err = readUInt16(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown country key: " + string(key))
		}
	}
	return offset, nil
}
