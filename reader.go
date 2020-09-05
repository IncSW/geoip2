package geoip2

import (
	"bytes"
	"errors"
	"net"
	"strconv"
)

var ErrNotFound = errors.New("not found")

type reader struct {
	metadata          *Metadata
	buffer            []byte
	decoderBuffer     []byte
	nodeBuffer        []byte
	ipV4Start         uint
	ipV4StartBitDepth uint
	nodeOffsetMult    uint
}

func (r *reader) getOffset(ip net.IP) (uint, error) {
	pointer, err := r.lookupPointer(ip)
	if err != nil {
		return 0, err
	}
	offset := pointer - uint(r.metadata.NodeCount) - uint(dataSectionSeparatorSize)
	if offset >= uint(len(r.buffer)) {
		return 0, errors.New("the MaxMind DB search tree is corrupt: " + strconv.Itoa(int(pointer)))
	}
	return offset, nil
}

func (r *reader) lookupPointer(ip net.IP) (uint, error) {
	if ip == nil {
		return 0, errors.New("IP cannot be nil")
	}
	ipV4 := ip.To4()
	if ipV4 != nil {
		ip = ipV4
	}
	if len(ip) == 16 && r.metadata.IPVersion == 4 {
		return 0, errors.New("cannot look up an IPv6 address in an IPv4-only database")
	}
	bitCount := uint(len(ip)) * 8
	node := uint(0)
	if bitCount == 32 {
		node = r.ipV4Start
	}
	nodeCount := uint(r.metadata.NodeCount)
	i := uint(0)
	for ; i < bitCount && node < nodeCount; i++ {
		bit := 1 & (ip[i>>3] >> (7 - (i % 8)))
		offset := node * r.nodeOffsetMult
		if bit == 0 {
			node = r.readLeft(offset)
		} else {
			node = r.readRight(offset)
		}
	}
	if node == nodeCount {
		return 0, ErrNotFound
	} else if node > nodeCount {
		return node, nil
	}
	return 0, errors.New("invalid node in search tree")
}

func (r *reader) readLeft(nodeNumber uint) uint {
	switch r.metadata.RecordSize {
	case 28:
		return ((uint(r.nodeBuffer[nodeNumber+3]) & 0xF0) << 20) | (uint(r.nodeBuffer[nodeNumber]) << 16) | (uint(r.nodeBuffer[nodeNumber+1]) << 8) | uint(r.nodeBuffer[nodeNumber+2])
	case 24:
		return (uint(r.nodeBuffer[nodeNumber]) << 16) | (uint(r.nodeBuffer[nodeNumber+1]) << 8) | uint(r.nodeBuffer[nodeNumber+2])
	default: // case 32:
		return (uint(r.nodeBuffer[nodeNumber]) << 24) | (uint(r.nodeBuffer[nodeNumber+1]) << 16) | (uint(r.nodeBuffer[nodeNumber+2]) << 8) | uint(r.nodeBuffer[nodeNumber+3])
	}
}

func (r *reader) readRight(nodeNumber uint) uint {
	switch r.metadata.RecordSize {
	case 28:
		return ((uint(r.nodeBuffer[nodeNumber+3]) & 0x0F) << 24) | (uint(r.nodeBuffer[nodeNumber+4]) << 16) | (uint(r.nodeBuffer[nodeNumber+5]) << 8) | uint(r.nodeBuffer[nodeNumber+6])
	case 24:
		return (uint(r.nodeBuffer[nodeNumber+3]) << 16) | (uint(r.nodeBuffer[nodeNumber+4]) << 8) | uint(r.nodeBuffer[nodeNumber+5])
	default: // case 32:
		return (uint(r.nodeBuffer[nodeNumber+4]) << 24) | (uint(r.nodeBuffer[nodeNumber+5]) << 16) | (uint(r.nodeBuffer[nodeNumber+6]) << 8) | uint(r.nodeBuffer[nodeNumber+7])
	}
}

func newReader(buffer []byte) (*reader, error) {
	metadataStart := bytes.LastIndex(buffer, metadataStartMarker)
	metadata, err := readMetadata(buffer[metadataStart+len(metadataStartMarker):])
	if err != nil {
		return nil, err
	}
	nodeOffsetMult := uint(metadata.RecordSize) / 4
	searchTreeSize := uint(metadata.NodeCount) * nodeOffsetMult
	dataSectionStart := searchTreeSize + dataSectionSeparatorSize
	if dataSectionStart > uint(metadataStart) {
		return nil, errors.New("the MaxMind DB contains invalid metadata")
	}
	reader := &reader{
		metadata:       metadata,
		buffer:         buffer,
		decoderBuffer:  buffer[searchTreeSize+dataSectionSeparatorSize : metadataStart],
		nodeBuffer:     buffer[:searchTreeSize],
		nodeOffsetMult: nodeOffsetMult,
	}
	if metadata.IPVersion == 6 {
		node := uint(0)
		i := uint(0)
		for ; i < 96 && node < uint(metadata.NodeCount); i++ {
			node = reader.readLeft(node * nodeOffsetMult)
		}
		reader.ipV4Start = node
		reader.ipV4StartBitDepth = i
	}
	return reader, nil
}
