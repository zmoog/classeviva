package spaggiari

import "net/http"

const (
	baseUrl   = "https://web.spaggiari.eu/rest/v1"
	apiKey    = "Tg1NWEwNGIgIC0K"
	userAgent = "CVVS/std/4.2.3 Android/12"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
