package system

import (
	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid/provider/bolt"
	"github.com/powerpuffpenguin/sessionid/provider/redis"
	"github.com/powerpuffpenguin/sessionid_server/configure"
	"github.com/powerpuffpenguin/sessionid_server/logger"
	"go.uber.org/zap"
)

var provider sessionid.Provider

func DefaultProvider() sessionid.Provider {
	return provider
}
func initProvider() {
	cnf := configure.DefaultConfigure().Provider
	switch cnf.Backend {
	case `Memory`:
		initProviderMemory(&cnf.Memory)
	case `Redis`:
		initProviderRedis(&cnf.Redis)
	case `Bolt`:
		initProviderBolt(&cnf.Bolt)
	default:
		panic(`unknow Provider.Backend : ` + cnf.Backend)
	}
}
func initProviderMemory(cnf *configure.ProviderMemory) {
	logger.Logger.Info(`memory provider`,
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Int(`max size`, cnf.MaxSize),
		zap.Int(`check batch`, cnf.Batch),
		zap.Duration(`clear`, cnf.Clear),
	)
	provider = sessionid.NewProvider(
		sessionid.WithProviderAccess(cnf.Access),
		sessionid.WithProviderRefresh(cnf.Refresh),
		sessionid.WithProviderMaxSize(cnf.MaxSize),
		sessionid.WithProviderCheckBatch(cnf.Batch),
		sessionid.WithProviderClear(cnf.Clear),
	)
}
func initProviderRedis(cnf *configure.ProviderRedis) {
	logger.Logger.Info(`redis provider`,
		zap.String(`url`, cnf.URL),
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Int(`check batch`, cnf.Batch),
		zap.String(`key prefix`, cnf.KeyPrefix),
		zap.String(`metadata key`, cnf.MetadataKey),
	)
	var e error
	provider, e = redis.New(
		redis.WithURL(cnf.URL),
		redis.WithAccess(cnf.Access),
		redis.WithRefresh(cnf.Refresh),
		redis.WithCheckBatch(cnf.Batch),
		redis.WithKeyPrefix(cnf.KeyPrefix),
		redis.WithMetadataKey(cnf.MetadataKey),
	)
	if e != nil {
		logger.Logger.Fatal(`redis provider`,
			zap.Error(e),
			zap.String(`url`, cnf.URL),
			zap.Duration(`access`, cnf.Access),
			zap.Duration(`refresh`, cnf.Refresh),
			zap.Int(`check batch`, cnf.Batch),
			zap.String(`key prefix`, cnf.KeyPrefix),
			zap.String(`metadata key`, cnf.MetadataKey),
		)
	}
}
func initProviderBolt(cnf *configure.ProviderBolt) {
	logger.Logger.Info(`bolt provider`,
		zap.String(`filename`, cnf.Filename),
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Int(`max size`, cnf.MaxSize),
		zap.Int(`check batch`, cnf.Batch),
		zap.Duration(`clear`, cnf.Clear),
	)
	var e error
	provider, e = bolt.New(
		bolt.WithFilename(cnf.Filename),
		bolt.WithAccess(cnf.Access),
		bolt.WithRefresh(cnf.Refresh),
		bolt.WithMaxSize(cnf.MaxSize),
		bolt.WithCheckBatch(cnf.Batch),
		bolt.WithClear(cnf.Clear),
	)
	if e != nil {
		logger.Logger.Fatal(`bolt provider`,
			zap.Error(e),
			zap.String(`filename`, cnf.Filename),
			zap.Duration(`access`, cnf.Access),
			zap.Duration(`refresh`, cnf.Refresh),
			zap.Int(`max size`, cnf.MaxSize),
			zap.Int(`check batch`, cnf.Batch),
			zap.Duration(`clear`, cnf.Clear),
		)
	}
}
