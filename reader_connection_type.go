package geoip2

import (
	"errors"
	"io/ioutil"
	"net"
	"strconv"
)

type ConnectionTypeReader struct {
	*reader
}

func (r *ConnectionTypeReader) Lookup(ip net.IP) (string, error) {
	offset, err := r.getOffset(ip)
	if err != nil {
		return "", err
	}
	dataType, size, offset, err := readControl(r.decoderBuffer, offset)
	if err != nil {
		return "", err
	}
	result := &ConnectionType{}
	switch dataType {
	case dataTypeMap:
		_, err = readConnectionTypeMap(result, r.decoderBuffer, size, offset)
		if err != nil {
			return "", err
		}
	case dataTypePointer:
		pointer, _, err := readPointer(r.decoderBuffer, size, offset)
		if err != nil {
			return "", err
		}
		dataType, size, offset, err := readControl(r.decoderBuffer, pointer)
		if err != nil {
			return "", err
		}
		if dataType != dataTypeMap {
			return "", errors.New("invalid Connection-Type pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readConnectionTypeMap(result, r.decoderBuffer, size, offset)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("invalid Connection-Type type: " + strconv.Itoa(int(dataType)))
	}
	return result.ConnectionType, nil
}

func NewConnectionTypeReaderType(buffer []byte, expectedTypes ...string) (*ConnectionTypeReader, error) {
	reader, err := newReader(buffer)
	if err != nil {
		return nil, err
	}
	if !isExpectedDatabaseType(reader.metadata.DatabaseType, expectedTypes...) {
		return nil, errors.New("wrong MaxMind DB Connection-Type type: " + reader.metadata.DatabaseType)
	}
	return &ConnectionTypeReader{
		reader: reader,
	}, nil
}

func NewConnectionTypeReader(buffer []byte) (*ConnectionTypeReader, error) {
	return NewConnectionTypeReaderType(buffer, "GeoIP2-Connection-Type")
}

func NewConnectionTypeReaderFromFile(filename string) (*ConnectionTypeReader, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewConnectionTypeReader(buffer)
}
