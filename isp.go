package geoip2

import "errors"

func readISPMap(result *ISP, buffer []byte, mapSize uint, offset uint) (uint, error) {
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
		case "isp":
			result.ISP, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "organization":
			result.Organization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "mobile_country_code":
			result.MobileCountryCode, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "mobile_network_code":
			result.MobileNetworkCode, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown isp key: " + string(key))
		}
	}
	return offset, nil
}
