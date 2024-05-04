package config

type Kafka struct {
	Broker           string
	AddressFamily    string
	SecurityProtocol string
	LogLevel         int
	MaxMessageBytes  int
}
