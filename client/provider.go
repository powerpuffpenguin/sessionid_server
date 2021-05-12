package client

import (
	"context"

	"github.com/powerpuffpenguin/sessionid"
	grpc_provider "github.com/powerpuffpenguin/sessionid_server/protocol/provider"
	"google.golang.org/grpc"
)

type Provider struct {
	provider grpc_provider.ProviderClient
}

func NewProvider(cc *grpc.ClientConn) *Provider {
	return &Provider{
		provider: grpc_provider.NewProviderClient(cc),
	}
}

// // Create new session
// Create(ctx context.Context,
// 	access, refresh string, // id.sessionid.signature
// 	pair []PairBytes,
// ) (e error)
// // Destroy a session by id
// Destroy(ctx context.Context, id string) (e error)
// // Destroy a session by token
// DestroyByToken(ctx context.Context, token string) (e error)
// Check token status
func (p *Provider) Check(ctx context.Context, token string) (e error) {
	_, e = p.provider.Verify(ctx, &grpc_provider.VerifyRequest{
		Access: token,
	})
	if e != nil {
		e = toError(e)
		return
	}
	return
}

// Put key value for token
func (p *Provider) Put(ctx context.Context, token string, pair []sessionid.PairBytes) (e error) {
	req := &grpc_provider.PutRequest{
		Access: token,
	}
	if len(pair) != 0 {
		req.Pairs = make([]*grpc_provider.Pair, len(pair))
		for i, kv := range pair {
			req.Pairs[i] = &grpc_provider.Pair{
				Key:   kv.Key,
				Value: kv.Value,
			}
		}
	}
	_, e = p.provider.Put(ctx, req)
	if e != nil {
		e = toError(e)
		return
	}
	return
}

// Get key's value from token
func (p *Provider) Get(ctx context.Context, token string, keys []string) (vals []sessionid.Value, e error) {
	resp, e := p.provider.Get(ctx, &grpc_provider.GetRequest{
		Access: token,
		Keys:   keys,
	})
	if e != nil {
		e = toError(e)
		return
	}
	vals = make([]sessionid.Value, len(resp.Value))
	for i, val := range resp.Value {
		vals[i].Bytes = val.Bytes
		vals[i].Exists = val.Exists
	}
	return
}

// Keys return all token's key
func (p *Provider) Keys(ctx context.Context, token string) (keys []string, e error) {
	resp, e := p.provider.Keys(ctx, &grpc_provider.KeysRequest{
		Access: token,
	})
	if e != nil {
		e = toError(e)
		return
	}
	keys = resp.Result
	return
}

// Delete keys
func (p *Provider) Delete(ctx context.Context, token string, keys []string) (e error) {
	_, e = p.provider.RemoveKeys(ctx, &grpc_provider.RemoveKeysRequest{
		Access: token,
		Keys:   keys,
	})
	if e != nil {
		e = toError(e)
		return
	}
	return
}

// Refresh a new access token
func (p *Provider) Refresh(ctx context.Context, access, refresh, newAccess, newRefresh string) (e error) {
	_, e = p.provider.Refresh(ctx, &grpc_provider.RefreshRequest{
		Access:     access,
		Refresh:    refresh,
		NewAccess:  newAccess,
		NewRefresh: newRefresh,
	})
	if e != nil {
		e = toError(e)
		return
	}
	return
}

// Close provider
func (p *Provider) Close() (e error) {
	return
}
func toError(e error) error {
	switch e.Error() {
	case sessionid.ErrTokenExpired.Error():
		return sessionid.ErrTokenExpired
	case sessionid.ErrTokenInvalid.Error():
		return sessionid.ErrTokenInvalid
	case sessionid.ErrTokenNotExists.Error():
		return sessionid.ErrTokenNotExists
	case sessionid.ErrTokenIDNotMatched.Error():
		return sessionid.ErrTokenIDNotMatched
	case sessionid.ErrRefreshTokenNotMatched.Error():
		return sessionid.ErrRefreshTokenNotMatched
	case sessionid.ErrProviderReturnNotMatch.Error():
		return sessionid.ErrProviderReturnNotMatch
	case sessionid.ErrProviderClosed.Error():
		return sessionid.ErrProviderClosed
	case sessionid.ErrNeedsPointer.Error():
		return sessionid.ErrNeedsPointer
	case sessionid.ErrPointerToPointer.Error():
		return sessionid.ErrPointerToPointer
	case sessionid.ErrKeyNotExists.Error():
		return sessionid.ErrKeyNotExists
	}
	return e
}
