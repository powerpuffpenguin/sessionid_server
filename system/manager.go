package system

import (
	"context"
	"encoding/base64"

	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid/cryptoer"
	"github.com/powerpuffpenguin/sessionid_server/configure"
)

var manager Manager

func DefaultManager() *Manager {
	return &manager
}

type Manager struct {
	method   cryptoer.SigningMethod
	key      []byte
	provider sessionid.Provider
}

func (m *Manager) init(provider sessionid.Provider) {
	cnf := configure.DefaultConfigure().Manager
	m.key = []byte(cnf.Key)
	switch cnf.Method {
	case `HMD5`:
		m.method = cryptoer.SigningMethodHMD5
	case `HS1`:
		m.method = cryptoer.SigningMethodHS1
	case `HS256`:
		m.method = cryptoer.SigningMethodHS256
	case `HS384`:
		m.method = cryptoer.SigningMethodHS384
	case `HS512`:
		m.method = cryptoer.SigningMethodHS512
	default:
		panic(`unknow Manager.Method : ` + cnf.Method)
	}
	m.provider = provider
}
func encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}
func (m *Manager) Create(ctx context.Context,
	id string,
	kv []sessionid.PairBytes,
) (eid, access, refresh string, e error) {
	eid = encode([]byte(id))
	access, refresh, e = sessionid.CreateToken(m.method, m.key, eid)
	if e != nil {
		return
	}
	provider := m.provider
	e = provider.Create(ctx, access, refresh, kv)
	if e != nil {
		return
	}
	return
}

// Destroy a session by id
func (m *Manager) Destroy(ctx context.Context, id string) error {
	return m.provider.Destroy(ctx, encode([]byte(id)))
}

// Destroy a session by token
func (m *Manager) DestroyByToken(ctx context.Context, token string) error {
	return m.provider.DestroyByToken(ctx, token)
}

// Verify token
func (m *Manager) Verify(ctx context.Context, token string) (id string, e error) {
	id, _, signature, e := sessionid.Split(token)
	if e != nil {
		return
	}
	e = m.method.Verify(token[:len(token)-len(signature)-1], signature, m.key)
	if e != nil {
		return
	}
	return
}

// Refresh a new access token
func (m *Manager) Refresh(ctx context.Context, access, refresh string) (newAccess, newRefresh string, e error) {
	id, _, signature, e := sessionid.Split(access)
	if e != nil {
		return
	}
	e = m.method.Verify(access[:len(access)-len(signature)-1], signature, m.key)
	if e != nil {
		return
	}
	id1, _, signature, e := sessionid.Split(refresh)
	if e != nil {
		return
	}
	e = m.method.Verify(refresh[:len(refresh)-len(signature)-1], signature, m.key)
	if e != nil {
		return
	}
	if id != id1 {
		e = sessionid.ErrTokenIDNotMatched
		return
	}

	newAccess, newRefresh, e = sessionid.CreateToken(m.method, m.key, id)
	if e != nil {
		return
	}
	e = m.provider.Refresh(ctx, access, refresh, newAccess, newRefresh)
	return
}
