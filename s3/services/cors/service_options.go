package cors

type Option func(svc *service)

func WithAllowOrigins(origins []string) Option {
	return func(svc *service) {
		svc.allowOrigins = origins
	}
}

func WithAllowMethods(methods []string) Option {
	return func(svc *service) {
		svc.allowMethods = methods
	}
}

func WithAllowHeaders(headers []string) Option {
	return func(svc *service) {
		svc.allowHeaders = headers
	}
}
