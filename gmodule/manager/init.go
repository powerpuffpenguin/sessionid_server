package manager

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	grpc_manager "github.com/powerpuffpenguin/sessionid_server/protocol/manager"
	"google.golang.org/grpc"
)

func init() {
	gmodule.RegisterModule(`manager`, Module(0))
}

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_manager.RegisterManagerServer(srv, server{})
}
func (Module) RegisterGateway(mux *runtime.ServeMux, cc *grpc.ClientConn) {
	grpc_manager.RegisterManagerHandler(context.Background(), mux, cc)
}
