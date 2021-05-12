package manager

import (
	"context"

	"github.com/powerpuffpenguin/sessionid"
	grpc_manager "github.com/powerpuffpenguin/sessionid_server/protocol/manager"
	"github.com/powerpuffpenguin/sessionid_server/system"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	grpc_manager.UnimplementedManagerServer
}

func (server) Create(ctx context.Context, req *grpc_manager.CreateRequest) (resp *grpc_manager.CreateResponse, e error) {
	kv := make([]sessionid.PairBytes, 0, len(req.Pairs))
	for _, pair := range req.Pairs {
		kv = append(kv, sessionid.PairBytes{
			Key:   pair.Key,
			Value: pair.Value,
		})
	}
	id, access, refresh, e := system.DefaultManager().Create(ctx, req.Id, kv)
	if e != nil {
		return
	}
	resp = &grpc_manager.CreateResponse{
		Id:      id,
		Access:  access,
		Refresh: refresh,
	}
	return
}
func (server) RemoveID(context.Context, *grpc_manager.RemoveIDRequest) (*grpc_manager.RemoveIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveID not implemented")
}
func (server) RemoveAccess(context.Context, *grpc_manager.RemoveAccessRequest) (*grpc_manager.RemoveAccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAccess not implemented")
}
func (server) Verify(context.Context, *grpc_manager.VerifyRequest) (*grpc_manager.VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (server) Refresh(context.Context, *grpc_manager.RefreshRequest) (*grpc_manager.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
