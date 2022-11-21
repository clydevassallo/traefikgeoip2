// Package traefikgeoip2 is a Traefik plugin for Maxmind GeoIP2.
package traefikgeoip2

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/IncSW/geoip2"
)

var lookupCache LookupGeoIP2

// Config the plugin configuration.
type Config struct {
	DBPath string `json:"dbPath,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		DBPath: DefaultDBPath,
	}
}

// TraefikGeoIP2 a traefik geoip2 plugin.
type TraefikGeoIP2 struct {
	next   http.Handler
	lookup LookupGeoIP2
	name   string
}

// New created a new TraefikGeoIP2 plugin.
func New(ctx context.Context, next http.Handler, cfg *Config, name string) (http.Handler, error) {
	lookup, err := getLookup(cfg)
	if err != nil {
		return nil, err
	}

	return &TraefikGeoIP2{
		lookup: lookup,
		next:   next,
		name:   name,
	}, nil
}

func (mw *TraefikGeoIP2) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	log.Printf("[geoip2] remoteAddr: %v", req.RemoteAddr)

	if mw.lookup == nil {
		req.Header.Set(CountryHeader, Unknown)
		req.Header.Set(RegionHeader, Unknown)
		req.Header.Set(CityHeader, Unknown)
		mw.next.ServeHTTP(reqWr, req)
		return
	}

	ipStr := req.RemoteAddr
	tmp, _, err := net.SplitHostPort(ipStr)
	if err == nil {
		ipStr = tmp
	}

	res, err := mw.lookup(net.ParseIP(ipStr))
	if err != nil {
		log.Printf("[geoip2] Unable to find for `%s', %v", ipStr, err)
		res = &GeoIPResult{
			country: Unknown,
			region:  Unknown,
			city:    Unknown,
		}
	}

	req.Header.Set(CountryHeader, res.country)
	req.Header.Set(RegionHeader, res.region)
	req.Header.Set(CityHeader, res.city)

	mw.next.ServeHTTP(reqWr, req)
}

// ClearCache resets GeoIP2 DB cache.
func ClearCache() {
	lookupCache = nil
}

func getLookup(cfg *Config) (LookupGeoIP2, error) {
	if lookupCache != nil {
		return lookupCache, nil
	}

	if _, err := os.Stat(cfg.DBPath); err != nil {
		log.Printf("[geoip2] DB `%s' not found: %v", cfg.DBPath, err)
		return lookupCache, nil
	}

	if strings.Contains(cfg.DBPath, "City") {
		rdr, err := geoip2.NewCityReaderFromFile(cfg.DBPath)
		if err != nil {
			log.Printf("[geoip2] DB `%s' not initialized: %v", cfg.DBPath, err)
			return nil, fmt.Errorf("%w", err)
		}

		lookupCache = CreateCityDBLookup(rdr)
		log.Printf("[geoip2] DB city `%s' was initialized", cfg.DBPath)
		return lookupCache, nil
	}

	if strings.Contains(cfg.DBPath, "Country") {
		rdr, err := geoip2.NewCountryReaderFromFile(cfg.DBPath)
		if err != nil {
			log.Printf("[geoip2] DB `%s' not initialized: %v", cfg.DBPath, err)
			return nil, fmt.Errorf("%w", err)
		}
		lookupCache = CreateCountryDBLookup(rdr)
		log.Printf("[geoip2] DB country `%s' was initialized", cfg.DBPath)
	}

	return lookupCache, nil
}
