package config

import "time"

type Grpc struct {
	Api             string
	ShutdownTimeout time.Duration
}
