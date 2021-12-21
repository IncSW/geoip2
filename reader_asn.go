package geoip2

import (
	"errors"
	"io/ioutil"
	"net"
	"strconv"
)

type ASNReader struct {
	*reader
}

func (r *ASNReader) Lookup(ip net.IP) (*ASN, error) {
	offset, err := r.getOffset(ip)
	if err != nil {
		return nil, err
	}
	dataType, size, offset, err := readControl(r.decoderBuffer, offset)
	if err != nil {
		return nil, err
	}
	result := &ASN{}
	switch dataType {
	case dataTypeMap:
		_, err = readASNMap(result, r.decoderBuffer, size, offset)
		if err != nil {
			return nil, err
		}
	case dataTypePointer:
		pointer, _, err := readPointer(r.decoderBuffer, size, offset)
		if err != nil {
			return nil, err
		}
		dataType, size, offset, err := readControl(r.decoderBuffer, pointer)
		if err != nil {
			return nil, err
		}
		if dataType != dataTypeMap {
			return nil, errors.New("invalid ASN pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readASNMap(result, r.decoderBuffer, size, offset)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid ASN type: " + strconv.Itoa(int(dataType)))
	}
	return result, nil
}

// NewASNReaderWithType creates a new ASNReader that accepts MMDB files with a custom database type. Note that
// ASNReader only implements the fields provided by MaxMind GeoLite2-ASN databases, and will ignore other fields.
// It is up to the developer to ensure that the database provides a compatible selection of fields.
func NewASNReaderWithType(buffer []byte, expectedTypes ...string) (*ASNReader, error) {
	reader, err := newReader(buffer)
	if err != nil {
		return nil, err
	}
	if !isExpectedDatabaseType(reader.metadata.DatabaseType, expectedTypes...) {
		return nil, errors.New("wrong MaxMind DB ASN type: " + reader.metadata.DatabaseType)
	}
	return &ASNReader{
		reader: reader,
	}, nil
}

func NewASNReader(buffer []byte) (*ASNReader, error) {
	return NewASNReaderWithType(buffer, "GeoLite2-ASN")
}

func NewASNReaderFromFile(filename string) (*ASNReader, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewASNReader(buffer)
}
