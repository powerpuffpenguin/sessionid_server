package client

import (
	"context"
	"fmt"

	"github.com/powerpuffpenguin/sessionid"
	grpc_manager "github.com/powerpuffpenguin/sessionid_server/protocol/manager"
	grpc_provider "github.com/powerpuffpenguin/sessionid_server/protocol/provider"
	"google.golang.org/grpc"
)

type Manager struct {
	manager  grpc_manager.ManagerClient
	provider grpc_provider.ProviderClient
	coder    sessionid.Coder
}

func NewManager(cc *grpc.ClientConn, coder sessionid.Coder) *Manager {
	return &Manager{
		manager: grpc_manager.NewManagerClient(cc),
		coder:   coder,
	}
}

// Create a session for the user
//
// * id uid or web-uid or mobile-uid or ...
//
// * pair session init key value
func (m *Manager) Create(ctx context.Context,
	id string,
	pair ...sessionid.Pair,
) (session *sessionid.Session, refresh string, e error) {
	req := &grpc_manager.CreateRequest{
		Id: id,
	}
	count := len(pair)
	if count != 0 {
		kv := make([]*grpc_manager.Pair, count)
		var v []byte
		for i := 0; i < count; i++ {
			v, e = m.coder.Encode(pair[i].Value)
			if e != nil {
				return
			}
			kv[i] = &grpc_manager.Pair{
				Key:   pair[i].Key,
				Value: v,
			}
		}
		req.Pairs = kv
	}
	resp, e := m.manager.Create(ctx, req)
	if e != nil {
		return
	}
	fmt.Println(
		resp.Id,
		resp.Access,
	)

	// eid := encode([]byte(id))
	// access, refresh, e := CreateToken(m.opts.method, m.opts.key, eid)
	// if e != nil {
	// 	return
	// }
	// opts := m.opts
	// provider := opts.provider
	// coder := opts.coder
	// var kv []PairBytes

	// e = provider.Create(ctx, access, refresh, kv)
	// if e != nil {
	// 	return
	// }
	// session = newSession(eid, access, provider, coder)
	return
}

// // Destroy a session by id
// func (m *LocalManager) Destroy(ctx context.Context, id string) error {
// 	return m.opts.provider.Destroy(ctx, encode([]byte(id)))
// }

// // Destroy a session by token
// func (m *LocalManager) DestroyByToken(ctx context.Context, token string) error {
// 	return m.opts.provider.DestroyByToken(ctx, token)
// }

// // Get session from token
// func (m *LocalManager) Get(ctx context.Context, token string) (s *Session, e error) {
// 	id, _, signature, e := Split(token)
// 	if e != nil {
// 		return
// 	}
// 	e = m.opts.method.Verify(token[:len(token)-len(signature)-1], signature, m.opts.key)
// 	if e != nil {
// 		return
// 	}
// 	s = newSession(id, token, m.opts.provider, m.opts.coder)
// 	return
// }

// // Refresh a new access token
// func (m *LocalManager) Refresh(ctx context.Context, access, refresh string) (newAccess, newRefresh string, e error) {
// 	id, _, signature, e := Split(access)
// 	if e != nil {
// 		return
// 	}
// 	e = m.opts.method.Verify(access[:len(access)-len(signature)-1], signature, m.opts.key)
// 	if e != nil {
// 		return
// 	}
// 	id1, _, signature, e := Split(refresh)
// 	if e != nil {
// 		return
// 	}
// 	e = m.opts.method.Verify(refresh[:len(refresh)-len(signature)-1], signature, m.opts.key)
// 	if e != nil {
// 		return
// 	}
// 	if id != id1 {
// 		e = ErrTokenIDNotMatched
// 		return
// 	}

// 	newAccess, newRefresh, e = CreateToken(m.opts.method, m.opts.key, id)
// 	if e != nil {
// 		return
// 	}
// 	e = m.opts.provider.Refresh(ctx, access, refresh, newAccess, newRefresh)
// 	return
// }
