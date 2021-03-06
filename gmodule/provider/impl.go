package provider

import (
	"context"

	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	"github.com/powerpuffpenguin/sessionid_server/logger"
	grpc_provider "github.com/powerpuffpenguin/sessionid_server/protocol/provider"
	"github.com/powerpuffpenguin/sessionid_server/system"
	"go.uber.org/zap"
)

type server struct {
	grpc_provider.UnimplementedProviderServer
	gmodule.Helper
}

var emptyCreateResponse grpc_provider.CreateResponse

func (server) Create(ctx context.Context, req *grpc_provider.CreateRequest) (resp *grpc_provider.CreateResponse, e error) {
	kv := make([]sessionid.PairBytes, 0, len(req.Pairs))
	for _, pair := range req.Pairs {
		kv = append(kv, sessionid.PairBytes{
			Key:   pair.Key,
			Value: pair.Value,
		})
	}
	e = system.DefaultProvider().Create(ctx,
		req.Access, req.Refresh,
		kv,
	)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/Create`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
				zap.String(`refresh`, req.Refresh),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/provider.Provider/Create`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
			zap.String(`refresh`, req.Refresh),
		)
	}
	resp = &emptyCreateResponse
	return
}

var emptyRemoveIDResponse grpc_provider.RemoveIDResponse

func (server) RemoveID(ctx context.Context, req *grpc_provider.RemoveIDRequest) (resp *grpc_provider.RemoveIDResponse, e error) {
	e = system.DefaultProvider().Destroy(ctx, req.Id)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/RemoveID`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`id`, req.Id),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/provider.Provider/RemoveID`); ce != nil {
		ce.Write(
			zap.String(`id`, req.Id),
		)
	}
	resp = &emptyRemoveIDResponse
	return
}

var emptyRemoveAccessResponse grpc_provider.RemoveAccessResponse

func (server) RemoveAccess(ctx context.Context, req *grpc_provider.RemoveAccessRequest) (resp *grpc_provider.RemoveAccessResponse, e error) {
	e = system.DefaultProvider().DestroyByToken(ctx, req.Access)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/RemoveAccess`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/provider.Provider/RemoveAccess`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
		)
	}
	resp = &emptyRemoveAccessResponse
	return
}

var emptyVerifyResponse grpc_provider.VerifyResponse

func (s server) Verify(ctx context.Context, req *grpc_provider.VerifyRequest) (resp *grpc_provider.VerifyResponse, e error) {
	e = system.DefaultProvider().Check(ctx, req.Access)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/Verify`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		e = s.ToError(e)
		return
	}
	resp = &emptyVerifyResponse
	return
}

var emptyPutResponse grpc_provider.PutResponse

func (s server) Put(ctx context.Context, req *grpc_provider.PutRequest) (resp *grpc_provider.PutResponse, e error) {
	kv := make([]sessionid.PairBytes, 0, len(req.Pairs))
	for _, pair := range req.Pairs {
		kv = append(kv, sessionid.PairBytes{
			Key:   pair.Key,
			Value: pair.Value,
		})
	}
	e = system.DefaultProvider().Put(ctx, req.Access, kv)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/Put`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		e = s.ToError(e)
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/provider.Provider/Put`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
		)
	}
	resp = &emptyPutResponse
	return
}
func (s server) Get(ctx context.Context, req *grpc_provider.GetRequest) (resp *grpc_provider.GetResponse, e error) {
	result, e := system.DefaultProvider().Get(ctx, req.Access, req.Keys)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/Get`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		e = s.ToError(e)
		return
	}
	resp = &grpc_provider.GetResponse{}
	if len(result) != 0 {
		resp.Value = make([]*grpc_provider.Value, len(result))
		for i, v := range result {
			resp.Value[i] = &grpc_provider.Value{
				Bytes:  v.Bytes,
				Exists: v.Exists,
			}
		}
	}
	return
}
func (s server) Keys(ctx context.Context, req *grpc_provider.KeysRequest) (resp *grpc_provider.KeysResponse, e error) {
	keys, e := system.DefaultProvider().Keys(ctx, req.Access)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/Keys`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		e = s.ToError(e)
		return
	}
	resp = &grpc_provider.KeysResponse{
		Result: keys,
	}
	return
}

var emptyRemoveKeysResponse grpc_provider.RemoveKeysResponse

func (server) RemoveKeys(ctx context.Context, req *grpc_provider.RemoveKeysRequest) (resp *grpc_provider.RemoveKeysResponse, e error) {
	e = system.DefaultProvider().Delete(ctx, req.Access, req.Keys)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/RemoveKeys`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/provider.Provider/RemoveKeys`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
		)
	}
	resp = &emptyRemoveKeysResponse
	return
}

var emptyRefreshResponse grpc_provider.RefreshResponse

func (s server) Refresh(ctx context.Context, req *grpc_provider.RefreshRequest) (resp *grpc_provider.RefreshResponse, e error) {
	e = system.DefaultProvider().Refresh(ctx,
		req.Access, req.Refresh,
		req.NewAccess, req.NewRefresh,
	)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `/provider.Provider/Refresh`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
				zap.String(`refresh`, req.Refresh),
				zap.String(`new access`, req.NewAccess),
				zap.String(`new refresh`, req.NewRefresh),
			)
		}
		e = s.ToError(e)
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/provider.Provider/Refresh`); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
			zap.String(`refresh`, req.Refresh),
			zap.String(`new access`, req.NewAccess),
			zap.String(`new refresh`, req.NewRefresh),
		)
	}
	resp = &emptyRefreshResponse
	return
}
