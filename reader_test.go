// Test DB https://github.com/maxmind/MaxMind-DB
package geoip2

import (
	"net"
	"testing"
)

func TestAnonymousIP(t *testing.T) {
	reader, err := NewAnonymousIPReaderFromFile("testdata/GeoIP2-Anonymous-IP-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Lookup(net.ParseIP("81.2.69.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.IsAnonymous != true {
		t.Fatal()
	}
	if record.IsAnonymousVPN != true {
		t.Fatal()
	}
	if record.IsHostingProvider != true {
		t.Fatal()
	}
	if record.IsPublicProxy != true {
		t.Fatal()
	}
	if record.IsTorExitNode != true {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("186.30.236.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.IsAnonymous != true {
		t.Fatal()
	}
	if record.IsAnonymousVPN != false {
		t.Fatal()
	}
	if record.IsHostingProvider != false {
		t.Fatal()
	}
	if record.IsPublicProxy != true {
		t.Fatal()
	}
	if record.IsTorExitNode != false {
		t.Fatal()
	}
}

func TestReaderZeroLength(t *testing.T) {
	_, err := newReader([]byte{})
	if err == nil {
		t.Fatal()
	}
}

func TestCity(t *testing.T) {
	reader, err := NewCityReaderFromFile("testdata/GeoIP2-City-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Lookup(net.ParseIP("81.2.69.142"))
	if err != nil {
		t.Fatal(err)
	}
	if record.City.GeoNameID != 2643743 {
		t.Fatal()
	}
	if record.City.Names["de"] != "London" ||
		record.City.Names["es"] != "Londres" {
		t.Fatal()
	}
	if record.Location.AccuracyRadius != 10 {
		t.Fatal()
	}
	if record.Location.Latitude != 51.5142 {
		t.Fatal()
	}
	if record.Location.Longitude != -0.0931 {
		t.Fatal()
	}
	if record.Location.TimeZone != "Europe/London" {
		t.Fatal()
	}
	if len(record.Subdivisions) != 1 {
		t.Fatal()
	}
	if record.Subdivisions[0].GeoNameID != 6269131 {
		t.Fatal()
	}
	if record.Subdivisions[0].ISOCode != "ENG" {
		t.Fatal()
	}
	if record.Subdivisions[0].Names["en"] != "England" ||
		record.Subdivisions[0].Names["pt-BR"] != "Inglaterra" {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("2a02:ff80::"))
	if err != nil {
		t.Fatal(err)
	}
	if record.City.GeoNameID != 0 {
		t.Fatal()
	}
	if record.Country.IsInEuropeanUnion != true {
		t.Fatal()
	}
	if record.Location.AccuracyRadius != 100 {
		t.Fatal()
	}
	if record.Location.Latitude != 51.5 {
		t.Fatal()
	}
	if record.Location.Longitude != 10.5 {
		t.Fatal()
	}
	if record.Location.TimeZone != "Europe/Berlin" {
		t.Fatal()
	}
	if len(record.Subdivisions) != 0 {
		t.Fatal()
	}
}

func TestConnectionType(t *testing.T) {
	reader, err := NewConnectionTypeReaderFromFile("testdata/GeoIP2-Connection-Type-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Lookup(net.ParseIP("1.0.0.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record != "Dialup" {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("1.0.1.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record != "Cable/DSL" {
		t.Fatal()
	}
}

func TestCountry(t *testing.T) {
	reader, err := NewCountryReaderFromFile("testdata/GeoIP2-Country-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Lookup(net.ParseIP("74.209.24.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.Continent.GeoNameID != 6255149 {
		t.Fatal()
	}
	if record.Continent.Code != "NA" {
		t.Fatal()
	}
	if record.Continent.Names["es"] != "Norteamérica" ||
		record.Continent.Names["ru"] != "Северная Америка" {
		t.Fatal()
	}
	if record.Country.GeoNameID != 6252001 {
		t.Fatal()
	}
	if record.Country.ISOCode != "US" {
		t.Fatal()
	}
	if record.Country.Names["fr"] != "États-Unis" ||
		record.Country.Names["pt-BR"] != "Estados Unidos" {
		t.Fatal()
	}
	if record.Country.IsInEuropeanUnion != false {
		t.Fatal()
	}
	if record.RegisteredCountry.GeoNameID != 6252001 {
		t.Fatal()
	}
	if record.RepresentedCountry.GeoNameID != 0 {
		t.Fatal()
	}
	if record.Traits.IsAnonymousProxy != true {
		t.Fatal()
	}
	if record.Traits.IsSatelliteProvider != true {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("2a02:ffc0::"))
	if err != nil {
		t.Fatal(err)
	}
	if record.Continent.GeoNameID != 6255148 {
		t.Fatal()
	}
	if record.Continent.Code != "EU" {
		t.Fatal()
	}
	if record.Continent.Names["en"] != "Europe" ||
		record.Continent.Names["zh-CN"] != "欧洲" {
		t.Fatal()
	}
	if record.Country.GeoNameID != 2411586 {
		t.Fatal()
	}
	if record.Country.ISOCode != "GI" {
		t.Fatal()
	}
	if record.Country.Names["en"] != "Gibraltar" ||
		record.Country.Names["ja"] != "ジブラルタル" {
		t.Fatal()
	}
	if record.Country.IsInEuropeanUnion != false {
		t.Fatal()
	}
	if record.RegisteredCountry.GeoNameID != 2411586 {
		t.Fatal()
	}
	if record.RepresentedCountry.GeoNameID != 0 {
		t.Fatal()
	}
	if record.Traits.IsAnonymousProxy != false {
		t.Fatal()
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

	record, err = reader.Lookup(net.ParseIP("186.30.236.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record != "replaced.com" {
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
	if record.Country.Confidence != 99 {
		t.Fatal()
	}
	if record.Postal.Code != "12037" {
		t.Fatal()
	}
	if record.Postal.Confidence != 11 {
		t.Fatal()
	}
	if len(record.Subdivisions) != 1 {
		t.Fatal()
	}
	if record.Subdivisions[0].Confidence != 93 {
		t.Fatal()
	}
	if record.Traits.AutonomousSystemNumber != 14671 {
		t.Fatal()
	}
	if record.Traits.AutonomousSystemOrganization != "FairPoint Communications" {
		t.Fatal()
	}
	if record.Traits.ISP != "Fairpoint Communications" {
		t.Fatal()
	}
	if record.Traits.Organization != "Fairpoint Communications" {
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
	if record.Traits.UserType != "residential" {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("81.2.69.160"))
	if err != nil {
		t.Fatal(err)
	}
	if record.Traits.ISP != "Andrews & Arnold Ltd" {
		t.Fatal()
	}
	if record.Traits.Organization != "STONEHOUSE office network" {
		t.Fatal()
	}
	if record.Traits.ConnectionType != "Corporate" {
		t.Fatal()
	}
	if record.Traits.Domain != "in-addr.arpa" {
		t.Fatal()
	}
	if record.Traits.StaticIPScore != 0.34 {
		t.Fatal()
	}
	if record.Traits.UserType != "government" {
		t.Fatal()
	}
}

func TestISP(t *testing.T) {
	reader, err := NewISPReaderFromFile("testdata/GeoIP2-ISP-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Lookup(net.ParseIP("1.128.0.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.AutonomousSystemNumber != 1221 {
		t.Fatal()
	}
	if record.AutonomousSystemOrganization != "Telstra Pty Ltd" {
		t.Fatal()
	}
	if record.ISP != "Telstra Internet" {
		t.Fatal()
	}
	if record.Organization != "Telstra Internet" {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("4.0.0.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.AutonomousSystemNumber != 0 {
		t.Fatal()
	}
	if record.AutonomousSystemOrganization != "" {
		t.Fatal()
	}
	if record.ISP != "Level 3 Communications" {
		t.Fatal()
	}
	if record.Organization != "Level 3 Communications" {
		t.Fatal()
	}
}

func TestASN(t *testing.T) {
	reader, err := NewASNReaderFromFile("testdata/GeoLite2-ASN-Test.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Lookup(net.ParseIP("1.128.0.0"))
	if err != nil {
		t.Fatal(err)
	}
	if record.AutonomousSystemNumber != 1221 {
		t.Fatal()
	}
	if record.AutonomousSystemOrganization != "Telstra Pty Ltd" {
		t.Fatal()
	}

	record, err = reader.Lookup(net.ParseIP("2600:6000::"))
	if err != nil {
		t.Fatal(err)
	}
	if record.AutonomousSystemNumber != 237 {
		t.Fatal()
	}
	if record.AutonomousSystemOrganization != "Merit Network Inc." {
		t.Fatal()
	}
}
