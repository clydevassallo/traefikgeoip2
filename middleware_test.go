package traefikgeoip2_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mw "github.com/clydevassallo/traefikgeoip2"
)

const (
	ValidIP       = "188.193.88.199"
)

func TestGeoIPFromRemoteAddr(t *testing.T) {
	mwCfg := mw.CreateConfig()
	
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	instance, _ := mw.New(context.TODO(), next, mwCfg, "traefik-geoip2")

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = fmt.Sprintf("%s:9999", ValidIP)
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, mw.CountryHeader, "DE")
	assertHeader(t, req, mw.RegionHeader, "BY")
	assertHeader(t, req, mw.CityHeader, "Munich")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	if req.Header.Get(key) != expected {
		t.Fatalf("invalid value of header [%s] != %s", key, req.Header.Get(key))
	}
}
