package cors

type Option func(svc *Service)

func WithAllowOrigins(origins []string) Option {
	return func(svc *Service) {
		svc.allowOrigins = origins
	}
}

func WithAllowMethods(methods []string) Option {
	return func(svc *Service) {
		svc.allowMethods = methods
	}
}

func WithAllowHeaders(headers []string) Option {
	return func(svc *Service) {
		svc.allowHeaders = headers
	}
}
