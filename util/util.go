package util

import (
	"fmt"
	"net/http"
)

func GetHttpUrl(schema, addr string) string {
	return fmt.Sprintf("%s://%s", schema, addr)
}

func GetFullUrl(r *http.Request) string{
	scheme := HTTP
	if r.TLS != nil {
		scheme = HTTPS
	}
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
}

func GetBaseUrl(r *http.Request) string {
	scheme := HTTP
	if r.TLS != nil {
		scheme = HTTPS
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}