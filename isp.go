package geoip2

import (
	"errors"
	"strconv"
)

func readISP(buffer []byte, offset uint) (*ISP, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, err
	}
	isp := &ISP{}
	switch dataType {
	case dataTypeMap:
		_, err = readISPMap(isp, buffer, size, offset)
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
			return nil, errors.New("invalid isp pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readISPMap(isp, buffer, size, offset)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid isp type: " + strconv.Itoa(int(dataType)))
	}
	return isp, nil
}

func readISPMap(isp *ISP, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "autonomous_system_number":
			isp.AutonomousSystemNumber, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "autonomous_system_organization":
			isp.AutonomousSystemOrganization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "isp":
			isp.ISP, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "organization":
			isp.Organization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown isp key: " + string(key))
		}
	}
	return offset, nil
}
