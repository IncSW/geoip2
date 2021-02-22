package geoip2

import "errors"

func readAnonymousIPMap(result *AnonymousIP, buffer []byte, mapSize uint, offset uint) (uint, error) {
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
		case "is_residential_proxy":
			result.IsResidentialProxy, offset, err = readBool(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown isp key: " + string(key))
		}
	}
	return offset, nil
}
