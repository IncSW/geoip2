package geoip2

import (
	"errors"
	"strconv"
)

func readLocation(location *Location, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readLocationMap(location, buffer, size, offset)
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
			return 0, errors.New("invalid location pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readLocationMap(location, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid location type: " + strconv.Itoa(int(dataType)))
	}
}

func readLocationMap(location *Location, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "latitude":
			location.Latitude, offset, err = readFloat64(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "longitude":
			location.Longitude, offset, err = readFloat64(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "accuracy_radius":
			location.AccuracyRadius, offset, err = readUInt16(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "time_zone":
			location.TimeZone, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "metro_code":
			location.MetroCode, offset, err = readUInt16(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown location key: " + string(key))
		}
	}
	return offset, nil
}
