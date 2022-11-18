// Package traefikgeoip2 is a Traefik plugin for Maxmind GeoIP2.
package traefikgeoip2

import (
	"context"
	"log"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
		Test string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Test: "hello",
	}
}

// TraefikGeoIP2 a traefik geoip2 plugin.
type TraefikGeoIP2 struct {
	next   http.Handler
	name   string
}

// New created a new TraefikGeoIP2 plugin.
func New(ctx context.Context, next http.Handler, cfg *Config, name string) (http.Handler, error) {
	log.Printf("[geoip2] newly created")

	return &TraefikGeoIP2{
		next:   next,
		name:   name,
	}, nil
}

func (mw *TraefikGeoIP2) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	log.Printf("[geoip2] remoteAddr: %v", req.RemoteAddr)

	req.Header.Set(CountryHeader, "DE")
	req.Header.Set(RegionHeader, "BY")
	req.Header.Set(CityHeader, "Munich")

	mw.next.ServeHTTP(reqWr, req)
}
