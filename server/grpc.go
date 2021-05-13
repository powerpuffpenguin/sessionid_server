package server

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/sessionid_server/configure"
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func newServer(mux *runtime.ServeMux, cc *grpc.ClientConn, cnf *configure.Auth) (srv *grpc.Server) {
	var opts []grpc.ServerOption

	if cnf.Auth {
		auth := make(serviceAuth)
		for _, rule := range cnf.Rule {
			for _, url := range rule.URL {
				m, ok := auth[url]
				if !ok {
					m = make(map[string]bool)
					auth[url] = m
				}
				for _, token := range rule.Bearer {
					m[token] = true
				}
			}
		}
		opts = append(opts,
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				UnaryServerInterceptor(auth.AuthFunc),
			)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				StreamServerInterceptor(auth.AuthFunc),
			)),
		)
	}
	srv = grpc.NewServer(opts...)
	gmodule.InitServer(srv, mux, cc)
	return
}

type serviceAuth map[string]map[string]bool

func (r serviceAuth) AuthFunc(ctx context.Context, fullMethod string) (context.Context, error) {
	if m, ok := r[fullMethod]; ok {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			strs := md.Get(`authorization`)
			for _, str := range strs {
				if strings.HasPrefix(str, `Bearer `) && m[str[7:]] {
					return ctx, nil
				}
			}
		}
		return ctx, status.Error(codes.PermissionDenied, `permission denied`)
	}
	return ctx, nil
}

// AuthFunc is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type AuthFunc func(ctx context.Context, fullMethod string) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			newCtx, err = authFunc(ctx, info.FullMethod)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			newCtx, err = authFunc(stream.Context(), info.FullMethod)
		}
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
