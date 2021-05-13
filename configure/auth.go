package configure

type Auth struct {
	Auth bool
	Rule []AuthRule
}
type AuthRule struct {
	URL    []string
	Bearer []string
}
