package logger

import (
	"context"

	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	"github.com/powerpuffpenguin/sessionid_server/logger"
	grpc_logger "github.com/powerpuffpenguin/sessionid_server/protocol/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	grpc_logger.UnimplementedLoggerServer
	gmodule.Helper
}

func (server) Level(ctx context.Context, req *grpc_logger.LevelRequest) (resp *grpc_logger.LevelResponse, e error) {
	file, e := logger.Logger.FileLevel().MarshalText()
	if e != nil {
		if ce := logger.Logger.Check(zap.ErrorLevel, `/logger.Logger/Level`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	console, e := logger.Logger.ConsoleLevel().MarshalText()
	if e != nil {
		if ce := logger.Logger.Check(zap.ErrorLevel, `/logger.Logger/Level`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}

	resp = &grpc_logger.LevelResponse{
		File:    string(file),
		Console: string(console),
	}
	return
}

var emptyFileResponse grpc_logger.FileResponse

func (server) File(ctx context.Context, req *grpc_logger.FileRequest) (resp *grpc_logger.FileResponse, e error) {
	lv := logger.Logger.FileLevel()
	e = lv.UnmarshalText([]byte(req.Level))
	if e != nil {
		e = status.Error(codes.InvalidArgument, e.Error())
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/logger.Logger/File`); ce != nil {
		ce.Write(
			zap.String(`level`, req.Level),
		)
	}
	resp = &emptyFileResponse
	return
}

var emptyConsoleResponse grpc_logger.ConsoleResponse

func (server) Console(ctx context.Context, req *grpc_logger.ConsoleRequest) (resp *grpc_logger.ConsoleResponse, e error) {
	lv := logger.Logger.ConsoleLevel()
	e = lv.UnmarshalText([]byte(req.Level))
	if e != nil {
		e = status.Error(codes.InvalidArgument, e.Error())
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `/logger.Logger/Console`); ce != nil {
		ce.Write(
			zap.String(`level`, req.Level),
		)
	}
	resp = &emptyConsoleResponse
	return
}
