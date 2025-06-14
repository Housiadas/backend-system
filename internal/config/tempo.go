package config

type Tempo struct {
	Host string
	// Shouldn't use a high Probability value in non-developer systems.
	// 0.05 should be enough for most systems. Some might want to have
	// this even lower.
	Probability float64
}
