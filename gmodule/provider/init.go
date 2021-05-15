package provider

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_provider "github.com/powerpuffpenguin/sessionid_server/protocol/provider"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_provider.RegisterProviderServer(srv, server{})
}
func (Module) RegisterGateway(mux *runtime.ServeMux, cc *grpc.ClientConn) {
	grpc_provider.RegisterProviderHandler(context.Background(), mux, cc)
}
