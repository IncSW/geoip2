package geoip2

import (
	"errors"
	"strconv"
)

func readTraits(traits *Traits, buffer []byte, offset uint) (uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readTraitsMap(traits, buffer, size, offset)
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
			return 0, errors.New("invalid traits pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readTraitsMap(traits, buffer, size, offset)
		if err != nil {
			return 0, err
		}
		return newOffset, nil
	default:
		return 0, errors.New("invalid traits type: " + strconv.Itoa(int(dataType)))
	}
}

func readTraitsMap(traits *Traits, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "is_anonymous_proxy":
			traits.IsAnonymousProxy, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_satellite_provider":
			traits.IsSatelliteProvider, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_legitimate_proxy":
			traits.IsLegitimateProxy, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "static_ip_score":
			traits.StaticIPScore, offset, err = readFloat64(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "autonomous_system_number":
			traits.AutonomousSystemNumber, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "autonomous_system_organization":
			traits.AutonomousSystemOrganization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "isp":
			traits.ISP, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "organization":
			traits.Organization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "connection_type":
			traits.ConnectionType, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "domain":
			traits.Domain, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "user_type":
			traits.UserType, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown traits key: " + string(key))
		}
	}
	return offset, nil
}
