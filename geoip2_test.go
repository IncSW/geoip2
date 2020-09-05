package geoip2

import (
	"net"
	"testing"
	"time"
)

func TestReader(t *testing.T) {
	ip := net.ParseIP("81.2.69.142")

	countryReader, err := NewCountryReaderFromFile("testdata/GeoIP2-Country.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = countryReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}
	countryLiteReader, err := NewCountryReaderFromFile("testdata/GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = countryLiteReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	cityReader, err := NewCityReaderFromFile("testdata/GeoIP2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = cityReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}
	cityLiteReader, err := NewCityReaderFromFile("testdata/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = cityLiteReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	ispReader, err := NewISPReaderFromFile("testdata/GeoIP2-ISP.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ispReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	connectionTypeReader, err := NewConnectionTypeReaderFromFile("testdata/GeoIP2-Connection-Type.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = connectionTypeReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	asnReader, err := NewASNReaderFromFile("testdata/GeoLite2-ASN.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = asnReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDomain(t *testing.T) {
	reader, err := NewDomainReaderFromFile("testdata/GeoIP2-Domain-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	record, err := reader.Lookup(net.ParseIP("1.2.0.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record != "maxmind.com" {
		t.Fatal()
	}
}

func TestAnonymousIP(t *testing.T) {
	reader, err := NewAnonymousIPReaderFromFile("testdata/GeoIP2-Anonymous-IP-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	record, err := reader.Lookup(net.ParseIP("1.2.0.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.IsAnonymous != true {
		t.Fatal()
	}
	if record.IsAnonymousVPN != true {
		t.Fatal()
	}
	if record.IsHostingProvider != false {
		t.Fatal()
	}
	if record.IsPublicProxy != false {
		t.Fatal()
	}
	if record.IsTorExitNode != false {
		t.Fatal()
	}
}

func TestEnterprise(t *testing.T) {
	reader, err := NewEnterpriseReaderFromFile("testdata/GeoIP2-Enterprise-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	record, err := reader.Lookup(net.ParseIP("74.209.24.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.City.Confidence != 11 {
		t.Fatal()
	}
	if record.Traits.AutonomousSystemNumber != 14671 {
		t.Fatal()
	}
	if record.Traits.AutonomousSystemOrganization != "FairPoint Communications" {
		t.Fatal()
	}
	if record.Traits.ConnectionType != "Cable/DSL" {
		t.Fatal()
	}
	if record.Traits.Domain != "frpt.net" {
		t.Fatal()
	}
	if record.Traits.StaticIPScore != 0.34 {
		t.Fatal()
	}
}

func TestBench(t *testing.T) {
	reader, err := NewCityReaderFromFile("testdata/GeoIP2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	ip := net.ParseIP("81.2.69.142")
	var minDuration time.Duration
	for i := 0; i < 200; i++ {
		start := time.Now()
		for j := 0; j < 10000; j++ {
			_, _ = reader.Lookup(ip)
		}
		duration := time.Since(start)
		if minDuration == 0 || minDuration > duration {
			minDuration = duration
		}
	}
	t.Log(int(minDuration/10000), "ns/op")
}

func BenchmarkGeoIP2(b *testing.B) {
	ip := net.ParseIP("81.2.69.142")
	b.ReportAllocs()

	b.Run("country", func(b *testing.B) {
		reader, err := NewCountryReaderFromFile("testdata/GeoIP2-Country.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.Lookup(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.Lookup(ip)
				}
			})
		})
	})

	b.Run("city", func(b *testing.B) {
		reader, err := NewCityReaderFromFile("testdata/GeoIP2-City.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.Lookup(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.Lookup(ip)
				}
			})
		})
	})

	b.Run("isp", func(b *testing.B) {
		reader, err := NewISPReaderFromFile("testdata/GeoIP2-ISP.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.Lookup(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.Lookup(ip)
				}
			})
		})
	})

	b.Run("connection_type", func(b *testing.B) {
		reader, err := NewConnectionTypeReaderFromFile("testdata/GeoIP2-Connection-Type.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.Lookup(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.Lookup(ip)
				}
			})
		})
	})

	b.Run("asn", func(b *testing.B) {
		reader, err := NewASNReaderFromFile("testdata/GeoLite2-ASN.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = reader.Lookup(ip)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = reader.Lookup(ip)
				}
			})
		})
	})
}
