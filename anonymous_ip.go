package geoip2

import (
	"errors"
	"strconv"
)

func readAnonymousIPResult(buffer []byte, offset uint) (*AnonymousIPResult, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, err
	}
	result := &AnonymousIPResult{}
	switch dataType {
	case dataTypeMap:
		_, err = readAnonymousIPMap(result, buffer, size, offset)
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
			return nil, errors.New("invalid anonymous ip pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readAnonymousIPMap(result, buffer, size, offset)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid anonymous ip type: " + strconv.Itoa(int(dataType)))
	}
	return result, nil
}

func readAnonymousIPMap(result *AnonymousIPResult, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "is_anonymous":
			result.IsAnonymous, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_anonymous_vpn":
			result.IsAnonymousVPN, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_hosting_provider":
			result.IsHostingProvider, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_public_proxy":
			result.IsPublicProxy, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "is_tor_exit_node":
			result.IsTorExitNode, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown isp key: " + string(key))
		}
	}
	return offset, nil
}
