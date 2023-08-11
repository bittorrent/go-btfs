package server

type Option func(*Server)

func WithAddress(address string) Option {
	return func(s *Server) {
		s.address = address
	}
}
