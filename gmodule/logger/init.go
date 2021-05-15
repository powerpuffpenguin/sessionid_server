package logger

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_logger "github.com/powerpuffpenguin/sessionid_server/protocol/logger"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_logger.RegisterLoggerServer(srv, server{})
}
func (Module) RegisterGateway(mux *runtime.ServeMux, cc *grpc.ClientConn) {
	grpc_logger.RegisterLoggerHandler(context.Background(), mux, cc)
}
