package geoip2

import (
	"net"
	"runtime"
	"testing"
)

func TestReader(t *testing.T) {
	ip := net.ParseIP("81.2.69.160")

	countryReader, err := NewCountryReaderFromFile("testdata/maxmind/test-data/GeoIP2-Country-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = countryReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}
	countryLiteReader, err := NewCountryReaderFromFile("testdata/maxmind/test-data/GeoLite2-Country-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = countryLiteReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	cityReader, err := NewCityReaderFromFile("testdata/maxmind/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = cityReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}
	cityLiteReader, err := NewCityReaderFromFile("testdata/maxmind/test-data/GeoLite2-City-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = cityLiteReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	ispReader, err := NewISPReaderFromFile("testdata/maxmind/test-data/GeoIP2-ISP-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ispReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	ip = net.ParseIP("2.125.160.216")
	connectionTypeReader, err := NewConnectionTypeReaderFromFile("testdata/maxmind/test-data/GeoIP2-Connection-Type-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = connectionTypeReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}

	ip = net.ParseIP("81.128.69.160")
	asnReader, err := NewASNReaderFromFile("testdata/maxmind/test-data/GeoLite2-ASN-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	_, err = asnReader.Lookup(ip)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkGeoIP2(b *testing.B) {
	ip := net.ParseIP("1.128.0.0")
	b.ReportAllocs()

	b.Run("country", func(b *testing.B) {
		reader, err := NewCountryReaderFromFile("testdata/maxmind/test-data/GeoIP2-Country-Test.mmdb")
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
		reader, err := NewCityReaderFromFile("testdata/maxmind/test-data/GeoIP2-City-Test.mmdb")
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
		reader, err := NewISPReaderFromFile("testdata/maxmind/test-data/GeoIP2-ISP-Test.mmdb")
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
		reader, err := NewConnectionTypeReaderFromFile("testdata/maxmind/test-data/GeoIP2-Connection-Type-Test.mmdb")
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
		reader, err := NewASNReaderFromFile("testdata/maxmind/test-data/GeoLite2-ASN-Test.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				asn, err := reader.Lookup(ip)
				if err != nil {
					b.Fatal(err)
				}
				runtime.KeepAlive(asn.Network)
			}
		})
		b.Run("parallel", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					asn, err := reader.Lookup(ip)
					if err != nil {
						b.Fatal(err)
					}
					runtime.KeepAlive(asn.Network)
				}
			})
		})
	})
}
