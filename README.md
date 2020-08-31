[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/IncSW/geoip2?style=flat-square)](https://goreportcard.com/report/github.com/IncSW/geoip2)

# GeoIP2 Reader for Go

This library reads MaxMind GeoIP2 databases.

## Installation

`go get github.com/IncSW/geoip2`

## Quick Start

```go
import "github.com/IncSW/geoip2"

reader, err := geoip2.NewReaderFromFile("path/to/GeoIP2-City.mmdb")
if err != nil {
	panic(err)
}
record, err := reader.LookupCity(net.ParseIP("1.1.1.1"))
if err != nil {
	panic(err)
}
data, err := json.Marshal(record)
if err != nil {
	panic(err)
}
println(string(data))
// {
// 	"Continent": {
// 		"GeoNameID": 6255151,
// 		"Code": "OC",
// 		"Names": {
// 			"de": "Ozeanien",
// 			"en": "Oceania",
// 			"es": "Oceanía",
// 			"fr": "Océanie",
// 			"ja": "オセアニア",
// 			"pt-BR": "Oceania",
// 			"ru": "Океания",
// 			"zh-CN": "大洋洲"
// 		}
// 	},
// 	"City": {
// 		"GeoNameID": 0,
// 		"Names": null
// 	},
// 	"Country": {
// 		"GeoNameID": 2077456,
// 		"ISOCode": "AU",
// 		"IsInEuropeanEnion": false,
// 		"Names": {
// 			"de": "Australien",
// 			"en": "Australia",
// 			"es": "Australia",
// 			"fr": "Australie",
// 			"ja": "オーストラリア",
// 			"pt-BR": "Austrália",
// 			"ru": "Австралия",
// 			"zh-CN": "澳大利亚"
// 		},
// 		"Type": ""
// 	},
// 	"Subdivisions": null,
// 	"Location": {
// 		"AccuracyRadius": 1000,
// 		"MetroCode": 0,
// 		"Latitude": -33.494,
// 		"Longitude": 143.2104,
// 		"TimeZone": "Australia/Sydney"
// 	},
// 	"Postal": {
// 		"Code": ""
// 	},
// 	"RegisteredCountry": {
// 		"GeoNameID": 2077456,
// 		"ISOCode": "AU",
// 		"IsInEuropeanEnion": false,
// 		"Names": {
// 			"de": "Australien",
// 			"en": "Australia",
// 			"es": "Australia",
// 			"fr": "Australie",
// 			"ja": "オーストラリア",
// 			"pt-BR": "Austrália",
// 			"ru": "Австралия",
// 			"zh-CN": "澳大利亚"
// 		},
// 		"Type": ""
// 	},
// 	"RepresentedCountry": {
// 		"GeoNameID": 0,
// 		"ISOCode": "",
// 		"IsInEuropeanEnion": false,
// 		"Names": null,
// 		"Type": ""
// 	},
// 	"Traits": {
// 		"IsAnonymousProxy": false,
// 		"IsSatelliteProvider": false,
// 		"StaticIPScore": 0
// 	}
// }
```

## Performance

TODO

## License

[MIT License](LICENSE).
