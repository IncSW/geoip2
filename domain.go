package geoip2

import (
	"errors"
	"strconv"
)

func readDomain(buffer []byte, offset uint) (string, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return "", err
	}
	result := &DomainResult{}
	switch dataType {
	case dataTypeMap:
		_, err = readDomainMap(result, buffer, size, offset)
		if err != nil {
			return "", err
		}
	case dataTypePointer:
		pointer, _, err := readPointer(buffer, size, offset)
		if err != nil {
			return "", err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return "", err
		}
		if dataType != dataTypeMap {
			return "", errors.New("invalid domain pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readDomainMap(result, buffer, size, offset)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("invalid domain type: " + strconv.Itoa(int(dataType)))
	}
	return result.Domain, nil
}

func readDomainMap(result *DomainResult, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "domain":
			result.Domain, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown domain key: " + string(key))
		}
	}
	return offset, nil
}
