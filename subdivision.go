package geoip2

import (
	"errors"
	"strconv"
)

func readSubdivisions(buffer []byte, offset uint) ([]Subdivision, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, 0, err
	}
	switch dataType {
	case dataTypeSlice:
		return readSubdivisionsSlice(buffer, size, offset)
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return nil, 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return nil, 0, err
		}
		if dataType != dataTypeSlice {
			return nil, 0, errors.New("invalid subdivisions pointer type: " + strconv.Itoa(int(dataType)))
		}
		subdivisions, _, err := readSubdivisionsSlice(buffer, size, offset)
		if err != nil {
			return nil, 0, err
		}
		return subdivisions, newOffset, nil
	default:
		return nil, 0, errors.New("invalid subdivisions type: " + strconv.Itoa(int(dataType)))
	}
}

func readSubdivisionsSlice(buffer []byte, subdivisionsSize uint, offset uint) ([]Subdivision, uint, error) {
	var err error
	subdivisions := make([]Subdivision, subdivisionsSize)
	for i := uint(0); i < subdivisionsSize; i++ {
		offset, err = readSubdivision(&subdivisions[i], buffer, offset)
		if err != nil {
			return nil, 0, err
		}
	}
	return subdivisions, offset, nil
}

func readSubdivision(subdivision *Subdivision, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readSubdivisionMap(subdivision, buffer, size, offset)
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
			return 0, errors.New("invalid subdivision pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readSubdivisionMap(subdivision, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid subdivision type: " + strconv.Itoa(int(dataType)))
	}
}

func readSubdivisionMap(subdivision *Subdivision, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "geoname_id":
			subdivision.GeoNameID, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "iso_code":
			subdivision.ISOCode, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "names":
			subdivision.Names, offset, err = readStringMap(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "confidence":
			subdivision.Confidence, offset, err = readUInt8(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown subdivision key: " + string(key))
		}
	}
	return offset, nil
}
