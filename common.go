package geoip2

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"unsafe"
)

func readControl(buffer []byte, offset uint) (byte, uint, uint, error) {
	controlByte := buffer[offset]
	offset++
	dataType := controlByte >> 5
	if dataType == dataTypeExtended {
		dataType = buffer[offset] + 7
		offset++
	}
	size := uint(controlByte & 0x1f)
	if dataType == dataTypeExtended || size < 29 {
		return dataType, size, offset, nil
	}
	bytesToRead := size - 28
	newOffset := offset + bytesToRead
	if newOffset > uint(len(buffer)) {
		return 0, 0, 0, errors.New("invalid offset")
	}
	size = uint(bytesToUInt64(buffer[offset:newOffset]))
	switch bytesToRead {
	case 1:
		size += 29
	case 2:
		size += 285
	default:
		size += 65821
	}
	return dataType, size, newOffset, nil
}

func readPointer(buffer []byte, size uint, offset uint) (uint, uint, error) {
	pointerSize := ((size >> 3) & 0x3) + 1
	newOffset := offset + pointerSize
	if newOffset > uint(len(buffer)) {
		return 0, 0, errors.New("invalid offset")
	}
	prefix := uint64(0)
	if pointerSize != 4 {
		prefix = uint64(size) & 0x7
	}
	unpacked := uint(bytesToUInt64WithPrefix(prefix, buffer[offset:newOffset]))
	pointerValueOffset := uint(0)
	switch pointerSize {
	case 2:
		pointerValueOffset = 2048
	case 3:
		pointerValueOffset = 526336
	case 4:
		pointerValueOffset = 0
	}
	return unpacked + pointerValueOffset, newOffset, nil
}

func readFloat64(buffer []byte, offset uint) (float64, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, 0, err
	}
	switch dataType {
	case dataTypeFloat64:
		newOffset := offset + size
		return bytesToFloat64(buffer[offset:newOffset]), newOffset, nil
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return 0, 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return 0, 0, err
		}
		if dataType != dataTypeFloat64 {
			return 0, 0, errors.New("invalid float64 pointer type: " + strconv.Itoa(int(dataType)))
		}
		return bytesToFloat64(buffer[offset : offset+size]), newOffset, nil
	default:
		return 0, 0, errors.New("invalid float64 type: " + strconv.Itoa(int(dataType)))
	}
}

func readUInt16(buffer []byte, offset uint) (uint16, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, 0, err
	}
	switch dataType {
	case dataTypeUint16:
		newOffset := offset + size
		return uint16(bytesToUInt64(buffer[offset:newOffset])), newOffset, nil
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return 0, 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return 0, 0, err
		}
		if dataType != dataTypeUint16 {
			return 0, 0, errors.New("invalid uint16 pointer type: " + strconv.Itoa(int(dataType)))
		}
		return uint16(bytesToUInt64(buffer[offset : offset+size])), newOffset, nil
	default:
		return 0, 0, errors.New("invalid uint16 type: " + strconv.Itoa(int(dataType)))
	}
}

func readUInt32(buffer []byte, offset uint) (uint32, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return 0, 0, err
	}
	switch dataType {
	case dataTypeUint32:
		newOffset := offset + size
		return uint32(bytesToUInt64(buffer[offset:newOffset])), newOffset, nil
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return 0, 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return 0, 0, err
		}
		if dataType != dataTypeUint32 {
			return 0, 0, errors.New("invalid uint32 pointer type: " + strconv.Itoa(int(dataType)))
		}
		return uint32(bytesToUInt64(buffer[offset : offset+size])), newOffset, nil
	default:
		return 0, 0, errors.New("invalid uint32 type: " + strconv.Itoa(int(dataType)))
	}
}

func readBool(buffer []byte, offset uint) (bool, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return false, 0, err
	}
	switch dataType {
	case dataTypeBool:
		return size != 0, offset, nil
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return false, 0, err
		}
		dataType, size, _, err := readControl(buffer, pointer)
		if err != nil {
			return false, 0, err
		}
		if dataType != dataTypeBool {
			return false, 0, errors.New("invalid bool pointer type: " + strconv.Itoa(int(dataType)))
		}
		return size != 0, newOffset, nil
	default:
		return false, 0, errors.New("invalid bool type: " + strconv.Itoa(int(dataType)))
	}
}

func readString(buffer []byte, offset uint) (string, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return "", 0, err
	}
	switch dataType {
	case dataTypeString:
		newOffset := offset + size
		return b2s(buffer[offset:newOffset]), newOffset, nil
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return "", 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return "", 0, err
		}
		if dataType != dataTypeString {
			return "", 0, errors.New("invalid string pointer type: " + strconv.Itoa(int(dataType)))
		}
		return b2s(buffer[offset : offset+size]), newOffset, nil
	default:
		return "", 0, errors.New("invalid string type: " + strconv.Itoa(int(dataType)))
	}
}

func readStringMap(buffer []byte, offset uint) (map[string]string, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, 0, err
	}
	switch dataType {
	case dataTypeMap:
		return readStringMapMap(buffer, size, offset)
	case dataTypePointer:
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return nil, 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return nil, 0, err
		}
		if dataType != dataTypeMap {
			return nil, 0, errors.New("invalid stringMap pointer type: " + strconv.Itoa(int(dataType)))
		}
		value, _, err := readStringMapMap(buffer, size, offset)
		if err != nil {
			return nil, 0, err
		}
		return value, newOffset, nil
	default:
		return nil, 0, errors.New("invalid stringMap type: " + strconv.Itoa(int(dataType)))
	}
}

