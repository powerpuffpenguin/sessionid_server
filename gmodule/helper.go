package gmodule

import (
	"errors"

	"github.com/powerpuffpenguin/sessionid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Helper struct {
}

func (Helper) ToError(e error) error {
	if errors.Is(e, sessionid.ErrTokenExpired) {
		e = status.Error(codes.Unauthenticated, e.Error())
	}
	return e
}
