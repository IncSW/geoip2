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

	record, err := reader.LookupCity(net.ParseIP("81.2.69.142"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func BenchmarkGeoIP2(b *testing.B) {
	ip := net.ParseIP("81.2.69.142")
	b.ReportAllocs()

	b.Run("city", func(b *testing.B) {
		reader, err := NewReaderFromFile("testdata/GeoIP2-City.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.LookupCity(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.LookupCity(ip)
				}
			})
		})
	})

	b.Run("isp", func(b *testing.B) {
		reader, err := NewReaderFromFile("testdata/GeoIP2-ISP.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.LookupISP(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.LookupISP(ip)
				}
			})
		})
	})

	b.Run("connection_type", func(b *testing.B) {
		reader, err := NewReaderFromFile("testdata/GeoIP2-Connection-Type.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.LookupConnectionType(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.LookupConnectionType(ip)
				}
			})
		})
	})
}
