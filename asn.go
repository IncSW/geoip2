package geoip2

import (
	"errors"
	"strconv"
)

func readASNResult(buffer []byte, offset uint) (*ASNResult, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, err
	}
	result := &ASNResult{}
	switch dataType {
	case dataTypeMap:
		_, err = readASNMap(result, buffer, size, offset)
		if err != nil {
			return nil, err
		}
	case dataTypePointer:
		pointer, _, err := readPointer(buffer, size, offset)
		if err != nil {
			return nil, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return nil, err
		}
		if dataType != dataTypeMap {
			return nil, errors.New("invalid asn pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readASNMap(result, buffer, size, offset)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid asn type: " + strconv.Itoa(int(dataType)))
	}
	return result, nil
}

func readASNMap(result *ASNResult, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "autonomous_system_number":
			result.AutonomousSystemNumber, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "autonomous_system_organization":
			result.AutonomousSystemOrganization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown isp key: " + string(key))
		}
	}
	return offset, nil
}
