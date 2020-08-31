package geoip2

import (
	"encoding/json"
	"net"
	"testing"
)

func TestDebug(t *testing.T) {
	reader, err := NewReaderFromFile("testdata/GeoIP2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.LookupCity(net.ParseIP("1.1.1.1"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
