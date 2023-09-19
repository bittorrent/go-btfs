package providers

import "time"

type BtfsAPIOption func(api *BtfsAPI)

const (
	defaultBtfsAPIEndpointUrl           = ""
	defaultBtfsAPITimeout               = 20 * time.Minute
	defaultBtfsAPIResponseHeaderTimeout = 1 * time.Minute
)

func BtfsAPIWithTimeout(timeout time.Duration) BtfsAPIOption {
	return func(api *BtfsAPI) {
		api.timeout = timeout
	}
}

func BtfsAPIWithBtfsAPIHeaderTimeout(timeout time.Duration) BtfsAPIOption {
	return func(api *BtfsAPI) {
		api.headerTimout = timeout
	}
}

func BtfsAPIWithEndpointUrl(url string) BtfsAPIOption {
	return func(api *BtfsAPI) {
		api.endpointUrl = url
	}
}
