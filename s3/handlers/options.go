package handlers

type Option func(handlers *Handlers)

func WithCorsAllowOrigins(origins []string) Option {
	return func(handlers *Handlers) {
		handlers.corsAllowOrigins = origins
	}
}
