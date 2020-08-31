package geoip2

import (
	"errors"
	"strconv"
)

func readContinent(continent *Continent, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readContinentMap(continent, buffer, size, offset)
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
			return 0, errors.New("invalid continent pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readContinentMap(continent, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid continent type: " + strconv.Itoa(int(dataType)))
	}
}

func readContinentMap(continent *Continent, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "geoname_id":
			continent.GeoNameID, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "code":
			continent.Code, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "names":
			continent.Names, offset, err = readStringMap(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown continent key: " + string(key))
		}
	}
	return offset, nil
}
