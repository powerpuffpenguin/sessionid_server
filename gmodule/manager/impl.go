package manager

import (
	"context"

	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	"github.com/powerpuffpenguin/sessionid_server/logger"
	grpc_manager "github.com/powerpuffpenguin/sessionid_server/protocol/manager"
	"github.com/powerpuffpenguin/sessionid_server/system"
	"go.uber.org/zap"
)

type server struct {
	grpc_manager.UnimplementedManagerServer
	gmodule.Helper
}

func (s server) Create(ctx context.Context, req *grpc_manager.CreateRequest) (resp *grpc_manager.CreateResponse, e error) {
	kv := make([]sessionid.PairBytes, 0, len(req.Pairs))
	for _, pair := range req.Pairs {
		kv = append(kv, sessionid.PairBytes{
			Key:   pair.Key,
			Value: pair.Value,
		})
	}
	id, access, refresh, e := system.DefaultManager().Create(ctx, req.Id, kv)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/manager.Manager/Create`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/manager.Manager/Create`); ce != nil {
		ce.Write(
			zap.String(`id`, req.Id),
			zap.String(`access`, access),
			zap.String(`refresh`, refresh),
		)
	}
	resp = &grpc_manager.CreateResponse{
		Id:      id,
		Access:  access,
		Refresh: refresh,
	}
	return
}

var emptyRemoveIDResponse grpc_manager.RemoveIDResponse

func (s server) RemoveID(ctx context.Context, req *grpc_manager.RemoveIDRequest) (resp *grpc_manager.RemoveIDResponse, e error) {
	e = system.DefaultManager().Destroy(ctx, req.Id)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/manager.Manager/RemoveID`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`id`, req.Id),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/manager.Manager/RemoveID`); ce != nil {
		ce.Write(
			zap.String(`id`, req.Id),
		)
	}
	resp = &emptyRemoveIDResponse
	return
}

var emptyRemoveAccessResponse grpc_manager.RemoveAccessResponse

func (server) RemoveAccess(ctx context.Context, req *grpc_manager.RemoveAccessRequest) (resp *grpc_manager.RemoveAccessResponse, e error) {
	e = system.DefaultManager().DestroyByToken(ctx, req.Access)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/manager.Manager/RemoveAccess`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/manager.Manager/RemoveAccess`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
		)
	}
	resp = &emptyRemoveAccessResponse
	return
}
func (s server) Verify(ctx context.Context, req *grpc_manager.VerifyRequest) (resp *grpc_manager.VerifyResponse, e error) {
	id, e := system.DefaultManager().Verify(ctx, req.Access)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/manager.Manager/Verify`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		e = s.ToError(e)
		return
	}
	resp = &grpc_manager.VerifyResponse{
		Id: id,
	}
	return
}
func (s server) Refresh(ctx context.Context, req *grpc_manager.RefreshRequest) (resp *grpc_manager.RefreshResponse, e error) {
	access, refresh, e := system.DefaultManager().Refresh(ctx,
		req.Access, req.Refresh,
	)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/manager.Manager/Refresh`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
				zap.String(`refresh`, req.Refresh),
				zap.String(`new access`, access),
				zap.String(`new refresh`, refresh),
			)
		}
		e = s.ToError(e)
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/manager.Manager/Refresh`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
			zap.String(`refresh`, req.Refresh),
			zap.String(`new access`, access),
			zap.String(`new refresh`, refresh),
		)
	}
	resp = &grpc_manager.RefreshResponse{
		Access:  access,
		Refresh: refresh,
	}
	return
}
