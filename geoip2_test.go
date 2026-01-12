package geoip2

import (
	"net"
	"runtime"
	"testing"
)

func TestReader(t *testing.T) {
	t.Parallel()

	dbPath := func(db string) string {
		return "testdata/maxmind/test-data/" + db + ".mmdb"
	}

	tests := map[string][]struct {
		db string
		ip string
	}{
		"country": {
			{db: "GeoIP2-Country-Test", ip: "81.2.69.160"},
			{db: "GeoLite2-Country-Test", ip: "81.2.69.160"},
		},
		"city": {
			{db: "GeoIP2-City-Test", ip: "81.2.69.160"},
			{db: "GeoLite2-City-Test", ip: "81.2.69.160"},
		},
		"isp": {
			{db: "GeoIP2-ISP-Test", ip: "81.2.69.160"},
		},
		"connection_type": {
			{db: "GeoIP2-Connection-Type-Test", ip: "2.125.160.216"},
		},
		"asn": {
			{db: "GeoLite2-ASN-Test", ip: "81.128.69.160"},
		},
	}
	for reader, readerTest := range tests {
		for _, test := range readerTest {
			t.Run(test.db, func(t *testing.T) {
				t.Parallel()
				switch reader {
				case "country":
					r, err := NewCountryReaderFromFile(dbPath(test.db))
					if err != nil {
						t.Fatal(err)
					}
					if _, err = r.Lookup(net.ParseIP(test.ip)); err != nil {
						t.Fatal(err)
					}
				case "city":
					r, err := NewCityReaderFromFile(dbPath(test.db))
					if err != nil {
						t.Fatal(err)
					}
					if _, err = r.Lookup(net.ParseIP(test.ip)); err != nil {
						t.Fatal(err)
					}
				case "isp":
					r, err := NewISPReaderFromFile(dbPath(test.db))
					if err != nil {
						t.Fatal(err)
					}
					if _, err = r.Lookup(net.ParseIP(test.ip)); err != nil {
						t.Fatal(err)
					}
				case "connection_type":
					r, err := NewConnectionTypeReaderFromFile(dbPath(test.db))
					if err != nil {
						t.Fatal(err)
					}
					if _, err = r.Lookup(net.ParseIP(test.ip)); err != nil {
						t.Fatal(err)
					}
				case "asn":
					r, err := NewASNReaderFromFile(dbPath(test.db))
					if err != nil {
						t.Fatal(err)
					}
					if _, err = r.Lookup(net.ParseIP(test.ip)); err != nil {
						t.Fatal(err)
					}
				default:
					t.Fatal("unknown reader type: " + reader)
				}
			})
		}
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
			for range b.N {
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

	b.Run("country-GeoLite2-Country", func(b *testing.B) {
		reader, err := NewCountryReaderFromFile("testdata/GeoLite2-Country.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for range b.N {
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
			for range b.N {
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

	b.Run("city-GeoLite2-City", func(b *testing.B) {
		reader, err := NewCityReaderFromFile("testdata/GeoLite2-City.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for range b.N {
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
			for range b.N {
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
			for range b.N {
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
			for range b.N {
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

	b.Run("asn-GeoLite2-ASN", func(b *testing.B) {
		reader, err := NewASNReaderFromFile("testdata/GeoLite2-ASN.mmdb")
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sync", func(b *testing.B) {
			for range b.N {
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
