package config

import "time"

type Server struct {
	Api                string
	Debug              string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	CorsAllowedOrigins []string
}
