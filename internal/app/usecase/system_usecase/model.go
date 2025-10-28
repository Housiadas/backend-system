package system_usecase

import "encoding/json"

type Status struct {
	Status string `json:"status"`
}

// Encode implements the encoder interface.
func (s Status) Encode() ([]byte, string, error) {
	data, err := json.Marshal(s)
	return data, "application/json", err
}

// Info represents information about the usecase.
type Info struct {
	Status     string `json:"status,omitempty"`
	Build      string `json:"build,omitempty"`
	Host       string `json:"host,omitempty"`
	Name       string `json:"name,omitempty"`
	PodIP      string `json:"podIP,omitempty"`
	Node       string `json:"node,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
}

// Encode implements the encoder interface.
func (info Info) Encode() ([]byte, string, error) {
	data, err := json.Marshal(info)
	return data, "application/json", err
}
