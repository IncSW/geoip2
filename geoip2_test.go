package geoip2

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	assert := assert.New(t)

	readerCity, err := NewReaderFromFile("testdata/GeoIP2-City.mmdb")
	if !assert.NoError(err) {
		return
	}

	cityResponse, err := readerCity.LookupCity(net.ParseIP("1.1.1.1"))
	if !assert.NoError(err) {
		return
	}

	data, _ := json.Marshal(cityResponse)
	t.Log(string(data))
}

func TestGeoIP2(t *testing.T) {
	assert := assert.New(t)

	readerCity, err := NewReaderFromFile("testdata/GeoIP2-City.mmdb")
	if !assert.NoError(err) {
		return
	}
	readerISP, err := NewReaderFromFile("testdata/GeoIP2-ISP.mmdb")
	if !assert.NoError(err) {
		return
	}
	readerConnectionType, err := NewReaderFromFile("testdata/GeoIP2-Connection-Type.mmdb")
	if !assert.NoError(err) {
		return
	}

	oschwaldReaderCity, err := geoip2.Open("testdata/GeoIP2-City.mmdb")
	if !assert.NoError(err) {
		return
	}
	oschwaldReaderISP, err := geoip2.Open("testdata/GeoIP2-ISP.mmdb")
	if !assert.NoError(err) {
		return
	}
	oschwaldReaderConnectionType, err := geoip2.Open("testdata/GeoIP2-Connection-Type.mmdb")
	if !assert.NoError(err) {
		return
	}

	for _, ip := range ips {
		cityResponse, err := readerCity.LookupCity(net.ParseIP(ip))
		if !assert.NoError(err) {
			return
		}
		isp, err := readerISP.LookupISP(net.ParseIP(ip))
		if !assert.NoError(err) {
			return
		}
		connectionType, err := readerConnectionType.LookupConnectionType(net.ParseIP(ip))
		if !assert.NoError(err) {
			return
		}
		oschwaldCity, err := oschwaldReaderCity.City(net.ParseIP(ip))
		if !assert.NoError(err) {
			return
		}
		oschwaldISP, err := oschwaldReaderISP.ISP(net.ParseIP(ip))
		if !assert.NoError(err) {
			return
		}
		oschwaldConnectionType, err := oschwaldReaderConnectionType.ConnectionType(net.ParseIP(ip))
		if !assert.NoError(err) {
			return
		}
		if cityResponse == nil {
			if oschwaldCity.Country.GeoNameID != 0 {
				t.Fatal(ip)
			}
			continue
		}
		if isp == nil {
			if oschwaldISP.ISP != "" {
				t.Fatal(ip)
			}
			continue
		}
		if connectionType == "" {
			if oschwaldConnectionType.ConnectionType != "" {
				t.Fatal(ip)
			}
			continue
		}
		if cityResponse.Continent.GeoNameID != uint32(oschwaldCity.Continent.GeoNameID) ||
			cityResponse.Continent.Code != oschwaldCity.Continent.Code ||
			cityResponse.Continent.Names["en"] != oschwaldCity.Continent.Names["en"] {
			t.Fatal(ip)
		}
		if cityResponse.Country.GeoNameID != uint32(oschwaldCity.Country.GeoNameID) ||
			cityResponse.Country.ISOCode != oschwaldCity.Country.IsoCode ||
			cityResponse.Country.Names["en"] != oschwaldCity.Country.Names["en"] {
			t.Fatal(ip)
		}
		if isp.AutonomousSystemNumber != uint32(oschwaldISP.AutonomousSystemNumber) ||
			isp.AutonomousSystemOrganization != oschwaldISP.AutonomousSystemOrganization ||
			isp.ISP != oschwaldISP.ISP ||
			isp.Organization != oschwaldISP.Organization {
			t.Fatal(ip)
		}
		if connectionType != oschwaldConnectionType.ConnectionType {
			t.Fatal(ip)
		}
	}
}

func BenchmarkGeoIP2(b *testing.B) {
	ip := net.ParseIP(ips[0])
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

func BenchmarkOschwaldGeoIP2(b *testing.B) {
	ip := net.ParseIP(ips[0])
	b.ReportAllocs()

	b.Run("city", func(b *testing.B) {
		reader, err := geoip2.Open("testdata/GeoIP2-City.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.City(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.City(ip)
				}
			})
		})
	})

	b.Run("isp", func(b *testing.B) {
		reader, err := geoip2.Open("testdata/GeoIP2-ISP.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.ISP(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.ISP(ip)
				}
			})
		})
	})

	b.Run("connection_type", func(b *testing.B) {
		reader, err := geoip2.Open("testdata/GeoIP2-Connection-Type.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.ConnectionType(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.ConnectionType(ip)
				}
			})
		})
	})
}
