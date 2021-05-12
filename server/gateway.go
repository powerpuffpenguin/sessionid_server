package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func newProxy() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithIncomingHeaderMatcher(headerToGRPC),
		runtime.WithOutgoingHeaderMatcher(grpcToHeader),
	)
}
func headerToGRPC(key string) (string, bool) {
	switch key {
	case `Connection`:
		fallthrough
	case `Content-Type`:
		fallthrough
	case `Content-Length`:
		return key, false
	default:
		return key, true
	}
}
func grpcToHeader(key string) (string, bool) {
	switch key {
	case `x-http-code`:
		return key, false
	default:
		return key, true
	}
}
func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
	}
	return nil
}
