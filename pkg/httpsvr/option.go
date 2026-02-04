package httpsvr

import "net/http"

type Option func(s *Server)

func WithConfig(cfg *Config) Option {
	return func(s *Server) {
		s.server.Addr = ":" + cfg.Port
		s.server.ReadTimeout = cfg.ReadTimeout
		s.server.ReadHeaderTimeout = cfg.ReadHeaderTimeout
		s.server.WriteTimeout = cfg.WriteTimeout
		s.server.IdleTimeout = cfg.IdleTimeout
	}
}

func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.server.Handler = handler
	}
}
