package server

const defaultServerAddress = "127.0.0.1:15001"

type Option func(*Server)

func WithAddress(address string) Option {
	return func(s *Server) {
		if address != "" {
			s.address = address
		}
	}
}
