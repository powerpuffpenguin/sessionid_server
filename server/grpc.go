package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	"google.golang.org/grpc"
)

func newServer(mux *runtime.ServeMux, cc *grpc.ClientConn) (srv *grpc.Server) {
	srv = grpc.NewServer()
	gmodule.InitServer(srv, mux, cc)
	return
}
