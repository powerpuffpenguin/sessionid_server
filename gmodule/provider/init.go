package provider

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	grpc_provider "github.com/powerpuffpenguin/sessionid_server/protocol/provider"
	"google.golang.org/grpc"
)

func init() {
	gmodule.RegisterModule(`provider`, Module(0))
}

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_provider.RegisterProviderServer(srv, server{})
}
func (Module) RegisterGateway(mux *runtime.ServeMux, cc *grpc.ClientConn) {
	grpc_provider.RegisterProviderHandler(context.Background(), mux, cc)
}
