package config

import "time"

type Otel struct {
	Host        string
	ReporterURI string
	Probability float64
}

type Prometheus struct {
	Host            string
	Route           string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type Collect struct {
	From string
}

type Publish struct {
	To       string
	Interval time.Duration
}
