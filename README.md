[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/IncSW/geoip2?style=flat-square)](https://goreportcard.com/report/github.com/IncSW/geoip2)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/IncSW/geoip2?tab=doc)

# GeoIP2 Reader for Go

This library reads MaxMind GeoIP2 databases.

Inspired by [oschwald/geoip2-golang](https://github.com/oschwald/geoip2-golang).

## Installation

`go get github.com/IncSW/geoip2`

## Quick Start

```go
import "github.com/IncSW/geoip2"

reader, err := geoip2.NewCityReaderFromFile("path/to/GeoIP2-City.mmdb")
if err != nil {
	panic(err)
}
record, err := reader.Lookup(net.ParseIP("81.2.69.142"))
if err != nil {
	panic(err)
}
println(record.Continent.Names["zh-CN"]) // 欧洲
println(record.City.Names["pt-BR"]) // Wimbledon
if len(record.Subdivisions) != 0 {
	println(record.Subdivisions[0].Names["en"]) // England
}
println(record.Country.Names["ru"]) // Великобритания
println(record.Country.ISOCode) // GB
println(record.Location.TimeZone) // Europe/London
println(record.Country.GeoNameID) // 2635167, https://www.geonames.org/2635167
```

## Performance

### [IncSW/geoip2](https://github.com/IncSW/geoip2)
```
city-24                          342847    2981 ns/op   2032 B/op    12 allocs/op
city_parallel-24                4477626     269 ns/op   2032 B/op    12 allocs/op
isp-24                          3539738     336 ns/op     64 B/op     1 allocs/op
isp_parallel-24                46938070    25.7 ns/op     64 B/op     1 allocs/op
connection_type-24              8759110     133 ns/op      0 B/op     0 allocs/op
connection_type_parallel-24   142261742    8.34 ns/op      0 B/op     0 allocs/op
```

### [oschwald/geoip2-golang](https://github.com/oschwald/geoip2-golang)
```
city-24                          109092   10717 ns/op   2848 B/op   103 allocs/op
city_parallel-24                 662510    1718 ns/op   2848 B/op   103 allocs/op
isp-24                          1688287     705 ns/op    112 B/op     4 allocs/op
isp_parallel-24                14285560    84.4 ns/op    112 B/op     4 allocs/op
connection_type-24              3883234     305 ns/op     32 B/op     2 allocs/op
connection_type_parallel-24    34284831    32.1 ns/op     32 B/op     2 allocs/op
```

## Supported databases types

### Country
- GeoIP2-Country
- GeoLite2-Country
- DBIP-Country
- DBIP-Country-Lite

### City
- GeoIP2-City
- GeoLite2-City
- GeoIP2-Enterprise
- DBIP-City-Lite

### ISP
- GeoIP2-ISP

### ASN
- GeoLite2-ASN
- DBIP-ASN-Lite
- DBIP-ASN-Lite (compat=GeoLite2-ASN)

### Connection Type
- GeoIP2-Connection-Type

### Anonymous IP
- GeoIP2-Anonymous-IP

### Domain
- GeoIP2-Domain

## MMDB files for tests

MMDB files for tests are organised in their respective directories based on the source within the `testdata` repository root directory.

### MaxMind

These are obtained using a submodule into a `testdata/maxmind` directory.

```
git submodule init
git submodule update
```

### DB-IP

You must obtain these files manually and place them into the `testdata/dbip` directory.

## License

[MIT License](LICENSE).
