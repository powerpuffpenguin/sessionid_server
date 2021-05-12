package gmodule

import (
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module interface {
	RegisterGRPC(*grpc.Server)
	RegisterGateway(mux *runtime.ServeMux, cc *grpc.ClientConn)
}

var keys = make(map[string]Module)

func RegisterModule(id string, m Module) {
	if _, ok := keys[id]; ok {
		panic(`module id already exists : ` + id)
	}
	log.Println(`register module :`, id)
	keys[id] = m
}
func InitServer(srv *grpc.Server, mux *runtime.ServeMux, cc *grpc.ClientConn) {
	for _, m := range keys {
		m.RegisterGRPC(srv)
		if mux != nil {
			m.RegisterGateway(mux, cc)
		}
	}
}
