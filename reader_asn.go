package geoip2

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

type ASNReader struct {
	*reader
}

func (r *ASNReader) Lookup(ip net.IP) (*ASN, error) {
	offset, prefix, err := r.getOffsetWithPrefix(ip)
	if err != nil {
		return nil, err
	}
	dataType, size, offset, err := readControl(r.decoderBuffer, offset)
	if err != nil {
		return nil, err
	}
	result := &ASN{}
	_, network, err := net.ParseCIDR(fmt.Sprintf("%s/%s", ip.String(), strconv.Itoa(int(prefix))))
	result.Network = network.String()
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

func NewASNReader(buffer []byte) (*ASNReader, error) {
	reader, err := newReader(buffer)
	if err != nil {
		return nil, err
	}
	if reader.metadata.DatabaseType != "GeoLite2-ASN" &&
		reader.metadata.DatabaseType != "DBIP-ASN-Lite" &&
		reader.metadata.DatabaseType != "DBIP-ASN-Lite (compat=GeoLite2-ASN)" {
		return nil, errors.New("wrong MaxMind DB ASN type: " + reader.metadata.DatabaseType)
	}
	return &ASNReader{
		reader: reader,
	}, nil
}

func NewASNReaderFromFile(filename string) (*ASNReader, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewASNReader(buffer)
}
