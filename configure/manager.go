package configure

type Manager struct {
	// Token signature algorithm
	Method string
	// Signing key
	Key string
	// Serialization coder
	Coder string
}
