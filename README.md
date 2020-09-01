[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/IncSW/geoip2?style=flat-square)](https://goreportcard.com/report/github.com/IncSW/geoip2)

# GeoIP2 Reader for Go

This library reads MaxMind GeoIP2 databases.

Inspired by [oschwald/geoip2-golang](https://github.com/oschwald/geoip2-golang).

## Installation

`go get github.com/IncSW/geoip2`

## Quick Start

```go
import "github.com/IncSW/geoip2"

reader, err := geoip2.NewReaderFromFile("path/to/GeoIP2-City.mmdb")
if err != nil {
	panic(err)
}
record, err := reader.LookupCity(net.ParseIP("81.2.69.142"))
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

TODO

## License

[MIT License](LICENSE).
