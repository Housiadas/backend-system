// Package kafka is a client for apache kafka
package kafka

type Event struct {
	Topic string
	Data  []byte
}

type EventData struct {
	Action string
	Params any
}
