package geoip2

import (
	"errors"
	"strconv"
)

type Metadata struct {
	NodeCount                uint32            // node_count This is an unsigned 32-bit integer indicating the number of nodes in the search tree.
	RecordSize               uint16            // record_size This is an unsigned 16-bit integer. It indicates the number of bits in a record in the search tree. Note that each node consists of two records.
	IPVersion                uint16            // ip_version This is an unsigned 16-bit integer which is always 4 or 6. It indicates whether the database contains IPv4 or IPv6 address data.
	DatabaseType             string            // database_type This is a string that indicates the structure of each data record associated with an IP address. The actual definition of these structures is left up to the database creator. Names starting with “GeoIP” are reserved for use by MaxMind (and “GeoIP” is a trademark anyway).
	Languages                []string          // languages An array of strings, each of which is a locale code. A given record may contain data items that have been localized to some or all of these locales. Records should not contain localized data for locales not included in this array. This is an optional key, as this may not be relevant for all types of data.
	BinaryFormatMajorVersion uint16            // binary_format_major_version This is an unsigned 16-bit integer indicating the major version number for the database’s binary format.
	BinaryFormatMinorVersion uint16            // binary_format_minor_version This is an unsigned 16-bit integer indicating the minor version number for the database’s binary format.
	BuildEpoch               uint64            // build_epoch This is an unsigned 64-bit integer that contains the database build timestamp as a Unix epoch value.
	Description              map[string]string // description This key will always point to a map. The keys of that map will be language codes, and the values will be a description in that language as a UTF-8 string. The codes may include additional information such as script or country identifiers, like “zh-TW” or “mn-Cyrl-MN”. The additional identifiers will be separated by a dash character (“-“).
}

var metadataStartMarker = []byte("\xAB\xCD\xEFMaxMind.com")

func readMetadata(buffer []byte) (*Metadata, error) {
	dataType, metadataSize, offset, err := readControl(buffer, 0)
	if err != nil {
		return nil, err
	}
	if dataType != dataTypeMap {
		return nil, errors.New("invalid metadata type: " + strconv.Itoa(int(dataType)))
	}
	var key []byte
	metadata := &Metadata{}
	for i := uint(0); i < metadataSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return nil, err
		}
		size := uint(0)
		dataType, size, offset, err = readControl(buffer, offset)
		if err != nil {
			return nil, err
		}
		newOffset := uint(0)
		switch b2s(key) {
		case "binary_format_major_version":
			if dataType != dataTypeUint16 {
				return nil, errors.New("invalid binary_format_major_version type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.BinaryFormatMajorVersion = uint16(bytesToUInt64(buffer[offset:newOffset]))
		case "binary_format_minor_version":
			if dataType != dataTypeUint16 {
				return nil, errors.New("invalid binary_format_minor_version type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.BinaryFormatMinorVersion = uint16(bytesToUInt64(buffer[offset:newOffset]))
		case "build_epoch":
			if dataType != dataTypeUint64 {
				return nil, errors.New("invalid build_epoch type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.BuildEpoch = bytesToUInt64(buffer[offset:newOffset])
		case "database_type":
			if dataType != dataTypeString {
				return nil, errors.New("invalid database_type type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.DatabaseType = b2s(buffer[offset:newOffset])
		case "description":
			if dataType != dataTypeMap {
				return nil, errors.New("invalid description type: " + strconv.Itoa(int(dataType)))
			}
			metadata.Description, newOffset, err = readStringMapMap(buffer, size, offset)
			if err != nil {
				return nil, err
			}
		case "ip_version":
			if dataType != dataTypeUint16 {
				return nil, errors.New("invalid ip_version type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.IPVersion = uint16(bytesToUInt64(buffer[offset:newOffset]))
		case "languages":
			if dataType != dataTypeSlice {
				return nil, errors.New("invalid languages type: " + strconv.Itoa(int(dataType)))
			}
			metadata.Languages, newOffset, err = readStringSlice(buffer, size, offset)
			if err != nil {
				return nil, err
			}
		case "node_count":
			if dataType != dataTypeUint32 {
				return nil, errors.New("invalid node_count type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.NodeCount = uint32(bytesToUInt64(buffer[offset:newOffset]))
		case "record_size":
			if dataType != dataTypeUint16 {
				return nil, errors.New("invalid record_size type: " + strconv.Itoa(int(dataType)))
			}
			newOffset = offset + size
			metadata.RecordSize = uint16(bytesToUInt64(buffer[offset:newOffset]))
		default:
			return nil, errors.New("unknown key: " + string(key) + ", type: " + strconv.Itoa(int(dataType)))
		}
		offset = newOffset
	}
	return metadata, nil
}
