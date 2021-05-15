package main

import (
	"github.com/powerpuffpenguin/sessionid_server/gmodule"
	grpc_logger "github.com/powerpuffpenguin/sessionid_server/gmodule/logger"
	grpc_manager "github.com/powerpuffpenguin/sessionid_server/gmodule/manager"
	grpc_provider "github.com/powerpuffpenguin/sessionid_server/gmodule/provider"
)

func registerModule() {
	gmodule.RegisterModule(`manager`, grpc_manager.Module(0))
	gmodule.RegisterModule(`provider`, grpc_provider.Module(0))
	gmodule.RegisterModule(`logger`, grpc_logger.Module(0))
}
