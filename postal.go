package geoip2

import (
	"errors"
	"strconv"
)

func readPostal(postal *Postal, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readPostalMap(postal, buffer, size, offset)
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
			return 0, errors.New("invalid postal pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readPostalMap(postal, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid postal type: " + strconv.Itoa(int(dataType)))
	}
}

func readPostalMap(postal *Postal, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "code":
			postal.Code, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "confidence":
			postal.Confidence, offset, err = readUInt8(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown postal key: " + string(key))
		}
	}
	return offset, nil
}
