package systemapp

// Info represents information about the service.
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

type Status struct {
	Status string `json:"status"`
}

type ApiError struct {
	StatusCode int                    `json:"status"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details"`
	Context    map[string]interface{} `json:"context"`
}