func readStringMapMap(buffer []byte, mapSize uint, offset uint) (map[string]string, uint, error) {
	var key []byte
	var err error
	var dataType byte
	var size uint
	result := map[string]string{}
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return nil, 0, err
		}
		dataType, size, offset, err = readControl(buffer, offset)
		if err != nil {
			return nil, 0, err
		}
		switch dataType {
		case dataTypePointer:
			pointer, newOffset, err := readPointer(buffer, size, offset)
			if err != nil {
				return nil, 0, err
			}
			dataType, size, valueOffset, err := readControl(buffer, pointer)
			if err != nil {
				return nil, 0, err
			}
			if dataType != dataTypeString {
				return nil, 0, errors.New("map key must be a string, got: " + strconv.Itoa(int(dataType)))
			}
			offset = newOffset
			result[b2s(key)] = b2s(buffer[valueOffset : valueOffset+size])
		case dataTypeString:
			newOffset := offset + size
			value := b2s(buffer[offset:newOffset])
			offset = newOffset
			result[b2s(key)] = value
		default:
			return nil, 0, errors.New("invalid data type of key " + string(key) + ": " + strconv.Itoa(int(dataType)))
		}
	}
	return result, offset, nil
}

func readMapKey(buffer []byte, offset uint) ([]byte, uint, error) {
	dataType, size, offset, err := readControl(buffer, offset)
	if err != nil {
		return nil, 0, err
	}
	if dataType == dataTypePointer {
		pointer, newOffset, err := readPointer(buffer, size, offset)
		if err != nil {
			return nil, 0, err
		}
		dataType, size, offset, err := readControl(buffer, pointer)
		if err != nil {
			return nil, 0, err
		}
		if dataType != dataTypeString {
			return nil, 0, errors.New("map key must be a string, got: " + strconv.Itoa(int(dataType)))
		}
		return buffer[offset : offset+size], newOffset, nil
	}
	if dataType != dataTypeString {
		return nil, 0, errors.New("map key must be a string, got: " + strconv.Itoa(int(dataType)))
	}
	newOffset := offset + size
	if newOffset > uint(len(buffer)) {
		return nil, 0, errors.New("invalid offset")
	}
	return buffer[offset:newOffset], newOffset, nil
}

func readStringSlice(buffer []byte, sliceSize uint, offset uint) ([]string, uint, error) {
	var err error
	var value string
	result := make([]string, sliceSize)
	for i := uint(0); i < sliceSize; i++ {
		value, offset, err = readString(buffer, offset)
		if err != nil {
			return nil, 0, err
		}
		result[i] = value
	}
	return result, offset, nil
}

func bytesToUInt64(buffer []byte) uint64 {
	switch len(buffer) {
	case 1:
		return uint64(buffer[0])
	case 2:
		return uint64(buffer[0])<<8 | uint64(buffer[1])
	case 3:
		return (uint64(buffer[0])<<8|uint64(buffer[1]))<<8 | uint64(buffer[2])
	case 4:
		return ((uint64(buffer[0])<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8 | uint64(buffer[3])
	case 5:
		return (((uint64(buffer[0])<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8 | uint64(buffer[4])
	case 6:
		return ((((uint64(buffer[0])<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8|uint64(buffer[4]))<<8 | uint64(buffer[5])
	case 7:
		return (((((uint64(buffer[0])<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8|uint64(buffer[4]))<<8|uint64(buffer[5]))<<8 | uint64(buffer[6])
	case 8:
		return ((((((uint64(buffer[0])<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8|uint64(buffer[4]))<<8|uint64(buffer[5]))<<8|uint64(buffer[6]))<<8 | uint64(buffer[7])
	}
	return 0
}

func bytesToUInt64WithPrefix(prefix uint64, buffer []byte) uint64 {
	switch len(buffer) {
	case 0:
		return prefix
	case 1:
		return prefix<<8 | uint64(buffer[0])
	case 2:
		return (prefix<<8|uint64(buffer[0]))<<8 | uint64(buffer[1])
	case 3:
		return ((prefix<<8|uint64(buffer[0]))<<8|uint64(buffer[1]))<<8 | uint64(buffer[2])
	case 4:
		return (((prefix<<8|uint64(buffer[0]))<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8 | uint64(buffer[3])
	case 5:
		return ((((prefix<<8|uint64(buffer[0]))<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8 | uint64(buffer[4])
	case 6:
		return (((((prefix<<8|uint64(buffer[0]))<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8|uint64(buffer[4]))<<8 | uint64(buffer[5])
	case 7:
		return ((((((prefix<<8|uint64(buffer[0]))<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8|uint64(buffer[4]))<<8|uint64(buffer[5]))<<8 | uint64(buffer[6])
	case 8:
		return (((((((prefix<<8|uint64(buffer[0]))<<8|uint64(buffer[1]))<<8|uint64(buffer[2]))<<8|uint64(buffer[3]))<<8|uint64(buffer[4]))<<8|uint64(buffer[5]))<<8|uint64(buffer[6]))<<8 | uint64(buffer[7])
	}
	return 0
}

func bytesToFloat32(buffer []byte) float32 {
	bits := binary.BigEndian.Uint32(buffer)
	return math.Float32frombits(bits)
}

func bytesToFloat64(buffer []byte) float64 {
	bits := binary.BigEndian.Uint64(buffer)
	return math.Float64frombits(bits)
}

func b2s(value []byte) string {
	return *(*string)(unsafe.Pointer(&value))
}
