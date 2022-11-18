package traefikgeoip2

import (
)

// Unknown constant for undefined data.
const Unknown = "XX"

const (
	// CountryHeader country header name.
	CountryHeader = "X-GeoIP2-Country"
	// RegionHeader region header name.
	RegionHeader = "X-GeoIP2-Region"
	// CityHeader city header name.
	CityHeader = "X-GeoIP2-City"
)

// GeoIPResult GeoIPResult.
type GeoIPResult struct {
	country string
	region  string
	city    string
}