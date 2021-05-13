package client

import (
	"context"

	"github.com/powerpuffpenguin/sessionid"
	grpc_manager "github.com/powerpuffpenguin/sessionid_server/protocol/manager"
	"google.golang.org/grpc"
)

type Manager struct {
	manager  grpc_manager.ManagerClient
	provider sessionid.Provider
	coder    sessionid.Coder
}

func NewManager(cc *grpc.ClientConn, coder sessionid.Coder) *Manager {
	return &Manager{
		manager:  grpc_manager.NewManagerClient(cc),
		provider: NewProvider(cc),
		coder:    coder,
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
	session, e = sessionid.NewSession(resp.Access, m.provider, m.coder)
	if e != nil {
		return
	}
	return
}

// Destroy a session by id
func (m *Manager) Destroy(ctx context.Context, id string) (e error) {
	_, e = m.manager.RemoveID(ctx, &grpc_manager.RemoveIDRequest{
		Id: id,
	})
	if e != nil {
		e = toError(e)
		return
	}
	return
}

// Destroy a session by token
func (m *Manager) DestroyByToken(ctx context.Context, token string) (e error) {
	_, e = m.manager.RemoveAccess(ctx, &grpc_manager.RemoveAccessRequest{
		Access: token,
	})
	if e != nil {
		e = toError(e)
		return
	}
	return
}

// Get session from token
func (m *Manager) Get(ctx context.Context, token string) (s *sessionid.Session, e error) {
	_, e = m.manager.Verify(ctx, &grpc_manager.VerifyRequest{
		Access: token,
	})
	if e != nil {
		e = toError(e)
		return
	}
	s, e = sessionid.NewSession(token, m.provider, m.coder)
	return
}

// Refresh a new access token
func (m *Manager) Refresh(ctx context.Context, access, refresh string) (newAccess, newRefresh string, e error) {
	resp, e := m.manager.Refresh(ctx, &grpc_manager.RefreshRequest{
		Access:  access,
		Refresh: refresh,
	})
	if e != nil {
		e = toError(e)
		return
	}
	newAccess, newRefresh = resp.Access, resp.Refresh
	return
}
