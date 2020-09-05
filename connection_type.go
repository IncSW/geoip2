package geoip2

import "errors"

func readConnectionTypeMap(result *ConnectionType, buffer []byte, mapSize uint, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch b2s(key) {
		case "connection_type":
			result.ConnectionType, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown connectionType key: " + string(key))
		}
	}
	return offset, nil
}
