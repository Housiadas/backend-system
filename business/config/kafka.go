package config

type Kafka struct {
	Brokers          string
	AddressFamily    string
	SecurityProtocol string
	LogLevel         int
	MaxMessageBytes  int
	SessionTimeout   int
}
