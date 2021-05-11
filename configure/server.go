package configure

type Server struct {
	// http addr
	Addr string
	// if not empty use https
	CertFile string
	KeyFile  string
}
